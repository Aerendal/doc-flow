package sections

import "docflow/pkg/parser"

// SectionNode represents a single heading node in the document.
type SectionNode struct {
	Level    int
	Text     string
	Line     int
	Children []*SectionNode
}

// SectionTree holds the root of a heading hierarchy.
type SectionTree struct {
	Root *SectionNode
}

// Parse builds a SectionTree from markdown content.
func Parse(content string) *SectionTree {
	ms := parser.ParseHeadingsString(content)

	root := &SectionNode{Level: 0, Text: "ROOT"}
	stack := []*SectionNode{root}

	for _, h := range ms.Headings {
		if h.Level == 1 {
			continue
		}
		node := &SectionNode{
			Level: h.Level,
			Text:  h.Text,
			Line:  h.Line,
		}

		// Find parent with level < current level.
		for len(stack) > 0 && stack[len(stack)-1].Level >= node.Level {
			stack = stack[:len(stack)-1]
		}
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, node)
		stack = append(stack, node)
	}

	return &SectionTree{Root: root}
}
