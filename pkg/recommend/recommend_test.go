package recommend

import "testing"

func TestRecommendRanksByWeights(t *testing.T) {
	templates := []TemplateInfo{
		{Path: "a.md", DocType: "guide", Language: "pl", Quality: 80, Usage: 0},
		{Path: "b.md", DocType: "guide", Language: "pl", Quality: 90, Usage: 1},
		{Path: "c.md", DocType: "spec", Language: "en", Quality: 95, Usage: 5},
	}

	recs := Recommend(templates, Params{PreferredDocType: "guide", PreferredLang: "pl", TopN: 2})
	if len(recs) != 2 {
		t.Fatalf("len=%d, want 2", len(recs))
	}
	if recs[0].Template.Path != "b.md" {
		t.Fatalf("top template = %s, want b.md", recs[0].Template.Path)
	}
	if recs[1].Template.Path != "a.md" {
		t.Fatalf("second template = %s, want a.md", recs[1].Template.Path)
	}
}
