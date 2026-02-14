package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateSectionsDryRunAlias(t *testing.T) {
	tmp := t.TempDir()
	docs := filepath.Join(tmp, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Overview\n\n## Intro\n"
	mdPath := filepath.Join(docs, "sample.md")
	if err := os.WriteFile(mdPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(tmp, "docflow.yaml")
	cfg := `docs_root: ` + docs + `
section_aliases:
  Overview: ["Overview"]
log:
  level: error
  format: text
cache:
  enabled: false
  dir: .docflow/cache
  max_size_mb: 10
  max_age_days: 1
`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	prevCfg := configFile
	prevLog := logLevel
	configFile = cfgPath
	logLevel = "error"
	t.Cleanup(func() {
		configFile = prevCfg
		logLevel = prevLog
	})

	cmd := migrateSectionsCmd()
	cmd.SetArgs([]string{"--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	// Ensure file unchanged (dry-run alias forces apply=false).
	after, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != content {
		t.Fatalf("expected file unchanged in dry-run, got diff")
	}
}
