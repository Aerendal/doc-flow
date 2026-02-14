package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"docflow/pkg/compliance"

	"github.com/spf13/cobra"
)

func baselineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "baseline",
		Short: "Narzędzia do baseline (migracja i utrzymanie)",
	}
	cmd.AddCommand(baselineMigrateCmd())
	return cmd
}

func baselineMigrateCmd() *cobra.Command {
	var inPath string
	var outPath string
	var kind string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migruj baseline identity v1 -> v2",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inPath == "" || outPath == "" {
				return fmt.Errorf("podaj --in i --out")
			}
			kindValue := kind
			if kindValue == "" {
				detected, err := detectBaselineKind(inPath)
				if err != nil {
					return err
				}
				kindValue = detected
			}

			switch kindValue {
			case "validate":
				report, err := loadValidateBaseline(inPath)
				if err != nil {
					return err
				}
				migrateValidateBaseline(&report)
				return writeJSONFile(outPath, report)
			case "compliance":
				sum, err := loadComplianceBaseline(inPath)
				if err != nil {
					return err
				}
				migrateComplianceBaseline(&sum)
				return writeJSONFile(outPath, sum)
			default:
				return fmt.Errorf("nieobsługiwany --kind: %s (dozwolone: validate|compliance)", kindValue)
			}
		},
	}
	cmd.Flags().StringVar(&inPath, "in", "", "wejściowy baseline JSON")
	cmd.Flags().StringVar(&outPath, "out", "", "wyjściowy baseline JSON")
	cmd.Flags().StringVar(&kind, "kind", "", "typ baseline: validate|compliance (opcjonalnie, auto-detect)")
	return cmd
}

func detectBaselineKind(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return "", err
	}
	if _, ok := probe["issues"]; ok {
		return "validate", nil
	}
	if _, ok := probe["docs"]; ok {
		return "compliance", nil
	}
	return "", fmt.Errorf("nie można wykryć typu baseline (brak klucza issues/docs)")
}

func loadValidateBaseline(path string) (validateReportJSON, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return validateReportJSON{}, err
	}
	var report validateReportJSON
	if err := json.Unmarshal(data, &report); err != nil {
		return validateReportJSON{}, err
	}
	return report, nil
}

func loadComplianceBaseline(path string) (compliance.Summary, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return compliance.Summary{}, err
	}
	var sum compliance.Summary
	if err := json.Unmarshal(data, &sum); err != nil {
		return compliance.Summary{}, err
	}
	return sum, nil
}

func migrateValidateBaseline(report *validateReportJSON) {
	if report == nil {
		return
	}
	if report.SchemaVersion == "" {
		report.SchemaVersion = "1.0"
	}
	report.IdentityVersion = currentIdentityVersion()
	for i := range report.Issues {
		report.Issues[i].Details = inferValidateIssueDetails(report.Issues[i])
	}
}

func migrateComplianceBaseline(sum *compliance.Summary) {
	if sum == nil {
		return
	}
	if sum.SchemaVersion == "" {
		sum.SchemaVersion = "1.0"
	}
	sum.IdentityVersion = currentIdentityVersion()
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
