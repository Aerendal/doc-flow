package content

import "testing"

func TestExtract(t *testing.T) {
	md := "# T\n" +
		"text with [link](http://example.com)\n" +
		"![img](a.png)\n\n" +
		"| a | b |\n|---|---|\n| 1 | 2 |\n\n" +
		"```\ncode\n```\n"
	m := Extract(md)
	if m.CodeBlocks != 1 || m.Tables != 3 { // counts table lines
		t.Fatalf("code=%d tables=%d", m.CodeBlocks, m.Tables)
	}
	if m.Links != 2 || m.Images != 1 {
		t.Fatalf("links=%d images=%d", m.Links, m.Images)
	}
	if m.Words == 0 {
		t.Fatalf("words should be counted")
	}
}
