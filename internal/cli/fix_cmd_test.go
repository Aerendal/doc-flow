package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"docflow/internal/validator"
)

func TestFixCmdApplyReducesErrors(t *testing.T) {
	tmp := t.TempDir()
	docs := filepath.Join(tmp, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}

	docPath := filepath.Join(docs, "a.md")
	doc := `---
doc_id: a_doc
title: A
doc_type: guide
status: draft
context_sources: [seed]
---
# A
## Overview
`
	if err := os.WriteFile(docPath, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(tmp, "docflow.yaml")
	cfg := "docs_root: " + docs + "\n" +
		"log:\n  level: error\n  format: text\n" +
		"cache:\n  enabled: false\n  dir: .docflow/cache\n  max_size_mb: 10\n  max_age_days: 1\n" +
		"section_aliases:\n  Przegląd:\n    - Overview\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	restore := saveAndSetCLIState(cfgPath)
	defer restore()

	before, err := collectIssuesForTest()
	if err != nil {
		t.Fatal(err)
	}
	if !containsIssue(before, "missing_field", "a.md") {
		t.Fatalf("expected missing_field before apply, got=%v", before)
	}
	if !containsIssue(before, "legacy_section_name", "a.md") {
		t.Fatalf("expected legacy_section_name before apply, got=%v", before)
	}

	out, err := runFixCommandForTest([]string{"--format", "diff", "--output", "-", "--apply"})
	if err != nil {
		t.Fatalf("fix --apply failed: %v", err)
	}
	if !strings.Contains(out, "version: v0.0.0") {
		t.Fatalf("expected version fix in diff output, got:\n%s", out)
	}

	after, err := collectIssuesForTest()
	if err != nil {
		t.Fatal(err)
	}
	if containsIssue(after, "missing_field", "a.md") {
		t.Fatalf("missing_field should be fixed, got=%v", after)
	}
	if containsIssue(after, "legacy_section_name", "a.md") {
		t.Fatalf("legacy_section_name should be fixed, got=%v", after)
	}
}

func TestFixCmdDiffDeterministicNoCacheVsCache(t *testing.T) {
	tmp := t.TempDir()
	docs := filepath.Join(tmp, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `---
doc_id: a_doc
title: A
doc_type: guide
status: draft
context_sources: [seed]
---
# A
`
	if err := os.WriteFile(filepath.Join(docs, "a.md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(tmp, "docflow.yaml")
	cfg := "docs_root: " + docs + "\n" +
		"log:\n  level: error\n  format: text\n" +
		"cache:\n  enabled: true\n  dir: .docflow/cache\n  max_size_mb: 10\n  max_age_days: 1\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	restore := saveAndSetCLIState(cfgPath)
	defer restore()
	cacheDir = filepath.Join(tmp, "cache")
	cacheEnabled = true

	noCache = true
	diffNoCache, err := runFixCommandForTest([]string{"--format", "diff", "--output", "-"})
	if err != nil {
		t.Fatalf("fix no-cache failed: %v", err)
	}

	noCache = false
	diffCache, err := runFixCommandForTest([]string{"--format", "diff", "--output", "-"})
	if err != nil {
		t.Fatalf("fix cache failed: %v", err)
	}

	if diffNoCache != diffCache {
		t.Fatalf("fix diff should be identical for no-cache/cache\nno-cache:\n%s\ncache:\n%s", diffNoCache, diffCache)
	}
}

func saveAndSetCLIState(cfgPath string) func() {
	prevConfig := configFile
	prevCacheEnabled := cacheEnabled
	prevNoCache := noCache
	prevCacheDir := cacheDir
	prevChangedOnly := changedOnly
	prevSinceRef := sinceRef
	prevLogLevel := logLevel
	prevLogFormat := logFormat

	configFile = cfgPath
	cacheEnabled = true
	noCache = false
	cacheDir = ""
	changedOnly = false
	sinceRef = ""
	logLevel = "error"
	logFormat = "text"

	return func() {
		configFile = prevConfig
		cacheEnabled = prevCacheEnabled
		noCache = prevNoCache
		cacheDir = prevCacheDir
		changedOnly = prevChangedOnly
		sinceRef = prevSinceRef
		logLevel = prevLogLevel
		logFormat = prevLogFormat
	}
}

func runFixCommandForTest(args []string) (string, error) {
	cmd := fixCmd()
	cmd.SetArgs(args)

	stdout := os.Stdout
	stderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	err := cmd.Execute()
	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = stdout
	os.Stderr = stderr

	var outBuf bytes.Buffer
	_, _ = outBuf.ReadFrom(rOut)
	var errBuf bytes.Buffer
	_, _ = errBuf.ReadFrom(rErr)
	if err == nil && errBuf.Len() > 0 {
		outBuf.Write(errBuf.Bytes())
	}
	return outBuf.String(), err
}

func collectIssuesForTest() ([]validator.Issue, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	execRes, err := collectExecutionResult("fix-test", cfg)
	if err != nil {
		return nil, err
	}
	report, err := validator.ValidateFacts(execRes.Facts, execRes.Index, validatorOptions(cfg))
	if err != nil {
		return nil, err
	}
	return report.SortedIssues(), nil
}

func containsIssue(issues []validator.Issue, typ, pathSuffix string) bool {
	for _, issue := range issues {
		if string(issue.Type) == typ && strings.HasSuffix(issue.File, pathSuffix) {
			return true
		}
	}
	return false
}
