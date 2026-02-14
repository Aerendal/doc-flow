package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"docflow/internal/validator"
	"docflow/pkg/index"
)

// generateDocs creates n synthetic markdown documents in dir with a simple chain dependency.
func generateDocs(t *testing.T, dir string, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		docID := fmt.Sprintf("doc_%04d", i)
		dep := ""
		if i > 0 {
			dep = fmt.Sprintf("depends_on: [\"doc_%04d\"]\n", i-1)
		}
		content := fmt.Sprintf(`---
title: "Doc %d"
doc_id: %s
doc_type: guide
status: published
version: "v1.0.0"
context_sources: ["seed"]
%s---

# Przegląd
Treść %d

## Kroki
- krok
`, i, docID, dep, i)
		path := filepath.Join(dir, fmt.Sprintf("%s.md", docID))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write doc %d: %v", i, err)
		}
	}
}

func TestHappyPathIndexAndValidate(t *testing.T) {
	tmp := t.TempDir()
	const docs = 1000
	generateDocs(t, tmp, docs)

	start := time.Now()
	idx, err := index.BuildIndex(tmp, nil)
	if err != nil {
		t.Fatalf("build index: %v", err)
	}
	if idx.Count != docs {
		t.Fatalf("expected %d docs, got %d", docs, idx.Count)
	}
	buildDur := time.Since(start)

	opts := &validator.Options{
		PromoteContextFor: map[string]bool{},
		RequiredDeps:      map[string][]string{},
		SectionAliases:    map[string][]string{},
	}
	start = time.Now()
	report, err := validator.ValidateDocs(tmp, nil, opts)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if report.ErrorCount() > 0 {
		t.Fatalf("validation errors: %v", report.SortedIssues())
	}
	valDur := time.Since(start)

	t.Logf("index build: %s, validate: %s", buildDur, valDur)
}
