package generate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DocumentManifest struct {
	Documents []DocItem `yaml:"documents"`
}

type DocItem struct {
	DocID           string   `yaml:"doc_id"`
	DocType         string   `yaml:"doc_type"`
	TemplatePath    string   `yaml:"template"`
	OutputPath      string   `yaml:"output"`
	TemplateVersion string   `yaml:"template_version,omitempty"`
	DependsOn       []string `yaml:"depends_on"`
	ContextSources  []string `yaml:"context_sources"`
}

func LoadManifest(path string) (*DocumentManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nie można odczytać manifestu %s: %w", path, err)
	}
	return loadManifestBytes(data)
}

// helper for tests
func loadManifestBytes(data []byte) (*DocumentManifest, error) {
	var m DocumentManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("błąd parsowania manifestu: %w", err)
	}

	if len(m.Documents) == 0 {
		return nil, fmt.Errorf("manifest nie zawiera dokumentów")
	}

	for i, d := range m.Documents {
		if d.DocID == "" || d.DocType == "" || d.TemplatePath == "" || d.OutputPath == "" {
			return nil, fmt.Errorf("dokument[%d] ma puste wymagane pola", i)
		}
	}

	return &m, nil
}
