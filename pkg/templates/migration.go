package templates

import "docflow/pkg/recommend"

// SuggestMigration selects best active replacement for a deprecated template.
// Preference: same doc_type, highest quality, then usage.
func SuggestMigration(deprecated recommend.TemplateInfo, candidates []recommend.TemplateInfo) *recommend.TemplateInfo {
	var best *recommend.TemplateInfo
	bestScore := -1
	for i := range candidates {
		c := &candidates[i]
		if c.Status == "deprecated" || c.Status == "archived" {
			continue
		}
		if deprecated.DocType != "" && c.DocType != deprecated.DocType {
			continue
		}
		score := c.Quality*2 + c.Usage
		if score > bestScore {
			bestScore = score
			best = c
		}
	}
	return best
}
