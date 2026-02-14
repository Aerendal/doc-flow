package util

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	dirs := []string{
		"docs",
		"docs/sub",
		".git",
		"vendor",
		"node_modules",
	}
	for _, d := range dirs {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}

	files := map[string]string{
		"README.md":            "# README",
		"docs/guide.md":        "# Guide",
		"docs/sub/nested.md":   "# Nested",
		"docs/image.png":       "binary",
		"docs/notes.txt":       "text",
		".git/config":          "[core]",
		"vendor/lib.md":        "# Vendor lib",
		"node_modules/pkg.md":  "# Node pkg",
	}
	for name, content := range files {
		os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
	}

	return dir
}

func TestWalkMarkdown(t *testing.T) {
	dir := setupTestDir(t)

	files, err := WalkMarkdown(dir, []string{".git", "vendor", "node_modules"})
	if err != nil {
		t.Fatalf("WalkMarkdown failed: %v", err)
	}

	paths := make(map[string]bool)
	for _, f := range files {
		paths[f.RelPath] = true
	}

	expected := []string{"README.md", "docs/guide.md", "docs/sub/nested.md"}
	for _, e := range expected {
		if !paths[e] {
			t.Errorf("expected file %q not found", e)
		}
	}

	notExpected := []string{"docs/image.png", "docs/notes.txt", ".git/config", "vendor/lib.md", "node_modules/pkg.md"}
	for _, ne := range notExpected {
		if paths[ne] {
			t.Errorf("file %q should be ignored", ne)
		}
	}
}

func TestWalkMarkdownEmpty(t *testing.T) {
	dir := t.TempDir()
	files, err := WalkMarkdown(dir, nil)
	if err != nil {
		t.Fatalf("WalkMarkdown failed: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestWalkResult(t *testing.T) {
	dir := setupTestDir(t)
	files, err := WalkMarkdown(dir, []string{".git", "vendor", "node_modules"})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.Path == "" {
			t.Error("Path should not be empty")
		}
		if f.RelPath == "" {
			t.Error("RelPath should not be empty")
		}
		if f.Size <= 0 {
			t.Errorf("Size should be > 0, got %d for %s", f.Size, f.RelPath)
		}
		if f.ModTime <= 0 {
			t.Errorf("ModTime should be > 0 for %s", f.RelPath)
		}
	}
}
