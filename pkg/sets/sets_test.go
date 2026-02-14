package sets

import "testing"

func TestCoOccurrence(t *testing.T) {
	usages := []TemplateUsage{
		{SetID: "p1", Template: "spec.md"},
		{SetID: "p1", Template: "guide.md"},
		{SetID: "p2", Template: "spec.md"},
		{SetID: "p2", Template: "runbook.md"},
		{SetID: "p3", Template: "spec.md"},
		{SetID: "p3", Template: "guide.md"},
	}

	co := CoOccurrence(usages, 2)
	if len(co["spec.md"]) == 0 {
		t.Fatalf("expected co-occurrence for spec.md")
	}
	foundGuide := false
	for _, t := range co["spec.md"] {
		if t == "guide.md" {
			foundGuide = true
		}
	}
	if !foundGuide {
		t.Fatalf("spec.md should co-occur with guide.md")
	}
}
