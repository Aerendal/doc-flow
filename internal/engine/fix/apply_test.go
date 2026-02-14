package fix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"docflow/internal/engine"
)

func TestApplyPlanPreservesBOMAndCRLF(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "doc.md")
	original := "\uFEFF---\r\ndoc_id: a_doc\r\ntitle: A\r\ndoc_type: guide\r\nstatus: draft\r\ncontext_sources: [seed]\r\n---\r\n# A\r\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	plan, err := BuildPlan(tmp, []engine.Issue{
		{
			Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
			Type:    "missing_field",
			Path:    "doc.md",
			Message: "brak version",
			Details: map[string]any{"field": "version"},
		},
	}, []FileInput{{Path: "doc.md", Content: original}}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Files) != 1 {
		t.Fatalf("expected one planned file, got %d", len(plan.Files))
	}
	if err := ApplyPlan(plan, ApplyOptions{Root: tmp}); err != nil {
		t.Fatal(err)
	}

	updatedBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := string(updatedBytes)
	if !strings.HasPrefix(updated, "\uFEFF") {
		t.Fatalf("expected UTF-8 BOM to be preserved")
	}
	if !strings.Contains(updated, "\r\nversion: v0.0.0\r\n") {
		t.Fatalf("expected CRLF+version in updated content, got: %q", updated)
	}
}

func TestApplyPlanDetectsConcurrentModification(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "doc.md")
	original := "---\ndoc_id: a_doc\ntitle: A\ndoc_type: guide\nstatus: draft\ncontext_sources: [seed]\n---\n# A\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	plan, err := BuildPlan(tmp, []engine.Issue{
		{
			Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
			Type:    "missing_field",
			Path:    "doc.md",
			Message: "brak version",
			Details: map[string]any{"field": "version"},
		},
	}, []FileInput{{Path: "doc.md", Content: original}}, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(original+"\nchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ApplyPlan(plan, ApplyOptions{Root: tmp}); err == nil {
		t.Fatalf("expected conflict error when file changes after plan creation")
	}
}

func TestApplyPlanFailsOnReadOnlyDirectory(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "ro")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "doc.md")
	original := "---\ndoc_id: a_doc\ntitle: A\ndoc_type: guide\nstatus: draft\ncontext_sources: [seed]\n---\n# A\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	plan, err := BuildPlan(tmp, []engine.Issue{
		{
			Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
			Type:    "missing_field",
			Path:    "ro/doc.md",
			Message: "brak version",
			Details: map[string]any{"field": "version"},
		},
	}, []FileInput{{Path: "ro/doc.md", Content: original}}, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0o755)

	if err := ApplyPlan(plan, ApplyOptions{Root: tmp}); err == nil {
		t.Fatalf("expected write error in read-only directory")
	}
}
