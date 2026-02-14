package generate

import "testing"

func TestLoadManifest(t *testing.T) {
	yaml := `documents:
  - doc_id: a
    doc_type: guide
    template: tmpl.md
    output: out.md
`
	m, err := loadManifestBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Documents) != 1 {
		t.Fatalf("docs len = %d, want 1", len(m.Documents))
	}
	if m.Documents[0].DocID != "a" {
		t.Fatalf("doc_id = %s, want a", m.Documents[0].DocID)
	}
}

func TestLoadManifestFails(t *testing.T) {
	_, err := loadManifestBytes([]byte(`documents: []`))
	if err == nil {
		t.Fatalf("expected error for empty manifest")
	}
}
