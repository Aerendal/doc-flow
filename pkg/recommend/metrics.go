package recommend

// PrecisionAtK computes precision for top-k recommendations given a set of relevant items.
func PrecisionAtK(recs []TemplateInfo, relevant map[string]bool, k int) float64 {
	if k <= 0 || len(recs) == 0 {
		return 0
	}
	if k > len(recs) {
		k = len(recs)
	}
	hits := 0
	for i := 0; i < k; i++ {
		if relevant[recs[i].Path] {
			hits++
		}
	}
	return float64(hits) / float64(k)
}
