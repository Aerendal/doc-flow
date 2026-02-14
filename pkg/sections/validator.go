package sections

import (
	"fmt"
	"strings"

	"docflow/pkg/schema"
)

// ValidationResult contains info about missing sections.
type ValidationResult struct {
	Missing []string
}

// ValidateTree compares a SectionTree with a SectionSchema and returns missing required sections.
func ValidateTree(tree *SectionTree, sch *schema.SectionSchema) *ValidationResult {
	return ValidateTreeWithAliases(tree, sch, nil)
}

// ValidateTreeWithAliases allows matching required section names via alias map (normalized).
func ValidateTreeWithAliases(tree *SectionTree, sch *schema.SectionSchema, aliases map[string]string) *ValidationResult {
	found := make(map[string]bool)

	var walk func(node *SectionNode)
	walk = func(node *SectionNode) {
		if node.Level > 0 { // skip root
			key := norm(node.Text)
			found[key] = true
		}
		for _, c := range node.Children {
			walk(c)
		}
	}
	walk(tree.Root)

	var missing []string
	for _, rule := range sch.Sections {
		if !rule.Required {
			continue
		}
		if matched(found, rule.Name, aliases) {
			continue
		}
		missing = append(missing, rule.Name)
	}

	return &ValidationResult{Missing: missing}
}

// ValidateContent parses content and validates against schema.
func ValidateContent(content string, sch *schema.SectionSchema) (*ValidationResult, error) {
	return ValidateContentWithAliases(content, sch, nil)
}

func ValidateContentWithAliases(content string, sch *schema.SectionSchema, aliases map[string]string) (*ValidationResult, error) {
	tree := Parse(content)
	if tree == nil || tree.Root == nil {
		return nil, fmt.Errorf("nie udało się zbudować SectionTree")
	}
	return ValidateTreeWithAliases(tree, sch, aliases), nil
}

func norm(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func matched(found map[string]bool, required string, aliases map[string]string) bool {
	req := norm(required)
	if found[req] {
		return true
	}
	if aliases == nil {
		return false
	}
	// legacy map[string]string kept for backward compatibility: k->canonical
	if alt, ok := aliases[req]; ok {
		if found[norm(alt)] {
			return true
		}
	}
	for k, v := range aliases {
		if v == required && found[norm(k)] {
			return true
		}
	}
	return false
}
