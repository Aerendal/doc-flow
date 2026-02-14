package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.DocsRoot != "." {
		t.Errorf("expected docs_root='.', got %q", cfg.DocsRoot)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected log level 'info', got %q", cfg.Log.Level)
	}
	if !cfg.Cache.Enabled {
		t.Error("expected cache enabled by default")
	}
	if len(cfg.IgnorePatterns) == 0 {
		t.Error("expected default ignore patterns")
	}
}

func TestValidate(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}

	cfg.DocsRoot = ""
	if err := cfg.Validate(); err == nil {
		t.Error("empty docs_root should fail validation")
	}

	cfg = Default()
	cfg.Log.Level = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid log level should fail validation")
	}

	cfg = Default()
	cfg.Cache.MaxSizeMB = -1
	if err := cfg.Validate(); err == nil {
		t.Error("negative cache size should fail validation")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "docflow.yaml")

	content := `docs_root: "testdir"
ignore_patterns:
  - ".git"
  - "vendor"
log:
  level: "debug"
  format: "json"
cache:
  enabled: false
  dir: ".cache"
  max_size_mb: 50
  max_age_days: 7
  checksum_algorithm: "sha256"
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.DocsRoot != "testdir" {
		t.Errorf("expected docs_root='testdir', got %q", cfg.DocsRoot)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level 'debug', got %q", cfg.Log.Level)
	}
	if cfg.Log.Format != "json" {
		t.Errorf("expected log format 'json', got %q", cfg.Log.Format)
	}
	if cfg.Cache.Enabled {
		t.Error("expected cache disabled")
	}
	if cfg.Cache.MaxSizeMB != 50 {
		t.Errorf("expected max_size_mb=50, got %d", cfg.Cache.MaxSizeMB)
	}
}

func TestLoadMissing(t *testing.T) {
	cfg, err := Load("/nonexistent/docflow.yaml")
	if err != nil {
		t.Fatalf("missing config should return defaults, got error: %v", err)
	}
	if cfg.DocsRoot != "." {
		t.Errorf("expected default docs_root, got %q", cfg.DocsRoot)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "bad.yaml")
	os.WriteFile(cfgPath, []byte("{{{{invalid"), 0644)

	_, err := Load(cfgPath)
	if err == nil {
		t.Error("invalid YAML should return error")
	}
}
