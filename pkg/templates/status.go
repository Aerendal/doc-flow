package templates

import "docflow/pkg/recommend"

// DeprecationReport lists templates that should be deprecated based on rules.
type DeprecationReport struct {
	ToDeprecate []recommend.TemplateInfo
}

type Rules struct {
	MinQuality int
	MaxVersion string
	UnusedOnly bool
}

// EvaluateDeprecation returns templates with quality below threshold or unused (if UnusedOnly).
func EvaluateDeprecation(templates []recommend.TemplateInfo, rules Rules) DeprecationReport {
	var list []recommend.TemplateInfo
	for _, t := range templates {
		if t.Status == "archived" || t.Status == "deprecated" {
			continue
		}
		if rules.UnusedOnly && t.Usage > 0 {
			continue
		}
		if t.Quality < rules.MinQuality {
			list = append(list, t)
			continue
		}
	}
	return DeprecationReport{ToDeprecate: list}
}
