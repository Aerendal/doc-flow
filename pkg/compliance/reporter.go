package compliance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"docflow/pkg/governance"
	"docflow/pkg/index"
	"docflow/pkg/sections"
)

type DocResult struct {
	DocID      string            `json:"doc_id"`
	Path       string            `json:"path"`
	Status     string            `json:"status"`
	DocType    string            `json:"doc_type"`
	Violations []string          `json:"violations"`
	Details    map[string]string `json:"details,omitempty"`
}

type Summary struct {
	SchemaVersion   string              `json:"schema_version"`
	IdentityVersion string              `json:"identity_version,omitempty"`
	RulesPath       string              `json:"rules_path,omitempty"`
	RulesChecksum   string              `json:"rules_checksum,omitempty"`
	FailOn          string              `json:"fail_on,omitempty"`
	Show            string              `json:"show,omitempty"`
	Baseline        *BaselineMeta       `json:"baseline,omitempty"`
	Documents       int                 `json:"documents"`
	Passed          int                 `json:"passed"`
	Failed          int                 `json:"failed"`
	PassRate        float64             `json:"pass_rate"`
	ViolationsCount map[string]int      `json:"violations_count"`
	NewFailed       int                 `json:"new_failed,omitempty"`
	ExistingFailed  int                 `json:"existing_failed,omitempty"`
	NewViolations   map[string]int      `json:"new_violations_count,omitempty"`
	DuplicateDocIDs map[string][]string `json:"duplicate_doc_ids,omitempty"`
	Docs            []DocResult         `json:"docs"`
}

type BaselineMeta struct {
	Path            string `json:"path"`
	SchemaVersion   string `json:"schema_version,omitempty"`
	IdentityVersion string `json:"identity_version,omitempty"`
	Loaded          bool   `json:"loaded"`
}

type ContentFact struct {
	Content   string
	ReadError string
}

// Report runs governance validation on all documents in index.
func Report(idx *index.DocumentIndex, rules *governance.Rules) (*Summary, error) {
	return ReportWithFacts(idx, rules, nil)
}

// ReportWithFacts runs governance validation using preloaded per-path content.
// If facts is nil or a path is missing, it falls back to reading from disk.
func ReportWithFacts(idx *index.DocumentIndex, rules *governance.Rules, facts map[string]ContentFact) (*Summary, error) {
	if rules == nil {
		return nil, fmt.Errorf("rules are nil")
	}
	sum := &Summary{
		SchemaVersion:   "1.0",
		IdentityVersion: "2",
		ViolationsCount: map[string]int{},
	}

	for _, rec := range idx.All() {
		sum.Documents++
		res := DocResult{
			DocID:   rec.DocID,
			Path:    rec.Path,
			Status:  rec.Status,
			DocType: rec.DocType,
		}

		content := ""
		if fact, ok := facts[rec.Path]; ok {
			if fact.ReadError != "" {
				res.Violations = append(res.Violations, "read_error")
				res.Details = map[string]string{"error": fact.ReadError}
				sum.Docs = append(sum.Docs, res)
				sum.ViolationsCount["read_error"]++
				continue
			}
			content = fact.Content
		} else {
			readPath := rec.Path
			if !filepath.IsAbs(readPath) && idx.Root != "" {
				readPath = filepath.Join(idx.Root, rec.Path)
			}

			data, err := os.ReadFile(readPath)
			if err != nil {
				res.Violations = append(res.Violations, "read_error")
				res.Details = map[string]string{"error": err.Error()}
				sum.Docs = append(sum.Docs, res)
				sum.ViolationsCount["read_error"]++
				continue
			}
			content = string(data)
		}

		violations := validateAgainstRules(content, rec, rules)
		if len(violations) == 0 {
			sum.Passed++
		} else {
			sum.Failed++
			for _, v := range violations {
				sum.ViolationsCount[v]++
			}
		}
		res.Violations = violations
		sum.Docs = append(sum.Docs, res)
	}

	if sum.Documents > 0 {
		sum.PassRate = float64(sum.Passed) / float64(sum.Documents)
	}
	finalizeSummary(sum)
	return sum, nil
}

func finalizeSummary(sum *Summary) {
	for i := range sum.Docs {
		if sum.Docs[i].Violations == nil {
			sum.Docs[i].Violations = []string{}
		}
		sort.Strings(sum.Docs[i].Violations)
	}
	sort.Slice(sum.Docs, func(i, j int) bool {
		return sum.Docs[i].Path < sum.Docs[j].Path
	})

	byDocID := map[string][]string{}
	for _, d := range sum.Docs {
		if d.DocID == "" {
			continue
		}
		byDocID[d.DocID] = append(byDocID[d.DocID], d.Path)
	}

	duplicates := map[string][]string{}
	for docID, paths := range byDocID {
		if len(paths) <= 1 {
			continue
		}
		sort.Strings(paths)
		duplicates[docID] = paths
	}
	if len(duplicates) > 0 {
		sum.DuplicateDocIDs = duplicates
	}
}

func validateAgainstRules(content string, rec *index.DocumentRecord, rules *governance.Rules) []string {
	var out []string

	// status rules
	sr, ok := rules.Statuses[rec.Status]
	if !ok {
		sr = rules.Statuses["default"]
	}
	if sr.RequiredFields != nil {
		for _, f := range sr.RequiredFields {
			if !hasField(f, content) {
				out = append(out, "missing_"+f)
			}
		}
	}
	if sr.MinQuality > 0 && rec.Checksum == "" {
		// placeholder; cannot compute quality
	}
	if !sr.AllowEmptySections {
		metrics := sections.ComputeMetrics(content)
		if metrics.Empty > 0 {
			out = append(out, "empty_sections")
		}
	}

	// family rules
	if rules != nil {
		if fr, ok := rules.Families[rec.DocType]; ok {
			tree := sections.Parse(content)
			found := map[string]bool{}
			var walk func(n *sections.SectionNode)
			walk = func(n *sections.SectionNode) {
				if n.Level > 0 {
					found[strings.ToLower(strings.TrimSpace(n.Text))] = true
				}
				for _, c := range n.Children {
					walk(c)
				}
			}
			walk(tree.Root)
			for _, rs := range fr.RequiredSections {
				if !found[strings.ToLower(rs)] {
					out = append(out, "missing_section_"+rs)
				}
			}
		}
	}
	return out
}

func hasField(field, content string) bool {
	pat := field + ":"
	return strings.Contains(content, pat)
}

// SaveJSON writes summary to path.
func SaveJSON(sum *Summary, path string) error {
	data, err := json.MarshalIndent(sum, "", "  ")
	if err != nil {
		return err
	}
	if path == "-" {
		_, err = os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// SaveHTML renders a minimal HTML dashboard (no external assets).
func SaveHTML(sum *Summary, path string) error {
	var b strings.Builder
	bar := func(label string, val, total int) string {
		width := 0
		if total > 0 {
			width = int(float64(val) / float64(total) * 100)
		}
		return fmt.Sprintf(`<div style="margin:4px 0;">%s: %d<div style="background:#eee;width:300px;"><div style="background:#4caf50;width:%d%%;height:10px;"></div></div></div>`, label, val, width)
	}

	b.WriteString("<html><head><meta charset=\"utf-8\"><title>Compliance Report</title>")
	b.WriteString(`<style>body{font-family:Arial, sans-serif;max-width:900px;margin:20px auto;}table{border-collapse:collapse;width:100%;}th,td{border:1px solid #ddd;padding:6px;text-align:left;}th{background:#f5f5f5;}code{background:#f2f2f2;padding:2px 4px;}</style></head><body>`)
	b.WriteString("<h1>Compliance Report</h1>")
	b.WriteString(fmt.Sprintf("<p>Documents: %d | Passed: %d | Failed: %d | Pass rate: %.1f%%</p>", sum.Documents, sum.Passed, sum.Failed, sum.PassRate*100))
	b.WriteString(bar("Passed", sum.Passed, sum.Documents))
	b.WriteString(bar("Failed", sum.Failed, sum.Documents))

	if len(sum.ViolationsCount) > 0 {
		b.WriteString("<h2>Violations by type</h2><table><tr><th>Type</th><th>Count</th></tr>")
		for k, v := range sum.ViolationsCount {
			b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", k, v))
		}
		b.WriteString("</table>")
	}

	if sum.Failed > 0 {
		b.WriteString("<h2>Non-compliant documents</h2><table><tr><th>Path</th><th>DocID</th><th>Status</th><th>DocType</th><th>Violations</th></tr>")
		for _, d := range sum.Docs {
			if len(d.Violations) == 0 {
				continue
			}
			b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				d.Path, d.DocID, d.Status, d.DocType, strings.Join(d.Violations, ", ")))
		}
		b.WriteString("</table>")
	}

	b.WriteString("</body></html>")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
