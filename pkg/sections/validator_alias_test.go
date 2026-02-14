package sections

import (
	"testing"

	"docflow/pkg/schema"
)

func TestValidateAliases(t *testing.T) {
	content := `# T
## Zakres
`
	sch := &schema.SectionSchema{
		ID: "guide",
		Sections: []schema.SectionRule{
			{Name: "Scope", Level: 2, Required: true},
		},
	}
	aliases := map[string]string{
		"scope": "zakres",
	}
	res, err := ValidateContentWithAliases(content, sch, aliases)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	if len(res.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", res.Missing)
	}
}
