package sections

import (
	"testing"

	"docflow/pkg/schema"
)

func TestValidateTreeMissing(t *testing.T) {
	content := `# T
## Cel
## Zakres
`
	tree := Parse(content)
	sch := &schema.SectionSchema{
		ID: "req",
		Sections: []schema.SectionRule{
			{Name: "Cel", Level: 2, Required: true},
			{Name: "Zakres", Level: 2, Required: true},
			{Name: "Wymagania", Level: 2, Required: true},
		},
	}

	res := ValidateTree(tree, sch)
	if len(res.Missing) != 1 || res.Missing[0] != "Wymagania" {
		t.Fatalf("missing = %v, want [Wymagania]", res.Missing)
	}
}

func TestValidateContentNoMissing(t *testing.T) {
	content := `# Doc
## Cel
## Zakres
## Wymagania
### Niefunkcjonalne
`
	sch := &schema.SectionSchema{
		ID: "req",
		Sections: []schema.SectionRule{
			{Name: "Cel", Level: 2, Required: true},
			{Name: "Zakres", Level: 2, Required: true},
			{Name: "Wymagania", Level: 2, Required: true},
		},
	}

	res, err := ValidateContent(content, sch)
	if err != nil {
		t.Fatalf("ValidateContent error: %v", err)
	}
	if len(res.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", res.Missing)
	}
}
