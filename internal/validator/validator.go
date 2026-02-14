package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"docflow/internal/util"
	"docflow/pkg/governance"
	"docflow/pkg/index"
	"docflow/pkg/parser"
	"docflow/pkg/sections"
)

// IssueLevel represents severity of a validation finding.
type IssueLevel string

const (
	LevelError IssueLevel = "error"
	LevelWarn  IssueLevel = "warn"
)

// IssueType enumerates validation rule identifiers.
type IssueType string

const (
	IssueMissingFrontmatter IssueType = "missing_frontmatter"
	IssueInvalidYAML        IssueType = "invalid_yaml"
	IssueMissingField       IssueType = "missing_field"
	IssueEmptyField         IssueType = "empty_field"
	IssueInvalidDocID       IssueType = "invalid_doc_id"
	IssueDuplicateDocID     IssueType = "duplicate_doc_id"
	IssueDocIDFilename      IssueType = "doc_id_filename_mismatch"
	IssueMissingContext     IssueType = "missing_context_sources"
	IssueMissingExpectedDep IssueType = "missing_expected_dependency"
	IssueLegacySectionName  IssueType = "legacy_section_name"
	IssueGovernance         IssueType = "governance_violation"
	IssueCycleDetected      IssueType = "cycle_detected"
)

// Issue describes a single validation problem.
type Issue struct {
	Level   IssueLevel
	Type    IssueType
	File    string
	DocID   string
	Line    int
	Message string
}

// Report aggregates findings for a validation run.
type Report struct {
	Issues    []Issue
	Files     int
	Documents int
}

func (r *Report) add(issue Issue) {
	r.Issues = append(r.Issues, issue)
}

// ErrorCount returns number of error-level issues.
func (r *Report) ErrorCount() int {
	count := 0
	for _, i := range r.Issues {
		if i.Level == LevelError {
			count++
		}
	}
	return count
}

// WarnCount returns number of warning-level issues.
func (r *Report) WarnCount() int {
	count := 0
	for _, i := range r.Issues {
		if i.Level == LevelWarn {
			count++
		}
	}
	return count
}

// HasErrors reports whether any error-level issues exist.
func (r *Report) HasErrors() bool {
	return r.ErrorCount() > 0
}

// SortedIssues returns a stable, file-ordered slice of issues.
func (r *Report) SortedIssues() []Issue {
	out := make([]Issue, len(r.Issues))
	copy(out, r.Issues)
	sort.Slice(out, func(i, j int) bool {
		if out[i].File == out[j].File {
			if out[i].Line == out[j].Line {
				return out[i].Type < out[j].Type
			}
			return out[i].Line < out[j].Line
		}
		return out[i].File < out[j].File
	})
	return out
}

var docIDPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,80}$`)

type docInfo struct {
	DocID string
	File  string
}

func missingDeps(expected, declared []string) []string {
	decl := make(map[string]bool)
	for _, d := range declared {
		decl[d] = true
	}
	var miss []string
	for _, e := range expected {
		if !decl[e] {
			miss = append(miss, e)
		}
	}
	return miss
}

func levelFor(strict bool) IssueLevel {
	if strict {
		return LevelError
	}
	return LevelWarn
}

// reject traversal elements like ".." to keep deps/context inside docs root namespace.
func containsTraversal(items []string) (string, bool) {
	for _, it := range items {
		if strings.Contains(it, "..") {
			return it, true
		}
	}
	return "", false
}

type Options struct {
	PromoteContextFor map[string]bool
	RequiredDeps      map[string][]string
	PublishStrict     bool
	SectionAliases    map[string][]string
	StrictMode        bool
	StatusAware       bool
	GovernanceRules   *governance.Rules
}

// ValidateDocs walks Markdown files and validates metadata & identifiers.
func ValidateDocs(docsRoot string, ignorePatterns []string, opts *Options) (*Report, error) {
	if opts == nil {
		opts = &Options{PromoteContextFor: map[string]bool{}, SectionAliases: map[string][]string{}, RequiredDeps: map[string][]string{}}
	}
	if opts.SectionAliases == nil {
		opts.SectionAliases = map[string][]string{}
	}

	idx, err := index.BuildIndex(docsRoot, ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("błąd skanowania: %w", err)
	}

	files, err := util.WalkMarkdown(docsRoot, ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("błąd skanowania: %w", err)
	}

	report := &Report{}
	var docs []docInfo

	// Cycles from index
	for _, rec := range idx.All() {
		if len(rec.Cycle) > 0 {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueCycleDetected,
				File:    rec.Path,
				DocID:   rec.DocID,
				Message: fmt.Sprintf("wykryto cykl: %s", strings.Join(rec.Cycle, " -> ")),
			})
		}
	}

	for _, f := range files {
		report.Files++
		rel := filepath.ToSlash(f.RelPath)

		data, err := os.ReadFile(f.Path)
		if err != nil {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueMissingField,
				File:    rel,
				Message: fmt.Sprintf("nie można odczytać pliku: %v", err),
			})
			continue
		}

		content := string(data)
		res, err := parser.ParseFrontmatterString(content, rel)
		if err != nil {
			if perr, ok := err.(*parser.ParseError); ok {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueInvalidYAML,
					File:    rel,
					Line:    perr.Line,
					Message: perr.Message,
				})
			} else {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueInvalidYAML,
					File:    rel,
					Message: err.Error(),
				})
			}
			continue
		}

		if res == nil || !res.HasFrontmatter {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueMissingFrontmatter,
				File:    rel,
				Message: "brak frontmatter YAML",
			})
			continue
		}

		fm := res.Frontmatter
		raw := res.RawFields
		if raw == nil {
			raw = map[string]bool{}
		}
		docID := strings.TrimSpace(fm.DocID)
		if docID == "" {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, Message: "brak doc_id"})
			docID = index.NormalizeDocID(rel)
		}

		if strings.TrimSpace(fm.Title) == "" || !raw["title"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak title"})
		}
		docType := strings.TrimSpace(fm.DocType)
		if docType == "" || !raw["doc_type"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak doc_type"})
		}
		if strings.TrimSpace(fm.Status) == "" || !raw["status"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak status"})
		}
		if strings.TrimSpace(fm.Version) == "" || !raw["version"] {
			report.add(Issue{Level: LevelError, Type: IssueMissingField, File: rel, DocID: docID, Message: "brak version"})
		}

		if fm.DocID != "" && !docIDPattern.MatchString(docID) {
			report.add(Issue{Level: LevelError, Type: IssueInvalidDocID, File: rel, DocID: docID, Message: "doc_id nie spełnia konwencji snake_case"})
		}

		expectedID := index.NormalizeDocID(rel)
		if fm.DocID != "" && fm.DocID != expectedID {
			report.add(Issue{Level: LevelWarn, Type: IssueDocIDFilename, File: rel, DocID: fm.DocID, Message: fmt.Sprintf("doc_id=%s nie zgadza się z nazwą pliku (%s)", fm.DocID, expectedID)})
		}

		if fm.DocID != "" {
			docs = append(docs, docInfo{DocID: fm.DocID, File: rel})
			report.Documents++
		}

		statusLower := strings.ToLower(fm.Status)
		relaxed := statusLower == "archived" || statusLower == "deprecated" || statusLower == "generated"
		emptyStrict, contextStrict := statusProfile(statusLower)
		if !opts.StatusAware {
			emptyStrict = opts.PublishStrict && statusLower == "published"
			contextStrict = opts.StrictMode || opts.PromoteContextFor[docType]
		}

		// context_sources check
		if len(fm.ContextSources) == 0 {
			level := LevelWarn
			if contextStrict {
				level = LevelError
			}
			if relaxed {
				level = LevelWarn
			}
			report.add(Issue{Level: level, Type: IssueMissingContext, File: rel, DocID: docID, Message: "brak context_sources"})
		} else if bad, ok := containsTraversal(fm.ContextSources); ok {
			report.add(Issue{Level: LevelError, Type: IssueInvalidDocID, File: rel, DocID: docID, Message: fmt.Sprintf("context_sources zawiera traversal: %s", bad)})
		}

		// expected deps per doc_type
		if reqDeps, ok := opts.RequiredDeps[docType]; ok && len(reqDeps) > 0 {
			missing := missingDeps(reqDeps, fm.DependsOn)
			for _, md := range missing {
				report.add(Issue{
					Level:   levelFor(opts.StrictMode),
					Type:    IssueMissingExpectedDep,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("brak wymaganej zależności dla %s: %s", docType, md),
				})
			}
		}
		if bad, ok := containsTraversal(fm.DependsOn); ok {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueInvalidDocID,
				File:    rel,
				DocID:   docID,
				Message: fmt.Sprintf("depends_on zawiera traversal: %s", bad),
			})
		}

		// status-based validation for sections (empty headings)
		if emptyStrict && statusLower == "published" {
			metrics := sections.ComputeMetrics(content)
			if metrics.Empty > 0 {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueMissingField,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("published: %d pustych sekcji (completeness %.2f)", metrics.Empty, metrics.Completeness),
				})
			}
		} else if emptyStrict && relaxed {
			metrics := sections.ComputeMetrics(content)
			if metrics.Empty > 0 {
				report.add(Issue{
					Level:   LevelWarn,
					Type:    IssueMissingField,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("%s: %d pustych sekcji (completeness %.2f)", statusLower, metrics.Empty, metrics.Completeness),
				})
			}
		}

		// legacy section names (alias matches)
		if len(opts.SectionAliases) > 0 {
			hits := sections.FindAliasHits(content, opts.SectionAliases)
			for _, h := range hits {
				report.add(Issue{
					Level:   LevelWarn,
					Type:    IssueLegacySectionName,
					File:    rel,
					DocID:   docID,
					Message: fmt.Sprintf("legacy section name '%s' → '%s'", h.From, h.To),
				})
			}
		}

		// governance rules
		if opts.GovernanceRules != nil {
			checkGovernance(report, opts.GovernanceRules, rel, docID, statusLower, docType, content)
		}
	}

	// Duplicate doc_id detection.
	seen := make(map[string][]string)
	for _, d := range docs {
		seen[d.DocID] = append(seen[d.DocID], d.File)
	}

	for docID, filesWithID := range seen {
		if len(filesWithID) <= 1 {
			continue
		}
		sort.Strings(filesWithID)
		for _, f := range filesWithID {
			report.add(Issue{
				Level:   LevelError,
				Type:    IssueDuplicateDocID,
				File:    f,
				DocID:   docID,
				Message: fmt.Sprintf("duplikat doc_id=%s (%d wystąpień)", docID, len(filesWithID)),
			})
		}
	}

	return report, nil
}
