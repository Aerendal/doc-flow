package cache

import "docflow/pkg/index"

// InvalidateDependents returns doc_ids that should be invalidated when changedIDs changed.
// dependents map: key doc_id -> list of docs that depend on it (reverse graph).
func InvalidateDependents(changedIDs []string, dependents map[string][]string) []string {
	queue := make([]string, 0, len(changedIDs))
	seen := map[string]bool{}
	for _, id := range changedIDs {
		queue = append(queue, id)
		seen[id] = true
	}

	var result []string
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		result = append(result, curr)
		for _, dep := range dependents[curr] {
			if !seen[dep] {
				seen[dep] = true
				queue = append(queue, dep)
			}
		}
	}
	return result
}

// BuildReverseDeps builds reverse dependency map from DocumentIndex.
func BuildReverseDeps(idx *index.DocumentIndex) map[string][]string {
	rev := make(map[string][]string)
	for _, rec := range idx.All() {
		for _, d := range rec.DependsOn {
			rev[d] = append(rev[d], rec.DocID)
		}
	}
	return rev
}
