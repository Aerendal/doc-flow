package rules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type FamilyRule struct {
	DocType      string   `yaml:"doc_type"`
	RequiredDeps []string `yaml:"required_deps"`
	OptionalDeps []string `yaml:"optional_deps"`
}

type FamilyRules struct {
	Families []FamilyRule `yaml:"families"`
}

// LoadFamilyRules reads YAML file with family rules.
func LoadFamilyRules(path string) (*FamilyRules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nie można odczytać %s: %w", path, err)
	}
	var fr FamilyRules
	if err := yaml.Unmarshal(data, &fr); err != nil {
		return nil, fmt.Errorf("błąd parsowania %s: %w", path, err)
	}
	return &fr, nil
}

// RequiredFor returns required deps for given docType.
func (fr *FamilyRules) RequiredFor(docType string) []string {
	for _, f := range fr.Families {
		if f.DocType == docType {
			return f.RequiredDeps
		}
	}
	return nil
}
