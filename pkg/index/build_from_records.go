package index

import "sort"

// BuildIndexFromRecords builds a deterministic index from precomputed records.
// This is used by execution-layer pipelines that obtain facts from cache/parsing.
func BuildIndexFromRecords(root string, records []*DocumentRecord) *DocumentIndex {
	idx := New()
	idx.Root = NormalizePath(root)

	adj := map[string][]string{}
	for _, rec := range records {
		if rec == nil {
			continue
		}
		cp := *rec
		cp.Path = NormalizePath(cp.Path)
		cp.DependsOn = append([]string(nil), cp.DependsOn...)
		cp.ContextSources = append([]string(nil), cp.ContextSources...)
		cp.Tags = append([]string(nil), cp.Tags...)
		idx.Add(&cp)
		adj[cp.DocID] = append([]string(nil), cp.DependsOn...)
	}

	markCycles(idx, adj)

	idx.mu.Lock()
	sort.Slice(idx.Documents, func(i, j int) bool {
		return idx.Documents[i].Path < idx.Documents[j].Path
	})
	idx.Count = len(idx.Documents)
	idx.Checksum = idx.computeChecksum()
	idx.mu.Unlock()

	return idx
}
