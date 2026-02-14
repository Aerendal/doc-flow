package plan

import (
	"docflow/pkg/deps"
	"docflow/pkg/index"
)

type DailyItem struct {
	DocID  string
	Path   string
	Effort int
}

// PlanDaily returns docs in topo order with heuristic effort, limited by max (0 = all).
func PlanDaily(idx *index.DocumentIndex, max int) ([]DailyItem, error) {
	depsMap := make(map[string][]string)
	for _, rec := range idx.All() {
		if rec.DocID == "" {
			continue
		}
		depsMap[rec.DocID] = rec.DependsOn
	}

	graph := deps.Build(depsMap)
	topo := graph.Topo()

	byID := make(map[string]*index.DocumentRecord)
	for _, r := range idx.All() {
		if r.DocID != "" {
			byID[r.DocID] = r
		}
	}

	var items []DailyItem
	for _, docID := range topo.Order {
		rec := byID[docID]
		if rec == nil {
			continue
		}
		items = append(items, DailyItem{
			DocID:  docID,
			Path:   rec.Path,
			Effort: effort(rec),
		})
		if max > 0 && len(items) >= max {
			break
		}
	}

	return items, nil
}

func effort(rec *index.DocumentRecord) int {
	e := rec.Lines/80 + rec.HeadingCount/5 + rec.MaxDepth
	if e < 1 {
		e = 1
	}
	return e
}
