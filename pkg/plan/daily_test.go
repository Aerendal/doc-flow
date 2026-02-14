package plan

import (
	"testing"

	"docflow/pkg/index"
)

func TestPlanDailyTopoOrder(t *testing.T) {
	idx := index.New()
	idx.Add(&index.DocumentRecord{DocID: "a", Path: "a.md", DependsOn: []string{}, Lines: 80, HeadingCount: 4, MaxDepth: 2})
	idx.Add(&index.DocumentRecord{DocID: "b", Path: "b.md", DependsOn: []string{"a"}, Lines: 120, HeadingCount: 6, MaxDepth: 3})

	items, err := PlanDaily(idx, 0)
	if err != nil {
		t.Fatalf("PlanDaily error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
	if items[0].DocID != "a" || items[1].DocID != "b" {
		t.Fatalf("order = [%s,%s], want [a,b]", items[0].DocID, items[1].DocID)
	}
}

func TestPlanDailyLimit(t *testing.T) {
	idx := index.New()
	idx.Add(&index.DocumentRecord{DocID: "a", Path: "a.md"})
	idx.Add(&index.DocumentRecord{DocID: "b", Path: "b.md"})

	items, _ := PlanDaily(idx, 1)
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
}
