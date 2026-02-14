package chaos

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"docflow/pkg/index"
	"docflow/internal/validator"
)

// Ensure YAML bomb / large frontmatter doesn't panic.
func TestYAMLBomb_NoPanic(t *testing.T) {
	tmp := t.TempDir()
	// 1000 nested keys (simple, not exponential) to avoid extreme memory but still large.
	body := "a: 1\n"
	for i := 0; i < 1000; i++ {
		body += "k" + string('a'+byte(i%26)) + ": 1\n"
	}
	content := "---\n" + body + "---\n# T\n"
	if err := os.WriteFile(filepath.Join(tmp, "bomb.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	// We only care that it doesn't panic/hang.
	_, err := validator.ValidateDocs(tmp, nil, &validator.Options{})
	if err != nil {
		t.Logf("validation returned error (acceptable): %v", err)
	}
}

// Ensure symlink loops do not hang scanning.
func TestSymlinkLoop_NoHang(t *testing.T) {
	tmp := t.TempDir()
	loop := filepath.Join(tmp, "loop")
	if err := os.Mkdir(loop, 0o755); err != nil {
		t.Fatal(err)
	}
	// symlink loop/inner -> loop (cycle)
	if err := os.Symlink("..", filepath.Join(loop, "back")); err != nil {
		t.Skipf("symlink not supported: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = index.BuildIndex(tmp, nil)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("scan hung on symlink loop")
	}
}

// Corrupted cache file should return error, not panic.
func TestCorruptedCacheLoadJSON(t *testing.T) {
	tmp := t.TempDir()
	bad := filepath.Join(tmp, "doc_index.json")
	if err := os.WriteFile(bad, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := index.LoadJSON(bad); err == nil {
		t.Fatalf("expected error on corrupted cache")
	}
}

// Path traversal in depends_on should not crash validation.
func TestPathTraversalDependsOn_NoPanic(t *testing.T) {
	tmp := t.TempDir()
	content := `---
title: "Traversal"
doc_id: "trav"
doc_type: "guide"
status: "draft"
depends_on: ["../outside"]
---
# T`
	if err := os.WriteFile(filepath.Join(tmp, "trav.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := validator.ValidateDocs(tmp, nil, &validator.Options{}); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
}
