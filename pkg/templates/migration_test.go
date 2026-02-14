package templates

import (
	"testing"

	"docflow/pkg/recommend"
)

func TestSuggestMigration(t *testing.T) {
	dep := recommend.TemplateInfo{Path: "old.md", DocType: "guide", Status: "deprecated"}
	candidates := []recommend.TemplateInfo{
		{Path: "a.md", DocType: "guide", Quality: 70, Usage: 1, Status: "active"},
		{Path: "b.md", DocType: "guide", Quality: 80, Usage: 0, Status: "active"},
		{Path: "c.md", DocType: "spec", Quality: 90, Usage: 10, Status: "active"},
	}
	best := SuggestMigration(dep, candidates)
	if best == nil || best.Path != "b.md" {
		t.Fatalf("expected b.md as migration, got %+v", best)
	}
}
