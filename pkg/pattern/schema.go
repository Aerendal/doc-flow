package pattern

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type SectionSchemaEntry struct {
	Name     string `yaml:"name"`
	Level    int    `yaml:"level"`
	Required bool   `yaml:"required"`
	Order    int    `yaml:"order"`
}

type SectionSchema struct {
	PatternID   string               `yaml:"pattern_id"`
	Fingerprint string               `yaml:"fingerprint"`
	Frequency   int                  `yaml:"frequency"`
	ExampleDocs []string             `yaml:"example_docs"`
	Sections    []SectionSchemaEntry `yaml:"sections"`
}

func GenerateSchema(group *PatternGroup) *SectionSchema {
	entries := make([]SectionSchemaEntry, len(group.Sections))
	for i, s := range group.Sections {
		entries[i] = SectionSchemaEntry{
			Name:     s.Text,
			Level:    s.Level,
			Required: true,
			Order:    i + 1,
		}
	}

	patternID := fmt.Sprintf("pattern_%s", group.Fingerprint[:8])

	return &SectionSchema{
		PatternID:   patternID,
		Fingerprint: group.Fingerprint,
		Frequency:   group.Count,
		ExampleDocs: group.Examples,
		Sections:    entries,
	}
}

func GenerateSchemas(report *PatternReport, topN int) []*SectionSchema {
	groups := TopN(report, topN)
	schemas := make([]*SectionSchema, len(groups))
	for i, g := range groups {
		schemas[i] = GenerateSchema(g)
	}
	return schemas
}

func SaveSchemaYAML(schema *SectionSchema, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("nie można utworzyć katalogu %s: %w", dir, err)
	}

	data, err := yaml.Marshal(schema)
	if err != nil {
		return fmt.Errorf("błąd serializacji schematu %s: %w", schema.PatternID, err)
	}

	path := filepath.Join(dir, schema.PatternID+".yaml")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("błąd zapisu schematu do %s: %w", path, err)
	}
	return nil
}

func SaveAllSchemas(schemas []*SectionSchema, dir string) error {
	for _, s := range schemas {
		if err := SaveSchemaYAML(s, dir); err != nil {
			return err
		}
	}
	return nil
}

func FormatReport(report *PatternReport, topN int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Wzorce sekcji — raport\n")
	fmt.Fprintf(&b, "======================\n\n")
	fmt.Fprintf(&b, "Szablonów: %d\n", report.TotalTemplates)
	fmt.Fprintf(&b, "Unikalnych wzorców: %d\n\n", report.UniquePatterns)

	groups := TopN(report, topN)
	for i, g := range groups {
		pct := float64(g.Count) / float64(report.TotalTemplates) * 100
		fmt.Fprintf(&b, "--- Wzorzec #%d (fingerprint: %s) ---\n", i+1, g.Fingerprint[:8])
		fmt.Fprintf(&b, "Częstotliwość: %d (%.1f%%)\n", g.Count, pct)
		fmt.Fprintf(&b, "Sekcji: %d\n", len(g.Sections))

		maxShow := 10
		if len(g.Sections) < maxShow {
			maxShow = len(g.Sections)
		}
		for j := 0; j < maxShow; j++ {
			s := g.Sections[j]
			prefix := strings.Repeat("#", s.Level)
			fmt.Fprintf(&b, "  %s %s\n", prefix, s.Text)
		}
		if len(g.Sections) > 10 {
			fmt.Fprintf(&b, "  ... (+%d sekcji)\n", len(g.Sections)-10)
		}

		fmt.Fprintf(&b, "Przykłady: %s\n\n", strings.Join(g.Examples, ", "))
	}
	return b.String()
}
