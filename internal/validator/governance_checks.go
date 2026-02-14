package validator

import (
	"fmt"
	"strings"

	"docflow/pkg/governance"
	"docflow/pkg/sections"
)

func checkGovernance(report *Report, rules *governance.Rules, file, docID, status, docType, content string) {
	// status rules
	if sr, ok := rules.Statuses[status]; ok {
		for _, f := range sr.RequiredFields {
			if !hasField(f, content) {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueGovernance,
					File:    file,
					DocID:   docID,
					Message: fmt.Sprintf("governance: missing required field '%s' for status %s", f, status),
				})
			}
		}
	}

	// family rules
	if fr, ok := rules.Families[docType]; ok {
		tree := sections.Parse(content)
		found := map[string]bool{}
		var walk func(n *sections.SectionNode)
		walk = func(n *sections.SectionNode) {
			if n.Level > 0 {
				found[strings.ToLower(strings.TrimSpace(n.Text))] = true
			}
			for _, c := range n.Children {
				walk(c)
			}
		}
		walk(tree.Root)
		for _, rs := range fr.RequiredSections {
			if !found[strings.ToLower(rs)] {
				report.add(Issue{
					Level:   LevelError,
					Type:    IssueGovernance,
					File:    file,
					DocID:   docID,
					Message: fmt.Sprintf("governance: missing required section '%s' for family %s", rs, docType),
				})
			}
		}
	}
}

func hasField(field, content string) bool {
	// crude check: look for "field:" in frontmatter/body; real parser already parsed,
	// but for governance we reuse content string to avoid passing structs.
	pat := fmt.Sprintf("%s:", field)
	return strings.Contains(content, pat)
}
