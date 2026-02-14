package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateWithPreview(t *testing.T) {
	dir := t.TempDir()
	tmpl := filepath.Join(dir, "tmpl.md")
	out := filepath.Join(dir, "out.md")
	os.WriteFile(tmpl, []byte("# {{DOC_ID}}\nType: {{DOC_TYPE}}\nBody"), 0o644)

	item := DocItem{
		DocID:        "doc-1",
		DocType:      "guide",
		TemplatePath: tmpl,
		OutputPath:   out,
	}
	prev, err := Generate(item, 2)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if !strings.Contains(prev, "doc-1") || !strings.Contains(prev, "guide") {
		t.Fatalf("preview missing replacements: %s", prev)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "doc-1") {
		t.Fatalf("output not written with doc_id")
	}
}
