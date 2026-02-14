package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"docflow/pkg/compliance"
)

func TestMigrateValidateBaselineAddsIdentityV2AndDetails(t *testing.T) {
	in := validateReportJSON{
		SchemaVersion: "1.0",
		Issues: []validateIssueJSON{
			{
				Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
				Level:   "error",
				Type:    "missing_field",
				Path:    "doc.md",
				DocID:   "doc",
				Message: "brak version",
			},
		},
	}
	migrateValidateBaseline(&in)
	if in.IdentityVersion != "2" {
		t.Fatalf("expected identity_version=2, got %q", in.IdentityVersion)
	}
	if len(in.Issues) != 1 || in.Issues[0].Details == nil {
		t.Fatalf("expected details after migration, got %#v", in.Issues)
	}
	if field, _ := in.Issues[0].Details["field"].(string); field != "version" {
		t.Fatalf("expected details.field=version, got %#v", in.Issues[0].Details["field"])
	}
}

func TestMigrateComplianceBaselineAddsIdentityV2(t *testing.T) {
	in := compliance.Summary{
		SchemaVersion: "1.0",
		Docs: []compliance.DocResult{
			{Path: "doc.md", Violations: []string{"missing_owner"}},
		},
	}
	migrateComplianceBaseline(&in)
	if in.IdentityVersion != "2" {
		t.Fatalf("expected identity_version=2, got %q", in.IdentityVersion)
	}
}

func TestDetectBaselineKind(t *testing.T) {
	tmp := t.TempDir()

	validatePath := filepath.Join(tmp, "validate.json")
	validateData, _ := json.Marshal(validateReportJSON{Issues: []validateIssueJSON{{Code: "C"}}})
	if err := os.WriteFile(validatePath, validateData, 0o644); err != nil {
		t.Fatal(err)
	}
	kind, err := detectBaselineKind(validatePath)
	if err != nil {
		t.Fatal(err)
	}
	if kind != "validate" {
		t.Fatalf("expected validate kind, got %q", kind)
	}

	compliancePath := filepath.Join(tmp, "compliance.json")
	complianceData, _ := json.Marshal(compliance.Summary{Docs: []compliance.DocResult{{Path: "a.md"}}})
	if err := os.WriteFile(compliancePath, complianceData, 0o644); err != nil {
		t.Fatal(err)
	}
	kind, err = detectBaselineKind(compliancePath)
	if err != nil {
		t.Fatal(err)
	}
	if kind != "compliance" {
		t.Fatalf("expected compliance kind, got %q", kind)
	}
}
