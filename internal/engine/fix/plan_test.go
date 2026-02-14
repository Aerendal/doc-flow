package fix

import (
	"strings"
	"testing"

	"docflow/internal/engine"
)

func TestBuildPlanDeterministicAndSafeFixes(t *testing.T) {
	file := FileInput{
		Path: "docs/a.md",
		Content: `---
doc_id: a_doc
title: A
doc_type: guide
status: draft
context_sources: [seed]
---
# A
## Overview
`,
	}
	issuesA := []engine.Issue{
		{
			Code:    "DOCFLOW.VALIDATE.LEGACY_SECTION_NAME",
			Type:    "legacy_section_name",
			Path:    "docs/a.md",
			Message: "legacy section name 'Overview' → 'Przegląd'",
		},
		{
			Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
			Type:    "missing_field",
			Path:    "docs/a.md",
			Message: "brak version",
			Details: map[string]any{
				"field": "version",
			},
		},
	}
	issuesB := []engine.Issue{issuesA[1], issuesA[0]}
	opts := Options{
		Scopes: map[Scope]bool{
			ScopeFrontmatter: true,
			ScopeSections:    true,
		},
		Aliases: map[string][]string{
			"Przegląd": {"Overview"},
		},
		MaxFiles:   10,
		MaxChanges: 20,
	}

	planA, err := BuildPlan(".", issuesA, []FileInput{file}, opts)
	if err != nil {
		t.Fatal(err)
	}
	planB, err := BuildPlan(".", issuesB, []FileInput{file}, opts)
	if err != nil {
		t.Fatal(err)
	}

	if RenderDiff(planA) != RenderDiff(planB) {
		t.Fatalf("diff should be deterministic")
	}
	diff := RenderDiff(planA)
	if !strings.Contains(diff, "+version: v0.0.0") {
		t.Fatalf("expected version fix in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "+## Przegląd") {
		t.Fatalf("expected section rename in diff, got:\n%s", diff)
	}
	if planA.Stats.FilesTouched != 1 {
		t.Fatalf("expected one touched file, got %d", planA.Stats.FilesTouched)
	}
}

func TestBuildPlanRespectsOnlyCodes(t *testing.T) {
	file := FileInput{
		Path: "docs/a.md",
		Content: `---
doc_id: a_doc
title: A
doc_type: guide
status: draft
context_sources: [seed]
---
# A
## Overview
`,
	}
	issues := []engine.Issue{
		{
			Code:    "DOCFLOW.VALIDATE.LEGACY_SECTION_NAME",
			Type:    "legacy_section_name",
			Path:    "docs/a.md",
			Message: "legacy section name 'Overview' → 'Przegląd'",
		},
		{
			Code:    "DOCFLOW.VALIDATE.MISSING_FIELD",
			Type:    "missing_field",
			Path:    "docs/a.md",
			Message: "brak version",
			Details: map[string]any{
				"field": "version",
			},
		},
	}
	opts := Options{
		Scopes: map[Scope]bool{
			ScopeFrontmatter: true,
			ScopeSections:    true,
		},
		OnlyCodes: map[string]bool{
			"DOCFLOW.VALIDATE.MISSING_FIELD": true,
		},
		Aliases: map[string][]string{
			"Przegląd": {"Overview"},
		},
	}

	plan, err := BuildPlan(".", issues, []FileInput{file}, opts)
	if err != nil {
		t.Fatal(err)
	}
	diff := RenderDiff(plan)
	if !strings.Contains(diff, "+version: v0.0.0") {
		t.Fatalf("expected version fix in diff, got:\n%s", diff)
	}
	if strings.Contains(diff, "Przegląd") {
		t.Fatalf("section rename should be excluded by --only code filter")
	}
}
