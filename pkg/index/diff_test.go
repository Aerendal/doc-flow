package index

import "testing"

func TestChangedTemplatesAndImpact(t *testing.T) {
	old := New()
	old.Add(&DocumentRecord{DocID: "tmpl_a", DocType: "template", Checksum: "aaa"})

	newIdx := New()
	newIdx.Add(&DocumentRecord{DocID: "tmpl_a", DocType: "template", Checksum: "bbb"})
	newIdx.Add(&DocumentRecord{DocID: "doc1", TemplateSource: "tmpl_a", Path: "doc1.md"})

	ch := ChangedTemplates(old, newIdx)
	if !ch["tmpl_a"] {
		t.Fatalf("tmpl_a should be changed")
	}
	impact := Impact(ch, newIdx)
	if len(impact["tmpl_a"]) != 1 {
		t.Fatalf("impact should include doc1")
	}
}
