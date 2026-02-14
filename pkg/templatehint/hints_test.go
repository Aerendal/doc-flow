package templatehint

import "testing"

func TestExtractHints(t *testing.T) {
	content := "```\ncode\n```\n| col | col |\n| --- | --- |\n"
	h := Extract(content)
	if h.CodeBlocks != 1 {
		t.Fatalf("codeblocks=%d, want 1", h.CodeBlocks)
	}
	if h.Tables == 0 {
		t.Fatalf("tables=%d, want >0", h.Tables)
	}
}
