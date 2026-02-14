package validator

import "testing"

func TestMissingExpectedDeps(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", `---
doc_id: spec-1
title: Spec
doc_type: specification
version: "1.0.0"
status: draft
depends_on: [product_vision_statement]
---
# Title
`)

	opts := &Options{
		PromoteContextFor: map[string]bool{},
		RequiredDeps: map[string][]string{
			"specification": {"product_vision_statement", "system_architecture"},
		},
	}

	report, err := ValidateDocs(dir, nil, opts)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	found := false
	for _, i := range report.Issues {
		if i.Type == IssueMissingExpectedDep && i.DocID == "spec-1" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected missing expected dep issue")
	}
}
