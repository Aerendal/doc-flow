package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"docflow/pkg/compliance"
)

func TestApplyValidateBaseline(t *testing.T) {
	tmp := t.TempDir()
	baselinePath := filepath.Join(tmp, "validate_baseline.json")

	baseline := validateReportJSON{
		SchemaVersion: "1.0",
		Issues: []validateIssueJSON{
			{Code: "DOCFLOW.VALIDATE.MISSING_FIELD", Level: "error", Path: "a.md", DocID: "a", Line: 3, Message: "old"},
		},
	}
	data, err := json.Marshal(baseline)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(baselinePath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	current := validateReportJSON{
		SchemaVersion: "1.0",
		Issues: []validateIssueJSON{
			{Code: "DOCFLOW.VALIDATE.MISSING_FIELD", Level: "error", Path: "a.md", DocID: "a", Line: 3, Message: "old"},
			{Code: "DOCFLOW.VALIDATE.MISSING_CONTEXT_SOURCES", Level: "warn", Path: "b.md", DocID: "b", Message: "new"},
		},
	}

	got, err := applyValidateBaseline(current, baselinePath, "new")
	if err != nil {
		t.Fatal(err)
	}
	if got.Baseline == nil || !got.Baseline.Loaded {
		t.Fatalf("expected baseline metadata, got %#v", got.Baseline)
	}
	if got.NewErrorCount != 0 || got.NewWarnCount != 1 {
		t.Fatalf("unexpected new counts: errors=%d warns=%d", got.NewErrorCount, got.NewWarnCount)
	}
	if got.ExistingErrorCount != 1 || got.ExistingWarnCount != 0 {
		t.Fatalf("unexpected existing counts: errors=%d warns=%d", got.ExistingErrorCount, got.ExistingWarnCount)
	}
	if len(got.Issues) != 1 || got.Issues[0].Path != "b.md" {
		t.Fatalf("unexpected issues after show=new: %#v", got.Issues)
	}
}

func TestApplyValidateBaselineDistinctMessagesDoNotMerge(t *testing.T) {
	tmp := t.TempDir()
	baselinePath := filepath.Join(tmp, "validate_baseline.json")

	baseline := validateReportJSON{
		SchemaVersion:   "1.0",
		IdentityVersion: "2",
		Issues: []validateIssueJSON{
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "error",
				Path:    "doc.md",
				DocID:   "doc",
				Message: "brak version",
				Details: map[string]any{"field": "version"},
			},
		},
	}
	data, err := json.Marshal(baseline)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(baselinePath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	current := validateReportJSON{
		SchemaVersion:   "1.0",
		IdentityVersion: "2",
		Issues: []validateIssueJSON{
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "error",
				Path:    "doc.md",
				DocID:   "doc",
				Message: "brak version",
				Details: map[string]any{"field": "version"},
			},
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "error",
				Path:    "doc.md",
				DocID:   "doc",
				Message: "brak title",
				Details: map[string]any{"field": "title"},
			},
		},
	}

	got, err := applyValidateBaseline(current, baselinePath, "all")
	if err != nil {
		t.Fatal(err)
	}
	if got.NewErrorCount != 1 || got.ExistingErrorCount != 1 {
		t.Fatalf("unexpected baseline split: new=%d existing=%d", got.NewErrorCount, got.ExistingErrorCount)
	}
}

func TestApplyValidateBaselineMessageChangeIgnoredInIdentityV2(t *testing.T) {
	tmp := t.TempDir()
	baselinePath := filepath.Join(tmp, "validate_baseline.json")

	baseline := validateReportJSON{
		SchemaVersion:   "1.0",
		IdentityVersion: "2",
		Issues: []validateIssueJSON{
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "error",
				Path:    "doc.md",
				DocID:   "doc",
				Message: "brak version",
				Details: map[string]any{"field": "version"},
			},
		},
	}
	data, err := json.Marshal(baseline)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(baselinePath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	current := validateReportJSON{
		SchemaVersion:   "1.0",
		IdentityVersion: "2",
		Issues: []validateIssueJSON{
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "error",
				Path:    "doc.md",
				DocID:   "doc",
				Message: "pole version jest wymagane",
				Details: map[string]any{"field": "version"},
			},
		},
	}

	got, err := applyValidateBaseline(current, baselinePath, "all")
	if err != nil {
		t.Fatal(err)
	}
	if got.NewErrorCount != 0 || got.ExistingErrorCount != 1 {
		t.Fatalf("expected wording-only change to be existing, new=%d existing=%d", got.NewErrorCount, got.ExistingErrorCount)
	}
}

func TestApplyComplianceBaseline(t *testing.T) {
	tmp := t.TempDir()
	baselinePath := filepath.Join(tmp, "compliance_baseline.json")

	baseline := compliance.Summary{
		SchemaVersion: "1.0",
		Docs: []compliance.DocResult{
			{Path: "a.md", Violations: []string{"missing_owner"}},
		},
	}
	data, err := json.Marshal(baseline)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(baselinePath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	sum := &compliance.Summary{
		SchemaVersion: "1.0",
		Docs: []compliance.DocResult{
			{Path: "a.md", Violations: []string{"missing_owner", "missing_section_FAQ"}},
			{Path: "b.md", Violations: []string{"missing_owner"}},
		},
	}

	if err := applyComplianceBaseline(sum, baselinePath, "new"); err != nil {
		t.Fatal(err)
	}
	if sum.Baseline == nil || !sum.Baseline.Loaded {
		t.Fatalf("expected baseline metadata, got %#v", sum.Baseline)
	}
	if sum.NewFailed != 2 {
		t.Fatalf("expected NewFailed=2, got %d", sum.NewFailed)
	}
	if sum.ExistingFailed != 0 {
		t.Fatalf("expected ExistingFailed=0, got %d", sum.ExistingFailed)
	}
	wantNewViolations := map[string]int{
		"missing_owner":       1,
		"missing_section_FAQ": 1,
	}
	if !reflect.DeepEqual(sum.NewViolations, wantNewViolations) {
		t.Fatalf("unexpected new violations map: got=%v want=%v", sum.NewViolations, wantNewViolations)
	}
	if len(sum.Docs) != 2 {
		t.Fatalf("expected 2 docs in show=new, got %d", len(sum.Docs))
	}
	if !reflect.DeepEqual(sum.Docs[0].Violations, []string{"missing_section_FAQ"}) {
		t.Fatalf("unexpected filtered violations for a.md: %v", sum.Docs[0].Violations)
	}
}
