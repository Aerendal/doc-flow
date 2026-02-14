package validator

import "testing"

func TestMissingContextWarnAndPromote(t *testing.T) {
	dir := t.TempDir()
	content := `---
doc_id: ctx-1
title: T
doc_type: guide
version: "1.0.0"
status: draft
---
# H1
`
	writeFile(t, dir, "a.md", content)

	// warn by default
	report, err := ValidateDocs(dir, nil, nil)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	foundWarn := false
	for _, i := range report.Issues {
		if i.Type == IssueMissingContext && i.Level == LevelWarn {
			foundWarn = true
		}
	}
	if !foundWarn {
		t.Fatalf("missing context should be warning by default")
	}

	// promote to error for doc_type=guide
	opts := &Options{PromoteContextFor: map[string]bool{"guide": true}}
	report2, err := ValidateDocs(dir, nil, opts)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	foundErr := false
	for _, i := range report2.Issues {
		if i.Type == IssueMissingContext && i.Level == LevelError {
			foundErr = true
		}
	}
	if !foundErr {
		t.Fatalf("missing context should be error when promoted")
	}
}
