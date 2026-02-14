package validator

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	execore "docflow/internal/exec"
)

func TestValidateFactsMatchesValidateDocs(t *testing.T) {
	root := t.TempDir()
	good := `---
title: Good
doc_id: good_doc
doc_type: guide
status: draft
version: v1
context_sources: ["ctx"]
---
# T
## Przegląd
Body
`
	bad := `---
title: Bad
doc_id: bad_doc
doc_type: guide
status: draft
context_sources: ["ctx"]
---
# T
## Przegląd
Body
`
	if err := os.WriteFile(filepath.Join(root, "good.md"), []byte(good), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bad.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &Options{}
	legacy, err := ValidateDocs(root, nil, opts)
	if err != nil {
		t.Fatal(err)
	}

	facts, _, err := execore.CollectFacts(execore.RunnerOptions{
		Mode:     "validate",
		Root:     root,
		UseCache: false,
		CacheDir: filepath.Join(root, ".docflow", "cache"),
	})
	if err != nil {
		t.Fatal(err)
	}
	idx := execore.BuildIndexFromFacts(root, facts)
	fromFacts, err := ValidateFacts(facts, idx, opts)
	if err != nil {
		t.Fatal(err)
	}

	if legacy.Files != fromFacts.Files || legacy.Documents != fromFacts.Documents {
		t.Fatalf("header mismatch: legacy=(%d,%d) facts=(%d,%d)",
			legacy.Files, legacy.Documents, fromFacts.Files, fromFacts.Documents)
	}
	if !reflect.DeepEqual(legacy.SortedIssues(), fromFacts.SortedIssues()) {
		t.Fatalf("issues mismatch\nlegacy=%#v\nfacts=%#v", legacy.SortedIssues(), fromFacts.SortedIssues())
	}
}
