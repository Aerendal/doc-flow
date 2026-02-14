package pattern

import "testing"

func TestSectionJaccard(t *testing.T) {
	a := []SectionEntry{
		{Level: 2, Text: "Cel"},
		{Level: 2, Text: "Zakres"},
		{Level: 3, Text: "Wymagania"},
	}
	b := []SectionEntry{
		{Level: 2, Text: "Cel"},
		{Level: 2, Text: "Zakres"},
		{Level: 3, Text: "Ryzyka"},
	}

	sim := SectionJaccard(a, b)
	if sim < 0.5 || sim >= 1.0 {
		t.Fatalf("similarity = %f, want >=0.5 and <1.0", sim)
	}
}

func TestFindSimilar(t *testing.T) {
	patterns := []*TemplatePattern{
		{Path: "a.md", Sections: []SectionEntry{{2, "A"}, {2, "B"}}},
		{Path: "b.md", Sections: []SectionEntry{{2, "A"}, {2, "B"}}},
		{Path: "c.md", Sections: []SectionEntry{{2, "X"}}},
	}

	pairs := FindSimilar(patterns, 0.9)
	if len(pairs) != 1 {
		t.Fatalf("pairs len = %d, want 1", len(pairs))
	}
	if pairs[0].APath != "a.md" || pairs[0].BPath != "b.md" {
		t.Fatalf("unexpected pair: %+v", pairs[0])
	}
}
