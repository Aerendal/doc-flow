package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestScanCmdUsesProvidedPath(t *testing.T) {
	// Prepare a temporary docs root with one valid document.
	tmp := t.TempDir()
	docs := filepath.Join(tmp, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := []byte(`---
doc_id: test_doc
doc_type: guide
status: draft
version: v1
context_sources: ["seed"]
---
`)
	if err := os.WriteFile(filepath.Join(docs, "doc.md"), doc, 0o644); err != nil {
		t.Fatal(err)
	}

	// Config points to a non-existing path; we expect CLI to override it when arg provided.
	cfgPath := filepath.Join(tmp, "docflow.yaml")
	if err := os.WriteFile(cfgPath, []byte("docs_root: /path/does/not/exist\nlog:\n  level: error\n  format: text\ncache:\n  enabled: false\n  dir: .docflow/cache\n  max_size_mb: 10\n  max_age_days: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Override global config path for the command.
	prevConfig := configFile
	prevLogLevel := logLevel
	configFile = cfgPath
	logLevel = "error"
	t.Cleanup(func() {
		configFile = prevConfig
		logLevel = prevLogLevel
	})

	cmd := scanCmd()
	cmd.SetArgs([]string{docs})

	// Capture stdout/stderr because scanCmd prints with fmt.Printf/log.
	stdout := os.Stdout
	stderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	if err := cmd.Execute(); err != nil {
		wOut.Close()
		wErr.Close()
		os.Stdout = stdout
		os.Stderr = stderr
		t.Fatalf("scanCmd.Execute() error = %v", err)
	}
	wOut.Close()
	wErr.Close()
	os.Stdout = stdout
	os.Stderr = stderr

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(rOut); err != nil {
		t.Fatal(err)
	}
	if _, err := buf.ReadFrom(rErr); err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Razem: 1 dokumentów")) {
		t.Fatalf("expected output to mention 1 document, got: %s", buf.String())
	}
}
