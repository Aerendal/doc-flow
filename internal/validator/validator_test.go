package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %s: %v", name, err)
	}
	return path
}

func TestMissingFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", "# No FM")

	report, err := ValidateDocs(dir, nil, nil)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	if report.ErrorCount() != 1 {
		t.Fatalf("errors = %d, want 1", report.ErrorCount())
	}
	if report.Issues[0].Type != IssueMissingFrontmatter {
		t.Fatalf("issue type = %s, want %s", report.Issues[0].Type, IssueMissingFrontmatter)
	}
}

func TestInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bad.md", "---\n title: test\n status: [oops\n---\n# h1\n")

	report, err := ValidateDocs(dir, nil, nil)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	if report.ErrorCount() != 1 {
		t.Fatalf("errors = %d, want 1", report.ErrorCount())
	}
	if report.Issues[0].Type != IssueInvalidYAML {
		t.Fatalf("issue type = %s, want %s", report.Issues[0].Type, IssueInvalidYAML)
	}
}

func TestMissingRequiredFields(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "req.md", `---
doc_id: sample_doc
status: draft
doc_type: spec
---
# Title
`)

	report, err := ValidateDocs(dir, nil, nil)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}

	// Missing: title, version
	if report.ErrorCount() != 2 {
		t.Fatalf("errors = %d, want 2", report.ErrorCount())
	}
}

func TestInvalidDocIDAndMismatch(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "my file.md", `---
doc_id: Bad-ID
title: Test
doc_type: specification
version: "1.0.0"
status: approved
context_sources: [ctx-1]
---
# Title
`)

	report, err := ValidateDocs(dir, nil, nil)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}

	if report.ErrorCount() != 1 {
		t.Fatalf("errors = %d, want 1", report.ErrorCount())
	}
	if report.Issues[0].Type != IssueInvalidDocID {
		t.Fatalf("issue type = %s, want %s", report.Issues[0].Type, IssueInvalidDocID)
	}

	// Expect filename mismatch warning as well.
	if report.WarnCount() != 1 {
		t.Fatalf("warns = %d, want 1", report.WarnCount())
	}
}

func TestDuplicateDocID(t *testing.T) {
	dir := t.TempDir()
	content := `---
doc_id: dup
title: T
doc_type: specification
version: "1.0.0"
status: draft
---
`
	writeFile(t, dir, "a.md", content)
	writeFile(t, dir, "b.md", content)

	report, err := ValidateDocs(dir, nil, nil)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}

	if report.ErrorCount() != 2 {
		t.Fatalf("errors = %d, want 2 (duplicate per file)", report.ErrorCount())
	}
	found := 0
	for _, i := range report.Issues {
		if i.Type == IssueDuplicateDocID {
			found++
		}
	}
	if found != 2 {
		t.Fatalf("duplicate issues = %d, want 2", found)
	}
}
