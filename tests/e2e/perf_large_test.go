package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"docflow/internal/validator"
	"docflow/pkg/index"
)

// generateLargeDocs creates n markdown docs with ~bodySize bytes of content each.
// If fanOut is true, the first document depends on the next 50 documents (or n-1 if smaller).
func generateLargeDocs(t *testing.T, dir string, n int, bodySize int, fanOut bool) {
	t.Helper()

	// Build a body string roughly bodySize bytes long.
	chunk := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus placerat.\n"
	repeats := bodySize/len(chunk) + 1
	body := strings.Repeat(chunk, repeats)
	if len(body) > bodySize {
		body = body[:bodySize]
	}

	for i := 0; i < n; i++ {
		docID := fmt.Sprintf("doc_%04d", i)
		dep := ""
		if fanOut && i == 0 {
			limit := 50
			if n-1 < limit {
				limit = n - 1
			}
			deps := make([]string, 0, limit)
			for d := 1; d <= limit; d++ {
				deps = append(deps, fmt.Sprintf("doc_%04d", d))
			}
			dep = fmt.Sprintf("depends_on: [\"%s\"]\n", strings.Join(deps, `","`))
		}

		content := fmt.Sprintf(`---
title: "Doc %d"
doc_id: %s
doc_type: guide
status: published
version: "v1.0.0"
context_sources: ["seed"]
%s---

# Wstęp
%s
`, i, docID, dep, body)

		path := filepath.Join(dir, fmt.Sprintf("%s.md", docID))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write doc %d: %v", i, err)
		}
	}
}

func perfRun(t *testing.T, dir string, label string) {
	start := time.Now()
	idx, err := index.BuildIndex(dir, nil)
	if err != nil {
		t.Fatalf("[%s] build index: %v", label, err)
	}
	buildDur := time.Since(start)

	opts := &validator.Options{
		PromoteContextFor: map[string]bool{},
		RequiredDeps:      map[string][]string{},
		SectionAliases:    map[string][]string{},
	}
	start = time.Now()
	report, err := validator.ValidateDocs(dir, nil, opts)
	if err != nil {
		t.Fatalf("[%s] validate: %v", label, err)
	}
	if report.ErrorCount() > 0 {
		t.Fatalf("[%s] validation issues: %v", label, report.SortedIssues())
	}
	valDur := time.Since(start)

	t.Logf("[%s] docs=%d index=%s validate=%s", label, idx.Count, buildDur, valDur)
}

func TestPerfLarge(t *testing.T) {
	t.Run("large_text_10k_each", func(t *testing.T) {
		tmp := t.TempDir()
		generateLargeDocs(t, tmp, 1000, 10_000, false)
		perfRun(t, tmp, "large_text")
	})

	t.Run("fan_out_50_deps", func(t *testing.T) {
		tmp := t.TempDir()
		generateLargeDocs(t, tmp, 1000, 8_000, true)
		perfRun(t, tmp, "fan_out")
	})
}
