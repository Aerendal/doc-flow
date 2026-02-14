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

// generateDocsSized creates n markdown documents with roughly bodySize bytes each.
func generateDocsSized(t *testing.T, dir string, n int, bodySize int) {
	t.Helper()
	body := make([]byte, bodySize)
	for i := range body {
		body[i] = 'a' + byte(i%26)
	}
	for i := 0; i < n; i++ {
		docID := fmt.Sprintf("doc_%05d", i)
		dep := ""
		if i > 0 {
			dep = fmt.Sprintf("depends_on: [\"doc_%05d\"]\n", i-1)
		}
		content := fmt.Sprintf(`---
title: "Doc %d"
doc_id: %s
doc_type: guide
status: draft
version: "v1.0.0"
context_sources: ["seed"]
%s---

# Przegląd
%s
`, i, docID, dep, string(body))
		path := filepath.Join(dir, fmt.Sprintf("%s.md", docID))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write doc %d: %v", i, err)
		}
	}
}

func TestPerf10k(t *testing.T) {
	if os.Getenv("PERF10K") == "" {
		t.Skip("set PERF10K=1 to run (heavy)")
	}

	tmp := t.TempDir()
	const docs = 10000
	generateDocsSized(t, tmp, docs, 2000) // ~2KB each -> ~20 MB total

	start := time.Now()
	idx, err := index.BuildIndex(tmp, nil)
	if err != nil {
		t.Fatalf("build index: %v", err)
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
		t.Fatalf("validation issues: %v", report.SortedIssues())
	}
	valDur := time.Since(start)

	t.Logf("[10k] docs=%d index=%s validate=%s", idx.Count, buildDur, valDur)
}
