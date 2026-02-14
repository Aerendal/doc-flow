package quality

import (
	"testing"

	"docflow/pkg/schema"
)

func TestScoreContentWithSchema(t *testing.T) {
	content := `# Doc
## Cel
## Zakres
`
	sch := &schema.SectionSchema{
		ID: "req",
		Sections: []schema.SectionRule{
			{Name: "Cel", Level: 2, Required: true},
			{Name: "Zakres", Level: 2, Required: true},
			{Name: "Wymagania", Level: 2, Required: true},
		},
	}

	score := ScoreContent(content, sch)
	if score.Structure == 0 || score.Total == 0 {
		t.Fatalf("expected positive score, got %+v", score)
	}
	if len(score.Missing) != 1 || score.Missing[0] != "Wymagania" {
		t.Fatalf("missing = %v, want [Wymagania]", score.Missing)
	}
}

func TestScoreContentNoSchema(t *testing.T) {
	content := `# Doc
## A
## B
`
	score := ScoreContent(content, nil)
	if score.Total == 0 {
		t.Fatal("expected non-zero score without schema")
	}
}
