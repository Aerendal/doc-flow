package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"docflow/internal/buildinfo"
)

type sarifReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema,omitempty"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool       sarifTool      `json:"tool"`
	Results    []sarifResult  `json:"results"`
	Properties map[string]any `json:"properties,omitempty"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name    string      `json:"name"`
	Version string      `json:"version,omitempty"`
	Rules   []sarifRule `json:"rules,omitempty"`
}

type sarifRule struct {
	ID               string           `json:"id"`
	Name             string           `json:"name,omitempty"`
	ShortDescription *sarifTextObject `json:"shortDescription,omitempty"`
}

type sarifResult struct {
	RuleID     string          `json:"ruleId"`
	Level      string          `json:"level,omitempty"`
	Message    sarifTextObject `json:"message"`
	Locations  []sarifLocation `json:"locations,omitempty"`
	Properties map[string]any  `json:"properties,omitempty"`
}

type sarifTextObject struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}

func buildValidateSARIFReport(report validateReportJSON) sarifReport {
	issues := append([]validateIssueJSON(nil), report.Issues...)
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		if issues[i].Code != issues[j].Code {
			return issues[i].Code < issues[j].Code
		}
		if issues[i].Line != issues[j].Line {
			return issues[i].Line < issues[j].Line
		}
		if issues[i].DocID != issues[j].DocID {
			return issues[i].DocID < issues[j].DocID
		}
		return issues[i].Message < issues[j].Message
	})

	ruleByCode := map[string]sarifRule{}
	results := make([]sarifResult, 0, len(issues))
	for _, issue := range issues {
		if _, ok := ruleByCode[issue.Code]; !ok {
			ruleByCode[issue.Code] = sarifRule{
				ID:   issue.Code,
				Name: issue.Type,
				ShortDescription: &sarifTextObject{
					Text: issue.Type,
				},
			}
		}

		result := sarifResult{
			RuleID:  issue.Code,
			Level:   validateLevelToSARIF(issue.Level),
			Message: sarifTextObject{Text: issue.Message},
			Properties: map[string]any{
				"type": issue.Type,
			},
		}
		if issue.DocID != "" {
			result.Properties["doc_id"] = issue.DocID
		}
		if issue.Path != "" {
			loc := sarifLocation{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: issue.Path},
				},
			}
			if issue.Line > 0 {
				loc.PhysicalLocation.Region = &sarifRegion{StartLine: issue.Line}
			}
			result.Locations = []sarifLocation{loc}
		}
		results = append(results, result)
	}

	ruleCodes := make([]string, 0, len(ruleByCode))
	for code := range ruleByCode {
		ruleCodes = append(ruleCodes, code)
	}
	sort.Strings(ruleCodes)
	rules := make([]sarifRule, 0, len(ruleCodes))
	for _, code := range ruleCodes {
		rules = append(rules, ruleByCode[code])
	}

	runProps := map[string]any{
		"schema_version": report.SchemaVersion,
	}
	if report.IdentityVersion != "" {
		runProps["identity_version"] = report.IdentityVersion
	}
	if report.FailOn != "" {
		runProps["fail_on"] = report.FailOn
	}
	if report.Show != "" {
		runProps["show"] = report.Show
	}
	if report.Baseline != nil {
		runProps["baseline"] = report.Baseline
	}

	return sarifReport{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:    "docflow",
						Version: buildinfo.FullVersion(),
						Rules:   rules,
					},
				},
				Results:    results,
				Properties: runProps,
			},
		},
	}
}

func validateLevelToSARIF(level string) string {
	switch level {
	case "error":
		return "error"
	case "warn":
		return "warning"
	default:
		return "note"
	}
}

func writeValidateSARIFReport(report validateReportJSON, output string) error {
	sarif := buildValidateSARIFReport(report)
	data, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		return fmt.Errorf("błąd serializacji raportu SARIF: %w", err)
	}
	if output == "" || output == "-" {
		_, err = os.Stdout.Write(append(data, '\n'))
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("błąd utworzenia katalogu dla raportu SARIF: %w", err)
	}
	return os.WriteFile(output, data, 0o644)
}
