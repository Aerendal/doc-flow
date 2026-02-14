package schema

import (
	"os"
	"testing"
)

func TestLoadSectionSchemaValidates(t *testing.T) {
	yaml := `id: requirements
description: "Schemat wymagań"
sections:
  - name: Cel
    level: 2
    required: true
  - name: Zakres
    level: 2
    required: true
  - name: Wymagania
    level: 2
    required: true
  - name: Ryzyka
    level: 2
    required: false
`
	tmp := t.TempDir() + "/schema.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	s, err := LoadSectionSchema(tmp)
	if err != nil {
		t.Fatalf("LoadSectionSchema error: %v", err)
	}
	if s.ID != "requirements" {
		t.Fatalf("id = %s, want requirements", s.ID)
	}
	if len(s.Sections) != 4 {
		t.Fatalf("sections len = %d, want 4", len(s.Sections))
	}
}

func TestLoadSectionSchemaFails(t *testing.T) {
	yaml := `description: brak id
sections: []`
	tmp := t.TempDir() + "/schema.yaml"
	os.WriteFile(tmp, []byte(yaml), 0o644)
	if _, err := LoadSectionSchema(tmp); err == nil {
		t.Fatalf("expected error for missing id")
	}
}
