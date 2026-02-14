package exec

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	execache "docflow/internal/exec/cache"
)

func TestCollectFactsWithSQLiteCacheHit(t *testing.T) {
	root := t.TempDir()
	cacheDir := filepath.Join(root, ".docflow", "cache")
	if err := os.WriteFile(filepath.Join(root, "a.md"), []byte("---\ntitle: A\ndoc_id: a\ndoc_type: guide\nstatus: draft\nversion: v1\n---\n# A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.md"), []byte("---\ntitle: B\ndoc_id: b\ndoc_type: guide\nstatus: draft\nversion: v1\n---\n# B\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	facts1, stats1, err := CollectFacts(RunnerOptions{
		Mode:           "validate",
		Root:           root,
		IgnorePatterns: nil,
		UseCache:       true,
		CacheDir:       cacheDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if stats1.CacheMisses != 2 {
		t.Fatalf("expected first run misses=2, got %+v", stats1)
	}

	facts2, stats2, err := CollectFacts(RunnerOptions{
		Mode:           "validate",
		Root:           root,
		IgnorePatterns: nil,
		UseCache:       true,
		CacheDir:       cacheDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if stats2.CacheHits != 2 {
		t.Fatalf("expected second run hits=2, got %+v", stats2)
	}

	if !reflect.DeepEqual(stripRuntimeFields(facts1), stripRuntimeFields(facts2)) {
		t.Fatalf("facts differ between miss/hit runs")
	}
}

func stripRuntimeFields(in []execache.DocumentFacts) []execache.DocumentFacts {
	out := make([]execache.DocumentFacts, 0, len(in))
	for _, f := range in {
		cp := f
		cp.AbsPath = ""
		cp.Root = ""
		out = append(out, cp)
	}
	return out
}
