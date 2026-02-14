package engine

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"docflow/internal/validator"
)

func BuildValidateReport(report *validator.Report) ValidateReport {
	if report == nil {
		return ValidateReport{
			ReportMeta:      ReportMeta{SchemaVersion: "1.0"},
			IdentityVersion: "2",
			Issues:          []Issue{},
		}
	}

	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.SortedIssues() {
		engineIssue := Issue{
			Code:     validateIssueCode(issue.Type),
			Severity: Severity(issue.Level),
			Type:     string(issue.Type),
			Kind:     IssueKindValidate,
			Path:     issue.File,
			DocID:    issue.DocID,
			Message:  issue.Message,
			Details:  validateIssueDetails(issue),
		}
		if issue.Line > 0 {
			engineIssue.Location = &Location{Line: issue.Line}
		}
		issues = append(issues, engineIssue)
	}

	// Engine-level ordering is canonical. Renderers must not apply their own sorting.
	sortIssues(issues)

	return ValidateReport{
		ReportMeta: ReportMeta{
			SchemaVersion: "1.0",
		},
		IdentityVersion: "2",
		Files:           report.Files,
		Documents:       report.Documents,
		ErrorCount:      report.ErrorCount(),
		WarnCount:       report.WarnCount(),
		Issues:          issues,
	}
}

func validateIssueCode(issueType validator.IssueType) string {
	code := strings.ToUpper(strings.ReplaceAll(string(issueType), "-", "_"))
	return "DOCFLOW.VALIDATE." + code
}

func validateIssueDetails(issue validator.Issue) map[string]any {
	// Keep details stable and machine-oriented; message remains UI-facing.
	msg := strings.TrimSpace(issue.Message)
	lower := strings.ToLower(msg)
	out := map[string]any{}

	switch issue.Type {
	case validator.IssueMissingField, validator.IssueEmptyField:
		if field, ok := fieldFromPolishMessage(msg); ok {
			out["field"] = field
		}
		if matches := emptySectionsMetricRe.FindStringSubmatch(lower); len(matches) == 3 {
			out["status"] = matches[1]
			if n, err := strconv.Atoi(matches[2]); err == nil {
				out["empty_sections"] = n
			}
			out["reason"] = "empty_sections"
		}
	case validator.IssueInvalidYAML:
		out["stage"] = "frontmatter"
	case validator.IssueDocIDFilename:
		if matches := docIDFilenameMismatchRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["actual_doc_id"] = matches[1]
			out["expected_doc_id"] = matches[2]
		}
	case validator.IssueMissingExpectedDep:
		if matches := missingExpectedDepRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["doc_type"] = strings.TrimSpace(matches[1])
			out["dependency"] = strings.TrimSpace(matches[2])
		}
	case validator.IssueInvalidDocID:
		if strings.HasPrefix(lower, "context_sources zawiera traversal:") {
			out["field"] = "context_sources"
			out["reason"] = "path_traversal"
			out["value"] = strings.TrimSpace(strings.TrimPrefix(msg, "context_sources zawiera traversal:"))
		} else if strings.HasPrefix(lower, "depends_on zawiera traversal:") {
			out["field"] = "depends_on"
			out["reason"] = "path_traversal"
			out["value"] = strings.TrimSpace(strings.TrimPrefix(msg, "depends_on zawiera traversal:"))
		}
	case validator.IssueLegacySectionName:
		if matches := legacySectionAliasRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["from"] = matches[1]
			out["to"] = matches[2]
		}
	case validator.IssueGovernance:
		if matches := governanceFieldMissingRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["rule"] = "missing_required_field"
			out["field"] = strings.TrimSpace(matches[1])
			out["status"] = strings.TrimSpace(matches[2])
		} else if matches := governanceSectionMissRe.FindStringSubmatch(msg); len(matches) == 3 {
			out["rule"] = "missing_required_section"
			out["section"] = strings.TrimSpace(matches[1])
			out["doc_type"] = strings.TrimSpace(matches[2])
		}
	case validator.IssueDuplicateDocID:
		if matches := duplicateDocIDCountRe.FindStringSubmatch(msg); len(matches) == 2 {
			if n, err := strconv.Atoi(matches[1]); err == nil {
				out["duplicate_count"] = n
			}
		}
	case validator.IssueCycleDetected:
		if strings.HasPrefix(lower, "wykryto cykl: ") {
			out["cycle"] = strings.TrimSpace(strings.TrimPrefix(msg, "wykryto cykl: "))
		}
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

var (
	missingExpectedDepRe     = regexp.MustCompile(`^brak wymaganej zależności dla ([^:]+):\s*(.+)$`)
	docIDFilenameMismatchRe  = regexp.MustCompile(`^doc_id=([^\s]+)\s+nie zgadza się z nazwą pliku\s+\(([^)]+)\)$`)
	legacySectionAliasRe     = regexp.MustCompile(`^legacy section name '(.+)' → '(.+)'$`)
	governanceFieldMissingRe = regexp.MustCompile(`^governance: missing required field '([^']+)' for status (.+)$`)
	governanceSectionMissRe  = regexp.MustCompile(`^governance: missing required section '([^']+)' for family (.+)$`)
	duplicateDocIDCountRe    = regexp.MustCompile(`\((\d+)\s+wystąpień\)$`)
	emptySectionsMetricRe    = regexp.MustCompile(`^([a-z_]+):\s+(\d+)\s+pustych sekcji`)
)

func fieldFromPolishMessage(message string) (string, bool) {
	const prefix = "brak "
	msg := strings.ToLower(strings.TrimSpace(message))
	if !strings.HasPrefix(msg, prefix) {
		return "", false
	}
	field := strings.TrimSpace(strings.TrimPrefix(msg, prefix))
	if field == "" {
		return "", false
	}
	return field, true
}

func sortIssues(issues []Issue) {
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		if issues[i].Code != issues[j].Code {
			return issues[i].Code < issues[j].Code
		}
		if issues[i].DocID != issues[j].DocID {
			return issues[i].DocID < issues[j].DocID
		}
		li := 0
		if issues[i].Location != nil {
			li = issues[i].Location.Line
		}
		lj := 0
		if issues[j].Location != nil {
			lj = issues[j].Location.Line
		}
		if li != lj {
			return li < lj
		}
		if issues[i].Type != issues[j].Type {
			return issues[i].Type < issues[j].Type
		}
		return issues[i].Message < issues[j].Message
	})
}
