package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestHealthOverallExitPriority(t *testing.T) {
	if got := healthOverallExit(0, 1, 0); got != 1 {
		t.Fatalf("healthOverallExit(0,1,0)=%d want 1", got)
	}
	if got := healthOverallExit(0, 1, 2, 0); got != 2 {
		t.Fatalf("healthOverallExit(0,1,2,0)=%d want 2", got)
	}
	if got := healthOverallExit(0, 1, 3, 2); got != 3 {
		t.Fatalf("healthOverallExit(0,1,3,2)=%d want 3", got)
	}
}

func TestShouldStopHealth(t *testing.T) {
	if shouldStopHealth(1, false) {
		t.Fatal("domain exit (1) should not stop by default")
	}
	if !shouldStopHealth(2, false) {
		t.Fatal("usage exit (2) should stop by default")
	}
	if !shouldStopHealth(3, false) {
		t.Fatal("runtime exit (3) should stop by default")
	}
	if shouldStopHealth(3, true) {
		t.Fatal("continue-on-error should disable stopping")
	}
}

func TestHealthCmdCIGeneratesBundleAndReturnsDomainExit(t *testing.T) {
	tmp := t.TempDir()
	docsDir := filepath.Join(tmp, "docs")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Missing version -> domain validation error in --ci strict mode.
	doc := `---
doc_id: a_doc
title: A
doc_type: guide
status: draft
owner: team-docs
language: pl
priority: P2
context_sources: [seed]
---
# A
`
	if err := os.WriteFile(filepath.Join(docsDir, "a.md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(tmp, "docflow.yaml")
	cfg := "docs_root: " + docsDir + "\n" +
		"log:\n  level: error\n  format: text\n" +
		"cache:\n  enabled: false\n  dir: .docflow/cache\n  max_size_mb: 10\n  max_age_days: 1\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	rulesPath := filepath.Join(tmp, "rules.yaml")
	rules := "statuses:\n  default:\n    required_fields: []\n    allow_empty_sections: true\nfamilies: {}\napprovals: {}\n"
	if err := os.WriteFile(rulesPath, []byte(rules), 0o644); err != nil {
		t.Fatal(err)
	}

	restore := saveAndSetCLIState(cfgPath)
	defer restore()

	bundleDir := filepath.Join(tmp, "out")
	runID := "run-ci-test"
	cmd := healthCmd()
	cmd.SetArgs([]string{
		"--ci",
		"--baseline-mode", "none",
		"--bundle-dir", bundleDir,
		"--run-id", runID,
		"--rules", rulesPath,
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected domain error for ci strict run, got nil")
	}
	if got := ExitCode(err); got != 1 {
		t.Fatalf("ExitCode()=%d want 1 err=%v", got, err)
	}

	runDir := filepath.Join(bundleDir, runID)
	required := []string{
		"validate.json",
		"validate.sarif",
		"compliance.json",
		"summary.json",
		"summary.md",
		"meta.json",
	}
	for _, name := range required {
		if _, statErr := os.Stat(filepath.Join(runDir, name)); statErr != nil {
			t.Fatalf("expected artifact %s: %v", name, statErr)
		}
	}

	summaryData, readErr := os.ReadFile(filepath.Join(runDir, "summary.json"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	var summary healthSummary
	if unmarshalErr := json.Unmarshal(summaryData, &summary); unmarshalErr != nil {
		t.Fatalf("summary.json invalid json: %v", unmarshalErr)
	}
	if summary.OverallExit != 1 {
		t.Fatalf("summary.overall_exit=%d want 1", summary.OverallExit)
	}
	if summary.Validate.JSONExit != 1 {
		t.Fatalf("summary.validate.json_exit=%d want 1", summary.Validate.JSONExit)
	}
	if summary.Compliance.JSONExit != 0 {
		t.Fatalf("summary.compliance.json_exit=%d want 0", summary.Compliance.JSONExit)
	}
}
