package sections

import (
	"testing"

	"docflow/pkg/schema"
)

func TestComputeCompletenessWithSchema(t *testing.T) {
	sch := &schema.SectionSchema{
		ID: "guide",
		Sections: []schema.SectionRule{
			{Name: "Intro", Level: 2, Required: true},
			{Name: "Steps", Level: 2, Required: true},
			{Name: "FAQ", Level: 2, Required: false},
		},
	}

	content := `# Title
## Intro
Some text
## Steps
Details
`

	res := ComputeCompleteness(content, sch, nil)
	if res.RequiredPresent != 2 || res.RequiredTotal != 2 {
		t.Fatalf("required %d/%d, want 2/2", res.RequiredPresent, res.RequiredTotal)
	}
	if res.OptionalPresent != 0 || res.OptionalTotal != 1 {
		t.Fatalf("optional %d/%d, want 0/1", res.OptionalPresent, res.OptionalTotal)
	}
	if res.Completeness != 1.0 {
		t.Fatalf("completeness=%f, want 1.0", res.Completeness)
	}
}

func TestComputeCompletenessMissing(t *testing.T) {
	sch := &schema.SectionSchema{
		ID: "guide",
		Sections: []schema.SectionRule{
			{Name: "Intro", Level: 2, Required: true},
			{Name: "Steps", Level: 2, Required: true},
		},
	}
	content := `# Title
## Intro
`
	res := ComputeCompleteness(content, sch, nil)
	if len(res.Missing) != 1 || res.Missing[0] != "Steps" {
		t.Fatalf("missing %+v, want [Steps]", res.Missing)
	}
	if res.Completeness >= 1.0 {
		t.Fatalf("completeness should be <1.0")
	}
}
