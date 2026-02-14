package pattern

import (
        "os"
        "path/filepath"
        "strings"
        "testing"
)

func TestExtractPattern(t *testing.T) {
        content := `---
title: Test
status: draft
---

# Main Title

## Cel dokumentu

## Metadane

## Zakres
`
        p := ExtractPattern(content, "test", "test.md")
        if p.DocID != "test" {
                t.Errorf("DocID = %q, want 'test'", p.DocID)
        }
        if len(p.Sections) != 4 {
                t.Fatalf("Sections = %d, want 4", len(p.Sections))
        }
        if p.Sections[0].Level != 1 || p.Sections[0].Text != "Main Title" {
                t.Errorf("Section[0] = %v, want {1, 'Main Title'}", p.Sections[0])
        }
        if p.Fingerprint == "" {
                t.Error("Fingerprint should not be empty")
        }
}

func TestFingerprint(t *testing.T) {
        sections1 := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}}
        sections2 := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}}
        sections3 := []SectionEntry{{1, "A"}, {2, "B"}, {2, "D"}}

        fp1 := Fingerprint(sections1)
        fp2 := Fingerprint(sections2)
        fp3 := Fingerprint(sections3)

        if fp1 != fp2 {
                t.Errorf("identical sections should have same fingerprint: %s vs %s", fp1, fp2)
        }
        if fp1 == fp3 {
                t.Error("different sections should have different fingerprints")
        }
}

func TestFingerprintCaseInsensitive(t *testing.T) {
        s1 := []SectionEntry{{2, "Cel Dokumentu"}}
        s2 := []SectionEntry{{2, "cel dokumentu"}}
        if Fingerprint(s1) != Fingerprint(s2) {
                t.Error("fingerprint should be case-insensitive")
        }
}

func TestGroupPatterns(t *testing.T) {
        patterns := []*TemplatePattern{
                {DocID: "a", Fingerprint: "fp1", Sections: []SectionEntry{{1, "A"}}},
                {DocID: "b", Fingerprint: "fp1", Sections: []SectionEntry{{1, "A"}}},
                {DocID: "c", Fingerprint: "fp2", Sections: []SectionEntry{{1, "B"}}},
                {DocID: "d", Fingerprint: "fp1", Sections: []SectionEntry{{1, "A"}}},
        }

        report := GroupPatterns(patterns)
        if report.TotalTemplates != 4 {
                t.Errorf("TotalTemplates = %d, want 4", report.TotalTemplates)
        }
        if report.UniquePatterns != 2 {
                t.Errorf("UniquePatterns = %d, want 2", report.UniquePatterns)
        }
        if report.Groups[0].Count != 3 {
                t.Errorf("top group count = %d, want 3", report.Groups[0].Count)
        }
}

func TestTopN(t *testing.T) {
        report := &PatternReport{
                Groups: []*PatternGroup{
                        {Count: 100}, {Count: 50}, {Count: 20}, {Count: 5},
                },
        }
        top := TopN(report, 2)
        if len(top) != 2 {
                t.Fatalf("TopN(2) returned %d groups", len(top))
        }
        if top[0].Count != 100 || top[1].Count != 50 {
                t.Errorf("TopN order wrong: %d, %d", top[0].Count, top[1].Count)
        }

        all := TopN(report, 100)
        if len(all) != 4 {
                t.Errorf("TopN(100) should return all %d groups", len(all))
        }
}

func TestLevenshtein(t *testing.T) {
        a := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}}
        b := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}}
        if d := Levenshtein(a, b); d != 0 {
                t.Errorf("identical sequences: distance = %d, want 0", d)
        }

        c := []SectionEntry{{1, "A"}, {2, "X"}, {2, "C"}}
        if d := Levenshtein(a, c); d != 1 {
                t.Errorf("one substitution: distance = %d, want 1", d)
        }

        d := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}, {2, "D"}}
        if dist := Levenshtein(a, d); dist != 1 {
                t.Errorf("one insertion: distance = %d, want 1", dist)
        }

        if dist := Levenshtein(nil, a); dist != 3 {
                t.Errorf("empty vs 3: distance = %d, want 3", dist)
        }
}

func TestSimilarity(t *testing.T) {
        a := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}}
        b := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}}
        if s := Similarity(a, b); s != 1.0 {
                t.Errorf("identical: similarity = %.2f, want 1.0", s)
        }

        c := []SectionEntry{{1, "A"}, {2, "X"}, {2, "C"}}
        sim := Similarity(a, c)
        if sim < 0.6 || sim > 0.7 {
                t.Errorf("one diff in 3: similarity = %.2f, want ~0.67", sim)
        }
}

func TestNGrams(t *testing.T) {
        s := []SectionEntry{{1, "A"}, {2, "B"}, {2, "C"}, {2, "D"}}
        ng := NGrams(s, 2)
        if len(ng) != 3 {
                t.Fatalf("2-grams of 4 entries = %d, want 3", len(ng))
        }
        if !strings.Contains(ng[0], "a") || !strings.Contains(ng[0], "b") {
                t.Errorf("first 2-gram should contain A and B: %s", ng[0])
        }

        empty := NGrams(s, 0)
        if empty != nil {
                t.Error("n=0 should return nil")
        }

        tooLong := NGrams(s, 10)
        if tooLong != nil {
                t.Error("n > len should return nil")
        }
}

func TestCommonNGrams(t *testing.T) {
        p1 := &TemplatePattern{Sections: []SectionEntry{{1, "Title"}, {2, "B"}, {2, "C"}}}
        p2 := &TemplatePattern{Sections: []SectionEntry{{1, "Title"}, {2, "B"}, {2, "D"}}}
        p3 := &TemplatePattern{Sections: []SectionEntry{{1, "Title"}, {2, "B"}, {2, "C"}}}

        common := CommonNGrams([]*TemplatePattern{p1, p2, p3}, 2, 2)
        if len(common) == 0 {
                t.Fatal("should have common 2-grams with freq >= 2 (h2+ sections)")
        }
        bcKey := "2:b|2:c"
        if common[bcKey] != 2 {
                t.Errorf("B|C n-gram frequency = %d, want 2", common[bcKey])
        }
}

func TestMergeNearDuplicates(t *testing.T) {
        groups := []*PatternGroup{
                {Fingerprint: "a", Count: 100, Sections: []SectionEntry{{1, "T1"}, {2, "Cel"}, {2, "Metadane"}, {2, "Zakres"}}, Examples: []string{"a.md"}},
                {Fingerprint: "b", Count: 10, Sections: []SectionEntry{{1, "T2"}, {2, "Cel"}, {2, "Metadane"}, {2, "Zakres"}}, Examples: []string{"b.md"}},
                {Fingerprint: "c", Count: 5, Sections: []SectionEntry{{1, "T3"}, {2, "X"}, {2, "Y"}, {2, "Z"}}, Examples: []string{"c.md"}},
        }

        merged := MergeNearDuplicates(groups, 0.85)
        if len(merged) != 2 {
                t.Fatalf("expected 2 groups after merge (a+b similar h2, c different), got %d", len(merged))
        }
        if merged[0].Count != 110 {
                t.Errorf("merged group count = %d, want 110", merged[0].Count)
        }
        if merged[1].Count != 5 {
                t.Errorf("unmerged group count = %d, want 5", merged[1].Count)
        }
}

func TestGroupPatternsWithMerge(t *testing.T) {
        p1 := &TemplatePattern{DocID: "a", Fingerprint: Fingerprint([]SectionEntry{{2, "Cel"}, {2, "Meta"}}), Sections: []SectionEntry{{1, "Title A"}, {2, "Cel"}, {2, "Meta"}}}
        p2 := &TemplatePattern{DocID: "b", Fingerprint: Fingerprint([]SectionEntry{{2, "Cel"}, {2, "Meta"}}), Sections: []SectionEntry{{1, "Title B"}, {2, "Cel"}, {2, "Meta"}}}
        p3 := &TemplatePattern{DocID: "c", Fingerprint: Fingerprint([]SectionEntry{{2, "X"}, {2, "Y"}}), Sections: []SectionEntry{{1, "Title C"}, {2, "X"}, {2, "Y"}}}

        report := GroupPatternsWithMerge([]*TemplatePattern{p1, p2, p3}, DefaultSimilarityThreshold)
        if report.UniquePatterns != 2 {
                t.Errorf("expected 2 unique patterns after merge, got %d", report.UniquePatterns)
        }
}

func TestScanPatternsSmall(t *testing.T) {
        tmpDir := t.TempDir()
        for _, name := range []string{"doc1.md", "doc2.md"} {
                content := "---\ntitle: Test\n---\n\n# Title\n\n## Sekcja A\n\n## Sekcja B\n"
                os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o644)
        }

        patterns, err := ScanPatterns(tmpDir, nil)
        if err != nil {
                t.Fatalf("ScanPatterns error: %v", err)
        }
        if len(patterns) != 2 {
                t.Fatalf("expected 2 patterns, got %d", len(patterns))
        }
        if patterns[0].Fingerprint != patterns[1].Fingerprint {
                t.Error("identical docs should have same fingerprint")
        }
}

func TestScanPatternsTemplates(t *testing.T) {
        templatesDir := "../../testdata/templates"
        if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
                t.Skip("testdata/templates not found")
        }

        patterns, err := ScanPatterns(templatesDir, nil)
        if err != nil {
                t.Fatalf("ScanPatterns error: %v", err)
        }
        if len(patterns) < 50 {
                t.Fatalf("expected 50+ patterns, got %d", len(patterns))
        }

        report := GroupPatterns(patterns)
        if report.Groups[0].Count < 100 {
                t.Errorf("dominant pattern should cover 100+ templates (88%% identical structure), got %d", report.Groups[0].Count)
        }

        pct := float64(report.Groups[0].Count) / float64(report.TotalTemplates) * 100
        t.Logf("Templates: %d, Unique patterns: %d, Top pattern: %d (%.1f%%)",
                report.TotalTemplates, report.UniquePatterns, report.Groups[0].Count, pct)
}

func TestFormatReport(t *testing.T) {
        report := &PatternReport{
                TotalTemplates: 100,
                UniquePatterns: 3,
                Groups: []*PatternGroup{
                        {
                                Fingerprint: "abcdef1234567890abcdef12",
                                Count:       80,
                                Sections:    []SectionEntry{{1, "Title"}, {2, "Sekcja A"}},
                                Examples:    []string{"doc1.md", "doc2.md"},
                        },
                },
        }
        out := FormatReport(report, 5)
        if !strings.Contains(out, "80") {
                t.Error("report should contain count 80")
        }
        if !strings.Contains(out, "doc1.md") {
                t.Error("report should contain example file")
        }
}
