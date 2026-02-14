package compliance

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"docflow/pkg/governance"
	"docflow/pkg/index"
)

func TestReport(t *testing.T) {
	tmpDir := t.TempDir()
	file := tmpDir + "/doc.md"
	os.WriteFile(file, []byte("---\ntitle: T\ndoc_id: a\ndoc_type: guide\nstatus: draft\n---\n# T\n## Przegląd\ntext\n"), 0o644)

	idx := index.New()
	idx.Add(&index.DocumentRecord{DocID: "a", Path: file, Status: "draft", DocType: "guide"})
	rules := &governance.Rules{
		Statuses: map[string]governance.StatusRule{
			"draft": {RequiredFields: []string{"title"}},
		},
		Families: map[string]governance.FamilyRule{
			"guide": {RequiredSections: []string{"Przegląd"}},
		},
	}

	sum, err := Report(idx, rules)
	if err != nil {
		t.Fatal(err)
	}
	if sum.Passed != 1 || sum.Failed != 0 {
		t.Fatalf("unexpected pass/fail: %+v", sum)
	}
}

func TestReportSortsAndDetectsDuplicateDocIDs(t *testing.T) {
	tmpDir := t.TempDir()
	fileB := filepath.Join(tmpDir, "b.md")
	fileA := filepath.Join(tmpDir, "a.md")
	content := "---\ntitle: T\ndoc_id: dup\ndoc_type: guide\nstatus: draft\n---\n# T\n"
	if err := os.WriteFile(fileB, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileA, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	idx.Add(&index.DocumentRecord{DocID: "dup", Path: fileB, Status: "draft", DocType: "guide"})
	idx.Add(&index.DocumentRecord{DocID: "dup", Path: fileA, Status: "draft", DocType: "guide"})
	rules := &governance.Rules{
		Statuses: map[string]governance.StatusRule{
			"draft": {RequiredFields: []string{"z_field", "a_field"}},
		},
		Families: map[string]governance.FamilyRule{
			"guide": {RequiredSections: []string{"Z", "A"}},
		},
	}

	sum, err := Report(idx, rules)
	if err != nil {
		t.Fatal(err)
	}

	if len(sum.Docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(sum.Docs))
	}
	if sum.Docs[0].Path != fileA || sum.Docs[1].Path != fileB {
		t.Fatalf("docs are not sorted by path: [%s, %s]", sum.Docs[0].Path, sum.Docs[1].Path)
	}
	wantViolations := []string{"missing_a_field", "missing_section_A", "missing_section_Z", "missing_z_field"}
	if !reflect.DeepEqual(sum.Docs[0].Violations, wantViolations) {
		t.Fatalf("unexpected sorted violations: got=%v want=%v", sum.Docs[0].Violations, wantViolations)
	}
	dupPaths, ok := sum.DuplicateDocIDs["dup"]
	if !ok {
		t.Fatal("expected duplicate doc_id entry for dup")
	}
	wantDupPaths := []string{fileA, fileB}
	if !reflect.DeepEqual(dupPaths, wantDupPaths) {
		t.Fatalf("unexpected duplicate paths: got=%v want=%v", dupPaths, wantDupPaths)
	}
}

func TestReportWithFactsMatchesReport(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: T\ndoc_id: a\ndoc_type: guide\nstatus: draft\n---\n# T\n## Przegląd\ntext\n"
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := index.New()
	idx.Add(&index.DocumentRecord{DocID: "a", Path: file, Status: "draft", DocType: "guide"})
	rules := &governance.Rules{
		Statuses: map[string]governance.StatusRule{
			"draft": {RequiredFields: []string{"title"}},
		},
		Families: map[string]governance.FamilyRule{
			"guide": {RequiredSections: []string{"Przegląd"}},
		},
	}

	gotA, err := Report(idx, rules)
	if err != nil {
		t.Fatal(err)
	}
	gotB, err := ReportWithFacts(idx, rules, map[string]ContentFact{
		file: {Content: content},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotA, gotB) {
		t.Fatalf("report mismatch\nA=%#v\nB=%#v", gotA, gotB)
	}
}
