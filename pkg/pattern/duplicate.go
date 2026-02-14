package pattern

import (
	"fmt"
	"strings"
)

// DuplicatePair represents similarity between two templates.
type DuplicatePair struct {
	APath string
	BPath string
	Score float64
}

// SectionJaccard computes Jaccard similarity on sets of "level:text" for h2+ sections.
func SectionJaccard(a, b []SectionEntry) float64 {
	setA := make(map[string]bool)
	setB := make(map[string]bool)

	for _, s := range h2PlusSections(a) {
		setA[key(s)] = true
	}
	for _, s := range h2PlusSections(b) {
		setB[key(s)] = true
	}

	if len(setA) == 0 && len(setB) == 0 {
		return 1.0
	}

	inter := 0
	union := make(map[string]bool)
	for k := range setA {
		union[k] = true
		if setB[k] {
			inter++
		}
	}
	for k := range setB {
		union[k] = true
	}

	if len(union) == 0 {
		return 0
	}
	return float64(inter) / float64(len(union))
}

// FindSimilar returns pairs of templates whose Jaccard similarity >= threshold.
func FindSimilar(patterns []*TemplatePattern, threshold float64) []DuplicatePair {
	var pairs []DuplicatePair
	for i := 0; i < len(patterns); i++ {
		for j := i + 1; j < len(patterns); j++ {
			sim := SectionJaccard(patterns[i].Sections, patterns[j].Sections)
			if sim >= threshold {
				pairs = append(pairs, DuplicatePair{
					APath: patterns[i].Path,
					BPath: patterns[j].Path,
					Score: sim,
				})
			}
		}
	}
	return pairs
}

func key(s SectionEntry) string {
	return fmt.Sprintf("%d:%s", s.Level, strings.ToLower(strings.TrimSpace(s.Text)))
}
