package governance

import "docflow/pkg/recommend"

type Rule struct {
	MinQuality     int
	RequiredFields []string
}

type Violation struct {
	Path   string
	Type   string
	Detail string
}

type Report struct {
	Violations []Violation
}

func Validate(templates []recommend.TemplateInfo, rule Rule) Report {
	var v []Violation
	for _, t := range templates {
		if t.Quality < rule.MinQuality {
			v = append(v, Violation{Path: t.Path, Type: "low_quality", Detail: "quality below threshold"})
		}
		// Required fields check (for demo we only check presence of doc_type/lang)
		if t.DocType == "" {
			v = append(v, Violation{Path: t.Path, Type: "missing_field", Detail: "doc_type"})
		}
		if t.Language == "" {
			v = append(v, Violation{Path: t.Path, Type: "missing_field", Detail: "language"})
		}
	}
	return Report{Violations: v}
}
