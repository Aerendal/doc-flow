package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"docflow/internal/engine"
	"docflow/pkg/compliance"
	"docflow/pkg/governance"

	"github.com/spf13/cobra"
)

func complianceCmd() *cobra.Command {
	var rulesPath string
	var output string
	var format string
	var htmlPath string
	var strict bool
	var warnOnly bool
	var againstPath string
	var failOn string
	var show string
	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "Raport zgodności dokumentów z regułami governance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if format != "text" && format != "json" && format != "html" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.INVALID_FORMAT",
					fmt.Sprintf("nieobsługiwany format: %s (dozwolone: text|json|html)", format),
					map[string]any{"format": format},
				)
			}
			if strict && warnOnly {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.CONFLICTING_FLAGS",
					"--strict i --warn nie mogą być użyte jednocześnie",
					map[string]any{"flags": []string{"--strict", "--warn"}},
				)
			}
			if failOn != "any" && failOn != "new" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.INVALID_FAIL_ON",
					fmt.Sprintf("nieobsługiwane --fail-on: %s (dozwolone: any|new)", failOn),
					map[string]any{"fail_on": failOn},
				)
			}
			if show != "all" && show != "new" && show != "existing" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.INVALID_SHOW",
					fmt.Sprintf("nieobsługiwane --show: %s (dozwolone: all|new|existing)", show),
					map[string]any{"show": show},
				)
			}
			if failOn == "new" && againstPath == "" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.MISSING_AGAINST",
					"podaj --against dla --fail-on new",
					map[string]any{"flag": "--against"},
				)
			}
			if rulesPath == "" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.MISSING_RULES",
					"podaj --rules (GOVERNANCE_RULES.yaml)",
					map[string]any{"flag": "--rules"},
				)
			}
			rulesPath = filepath.Clean(rulesPath)
			rulesChecksum, err := calculateRulesChecksum(rulesPath)
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.RULES_CHECKSUM_FAILED",
					fmt.Sprintf("nie można obliczyć checksum dla %s", rulesPath),
					err,
					map[string]any{"path": rulesPath},
				)
			}
			rules, err := governance.Load(rulesPath)
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.RULES_LOAD_FAILED",
					fmt.Sprintf("nie można wczytać reguł governance z %s", rulesPath),
					err,
					map[string]any{"path": rulesPath},
				)
			}
			cfg, err := loadConfig()
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.CONFIG_LOAD_FAILED",
					"nie można wczytać konfiguracji",
					err,
					nil,
				)
			}
			execRes, err := collectExecutionResult("compliance", cfg)
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.INDEX_BUILD_FAILED",
					"nie można przygotować facts execution layer",
					err,
					map[string]any{"docs_root": cfg.DocsRoot},
				)
			}
			sum, err := compliance.ReportWithFacts(execRes.Index, rules, contentFactsByPath(execRes.Facts))
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.COMPLIANCE_REPORT_FAILED",
					"nie można wygenerować raportu compliance",
					err,
					nil,
				)
			}
			sum.RulesPath = rulesPath
			sum.RulesChecksum = "sha256:" + rulesChecksum
			sum.FailOn = failOn
			sum.Show = show
			if againstPath != "" {
				if err := applyComplianceBaseline(sum, againstPath, show); err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.BASELINE_LOAD_FAILED",
						fmt.Sprintf("nie można wczytać baseline compliance z %s", againstPath),
						err,
						map[string]any{"path": againstPath},
					)
				}
			}
			canonical := engine.BuildComplianceReport(sum)
			engine.ApplyComplianceReport(sum, canonical)

			switch format {
			case "json":
				if output == "" {
					output = "-"
				}
				if err := compliance.SaveJSON(sum, output); err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.OUTPUT_WRITE_FAILED",
						fmt.Sprintf("nie można zapisać raportu compliance do %s", output),
						err,
						map[string]any{"output": output},
					)
				}
				if output != "-" {
					fmt.Printf("Raport zapisany do %s\n", output)
				}
			case "html":
				if htmlPath == "" {
					htmlPath = ".docflow/output/compliance.html"
				}
				if err := compliance.SaveHTML(sum, htmlPath); err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.OUTPUT_WRITE_FAILED",
						fmt.Sprintf("nie można zapisać raportu HTML do %s", htmlPath),
						err,
						map[string]any{"output": htmlPath},
					)
				}
				fmt.Printf("Raport HTML zapisany do %s\n", htmlPath)
			default:
				printComplianceText(sum)
			}

			failedForGate := sum.Failed
			if failOn == "new" {
				failedForGate = sum.NewFailed
			}

			if strict && failedForGate > 0 && !warnOnly {
				if failOn == "new" {
					return fmt.Errorf("compliance zakończony: %d dokumentów z nowymi naruszeniami", failedForGate)
				}
				return fmt.Errorf("compliance zakończony: %d dokumentów niezgodnych", failedForGate)
			}
			if warnOnly && failedForGate > 0 {
				if failOn == "new" {
					fmt.Fprintf(os.Stderr, "WARN: compliance wykrył %d dokumentów z nowymi naruszeniami (tryb --warn, brak blokady)\n", failedForGate)
				} else {
					fmt.Fprintf(os.Stderr, "WARN: compliance wykrył %d dokumentów niezgodnych (tryb --warn, brak blokady)\n", failedForGate)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&rulesPath, "rules", "", "ścieżka do GOVERNANCE_RULES.yaml")
	cmd.Flags().StringVar(&output, "output", "", "plik wyjściowy (json)")
	cmd.Flags().StringVar(&format, "format", "text", "format: text|json|html")
	cmd.Flags().StringVar(&htmlPath, "html", "", "ścieżka HTML (opcjonalnie, gdy --format html)")
	cmd.Flags().BoolVar(&strict, "strict", false, "zwróć błąd exit gdy wykryto naruszenia compliance")
	cmd.Flags().BoolVar(&warnOnly, "warn", false, "raportuj, ale nie blokuj nawet przy naruszeniach")
	cmd.Flags().StringVar(&againstPath, "against", "", "ścieżka do baseline compliance JSON")
	cmd.Flags().StringVar(&failOn, "fail-on", "any", "gating: any|new")
	cmd.Flags().StringVar(&show, "show", "all", "które dokumenty zwrócić w JSON: all|new|existing")
	return cmd
}

func calculateRulesChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func printComplianceText(sum *compliance.Summary) {
	fmt.Printf("Dokumenty: %d | Passed: %d | Failed: %d | PassRate: %.1f%%\n",
		sum.Documents, sum.Passed, sum.Failed, sum.PassRate*100)
	fmt.Println(strings.Repeat("─", 60))
	if len(sum.ViolationsCount) > 0 {
		fmt.Println("Violations:")
		keys := make([]string, 0, len(sum.ViolationsCount))
		for k := range sum.ViolationsCount {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := sum.ViolationsCount[k]
			fmt.Printf("  %-30s %d\n", k, v)
		}
	}
	if sum.Failed > 0 {
		fmt.Println("\nNon-compliant docs:")
		for _, d := range sum.Docs {
			if len(d.Violations) == 0 {
				continue
			}
			fmt.Printf("- %s (%s): %v\n", d.Path, d.DocID, d.Violations)
		}
	}
}
