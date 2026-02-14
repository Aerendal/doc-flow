package pattern

// MasterSelection selects a master template from a group of near-duplicates.
// Returns index of master in the slice.
func MasterSelection(group []*TemplatePattern, quality map[string]int, usage map[string]int) int {
	best := 0
	bestScore := -1
	for i, g := range group {
		q := quality[g.Path]
		u := usage[g.Path]
		score := q*2 + u // simple heuristic
		if score > bestScore || (score == bestScore && g.Path < group[best].Path) {
			bestScore = score
			best = i
		}
	}
	return best
}
