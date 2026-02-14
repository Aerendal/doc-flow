package engine

import (
	"testing"

	"docflow/internal/validator"
)

func TestBuildValidateReportCanonicalSortAndDetails(t *testing.T) {
	in := &validator.Report{
		Files:     2,
		Documents: 2,
		Issues: []validator.Issue{
			{
				Level:   validator.LevelError,
				Type:    validator.IssueMissingField,
				File:    "b.md",
				DocID:   "b",
				Message: "brak version",
			},
			{
				Level:   validator.LevelWarn,
				Type:    validator.IssueInvalidYAML,
				File:    "a.md",
				Line:    12,
				Message: "invalid yaml",
			},
		},
	}

	got := BuildValidateReport(in)
	if got.SchemaVersion != "1.0" {
		t.Fatalf("unexpected schema version: %s", got.SchemaVersion)
	}
	if len(got.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(got.Issues))
	}
	if got.Issues[0].Path != "a.md" {
		t.Fatalf("expected sorted issues by path, got first=%s", got.Issues[0].Path)
	}
	if got.Issues[1].Details == nil {
		t.Fatalf("expected details for missing_field issue")
	}
	if field, _ := got.Issues[1].Details["field"].(string); field != "version" {
		t.Fatalf("expected details.field=version, got %#v", got.Issues[1].Details["field"])
	}
}
