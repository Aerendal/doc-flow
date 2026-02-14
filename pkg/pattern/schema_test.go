package pattern

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGenerateSchema(t *testing.T) {
	group := &PatternGroup{
		Fingerprint: "abcdef1234567890abcdef12",
		Count:       50,
		Sections:    []SectionEntry{{1, "Title"}, {2, "Cel dokumentu"}, {2, "Metadane"}},
		Examples:    []string{"doc1.md", "doc2.md"},
	}

	schema := GenerateSchema(group)
	if schema.PatternID != "pattern_abcdef12" {
		t.Errorf("PatternID = %q, want 'pattern_abcdef12'", schema.PatternID)
	}
	if schema.Frequency != 50 {
		t.Errorf("Frequency = %d, want 50", schema.Frequency)
	}
	if len(schema.Sections) != 3 {
		t.Fatalf("Sections = %d, want 3", len(schema.Sections))
	}
	if schema.Sections[0].Name != "Title" || schema.Sections[0].Level != 1 {
		t.Errorf("Section[0] = %v, want {Title, 1}", schema.Sections[0])
	}
	if !schema.Sections[0].Required {
		t.Error("all sections should be required")
	}
	if schema.Sections[1].Order != 2 {
		t.Errorf("Section[1].Order = %d, want 2", schema.Sections[1].Order)
	}
}

func TestSaveSchemaYAML(t *testing.T) {
	tmpDir := t.TempDir()
	schema := &SectionSchema{
		PatternID:   "pattern_test",
		Fingerprint: "abc123",
		Frequency:   10,
		ExampleDocs: []string{"doc.md"},
		Sections: []SectionSchemaEntry{
			{Name: "Title", Level: 1, Required: true, Order: 1},
			{Name: "Cel", Level: 2, Required: true, Order: 2},
		},
	}

	if err := SaveSchemaYAML(schema, tmpDir); err != nil {
		t.Fatalf("SaveSchemaYAML error: %v", err)
	}

	path := filepath.Join(tmpDir, "pattern_test.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read schema file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "pattern_test") {
		t.Error("schema file should contain pattern_id")
	}
	if !strings.Contains(content, "Cel") {
		t.Error("schema file should contain section name")
	}

	var loaded SectionSchema
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("cannot parse YAML: %v", err)
	}
	if loaded.PatternID != "pattern_test" {
		t.Errorf("loaded PatternID = %q", loaded.PatternID)
	}
	if len(loaded.Sections) != 2 {
		t.Errorf("loaded sections = %d, want 2", len(loaded.Sections))
	}
}

func TestGenerateSchemas(t *testing.T) {
	report := &PatternReport{
		TotalTemplates: 100,
		UniquePatterns: 3,
		Groups: []*PatternGroup{
			{Fingerprint: "aaa1234567890000aaa12300", Count: 80, Sections: []SectionEntry{{1, "A"}}, Examples: []string{"a.md"}},
			{Fingerprint: "bbb1234567890000bbb12300", Count: 15, Sections: []SectionEntry{{1, "B"}}, Examples: []string{"b.md"}},
			{Fingerprint: "ccc1234567890000ccc12300", Count: 5, Sections: []SectionEntry{{1, "C"}}, Examples: []string{"c.md"}},
		},
	}

	schemas := GenerateSchemas(report, 2)
	if len(schemas) != 2 {
		t.Fatalf("GenerateSchemas(2) = %d schemas, want 2", len(schemas))
	}
	if schemas[0].Frequency != 80 {
		t.Errorf("schema[0].Frequency = %d, want 80", schemas[0].Frequency)
	}
}

func TestSaveAllSchemas(t *testing.T) {
	tmpDir := t.TempDir()
	schemas := []*SectionSchema{
		{PatternID: "p1", Sections: []SectionSchemaEntry{{Name: "A", Level: 1, Required: true, Order: 1}}},
		{PatternID: "p2", Sections: []SectionSchemaEntry{{Name: "B", Level: 1, Required: true, Order: 1}}},
	}

	if err := SaveAllSchemas(schemas, tmpDir); err != nil {
		t.Fatalf("SaveAllSchemas error: %v", err)
	}

	for _, s := range schemas {
		path := filepath.Join(tmpDir, s.PatternID+".yaml")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("schema file %s not found", path)
		}
	}
}
