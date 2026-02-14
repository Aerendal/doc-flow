package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generate copies template to output with simple placeholder fill (doc_id, doc_type).
// Returns preview (first N lines) of generated content.
func Generate(item DocItem, previewLines int) (string, error) {
	data, err := os.ReadFile(item.TemplatePath)
	if err != nil {
		return "", fmt.Errorf("cannot read template %s: %w", item.TemplatePath, err)
	}
	content := string(data)
	content = replacePlaceholder(content, "DOC_ID", item.DocID)
	content = replacePlaceholder(content, "DOC_TYPE", item.DocType)
	if item.TemplateVersion != "" {
		content = replacePlaceholder(content, "TEMPLATE_VERSION", item.TemplateVersion)
	}

	if err := os.MkdirAll(filepath.Dir(item.OutputPath), 0o755); err != nil {
		return "", fmt.Errorf("cannot create output dir: %w", err)
	}
	if err := os.WriteFile(item.OutputPath, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("cannot write output %s: %w", item.OutputPath, err)
	}
	return preview(content, previewLines), nil
}

func replacePlaceholder(content, key, value string) string {
	// very simple placeholder format: {{KEY}}
	placeholder := fmt.Sprintf("{{%s}}", key)
	return strings.ReplaceAll(content, placeholder, value)
}

func preview(content string, lines int) string {
	if lines <= 0 {
		return ""
	}
	parts := strings.Split(content, "\n")
	if len(parts) > lines {
		parts = parts[:lines]
	}
	return strings.Join(parts, "\n")
}
