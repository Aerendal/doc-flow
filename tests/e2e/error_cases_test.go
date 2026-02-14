package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"docflow/internal/validator"
)

func TestMissingDocID(t *testing.T) {
	tmp := t.TempDir()
	content := `---
title: "No ID"
doc_type: guide
status: draft
version: "v1"
---
# T`
	path := filepath.Join(tmp, "no_id.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := validator.ValidateDocs(tmp, nil, &validator.Options{})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if report.ErrorCount() == 0 {
		t.Fatalf("expected validation error for missing doc_id")
	}
}

func TestInvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	content := `---
title: "Bad
doc_id: bad
---`
	path := filepath.Join(tmp, "bad.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := validator.ValidateDocs(tmp, nil, &validator.Options{})
	if err != nil {
		t.Fatalf("validate returned error, expected issue report: %v", err)
	}
	if report.ErrorCount() == 0 {
		t.Fatalf("expected issues for invalid YAML, got 0")
	}
}
