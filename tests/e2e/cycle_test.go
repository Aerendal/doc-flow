package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"docflow/internal/validator"
	"docflow/pkg/index"
)

func TestCycleDetectionChecksumChange(t *testing.T) {
	tmp := t.TempDir()
	write := func(name, dep string) {
		content := `---
title: "` + name + `"
doc_id: ` + name + `
doc_type: guide
status: draft
depends_on: [` + dep + `]
---
# T`
		if err := os.WriteFile(filepath.Join(tmp, name+".md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a", `"b"`)
	write("b", `"a"`)

	idx, err := index.BuildIndex(tmp, nil)
	if err != nil {
		t.Fatalf("build index: %v", err)
	}
	rec := idx.GetByID("a")
	if rec == nil || len(rec.Cycle) == 0 {
		t.Fatalf("expected cycle annotation for doc a")
	}

	// Validate should surface cycle as error.
	opts := &validator.Options{}
	report, err := validator.ValidateDocs(tmp, nil, opts)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if report.ErrorCount() == 0 {
		t.Fatalf("expected errors for cycle, got none")
	}
}
