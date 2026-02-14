package templates

import (
	"testing"

	"docflow/pkg/recommend"
)

func TestEvaluateDeprecation(t *testing.T) {
	templates := []recommend.TemplateInfo{
		{Path: "good.md", Quality: 80, Usage: 2, Status: "active"},
		{Path: "low.md", Quality: 40, Usage: 0, Status: "active"},
		{Path: "used.md", Quality: 45, Usage: 5, Status: "active"},
	}
	rules := Rules{MinQuality: 50, UnusedOnly: false}
	rep := EvaluateDeprecation(templates, rules)
	if len(rep.ToDeprecate) != 2 {
		t.Fatalf("to deprecate = %d, want 2", len(rep.ToDeprecate))
	}

	rulesUnused := Rules{MinQuality: 50, UnusedOnly: true}
	rep2 := EvaluateDeprecation(templates, rulesUnused)
	if len(rep2.ToDeprecate) != 1 || rep2.ToDeprecate[0].Path != "low.md" {
		t.Fatalf("unused filter failed: %+v", rep2.ToDeprecate)
	}
}
