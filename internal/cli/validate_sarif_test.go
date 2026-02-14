package cli

import (
	"encoding/json"
	"testing"
)

func TestBuildValidateSARIFReportSortedAndMapped(t *testing.T) {
	report := validateReportJSON{
		SchemaVersion: "1.0",
		FailOn:        "new",
		Show:          "new",
		Issues: []validateIssueJSON{
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "warn",
				Type:    "missing_field",
				Path:    "b.md",
				DocID:   "b",
				Message: "brak version",
			},
			{
				Code:    "DOCFLOW.VALIDATE.INVALID_YAML",
				Level:   "error",
				Type:    "invalid_yaml",
				Path:    "a.md",
				Line:    12,
				Message: "invalid yaml",
			},
		},
	}

	sarif := buildValidateSARIFReport(report)
	if sarif.Version != "2.1.0" {
		t.Fatalf("unexpected sarif version: %s", sarif.Version)
	}
	if len(sarif.Runs) != 1 {
		t.Fatalf("expected one run, got %d", len(sarif.Runs))
	}
	run := sarif.Runs[0]
	if len(run.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(run.Results))
	}
	// Sorted by path => a.md first.
	if run.Results[0].RuleID != "DOCFLOW.VALIDATE.INVALID_YAML" {
		t.Fatalf("unexpected first rule: %s", run.Results[0].RuleID)
	}
	if run.Results[0].Level != "error" {
		t.Fatalf("unexpected level mapping: %s", run.Results[0].Level)
	}
	if len(run.Results[0].Locations) != 1 || run.Results[0].Locations[0].PhysicalLocation.Region == nil {
		t.Fatal("expected line-based location for issue with line")
	}
	if run.Results[1].Level != "warning" {
		t.Fatalf("unexpected warning mapping: %s", run.Results[1].Level)
	}
	if run.Properties["fail_on"] != "new" || run.Properties["show"] != "new" {
		t.Fatalf("unexpected run properties: %+v", run.Properties)
	}
}

func TestBuildValidateSARIFReportDeterministic(t *testing.T) {
	report := validateReportJSON{
		SchemaVersion: "1.0",
		Issues: []validateIssueJSON{
			{Code: "DOCFLOW.VALIDATE.MISSING_FIELD", Level: "error", Type: "missing_field", Path: "doc.md", Message: "m1"},
			{Code: "DOCFLOW.VALIDATE.MISSING_CONTEXT_SOURCES", Level: "warn", Type: "missing_context_sources", Path: "doc.md", Message: "m2"},
		},
	}

	a, err := json.Marshal(buildValidateSARIFReport(report))
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(buildValidateSARIFReport(report))
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("sarif output is not deterministic")
	}
}
