package templatehint

import "strings"

type Hints struct {
	CodeBlocks int
	Tables     int
}

// Extract scans markdown content for code fences and simple tables.
func Extract(content string) Hints {
	lines := strings.Split(content, "\n")
	codeBlocks := 0
	tables := 0
	inFence := false
	for _, ln := range lines {
		trim := strings.TrimSpace(ln)
		if strings.HasPrefix(trim, "```") {
			if inFence {
				codeBlocks++
				inFence = false
			} else {
				inFence = true
			}
			continue
		}
		if strings.HasPrefix(trim, "|") && strings.Contains(trim, "|") {
			tables++
		}
	}
	return Hints{CodeBlocks: codeBlocks, Tables: tables}
}
