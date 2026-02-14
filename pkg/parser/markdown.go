package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type Heading struct {
	Level int
	Text  string
	Line  int
}

type MarkdownStructure struct {
	Headings []Heading
	Lines    int
}

var headingRegex = regexp.MustCompile(`^(#{1,6})\s+(.+)`)

func ParseHeadings(r io.Reader) *MarkdownStructure {
	scanner := bufio.NewScanner(r)
	ms := &MarkdownStructure{}
	lineNum := 0
	inFrontmatter := false
	fmDelimCount := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			fmDelimCount++
			if fmDelimCount == 1 && lineNum <= 1 {
				inFrontmatter = true
				continue
			}
			if fmDelimCount == 2 && inFrontmatter {
				inFrontmatter = false
				continue
			}
		}

		if inFrontmatter {
			continue
		}

		if m := headingRegex.FindStringSubmatch(line); m != nil {
			ms.Headings = append(ms.Headings, Heading{
				Level: len(m[1]),
				Text:  strings.TrimSpace(m[2]),
				Line:  lineNum,
			})
		}
	}

	ms.Lines = lineNum
	return ms
}

func ParseHeadingsString(content string) *MarkdownStructure {
	return ParseHeadings(strings.NewReader(content))
}

func (ms *MarkdownStructure) HeadingTexts() []string {
	texts := make([]string, len(ms.Headings))
	for i, h := range ms.Headings {
		texts[i] = h.Text
	}
	return texts
}

func (ms *MarkdownStructure) MaxDepth() int {
	max := 0
	for _, h := range ms.Headings {
		if h.Level > max {
			max = h.Level
		}
	}
	return max
}

func (ms *MarkdownStructure) HeadingsAtLevel(level int) []Heading {
	var result []Heading
	for _, h := range ms.Headings {
		if h.Level == level {
			result = append(result, h)
		}
	}
	return result
}
