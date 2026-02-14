package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SectionRule describes a single expected section in a document template.
type SectionRule struct {
	Name     string `yaml:"name"`
	Level    int    `yaml:"level"`
	Required bool   `yaml:"required"`
	Order    int    `yaml:"order,omitempty"`
}

// SectionSchema represents an ordered list of sections for a document family.
type SectionSchema struct {
	ID          string        `yaml:"id"`
	Description string        `yaml:"description,omitempty"`
	Sections    []SectionRule `yaml:"sections"`
}

// LoadSectionSchema reads YAML schema from a file path.
func LoadSectionSchema(path string) (*SectionSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nie można odczytać schematu %s: %w", path, err)
	}

	var s SectionSchema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("błąd parsowania schematu %s: %w", path, err)
	}

	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("schemat %s nieprawidłowy: %w", path, err)
	}
	return &s, nil
}

// Validate checks minimal correctness of the schema.
func (s *SectionSchema) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("pole 'id' jest wymagane")
	}
	if len(s.Sections) == 0 {
		return fmt.Errorf("schemat musi mieć co najmniej jedną sekcję")
	}
	for i, sec := range s.Sections {
		if sec.Name == "" {
			return fmt.Errorf("sekcja[%d]: 'name' nie może być puste", i)
		}
		if sec.Level <= 0 {
			return fmt.Errorf("sekcja[%d]: 'level' musi być > 0", i)
		}
	}
	return nil
}
