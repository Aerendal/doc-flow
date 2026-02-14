package parser

import (
        "os"
        "path/filepath"
        "strings"
        "testing"
)

func TestParseHeadingsBasic(t *testing.T) {
        content := `---
title: Test
status: draft
---

# Main Title

## Section One

Some text.

## Section Two

### Subsection

## Section Three
`
        ms := ParseHeadingsString(content)

        if len(ms.Headings) != 5 {
                t.Fatalf("expected 5 headings, got %d", len(ms.Headings))
        }

        expected := []struct {
                level int
                text  string
        }{
                {1, "Main Title"},
                {2, "Section One"},
                {2, "Section Two"},
                {3, "Subsection"},
                {2, "Section Three"},
        }

        for i, exp := range expected {
                if i >= len(ms.Headings) {
                        break
                }
                h := ms.Headings[i]
                if h.Level != exp.level {
                        t.Errorf("heading[%d] level = %d, want %d", i, h.Level, exp.level)
                }
                if h.Text != exp.text {
                        t.Errorf("heading[%d] text = %q, want %q", i, h.Text, exp.text)
                }
        }
}

func TestParseHeadingsSkipsFrontmatter(t *testing.T) {
        content := `---
title: Test
---

# Real Heading
`
        ms := ParseHeadingsString(content)

        if len(ms.Headings) != 1 {
                t.Fatalf("expected 1 heading (skip frontmatter), got %d", len(ms.Headings))
        }
        if ms.Headings[0].Text != "Real Heading" {
                t.Errorf("heading text = %q, want 'Real Heading'", ms.Headings[0].Text)
        }
}

func TestParseHeadingsNoFrontmatter(t *testing.T) {
        content := `# Title Without Frontmatter

## Section A

## Section B
`
        ms := ParseHeadingsString(content)

        if len(ms.Headings) != 3 {
                t.Fatalf("expected 3 headings, got %d", len(ms.Headings))
        }
}

func TestParseHeadingsEmpty(t *testing.T) {
        ms := ParseHeadingsString("")
        if len(ms.Headings) != 0 {
                t.Errorf("expected 0 headings for empty content, got %d", len(ms.Headings))
        }
        if ms.Lines != 0 {
                t.Errorf("expected 0 lines for empty content, got %d", ms.Lines)
        }
}

func TestHeadingTexts(t *testing.T) {
        content := `# A
## B
## C
`
        ms := ParseHeadingsString(content)
        texts := ms.HeadingTexts()

        if len(texts) != 3 {
                t.Fatalf("expected 3 texts, got %d", len(texts))
        }
        if texts[0] != "A" || texts[1] != "B" || texts[2] != "C" {
                t.Errorf("texts = %v, want [A B C]", texts)
        }
}

func TestMaxDepth(t *testing.T) {
        content := `# H1
## H2
### H3
`
        ms := ParseHeadingsString(content)
        if ms.MaxDepth() != 3 {
                t.Errorf("max depth = %d, want 3", ms.MaxDepth())
        }
}

func TestHeadingsAtLevel(t *testing.T) {
        content := `# Title
## A
## B
### C
## D
`
        ms := ParseHeadingsString(content)
        h2 := ms.HeadingsAtLevel(2)
        if len(h2) != 3 {
                t.Errorf("h2 count = %d, want 3", len(h2))
        }
}

func TestParseHeadingsLineNumbers(t *testing.T) {
        content := `---
title: Test
---

# Title

## Section
`
        ms := ParseHeadingsString(content)
        if len(ms.Headings) < 2 {
                t.Fatalf("expected at least 2 headings, got %d", len(ms.Headings))
        }

        if ms.Headings[0].Line != 5 {
                t.Errorf("first heading line = %d, want 5", ms.Headings[0].Line)
        }
}

func TestParseHeadingsTemplateFile(t *testing.T) {
        templateDir := filepath.Join("..", "..", "testdata", "templates")
        entries, err := os.ReadDir(templateDir)
        if err != nil {
                t.Skipf("cannot read templates directory: %v", err)
        }

        checked := 0
        for _, e := range entries {
                if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
                        continue
                }
                if checked >= 3 {
                        break
                }

                fpath := filepath.Join(templateDir, e.Name())
                data, err := os.ReadFile(fpath)
                if err != nil {
                        t.Errorf("cannot read %s: %v", e.Name(), err)
                        continue
                }

                ms := ParseHeadingsString(string(data))
                if len(ms.Headings) == 0 {
                        t.Errorf("%s: expected headings, got none", e.Name())
                        continue
                }

                if ms.Headings[0].Level != 1 {
                        t.Errorf("%s: first heading level = %d, want 1", e.Name(), ms.Headings[0].Level)
                }

                h2Count := len(ms.HeadingsAtLevel(2))
                t.Logf("%s: %d headings total, %d at h2, max depth %d", e.Name(), len(ms.Headings), h2Count, ms.MaxDepth())
                checked++
        }

        if checked < 3 {
                t.Errorf("expected to check at least 3 template files, got %d", checked)
        }
}
