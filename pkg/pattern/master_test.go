package pattern

import "testing"

func TestMasterSelection(t *testing.T) {
	group := []*TemplatePattern{
		{Path: "a.md"},
		{Path: "b.md"},
	}
	quality := map[string]int{"a.md": 80, "b.md": 90}
	usage := map[string]int{"a.md": 1, "b.md": 0}

	m := MasterSelection(group, quality, usage)
	if m != 1 {
		t.Fatalf("master = %d, want 1", m)
	}
}
