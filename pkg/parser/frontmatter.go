package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title          string   `yaml:"title"`
	DocID          string   `yaml:"doc_id"`
	DocType        string   `yaml:"doc_type"`
	Version        string   `yaml:"version"`
	Status         string   `yaml:"status"`
	Priority       string   `yaml:"priority"`
	Owner          string   `yaml:"owner"`
	Created        string   `yaml:"created"`
	Updated        string   `yaml:"updated"`
	Tags           []string `yaml:"tags"`
	DependsOn      []string `yaml:"depends_on"`
	ContextSources []string `yaml:"context_sources"`
	TemplateSource string   `yaml:"template_source"`
	Language       string   `yaml:"language"`
	Audience       string   `yaml:"audience"`
	ReviewCycle    string   `yaml:"review_cycle"`

	Extra map[string]interface{} `yaml:"-"`
}

type ParseResult struct {
	Frontmatter    *Frontmatter
	Body           string
	HasFrontmatter bool
	Warnings       []string
	RawFields      map[string]bool
	RawYAML        string
}

type ParseError struct {
	File    string
	Line    int
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("linia %d: %s", e.Line, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func ParseFrontmatter(r io.Reader, filename string) (*ParseResult, error) {
	scanner := bufio.NewScanner(r)
	result := &ParseResult{}

	if !scanner.Scan() {
		return result, nil
	}
	firstLine := scanner.Text()
	if strings.TrimSpace(firstLine) != "---" {
		var bodyLines []string
		bodyLines = append(bodyLines, firstLine)
		for scanner.Scan() {
			bodyLines = append(bodyLines, scanner.Text())
		}
		result.Body = strings.Join(bodyLines, "\n")
		return result, nil
	}

	var yamlLines []string
	foundEnd := false
	lineNum := 1

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			foundEnd = true
			break
		}
		yamlLines = append(yamlLines, line)
	}

	if !foundEnd {
		return nil, &ParseError{
			File:    filename,
			Line:    1,
			Message: "brak zamykającego '---' dla frontmatter",
		}
	}

	result.HasFrontmatter = true
	yamlContent := strings.Join(yamlLines, "\n")
	result.RawYAML = yamlContent

	// Capture raw keys before mapping to struct.
	var raw map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &raw); err != nil {
		return nil, &ParseError{
			File:    filename,
			Line:    2,
			Message: fmt.Sprintf("błąd parsowania YAML: %v", err),
			Err:     err,
		}
	}

	fm := &Frontmatter{}
	if err := yaml.Unmarshal([]byte(yamlContent), fm); err != nil {
		return nil, &ParseError{
			File:    filename,
			Line:    2,
			Message: fmt.Sprintf("błąd parsowania YAML: %v", err),
			Err:     err,
		}
	}

	if raw != nil {
		result.RawFields = make(map[string]bool, len(raw))
		knownKeys := map[string]bool{
			"title": true, "doc_id": true, "doc_type": true, "version": true,
			"status": true, "priority": true, "owner": true, "created": true,
			"updated": true, "tags": true, "depends_on": true, "context_sources": true,
			"language": true, "audience": true, "review_cycle": true,
		}
		extras := make(map[string]interface{})
		for k, v := range raw {
			result.RawFields[k] = true
			if !knownKeys[k] {
				extras[k] = v
			}
		}
		if len(extras) > 0 {
			fm.Extra = extras
		}
	}

	applyDefaults(fm, result)
	result.Frontmatter = fm

	var bodyLines []string
	for scanner.Scan() {
		bodyLines = append(bodyLines, scanner.Text())
	}
	result.Body = strings.Join(bodyLines, "\n")

	return result, nil
}

func ParseFrontmatterString(content string, filename string) (*ParseResult, error) {
	return ParseFrontmatter(strings.NewReader(content), filename)
}

type MetaBlock struct {
	Title   string
	Status  string
	Owner   string
	Version string
}

func ExtractMetaBlock(content string) *MetaBlock {
	lines := strings.Split(content, "\n")
	inMetaSection := false
	meta := &MetaBlock{}
	found := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "## Metadane" || trimmed == "## Meta" {
			inMetaSection = true
			continue
		}

		if inMetaSection && strings.HasPrefix(trimmed, "## ") {
			break
		}

		if !inMetaSection {
			continue
		}

		if strings.HasPrefix(trimmed, "- ") {
			kv := strings.TrimPrefix(trimmed, "- ")
			parts := strings.SplitN(kv, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			val := strings.TrimSpace(parts[1])

			if val == "" || val == "[osoba/rola]" || val == "RRRR-MM-DD" {
				continue
			}

			switch key {
			case "właściciel", "owner":
				meta.Owner = val
				found = true
			case "wersja", "version":
				meta.Version = val
				found = true
			case "status":
				meta.Status = val
				found = true
			case "tytuł", "title":
				meta.Title = val
				found = true
			}
		}
	}

	if !found {
		return nil
	}
	return meta
}

func applyDefaults(fm *Frontmatter, result *ParseResult) {
	if fm.Status == "" {
		fm.Status = "draft"
		result.Warnings = append(result.Warnings, "brak pola 'status', ustawiono domyślne: draft")
	}
	if fm.DocType == "" {
		fm.DocType = "unknown"
		result.Warnings = append(result.Warnings, "brak pola 'doc_type', ustawiono domyślne: unknown")
	}
	if fm.Version == "" {
		fm.Version = "0.1.0"
	}
	if fm.Priority == "" {
		fm.Priority = "normal"
	}
	if fm.Language == "" {
		fm.Language = "pl"
	}
	if fm.Created == "" {
		fm.Created = time.Now().Format("2006-01-02")
	}
}
