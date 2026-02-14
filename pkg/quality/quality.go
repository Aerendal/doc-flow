package quality

import (
	"docflow/pkg/schema"
	"docflow/pkg/sections"
)

// Score holds component scores and total 0-100.
type Score struct {
	Total        int
	Structure    int
	Completeness int
	Meta         int
	Missing      []string
}

// ScoreContent calculates a simple quality score for a template/document content
// against a section schema. If schema is nil, only structural heuristics based
// on headings depth/amount are used.
func ScoreContent(content string, sch *schema.SectionSchema) *Score {
	const (
		maxStructure    = 60
		maxCompleteness = 30
		maxMeta         = 10
	)

	score := &Score{}

	if sch != nil {
		res, _ := sections.ValidateContent(content, sch)
		if res != nil {
			score.Missing = res.Missing
			required := requiredCount(sch)
			missing := len(res.Missing)
			if required == 0 {
				score.Structure = maxStructure
			} else {
				score.Structure = int(float64(maxStructure) * (1.0 - float64(missing)/float64(required)))
				if score.Structure < 0 {
					score.Structure = 0
				}
			}
			score.Completeness = maxCompleteness - missing*5
			if score.Completeness < 0 {
				score.Completeness = 0
			}
		}
	} else {
		// Heuristic fallback: more headings → better structure.
		tree := sections.Parse(content)
		count := len(tree.Root.Children)
		switch {
		case count >= 8:
			score.Structure = maxStructure
			score.Completeness = maxCompleteness
		case count >= 4:
			score.Structure = int(0.7 * float64(maxStructure))
			score.Completeness = int(0.7 * float64(maxCompleteness))
		case count >= 2:
			score.Structure = int(0.4 * float64(maxStructure))
			score.Completeness = int(0.4 * float64(maxCompleteness))
		default:
			score.Structure = int(0.2 * float64(maxStructure))
			score.Completeness = int(0.2 * float64(maxCompleteness))
		}
	}

	score.Meta = maxMeta // placeholder; brak metryk meta na tym etapie
	score.Total = clamp(score.Structure+score.Completeness+score.Meta, 0, 100)
	return score
}

func requiredCount(sch *schema.SectionSchema) int {
	n := 0
	for _, s := range sch.Sections {
		if s.Required {
			n++
		}
	}
	return n
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
