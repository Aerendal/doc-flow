package engine

import (
	"reflect"
	"testing"

	"docflow/pkg/compliance"
)

func TestComplianceRoundTripAdapter(t *testing.T) {
	sum := &compliance.Summary{
		SchemaVersion: "1.0",
		FailOn:        "new",
		Show:          "new",
		RulesPath:     "rules.yaml",
		RulesChecksum: "sha256:abc",
		Documents:     2,
		Passed:        1,
		Failed:        1,
		PassRate:      0.5,
		ViolationsCount: map[string]int{
			"missing_owner": 1,
		},
		NewFailed:       1,
		ExistingFailed:  0,
		NewViolations:   map[string]int{"missing_owner": 1},
		DuplicateDocIDs: map[string][]string{"dup": {"a.md", "b.md"}},
		Docs: []compliance.DocResult{
			{DocID: "dup", Path: "b.md", Status: "draft", DocType: "guide", Violations: []string{"missing_owner"}},
			{DocID: "dup", Path: "a.md", Status: "draft", DocType: "guide", Violations: []string{"missing_owner"}},
		},
	}

	report := BuildComplianceReport(sum)
	if len(report.Docs) != 2 {
		t.Fatalf("expected docs=2, got %d", len(report.Docs))
	}
	if report.Docs[0].Path != "a.md" {
		t.Fatalf("expected sorted docs, got first=%s", report.Docs[0].Path)
	}

	var out compliance.Summary
	ApplyComplianceReport(&out, report)

	if out.RulesPath != sum.RulesPath || out.RulesChecksum != sum.RulesChecksum {
		t.Fatalf("metadata mismatch after roundtrip")
	}
	if !reflect.DeepEqual(out.DuplicateDocIDs, map[string][]string{"dup": {"a.md", "b.md"}}) {
		t.Fatalf("unexpected duplicate ids map: %#v", out.DuplicateDocIDs)
	}
}
