package governance

import (
	"os"
	"testing"

	"docflow/pkg/recommend"
)

func TestValidateGovernance(t *testing.T) {
	templates := []recommend.TemplateInfo{
		{Path: "good.md", Quality: 80, DocType: "spec", Language: "pl"},
		{Path: "bad.md", Quality: 40, DocType: "", Language: ""},
	}
	rule := Rule{MinQuality: 50}
	rep := Validate(templates, rule)
	if len(rep.Violations) < 2 {
		t.Fatalf("expected violations for bad template")
	}
}

func TestLoadRules(t *testing.T) {
	yaml := []byte(`
statuses:
  draft:
    required_fields: [title]
families:
  guide:
    required_sections: ["Intro"]
`)
	tmp := t.TempDir() + "/rules.yaml"
	if err := os.WriteFile(tmp, yaml, 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Load(tmp)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(r.Statuses) != 1 || len(r.Families) != 1 {
		t.Fatalf("unexpected rules content")
	}
}
