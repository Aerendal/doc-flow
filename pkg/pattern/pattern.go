package pattern

import (
        "crypto/sha256"
        "fmt"
        "os"
        "sort"
        "strings"

        "docflow/internal/util"
        "docflow/pkg/parser"
)

type SectionEntry struct {
        Level int    `json:"level" yaml:"level"`
        Text  string `json:"text" yaml:"text"`
}

type TemplatePattern struct {
        DocID    string         `json:"doc_id"`
        Path     string         `json:"path"`
        Sections []SectionEntry `json:"sections"`
        Fingerprint string     `json:"fingerprint"`
}

type PatternGroup struct {
        Fingerprint string            `json:"fingerprint"`
        Count       int               `json:"count"`
        Sections    []SectionEntry    `json:"sections"`
        Examples    []string          `json:"examples"`
}

type PatternReport struct {
        TotalTemplates int             `json:"total_templates"`
        UniquePatterns int             `json:"unique_patterns"`
        Groups         []*PatternGroup `json:"groups"`
}

func ExtractPattern(content string, docID, path string) *TemplatePattern {
        ms := parser.ParseHeadingsString(content)
        sections := make([]SectionEntry, 0, len(ms.Headings))
        for _, h := range ms.Headings {
                sections = append(sections, SectionEntry{Level: h.Level, Text: h.Text})
        }
        structSections := make([]SectionEntry, 0, len(sections))
        for _, s := range sections {
                if s.Level >= 2 {
                        structSections = append(structSections, s)
                }
        }
        fp := Fingerprint(structSections)
        return &TemplatePattern{
                DocID:       docID,
                Path:        path,
                Sections:    sections,
                Fingerprint: fp,
        }
}

func Fingerprint(sections []SectionEntry) string {
        var b strings.Builder
        for i, s := range sections {
                if i > 0 {
                        b.WriteByte('|')
                }
                fmt.Fprintf(&b, "%d:%s", s.Level, strings.ToLower(strings.TrimSpace(s.Text)))
        }
        hash := sha256.Sum256([]byte(b.String()))
        return fmt.Sprintf("%x", hash[:12])
}

func ScanPatterns(docsRoot string, ignorePatterns []string) ([]*TemplatePattern, error) {
        files, err := util.WalkMarkdown(docsRoot, ignorePatterns)
        if err != nil {
                return nil, fmt.Errorf("błąd skanowania szablonów: %w", err)
        }

        patterns := make([]*TemplatePattern, 0, len(files))
        for _, f := range files {
                data, err := os.ReadFile(f.Path)
                if err != nil {
                        continue
                }
                base := f.RelPath
                if idx := strings.LastIndex(base, "/"); idx >= 0 {
                        base = base[idx+1:]
                }
                docID := strings.TrimSuffix(base, ".md")
                tp := ExtractPattern(string(data), docID, f.RelPath)
                patterns = append(patterns, tp)
        }
        return patterns, nil
}

const DefaultSimilarityThreshold = 0.85

func GroupPatterns(patterns []*TemplatePattern) *PatternReport {
        return GroupPatternsWithMerge(patterns, 0)
}

func GroupPatternsWithMerge(patterns []*TemplatePattern, similarityThreshold float64) *PatternReport {
        groups := make(map[string]*PatternGroup)

        for _, p := range patterns {
                g, ok := groups[p.Fingerprint]
                if !ok {
                        g = &PatternGroup{
                                Fingerprint: p.Fingerprint,
                                Sections:    p.Sections,
                        }
                        groups[p.Fingerprint] = g
                }
                g.Count++
                if len(g.Examples) < 5 {
                        g.Examples = append(g.Examples, p.Path)
                }
        }

        sorted := make([]*PatternGroup, 0, len(groups))
        for _, g := range groups {
                sorted = append(sorted, g)
        }
        sort.Slice(sorted, func(i, j int) bool {
                return sorted[i].Count > sorted[j].Count
        })

        if similarityThreshold > 0 {
                sorted = MergeNearDuplicates(sorted, similarityThreshold)
        }

        return &PatternReport{
                TotalTemplates: len(patterns),
                UniquePatterns: len(sorted),
                Groups:         sorted,
        }
}

func h2PlusSections(sections []SectionEntry) []SectionEntry {
        result := make([]SectionEntry, 0, len(sections))
        for _, s := range sections {
                if s.Level >= 2 {
                        result = append(result, s)
                }
        }
        return result
}

func MergeNearDuplicates(groups []*PatternGroup, threshold float64) []*PatternGroup {
        if len(groups) <= 1 {
                return groups
        }

        merged := make([]bool, len(groups))
        result := make([]*PatternGroup, 0, len(groups))

        for i := 0; i < len(groups); i++ {
                if merged[i] {
                        continue
                }

                primary := groups[i]
                primaryH2 := h2PlusSections(primary.Sections)

                for j := i + 1; j < len(groups); j++ {
                        if merged[j] {
                                continue
                        }
                        candidateH2 := h2PlusSections(groups[j].Sections)
                        sim := Similarity(primaryH2, candidateH2)
                        if sim >= threshold {
                                primary.Count += groups[j].Count
                                for _, ex := range groups[j].Examples {
                                        if len(primary.Examples) < 5 {
                                                primary.Examples = append(primary.Examples, ex)
                                        }
                                }
                                merged[j] = true
                        }
                }

                result = append(result, primary)
        }

        sort.Slice(result, func(i, j int) bool {
                return result[i].Count > result[j].Count
        })

        return result
}

func TopN(report *PatternReport, n int) []*PatternGroup {
        if n > len(report.Groups) {
                n = len(report.Groups)
        }
        return report.Groups[:n]
}

func Levenshtein(a, b []SectionEntry) int {
        la, lb := len(a), len(b)
        if la == 0 {
                return lb
        }
        if lb == 0 {
                return la
        }

        prev := make([]int, lb+1)
        curr := make([]int, lb+1)
        for j := 0; j <= lb; j++ {
                prev[j] = j
        }

        for i := 1; i <= la; i++ {
                curr[0] = i
                for j := 1; j <= lb; j++ {
                        cost := 1
                        if a[i-1].Level == b[j-1].Level && strings.EqualFold(a[i-1].Text, b[j-1].Text) {
                                cost = 0
                        }
                        del := prev[j] + 1
                        ins := curr[j-1] + 1
                        sub := prev[j-1] + cost
                        curr[j] = min3(del, ins, sub)
                }
                prev, curr = curr, prev
        }
        return prev[lb]
}

func min3(a, b, c int) int {
        if a < b {
                if a < c {
                        return a
                }
                return c
        }
        if b < c {
                return b
        }
        return c
}

func Similarity(a, b []SectionEntry) float64 {
        maxLen := len(a)
        if len(b) > maxLen {
                maxLen = len(b)
        }
        if maxLen == 0 {
                return 1.0
        }
        dist := Levenshtein(a, b)
        return 1.0 - float64(dist)/float64(maxLen)
}

func NGrams(sections []SectionEntry, n int) []string {
        if n <= 0 || len(sections) < n {
                return nil
        }
        result := make([]string, 0, len(sections)-n+1)
        for i := 0; i <= len(sections)-n; i++ {
                parts := make([]string, n)
                for j := 0; j < n; j++ {
                        parts[j] = fmt.Sprintf("%d:%s", sections[i+j].Level, strings.ToLower(sections[i+j].Text))
                }
                result = append(result, strings.Join(parts, "|"))
        }
        return result
}

func CommonNGrams(patterns []*TemplatePattern, n int, minFreq int) map[string]int {
        freq := make(map[string]int)
        for _, p := range patterns {
                seen := make(map[string]bool)
                h2plus := h2PlusSections(p.Sections)
                for _, ng := range NGrams(h2plus, n) {
                        if !seen[ng] {
                                freq[ng]++
                                seen[ng] = true
                        }
                }
        }
        result := make(map[string]int)
        for ng, count := range freq {
                if count >= minFreq {
                        result[ng] = count
                }
        }
        return result
}
