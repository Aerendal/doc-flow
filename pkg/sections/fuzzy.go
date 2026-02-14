package sections

import (
	"regexp"
	"strings"
)

// NormalizeSectionName lowercases, trims, removes extra spaces and punctuation.
func NormalizeSectionName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	spaceRe := regexp.MustCompile(`\s+`)
	s = spaceRe.ReplaceAllString(s, " ")
	punctRe := regexp.MustCompile(`[[:punct:]]`)
	s = punctRe.ReplaceAllString(s, "")
	return s
}

// CanonicalFor returns canonical name if the provided name matches canonical or its aliases.
// exact=true when name already equals canonical (normalized).
func CanonicalFor(name string, aliases map[string][]string) (canonical string, ok bool, exact bool) {
	norm := NormalizeSectionName(name)
	for canon, list := range aliases {
		cn := NormalizeSectionName(canon)
		if norm == cn {
			return canon, true, true
		}
		for _, a := range list {
			if norm == NormalizeSectionName(a) {
				return canon, true, false
			}
		}
	}
	return "", false, false
}

type AliasHit struct {
	From string
	To   string
}

// FindAliasHits scans headings and returns alias→canonical replacements to suggest.
func FindAliasHits(content string, aliases map[string][]string) []AliasHit {
	var hits []AliasHit
	tree := Parse(content)
	var walk func(n *SectionNode)
	walk = func(n *SectionNode) {
		if n.Level > 0 {
			if to, ok, exact := CanonicalFor(n.Text, aliases); ok && !exact {
				hits = append(hits, AliasHit{From: n.Text, To: to})
			}
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(tree.Root)
	return hits
}

// RenameAliases rewrites heading lines using aliases map; returns new content and hits performed.
func RenameAliases(content string, aliases map[string][]string) (string, []AliasHit) {
	lines := strings.Split(content, "\n")
	var hits []AliasHit
	headingRe := regexp.MustCompile(`^(#{1,6})(\s+)(.+)$`)
	for i, ln := range lines {
		m := headingRe.FindStringSubmatch(ln)
		if len(m) == 0 {
			continue
		}
		text := m[3]
		if to, ok, exact := CanonicalFor(text, aliases); ok && !exact {
			hits = append(hits, AliasHit{From: text, To: to})
			lines[i] = m[1] + m[2] + to
		}
	}
	return strings.Join(lines, "\n"), hits
}
