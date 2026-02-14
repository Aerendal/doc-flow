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

// generateFanOutDocs creates n documents where doc_0 depends on up to fanOut leaves.
func generateFanOutDocs(t *testing.T, dir string, n int, fanOut int, bodySize int) {
	t.Helper()
	if fanOut >= n {
		fanOut = n - 1
	}
	body := make([]byte, bodySize)
	for i := range body {
		body[i] = 'a' + byte(i%26)
	}

	for i := 0; i < n; i++ {
		docID := fmt.Sprintf("doc_%05d", i)
		dep := ""
		if i == 0 {
			deps := make([]string, 0, fanOut)
			for d := 1; d <= fanOut; d++ {
				deps = append(deps, fmt.Sprintf("doc_%05d", d))
			}
			dep = fmt.Sprintf("depends_on: [\"%s\"]\n", join(deps, "\",\""))
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

// join helper to avoid importing strings (to keep small deps).
func join(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	out := items[0]
	for _, s := range items[1:] {
		out += sep + s
	}
	return out
}

func TestPerf20kFanOut(t *testing.T) {
	if os.Getenv("PERF20K") == "" {
		t.Skip("set PERF20K=1 to run (heavy)")
	}

	tmp := t.TempDir()
	const docs = 20000
	generateFanOutDocs(t, tmp, docs, 50, 1200) // ~1.2KB each

	start := time.Now()
	idx, err := index.BuildIndex(tmp, nil)
	if err != nil {
		t.Fatalf("build index: %v", err)
	}
	buildDur := time.Since(start)

	opts := &validator.Options{}
	start = time.Now()
	report, err := validator.ValidateDocs(tmp, nil, opts)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if report.ErrorCount() > 0 {
		t.Fatalf("validation issues: %v", report.SortedIssues())
	}
	valDur := time.Since(start)

	t.Logf("[20k fan-out] docs=%d index=%s validate=%s", idx.Count, buildDur, valDur)
}
