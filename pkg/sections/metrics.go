package sections

import (
	"regexp"
	"strings"

	"docflow/pkg/schema"
)

type SectionMetrics struct {
	Total        int
	Empty        int
	Filled       int
	Completeness float64 // 0-1
}

type CompletenessResult struct {
	RequiredPresent int
	RequiredTotal   int
	OptionalPresent int
	OptionalTotal   int
	Missing         []string
	Empty           int
	Completeness    float64 // requiredPresent/requiredTotal
}

// ComputeMetrics counts filled vs empty headings (text-only heuristic).
func ComputeMetrics(content string) *SectionMetrics {
	tree := Parse(content)
	total := 0
	empty := 0

	var walk func(n *SectionNode)
	walk = func(n *SectionNode) {
		if n.Level > 0 {
			total++
			if strings.TrimSpace(n.Text) == "" {
				empty++
			}
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(tree.Root)

	// blank headings that parser nie łapie (np. "## " bez tekstu)
	blank := countBlankHeadings(content)
	total += blank
	empty += blank

	filled := total - empty
	comp := 0.0
	if total > 0 {
		comp = float64(filled) / float64(total)
	}

	return &SectionMetrics{
		Total:        total,
		Empty:        empty,
		Filled:       filled,
		Completeness: comp,
	}
}

// ComputeCompleteness measures presence of required/optional sections per schema.
// If schema is nil, falls back to ComputeMetrics-derived completeness.
func ComputeCompleteness(content string, sch *schema.SectionSchema, aliases map[string]string) *CompletenessResult {
	if sch == nil {
		m := ComputeMetrics(content)
		return &CompletenessResult{
			RequiredPresent: m.Filled,
			RequiredTotal:   m.Total,
			OptionalPresent: 0,
			OptionalTotal:   0,
			Missing:         nil,
			Empty:           m.Empty,
			Completeness:    m.Completeness,
		}
	}

	tree := Parse(content)
	found := make(map[string]bool)
	var walk func(n *SectionNode)
	walk = func(n *SectionNode) {
		if n.Level > 0 {
			key := norm(n.Text)
			found[key] = true
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(tree.Root)

	reqTotal, reqPresent := 0, 0
	optTotal, optPresent := 0, 0
	var missing []string

	for _, rule := range sch.Sections {
		if rule.Required {
			reqTotal++
			if matched(found, rule.Name, aliases) {
				reqPresent++
			} else {
				missing = append(missing, rule.Name)
			}
		} else {
			optTotal++
			if matched(found, rule.Name, aliases) {
				optPresent++
			}
		}
	}

	m := ComputeMetrics(content)
	comp := 0.0
	if reqTotal > 0 {
		comp = float64(reqPresent) / float64(reqTotal)
	} else {
		comp = m.Completeness
	}

	return &CompletenessResult{
		RequiredPresent: reqPresent,
		RequiredTotal:   reqTotal,
		OptionalPresent: optPresent,
		OptionalTotal:   optTotal,
		Missing:         missing,
		Empty:           m.Empty,
		Completeness:    comp,
	}
}

var blankHeadingRe = regexp.MustCompile(`^(#{2,6})\s*$`)

func countBlankHeadings(content string) int {
	count := 0
	lines := strings.Split(content, "\n")
	for _, ln := range lines {
		if blankHeadingRe.MatchString(strings.TrimSpace(ln)) {
			count++
		}
	}
	return count
}
