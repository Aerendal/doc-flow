package validator

import "testing"

func TestPublishedFailsOnEmptySections(t *testing.T) {
	dir := t.TempDir()
	content := `---
doc_id: pub-1
title: T
doc_type: guide
version: "1.0.0"
status: published
---
# T
## 
`
	writeFile(t, dir, "doc.md", content)

	opts := &Options{
		PromoteContextFor: map[string]bool{},
		RequiredDeps:      map[string][]string{},
		PublishStrict:     true,
	}
	report, err := ValidateDocs(dir, nil, opts)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	found := false
	for _, i := range report.Issues {
		if i.Level == LevelError && i.Type == IssueMissingField {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected error for empty sections in published document")
	}
}

func TestPublishedPassesWhenNotStrict(t *testing.T) {
	dir := t.TempDir()
	content := `---
doc_id: pub-2
title: T
doc_type: guide
version: "1.0.0"
status: published
---
# T
## 
`
	writeFile(t, dir, "doc.md", content)

	opts := &Options{
		PromoteContextFor: map[string]bool{},
		RequiredDeps:      map[string][]string{},
		PublishStrict:     false,
	}
	report, err := ValidateDocs(dir, nil, opts)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	for _, i := range report.Issues {
		if i.Level == LevelError && i.Type == IssueMissingField {
			t.Fatalf("expected no error for published when PublishStrict=false")
		}
	}
}
