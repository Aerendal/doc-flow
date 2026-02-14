package content

import (
	"regexp"
)

type Metrics struct {
	CodeBlocks int
	Tables     int
	Links      int
	Images     int
	Words      int
}

var (
	reCodeFence = regexp.MustCompile("(?m)^```")
	reTable     = regexp.MustCompile(`(?m)^\|.*\|`)
	reLink      = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
	reImage     = regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)
	reWord      = regexp.MustCompile(`\pL+`)
)

// Extract returns simple content metrics from markdown.
func Extract(md string) Metrics {
	code := len(reCodeFence.FindAllString(md, -1)) / 2
	tables := len(reTable.FindAllString(md, -1))
	links := len(reLink.FindAllString(md, -1))
	images := len(reImage.FindAllString(md, -1))
	words := len(reWord.FindAllString(md, -1))

	return Metrics{
		CodeBlocks: code,
		Tables:     tables,
		Links:      links,
		Images:     images,
		Words:      words,
	}
}
