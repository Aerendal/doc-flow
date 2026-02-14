package recommend

import "testing"

func TestPrecisionAtK(t *testing.T) {
	recs := []TemplateInfo{
		{Path: "a"},
		{Path: "b"},
		{Path: "c"},
	}
	rel := map[string]bool{"a": true, "c": true}
	p := PrecisionAtK(recs, rel, 2)
	if p != 0.5 {
		t.Fatalf("precision=%f, want 0.5", p)
	}
}
