package fix

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"docflow/internal/engine"
	"docflow/pkg/index"
	"docflow/pkg/sections"
)

type textFormat struct {
	hasBOM      bool
	eol         string
	trailingEOL bool
}

func BuildPlan(root string, issues []engine.Issue, files []FileInput, opts Options) (Plan, error) {
	opts = defaultOptions(opts)

	issuesCopy := append([]engine.Issue(nil), issues...)
	sortIssuesCanonical(issuesCopy)

	fileByPath := make(map[string]string, len(files))
	for _, f := range files {
		p := index.NormalizePath(f.Path)
		if p == "" {
			continue
		}
		fileByPath[p] = f.Content
	}

	issuesByPath := map[string][]engine.Issue{}
	for _, issue := range issuesCopy {
		if !allowedByCode(issue.Code, opts.OnlyCodes) {
			continue
		}
		p := index.NormalizePath(issue.Path)
		if p == "" {
			continue
		}
		issuesByPath[p] = append(issuesByPath[p], issue)
	}

	paths := make([]string, 0, len(issuesByPath))
	for path := range issuesByPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	plan := Plan{
		SchemaVersion:   "1.0",
		IdentityVersion: "2",
		Root:            root,
		Files:           []PlannedFile{},
	}

	for _, path := range paths {
		original, ok := fileByPath[path]
		if !ok {
			continue
		}

		format := detectTextFormat(original)
		normalized := normalizeContentForEdit(original, format)
		current := normalized
		changes := 0

		codeSet := map[string]bool{}
		summarySet := map[string]bool{}
		sectionRenameApplied := false

		for _, issue := range issuesByPath[path] {
			if issue.Code == "DOCFLOW.VALIDATE.MISSING_FIELD" && opts.Scopes[ScopeFrontmatter] {
				field, _ := stringDetail(issue.Details, "field")
				if field == "version" {
					next, changed, err := ensureFrontmatterField(current, "version", "v0.0.0")
					if err != nil {
						return Plan{}, fmt.Errorf("fix %s: %w", path, err)
					}
					if changed {
						current = next
						changes++
						codeSet[issue.Code] = true
						summarySet["add frontmatter field version=v0.0.0"] = true
					}
				}
				continue
			}

			if issue.Code == "DOCFLOW.VALIDATE.LEGACY_SECTION_NAME" && opts.Scopes[ScopeSections] && !sectionRenameApplied {
				renamed, hits := sections.RenameAliases(current, opts.Aliases)
				sectionRenameApplied = true
				if len(hits) > 0 && renamed != current {
					current = renamed
					changes += len(hits)
					codeSet[issue.Code] = true
					summarySet[fmt.Sprintf("rename %d legacy section headings", len(hits))] = true
				}
			}
		}

		if current == normalized {
			continue
		}
		updated := restoreContentAfterEdit(current, format)

		codes := mapKeysSorted(codeSet)
		summaries := mapKeysSorted(summarySet)
		diff := UnifiedDiff(path, original, updated)

		plan.Files = append(plan.Files, PlannedFile{
			Path:             path,
			Codes:            codes,
			Summaries:        summaries,
			Changes:          changes,
			Diff:             diff,
			OriginalChecksum: checksum(original),
			UpdatedContent:   updated,
		})
		plan.Stats.FilesTouched++
		plan.Stats.TotalEdits += changes
	}

	if plan.Stats.FilesTouched > opts.MaxFiles {
		return Plan{}, fmt.Errorf("plan exceeds --max-files (%d > %d)", plan.Stats.FilesTouched, opts.MaxFiles)
	}
	if plan.Stats.TotalEdits > opts.MaxChanges {
		return Plan{}, fmt.Errorf("plan exceeds --max-changes (%d > %d)", plan.Stats.TotalEdits, opts.MaxChanges)
	}
	return plan, nil
}

func ParseScopesCSV(scopeCSV string) (map[Scope]bool, error) {
	scopeCSV = strings.TrimSpace(scopeCSV)
	if scopeCSV == "" {
		return map[Scope]bool{
			ScopeFrontmatter: true,
			ScopeSections:    true,
		}, nil
	}

	allowed := map[Scope]bool{
		ScopeFrontmatter: true,
		ScopeSections:    true,
	}

	scopes := map[Scope]bool{}
	for _, raw := range strings.Split(scopeCSV, ",") {
		token := Scope(strings.TrimSpace(raw))
		if token == "" {
			continue
		}
		if !allowed[token] {
			return nil, fmt.Errorf("unsupported scope: %s", token)
		}
		scopes[token] = true
	}
	if len(scopes) == 0 {
		return nil, fmt.Errorf("no valid scopes in --scope")
	}
	return scopes, nil
}

func ParseOnlyCodesCSV(csv string) map[string]bool {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil
	}
	out := map[string]bool{}
	for _, raw := range strings.Split(csv, ",") {
		code := strings.ToUpper(strings.TrimSpace(raw))
		if code == "" {
			continue
		}
		if !strings.HasPrefix(code, "DOCFLOW.") {
			code = "DOCFLOW.VALIDATE." + code
		}
		out[code] = true
	}
	return out
}

func RenderDiff(plan Plan) string {
	if len(plan.Files) == 0 {
		return ""
	}
	parts := make([]string, 0, len(plan.Files))
	for _, f := range plan.Files {
		if strings.TrimSpace(f.Diff) == "" {
			continue
		}
		parts = append(parts, strings.TrimSuffix(f.Diff, "\n"))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "\n\n") + "\n"
}

func mapKeysSorted(values map[string]bool) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for k := range values {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func stringDetail(details map[string]any, key string) (string, bool) {
	if len(details) == 0 {
		return "", false
	}
	v, ok := details[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	return s, true
}

func ensureFrontmatterField(content, key, value string) (string, bool, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return content, false, nil
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end <= 0 {
		return content, false, fmt.Errorf("invalid frontmatter delimiters")
	}

	prefix := strings.ToLower(key) + ":"
	for i := 1; i < end; i++ {
		l := strings.TrimSpace(lines[i])
		if strings.HasPrefix(strings.ToLower(l), prefix) {
			return content, false, nil
		}
	}

	insert := key + ": " + value
	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:end]...)
	out = append(out, insert)
	out = append(out, lines[end:]...)
	return strings.Join(out, "\n"), true, nil
}

func detectTextFormat(content string) textFormat {
	f := textFormat{
		eol: "\n",
	}
	if strings.HasPrefix(content, "\uFEFF") {
		f.hasBOM = true
		content = strings.TrimPrefix(content, "\uFEFF")
	}
	if strings.Contains(content, "\r\n") {
		f.eol = "\r\n"
	}
	f.trailingEOL = strings.HasSuffix(content, "\n")
	return f
}

func normalizeContentForEdit(content string, format textFormat) string {
	if format.hasBOM {
		content = strings.TrimPrefix(content, "\uFEFF")
	}
	if format.eol == "\r\n" {
		content = strings.ReplaceAll(content, "\r\n", "\n")
	}
	return content
}

func restoreContentAfterEdit(content string, format textFormat) string {
	if format.trailingEOL {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
	} else {
		content = strings.TrimSuffix(content, "\n")
	}
	if format.eol == "\r\n" {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if format.hasBOM {
		content = "\uFEFF" + content
	}
	return content
}

func absPath(root, rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	if root == "" {
		return rel
	}
	return filepath.Join(root, rel)
}
