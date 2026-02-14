package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGraphCmdJSONIdenticalNoCacheVsCache(t *testing.T) {
	tmp := t.TempDir()
	docs := filepath.Join(tmp, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}

	docA := []byte(`---
doc_id: A
doc_type: guide
status: draft
depends_on: [B]
---
# A
`)
	docB := []byte(`---
doc_id: B
doc_type: guide
status: draft
---
# B
`)
	if err := os.WriteFile(filepath.Join(docs, "a.md"), docA, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docs, "b.md"), docB, 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(tmp, "docflow.yaml")
	cfg := "docs_root: " + docs + "\n" +
		"log:\n  level: error\n  format: text\n" +
		"cache:\n  enabled: true\n  dir: .docflow/cache\n  max_size_mb: 10\n  max_age_days: 1\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	prevConfig := configFile
	prevCacheEnabled := cacheEnabled
	prevNoCache := noCache
	prevCacheDir := cacheDir
	prevChangedOnly := changedOnly
	prevSinceRef := sinceRef
	t.Cleanup(func() {
		configFile = prevConfig
		cacheEnabled = prevCacheEnabled
		noCache = prevNoCache
		cacheDir = prevCacheDir
		changedOnly = prevChangedOnly
		sinceRef = prevSinceRef
	})

	configFile = cfgPath
	cacheEnabled = true
	cacheDir = filepath.Join(tmp, "cache")
	changedOnly = false
	sinceRef = ""

	noCache = true
	noCacheJSON, err := runGraphCommandForTest(t, []string{"--format", "json", "--output", "-"})
	if err != nil {
		t.Fatalf("graph no-cache run failed: %v", err)
	}

	noCache = false
	cacheJSON, err := runGraphCommandForTest(t, []string{"--format", "json", "--output", "-"})
	if err != nil {
		t.Fatalf("graph cache run failed: %v", err)
	}

	if noCacheJSON != cacheJSON {
		t.Fatalf("graph output differs between no-cache and cache\nno-cache:\n%s\ncache:\n%s", noCacheJSON, cacheJSON)
	}
}

func runGraphCommandForTest(t *testing.T, args []string) (string, error) {
	t.Helper()
	cmd := graphCmd()
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
	if _, readErr := outBuf.ReadFrom(rOut); readErr != nil {
		t.Fatalf("read stdout: %v", readErr)
	}
	var errBuf bytes.Buffer
	if _, readErr := errBuf.ReadFrom(rErr); readErr != nil {
		t.Fatalf("read stderr: %v", readErr)
	}
	if err == nil && errBuf.Len() > 0 {
		t.Logf("graph stderr: %s", errBuf.String())
	}
	return outBuf.String(), err
}
