package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFrontmatterFull(t *testing.T) {
	content := `---
title: Test Document
doc_id: DOC-001
doc_type: specification
version: "1.0.0"
status: approved
priority: high
owner: jan.kowalski
created: "2026-01-15"
updated: "2026-02-01"
tags: [api, design]
depends_on: [DOC-002, DOC-003]
context_sources: [DOC-010]
language: pl
audience: developers
review_cycle: quarterly
---

# Test Document

Body content here.
`
	result, err := ParseFrontmatterString(content, "test.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasFrontmatter {
		t.Fatal("expected HasFrontmatter=true")
	}

	fm := result.Frontmatter
	if fm.Title != "Test Document" {
		t.Errorf("title = %q, want %q", fm.Title, "Test Document")
	}
	if fm.DocID != "DOC-001" {
		t.Errorf("doc_id = %q, want %q", fm.DocID, "DOC-001")
	}
	if fm.DocType != "specification" {
		t.Errorf("doc_type = %q, want %q", fm.DocType, "specification")
	}
	if fm.Status != "approved" {
		t.Errorf("status = %q, want %q", fm.Status, "approved")
	}
	if fm.Priority != "high" {
		t.Errorf("priority = %q, want %q", fm.Priority, "high")
	}
	if len(fm.DependsOn) != 2 {
		t.Errorf("depends_on len = %d, want 2", len(fm.DependsOn))
	}
	if len(fm.ContextSources) != 1 {
		t.Errorf("context_sources len = %d, want 1", len(fm.ContextSources))
	}
	if len(fm.Tags) != 2 {
		t.Errorf("tags len = %d, want 2", len(fm.Tags))
	}

	if !strings.Contains(result.Body, "Body content here.") {
		t.Error("body should contain content after frontmatter")
	}
}

func TestParseFrontmatterMinimal(t *testing.T) {
	content := `---
title: Minimal Doc
status: needs_content
---

# Minimal Doc
`
	result, err := ParseFrontmatterString(content, "minimal.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := result.Frontmatter
	if fm.Title != "Minimal Doc" {
		t.Errorf("title = %q, want %q", fm.Title, "Minimal Doc")
	}
	if fm.DocType != "unknown" {
		t.Errorf("doc_type should default to 'unknown', got %q", fm.DocType)
	}
	if fm.Version != "0.1.0" {
		t.Errorf("version should default to '0.1.0', got %q", fm.Version)
	}
	if fm.Priority != "normal" {
		t.Errorf("priority should default to 'normal', got %q", fm.Priority)
	}
	if fm.Language != "pl" {
		t.Errorf("language should default to 'pl', got %q", fm.Language)
	}

	if len(result.Warnings) < 1 {
		t.Error("expected warnings about missing fields")
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	content := `# Just a heading

Some body text without frontmatter.
`
	result, err := ParseFrontmatterString(content, "nofm.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HasFrontmatter {
		t.Error("expected HasFrontmatter=false")
	}
	if result.Frontmatter != nil {
		t.Error("expected nil frontmatter")
	}
	if !strings.Contains(result.Body, "Just a heading") {
		t.Error("body should contain the full content")
	}
}

func TestParseFrontmatterInvalidYAML(t *testing.T) {
	content := `---
title: Bad YAML
status: [invalid
  broken: yaml
---

# Bad
`
	_, err := ParseFrontmatterString(content, "bad.md")
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}

	perr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if perr.File != "bad.md" {
		t.Errorf("error file = %q, want %q", perr.File, "bad.md")
	}
	if !strings.Contains(perr.Message, "YAML") {
		t.Errorf("error message should mention YAML: %s", perr.Message)
	}
}

func TestParseFrontmatterUnclosed(t *testing.T) {
	content := `---
title: Unclosed
status: draft
`
	_, err := ParseFrontmatterString(content, "unclosed.md")
	if err == nil {
		t.Fatal("expected error for unclosed frontmatter")
	}

	perr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if !strings.Contains(perr.Message, "zamykającego") {
		t.Errorf("error should mention missing closing delimiter: %s", perr.Message)
	}
}

func TestParseFrontmatterExtraFields(t *testing.T) {
	content := `---
title: Extra Fields
status: draft
custom_field: custom_value
another_extra: 42
---

# Extra
`
	result, err := ParseFrontmatterString(content, "extra.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm := result.Frontmatter
	if fm.Extra == nil {
		t.Fatal("expected Extra map to be populated")
	}
	if fm.Extra["custom_field"] != "custom_value" {
		t.Errorf("extra[custom_field] = %v, want 'custom_value'", fm.Extra["custom_field"])
	}
	if fm.Extra["another_extra"] != 42 {
		t.Errorf("extra[another_extra] = %v, want 42", fm.Extra["another_extra"])
	}
}

func TestParseFrontmatterEmptyContent(t *testing.T) {
	result, err := ParseFrontmatterString("", "empty.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasFrontmatter {
		t.Error("expected HasFrontmatter=false for empty content")
	}
}

func TestParseTemplateFiles(t *testing.T) {
	templateDir := filepath.Join("..", "..", "testdata", "templates")
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		t.Skipf("cannot read templates directory: %v", err)
	}

	parsed := 0
	errors := 0
	withFM := 0

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if parsed >= 10 {
			break
		}

		fpath := filepath.Join(templateDir, e.Name())
		data, err := os.ReadFile(fpath)
		if err != nil {
			t.Errorf("cannot read %s: %v", e.Name(), err)
			errors++
			continue
		}

		result, err := ParseFrontmatterString(string(data), e.Name())
		if err != nil {
			t.Errorf("parse error in %s: %v", e.Name(), err)
			errors++
			continue
		}

		parsed++
		if result.HasFrontmatter {
			withFM++
			if result.Frontmatter.Title == "" {
				t.Errorf("%s: has frontmatter but no title", e.Name())
			}
		}
	}

	t.Logf("Parsed %d templates: %d with frontmatter, %d errors", parsed, withFM, errors)
	if parsed < 10 {
		t.Errorf("expected to parse at least 10 templates, got %d", parsed)
	}
}

func TestParseGroundTruthFiles(t *testing.T) {
	gtDir := filepath.Join("..", "..", "testdata", "ground_truth")
	entries, err := os.ReadDir(gtDir)
	if err != nil {
		t.Skipf("cannot read ground_truth directory: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		fpath := filepath.Join(gtDir, e.Name())
		data, err := os.ReadFile(fpath)
		if err != nil {
			t.Errorf("cannot read %s: %v", e.Name(), err)
			continue
		}

		result, err := ParseFrontmatterString(string(data), e.Name())
		if err != nil {
			t.Errorf("parse error in %s: %v", e.Name(), err)
			continue
		}

		if !result.HasFrontmatter {
			t.Errorf("%s: expected frontmatter in ground_truth file", e.Name())
			continue
		}

		fm := result.Frontmatter
		if fm.Title == "" {
			t.Errorf("%s: missing title", e.Name())
		}
		if fm.DocID == "" {
			t.Errorf("%s: missing doc_id", e.Name())
		}
	}
}
