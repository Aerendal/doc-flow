package recommend

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUsageStore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "usage.json")
	store, err := LoadUsage(path)
	if err != nil {
		t.Fatalf("LoadUsage error: %v", err)
	}
	store.Inc("a.md")
	store.Inc("a.md")
	if err := store.Save(); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if len(data) == 0 {
		t.Fatalf("usage file empty")
	}
	store2, err := LoadUsage(path)
	if err != nil {
		t.Fatalf("LoadUsage2 error: %v", err)
	}
	if store2.Usage["a.md"] != 2 {
		t.Fatalf("usage count = %d, want 2", store2.Usage["a.md"])
	}
}
