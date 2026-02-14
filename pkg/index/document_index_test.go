package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewIndex(t *testing.T) {
	idx := New()
	if idx.Version != "1.0.0" {
		t.Errorf("version = %q, want %q", idx.Version, "1.0.0")
	}
	if idx.Count != 0 {
		t.Errorf("count = %d, want 0", idx.Count)
	}
}

func TestAddAndGet(t *testing.T) {
	idx := New()

	rec := &DocumentRecord{
		DocID:   "test-doc-001",
		Path:    "docs/test.md",
		Title:   "Test Document",
		DocType: "specification",
		Status:  "draft",
	}

	idx.Add(rec)

	if idx.Count != 1 {
		t.Errorf("count = %d, want 1", idx.Count)
	}

	byID := idx.GetByID("test-doc-001")
	if byID == nil {
		t.Fatal("GetByID returned nil")
	}
	if byID.Title != "Test Document" {
		t.Errorf("title = %q, want %q", byID.Title, "Test Document")
	}

	byPath := idx.GetByPath("docs/test.md")
	if byPath == nil {
		t.Fatal("GetByPath returned nil")
	}
	if byPath.DocID != "test-doc-001" {
		t.Errorf("doc_id = %q, want %q", byPath.DocID, "test-doc-001")
	}
}

func TestSaveAndLoadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "cache", "doc_index.json")

	idx := New()
	for i := 0; i < 50; i++ {
		rec := &DocumentRecord{
			DocID:   NormalizeDocID(filepath.Join("docs", "doc_"+string(rune('a'+i%26))+".md")),
			Path:    filepath.Join("docs", "doc_"+string(rune('a'+i%26))+".md"),
			Title:   "Document " + string(rune('A'+i%26)),
			DocType: "specification",
			Status:  "draft",
			Lines:   100 + i,
			Size:    int64(1000 + i*100),
			ModTime: time.Now().Format(time.RFC3339),
		}
		if i%5 == 0 {
			rec.DependsOn = []string{"base-doc"}
		}
		idx.Add(rec)
	}

	if err := idx.SaveJSON(indexPath); err != nil {
		t.Fatalf("SaveJSON error: %v", err)
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatal("index file not created")
	}

	loaded, err := LoadJSON(indexPath)
	if err != nil {
		t.Fatalf("LoadJSON error: %v", err)
	}

	if loaded.Count != idx.Count {
		t.Errorf("loaded count = %d, want %d", loaded.Count, idx.Count)
	}

	if loaded.Checksum != idx.Checksum {
		t.Errorf("checksum mismatch: loaded=%s, original=%s", loaded.Checksum, idx.Checksum)
	}

	if loaded.Version != "1.0.0" {
		t.Errorf("loaded version = %q, want %q", loaded.Version, "1.0.0")
	}

	for _, rec := range loaded.Documents {
		if rec.DocID == "" {
			t.Error("loaded document has empty doc_id")
		}
		found := loaded.GetByID(rec.DocID)
		if found == nil {
			t.Errorf("GetByID(%q) returned nil after load", rec.DocID)
		}
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "index.json")

	idx := New()
	rec := &DocumentRecord{
		DocID:          "roundtrip-001",
		Path:           "docs/roundtrip.md",
		Title:          "Roundtrip Test",
		DocType:        "guide",
		Status:         "approved",
		Version:        "1.2.0",
		Priority:       "high",
		Owner:          "jan.kowalski",
		Language:       "pl",
		DependsOn:      []string{"dep-001", "dep-002"},
		ContextSources: []string{"ctx-001"},
		Tags:           []string{"test", "roundtrip"},
		HeadingCount:   15,
		MaxDepth:       3,
		Lines:          200,
		Size:           5000,
		ModTime:        "2026-02-09T10:00:00Z",
	}
	idx.Add(rec)

	if err := idx.SaveJSON(indexPath); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := LoadJSON(indexPath)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	lr := loaded.GetByID("roundtrip-001")
	if lr == nil {
		t.Fatal("loaded record not found")
	}

	if lr.Title != rec.Title {
		t.Errorf("title = %q, want %q", lr.Title, rec.Title)
	}
	if lr.DocType != rec.DocType {
		t.Errorf("doc_type = %q, want %q", lr.DocType, rec.DocType)
	}
	if lr.Status != rec.Status {
		t.Errorf("status = %q, want %q", lr.Status, rec.Status)
	}
	if lr.Owner != rec.Owner {
		t.Errorf("owner = %q, want %q", lr.Owner, rec.Owner)
	}
	if len(lr.DependsOn) != 2 {
		t.Errorf("depends_on len = %d, want 2", len(lr.DependsOn))
	}
	if len(lr.ContextSources) != 1 {
		t.Errorf("context_sources len = %d, want 1", len(lr.ContextSources))
	}
	if len(lr.Tags) != 2 {
		t.Errorf("tags len = %d, want 2", len(lr.Tags))
	}
	if lr.HeadingCount != 15 {
		t.Errorf("heading_count = %d, want 15", lr.HeadingCount)
	}
	if lr.Lines != 200 {
		t.Errorf("lines = %d, want 200", lr.Lines)
	}
}

func TestNormalizeDocID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"docs/My Document.md", "my_document"},
		{"test/API Design Specification.md", "api_design_specification"},
		{"a_b_testing_guide.md", "a_b_testing_guide"},
		{"UPPERCASE.md", "uppercase"},
		{"path/to/file with spaces.md", "file_with_spaces"},
	}

	for _, tt := range tests {
		got := NormalizeDocID(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeDocID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"docs//test.md", "docs/test.md"},
		{"./docs/test.md", "docs/test.md"},
		{"docs/sub/../test.md", "docs/test.md"},
	}

	for _, tt := range tests {
		got := NormalizePath(tt.input)
		if got != tt.want {
			t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildIndexSmall(t *testing.T) {
	gtDir := filepath.Join("..", "..", "testdata", "ground_truth")
	if _, err := os.Stat(gtDir); os.IsNotExist(err) {
		t.Skipf("ground_truth directory not found: %s", gtDir)
	}

	idx, err := BuildIndex(gtDir, nil)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	if idx.Count == 0 {
		t.Fatal("expected non-zero document count")
	}

	for _, rec := range idx.All() {
		if rec.DocID == "" {
			t.Errorf("document %s has empty doc_id", rec.Path)
		}
		if rec.Title == "" {
			t.Errorf("document %s has empty title", rec.Path)
		}
		if rec.HeadingCount == 0 {
			t.Errorf("document %s has 0 headings", rec.Path)
		}
	}
}

func TestBuildIndex50Templates(t *testing.T) {
	templateDir := filepath.Join("..", "..", "testdata", "templates")
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Skipf("templates directory not found: %s", templateDir)
	}

	start := time.Now()
	idx, err := BuildIndex(templateDir, nil)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	if idx.Count < 50 {
		t.Errorf("expected at least 50 documents, got %d", idx.Count)
	}

	t.Logf("Indexed %d documents in %v", idx.Count, elapsed)

	perDoc := elapsed / time.Duration(idx.Count)
	t.Logf("Per-document indexing: %v", perDoc)
	if perDoc > 10*time.Millisecond {
		t.Errorf("per-document indexing too slow: %v (> 10ms per doc)", perDoc)
	}

	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "doc_index.json")

	if err := idx.SaveJSON(cachePath); err != nil {
		t.Fatalf("SaveJSON error: %v", err)
	}

	loadStart := time.Now()
	loaded, err := LoadJSON(cachePath)
	loadElapsed := time.Since(loadStart)

	if err != nil {
		t.Fatalf("LoadJSON error: %v", err)
	}

	t.Logf("Loaded %d documents from cache in %v", loaded.Count, loadElapsed)

	if loadElapsed > 100*time.Millisecond {
		t.Errorf("cache load took %v, expected < 100ms", loadElapsed)
	}

	if loaded.Count != idx.Count {
		t.Errorf("loaded count = %d, saved count = %d", loaded.Count, idx.Count)
	}

	if loaded.Checksum != idx.Checksum {
		t.Errorf("checksum mismatch after roundtrip")
	}
}

func TestDocumentsWithoutFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	doc1 := `# Document Without Frontmatter

## Metadane
- Właściciel: jan.kowalski
- Wersja: v1.0
- Status: approved

## Cel dokumentu
This is a test document without YAML frontmatter.
`
	doc2 := `# Another No-FM Doc

Just plain markdown content.
`

	os.WriteFile(filepath.Join(tmpDir, "no_fm_doc1.md"), []byte(doc1), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "no_fm_doc2.md"), []byte(doc2), 0o644)

	idx, err := BuildIndex(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildIndex error: %v", err)
	}

	if idx.Count != 2 {
		t.Fatalf("expected 2 documents, got %d", idx.Count)
	}

	rec1 := idx.GetByID("no_fm_doc1")
	if rec1 == nil {
		t.Fatal("no_fm_doc1 not found in index")
	}
	if rec1.Owner != "jan.kowalski" {
		t.Errorf("doc1 owner = %q, want %q (from meta block fallback)", rec1.Owner, "jan.kowalski")
	}
	if rec1.DocType != "unknown" {
		t.Errorf("doc1 doc_type = %q, want 'unknown'", rec1.DocType)
	}

	rec2 := idx.GetByID("no_fm_doc2")
	if rec2 == nil {
		t.Fatal("no_fm_doc2 not found in index")
	}
	if rec2.Status != "draft" {
		t.Errorf("doc2 status = %q, want 'draft' (default)", rec2.Status)
	}
}

func TestLoadCorruptedChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "corrupted.json")

	idx := New()
	idx.Add(&DocumentRecord{
		DocID: "test-001",
		Path:  "test.md",
		Title: "Test",
	})
	if err := idx.SaveJSON(indexPath); err != nil {
		t.Fatalf("save error: %v", err)
	}

	data, _ := os.ReadFile(indexPath)
	corrupted := strings.Replace(string(data), idx.Checksum, "deadbeefdeadbeef", 1)
	os.WriteFile(indexPath, []byte(corrupted), 0o644)

	_, err := LoadJSON(indexPath)
	if err == nil {
		t.Fatal("expected error for corrupted checksum")
	}
	if !strings.Contains(err.Error(), "checksum") {
		t.Errorf("error should mention checksum: %v", err)
	}
}

func TestIndexJSONFormat(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "test.json")

	idx := New()
	idx.Add(&DocumentRecord{
		DocID: "test-001",
		Path:  "test.md",
		Title: "Test",
	})

	if err := idx.SaveJSON(indexPath); err != nil {
		t.Fatalf("save error: %v", err)
	}

	data, _ := os.ReadFile(indexPath)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	requiredKeys := []string{"version", "created_at", "count", "checksum", "documents"}
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing key %q in index JSON", key)
		}
	}
}

func TestSaveJSONWithNoTimestamps(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "deterministic.json")

	idx := New()
	idx.CreatedAt = "2026-02-13T13:00:00Z"
	idx.Add(&DocumentRecord{
		DocID: "test-001",
		Path:  "docs/test.md",
		Title: "Test",
	})

	if err := idx.SaveJSONWithOptions(indexPath, SaveOptions{NoTimestamps: true}); err != nil {
		t.Fatalf("save error: %v", err)
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got := raw["created_at"]; got != "0001-01-01T00:00:00Z" {
		t.Fatalf("created_at = %v, want fixed deterministic value", got)
	}
}
