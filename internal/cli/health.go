package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"docflow/internal/buildinfo"
	"docflow/pkg/compliance"
	"docflow/pkg/config"

	"github.com/spf13/cobra"
)

func healthCmd() *cobra.Command {
	var (
		ciMode          bool
		bundleDir       string
		runID           string
		baselineMode    string
		baselineDir     string
		portable        bool
		continueOnError bool
		rulesPath       string
	)

	cmd := &cobra.Command{
		Use:   "health",
		Short: "One-command orchestration: validate + compliance + bundle artefaktów",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if baselineMode != "repo" && baselineMode != "artifact" && baselineMode != "none" {
				return usageError(
					"DOCFLOW.CLI.INVALID_BASELINE_MODE",
					fmt.Sprintf("nieobsługiwane --baseline-mode: %s (dozwolone: repo|artifact|none)", baselineMode),
					map[string]any{"baseline_mode": baselineMode},
				)
			}
			if runID == "" {
				runID = defaultHealthRunID()
			}
			if strings.Contains(runID, "/") || strings.Contains(runID, string(os.PathSeparator)) {
				return usageError(
					"DOCFLOW.CLI.INVALID_RUN_ID",
					"run_id nie może zawierać separatorów katalogów",
					map[string]any{"run_id": runID},
				)
			}
			runDir := filepath.Join(bundleDir, runID)
			if err := os.MkdirAll(runDir, 0o755); err != nil {
				return runtimeError(
					"DOCFLOW.RUNTIME.BUNDLE_DIR_CREATE_FAILED",
					fmt.Sprintf("nie można utworzyć katalogu bundle %s", runDir),
					err,
					map[string]any{"bundle_dir": runDir},
				)
			}

			mode := "local"
			if ciMode {
				mode = "ci"
			}
			validateAgainst, complianceAgainst, missing := resolveHealthBaselinePaths(baselineMode, baselineDir)
			if len(missing) > 0 {
				if ciMode {
					return runtimeError(
						"DOCFLOW.RUNTIME.BASELINE_NOT_FOUND",
						fmt.Sprintf("brakuje plików baseline: %s", strings.Join(missing, ", ")),
						fmt.Errorf("missing baseline"),
						map[string]any{"missing": missing, "baseline_mode": baselineMode, "baseline_dir": baselineDir},
					)
				}
				for _, miss := range missing {
					fmt.Fprintf(os.Stderr, "WARN: baseline %s nie istnieje; health przechodzi bez tego baseline\n", miss)
				}
				if containsString(missing, "validate.json") {
					validateAgainst = ""
				}
				if containsString(missing, "compliance.json") {
					complianceAgainst = ""
				}
			}

			validateJSONPath := filepath.Join(runDir, "validate.json")
			validateSARIFPath := filepath.Join(runDir, "validate.sarif")
			complianceJSONPath := filepath.Join(runDir, "compliance.json")
			summaryJSONPath := filepath.Join(runDir, "summary.json")
			summaryMDPath := filepath.Join(runDir, "summary.md")
			metaJSONPath := filepath.Join(runDir, "meta.json")

			summary := healthSummary{
				SchemaVersion: "1.0",
				RunID:         runID,
				Mode:          mode,
				Baseline: healthSummaryBaseline{
					Mode:           baselineMode,
					ValidatePath:   validateAgainst,
					CompliancePath: complianceAgainst,
				},
				Validate: healthSummaryValidate{
					JSON:  "validate.json",
					SARIF: "validate.sarif",
				},
				Compliance: healthSummaryCompliance{
					JSON: "compliance.json",
				},
			}
			meta := healthMeta{
				SchemaVersion: "1.0",
				RunID:         runID,
				Mode:          mode,
				ToolVersion:   buildinfo.Version,
				ToolCommit:    buildinfo.Commit,
				ToolDate:      buildinfo.Date,
				BundleDir:     runDir,
				BaselineMode:  baselineMode,
				BaselineDir:   baselineDir,
				Portable:      portable,
				ConfigPath:    resolvedConfigPathForHealth(),
				RulesPath:     rulesPath,
			}
			cfg, cfgErr := loadConfig()
			if cfgErr == nil {
				cacheOn, cachePath := effectiveCacheConfig(cfg)
				meta.CacheEnabled = cacheOn
				meta.CacheDir = cachePath
			}

			validateJSONRes := runHealthSubcommand(
				validateCmd,
				healthValidateArgs(validateJSONPath, "json", ciMode, validateAgainst),
			)
			summary.Validate.JSONExit = validateJSONRes.Exit
			if validateJSONRes.Error != "" {
				summary.Validate.Error = validateJSONRes.Error
			}

			if report, err := loadValidateReport(validateJSONPath); err == nil {
				summary.Baseline.ValidateLoaded = report.Baseline != nil && report.Baseline.Loaded
				summary.Validate.ErrorCount = report.ErrorCount
				summary.Validate.WarnCount = report.WarnCount
				if summary.Baseline.ValidateLoaded {
					summary.Validate.NewErrors = report.NewErrorCount
					summary.Validate.NewWarnings = report.NewWarnCount
				} else {
					summary.Validate.NewErrors = report.ErrorCount
					summary.Validate.NewWarnings = report.WarnCount
				}
			} else {
				summary.Validate.ParseError = err.Error()
			}

			validateSARIFRes := healthCommandResult{Skipped: true}
			if !shouldStopHealth(validateJSONRes.Exit, continueOnError) {
				validateSARIFRes = runHealthSubcommand(
					validateCmd,
					healthValidateArgs(validateSARIFPath, "sarif", ciMode, validateAgainst),
				)
			}
			summary.Validate.SARIFExit = validateSARIFRes.Exit
			summary.Validate.SARIFSkipped = validateSARIFRes.Skipped

			complianceRes := healthCommandResult{Skipped: true}
			if !shouldStopHealth(validateJSONRes.Exit, continueOnError) &&
				!shouldStopHealth(validateSARIFRes.Exit, continueOnError) {
				complianceRes = runHealthSubcommand(
					complianceCmd,
					healthComplianceArgs(complianceJSONPath, rulesPath, ciMode, complianceAgainst),
				)
			}
			summary.Compliance.JSONExit = complianceRes.Exit
			summary.Compliance.Skipped = complianceRes.Skipped
			if complianceRes.Error != "" {
				summary.Compliance.Error = complianceRes.Error
			}
			if !complianceRes.Skipped {
				if comp, err := loadComplianceReport(complianceJSONPath); err == nil {
					summary.Baseline.ComplianceLoaded = comp.Baseline != nil && comp.Baseline.Loaded
					summary.Compliance.Failed = comp.Failed
					if summary.Baseline.ComplianceLoaded {
						summary.Compliance.NewFailed = comp.NewFailed
						if len(comp.NewViolations) > 0 {
							summary.Compliance.NewViolations = comp.NewViolations
						}
					} else {
						summary.Compliance.NewFailed = comp.Failed
						if len(comp.ViolationsCount) > 0 {
							summary.Compliance.NewViolations = comp.ViolationsCount
						}
					}
				} else {
					summary.Compliance.ParseError = err.Error()
				}
			}

			summary.Validate.Exit = healthOverallExit(validateJSONRes.Exit, validateSARIFRes.Exit)
			summary.Compliance.Exit = complianceRes.Exit
			summary.OverallExit = healthOverallExit(
				validateJSONRes.Exit,
				validateSARIFRes.Exit,
				complianceRes.Exit,
			)

			meta.ValidateJSONExit = validateJSONRes.Exit
			meta.ValidateSARIFExit = validateSARIFRes.Exit
			meta.ComplianceJSONExit = complianceRes.Exit
			meta.OverallExit = summary.OverallExit

			if err := writeHealthJSONFile(summaryJSONPath, summary); err != nil {
				return runtimeError(
					"DOCFLOW.RUNTIME.SUMMARY_WRITE_FAILED",
					fmt.Sprintf("nie można zapisać %s", summaryJSONPath),
					err,
					map[string]any{"path": summaryJSONPath},
				)
			}
			if err := os.WriteFile(summaryMDPath, []byte(renderHealthSummaryMarkdown(summary, runDir)), 0o644); err != nil {
				return runtimeError(
					"DOCFLOW.RUNTIME.SUMMARY_WRITE_FAILED",
					fmt.Sprintf("nie można zapisać %s", summaryMDPath),
					err,
					map[string]any{"path": summaryMDPath},
				)
			}
			if err := writeHealthJSONFile(metaJSONPath, meta); err != nil {
				return runtimeError(
					"DOCFLOW.RUNTIME.META_WRITE_FAILED",
					fmt.Sprintf("nie można zapisać %s", metaJSONPath),
					err,
					map[string]any{"path": metaJSONPath},
				)
			}

			if summary.OverallExit == 0 {
				fmt.Printf("Health OK. Bundle: %s\n", runDir)
				return nil
			}

			primary := firstHealthErrorByExit(summary.OverallExit, validateJSONRes, validateSARIFRes, complianceRes)
			switch summary.OverallExit {
			case exitCodeUsage:
				return usageError(
					"DOCFLOW.CLI.HEALTH_SUBCOMMAND_FAILED",
					fmt.Sprintf("health zakończony błędem usage: %s", primary),
					map[string]any{"summary": summaryJSONPath, "run_id": runID},
				)
			case exitCodeRuntime:
				return runtimeError(
					"DOCFLOW.RUNTIME.HEALTH_SUBCOMMAND_FAILED",
					fmt.Sprintf("health zakończony błędem runtime: %s", primary),
					fmt.Errorf("%s", primary),
					map[string]any{"summary": summaryJSONPath, "run_id": runID},
				)
			default:
				return fmt.Errorf("health zakończony: wykryto problemy (summary: %s)", summaryJSONPath)
			}
		},
	}

	cmd.Flags().BoolVar(&ciMode, "ci", false, "tryb CI: strict + fail-on new (jeśli baseline dostępny)")
	cmd.Flags().StringVar(&bundleDir, "bundle-dir", ".docflow/out", "katalog na bundle artefaktów")
	cmd.Flags().StringVar(&runID, "run-id", "", "identyfikator runu (domyślnie UTC timestamp)")
	cmd.Flags().StringVar(&baselineMode, "baseline-mode", "repo", "źródło baseline: repo|artifact|none")
	cmd.Flags().StringVar(&baselineDir, "baseline-dir", ".docflow/baseline", "katalog baseline (validate.json/compliance.json)")
	cmd.Flags().BoolVar(&portable, "portable", false, "portable metadata (bez ścieżek zależnych od hosta)")
	cmd.Flags().BoolVar(&continueOnError, "continue-on-error", false, "kontynuuj kolejne kroki mimo błędów usage/runtime")
	cmd.Flags().StringVar(&rulesPath, "rules", "docs/_meta/GOVERNANCE_RULES.yaml", "ścieżka do GOVERNANCE_RULES.yaml dla compliance")
	return cmd
}

type healthSummary struct {
	SchemaVersion string                  `json:"schema_version"`
	RunID         string                  `json:"run_id"`
	Mode          string                  `json:"mode"`
	Baseline      healthSummaryBaseline   `json:"baseline"`
	Validate      healthSummaryValidate   `json:"validate"`
	Compliance    healthSummaryCompliance `json:"compliance"`
	OverallExit   int                     `json:"overall_exit"`
}

type healthSummaryBaseline struct {
	Mode             string `json:"mode"`
	ValidatePath     string `json:"validate_path,omitempty"`
	CompliancePath   string `json:"compliance_path,omitempty"`
	ValidateLoaded   bool   `json:"validate_loaded"`
	ComplianceLoaded bool   `json:"compliance_loaded"`
}

type healthSummaryValidate struct {
	Exit         int    `json:"exit"`
	JSONExit     int    `json:"json_exit"`
	SARIFExit    int    `json:"sarif_exit"`
	SARIFSkipped bool   `json:"sarif_skipped,omitempty"`
	ErrorCount   int    `json:"error_count,omitempty"`
	WarnCount    int    `json:"warn_count,omitempty"`
	NewErrors    int    `json:"new_errors,omitempty"`
	NewWarnings  int    `json:"new_warnings,omitempty"`
	JSON         string `json:"json"`
	SARIF        string `json:"sarif"`
	Error        string `json:"error,omitempty"`
	ParseError   string `json:"parse_error,omitempty"`
}

type healthSummaryCompliance struct {
	Exit          int            `json:"exit"`
	JSONExit      int            `json:"json_exit"`
	Skipped       bool           `json:"skipped,omitempty"`
	Failed        int            `json:"failed,omitempty"`
	NewFailed     int            `json:"new_failed,omitempty"`
	NewViolations map[string]int `json:"new_violations,omitempty"`
	JSON          string         `json:"json"`
	Error         string         `json:"error,omitempty"`
	ParseError    string         `json:"parse_error,omitempty"`
}

type healthMeta struct {
	SchemaVersion      string `json:"schema_version"`
	RunID              string `json:"run_id"`
	Mode               string `json:"mode"`
	ToolVersion        string `json:"tool_version"`
	ToolCommit         string `json:"tool_commit"`
	ToolDate           string `json:"tool_date"`
	BundleDir          string `json:"bundle_dir"`
	ConfigPath         string `json:"config_path,omitempty"`
	RulesPath          string `json:"rules_path,omitempty"`
	BaselineMode       string `json:"baseline_mode"`
	BaselineDir        string `json:"baseline_dir"`
	Portable           bool   `json:"portable"`
	CacheEnabled       bool   `json:"cache_enabled"`
	CacheDir           string `json:"cache_dir,omitempty"`
	ValidateJSONExit   int    `json:"validate_json_exit"`
	ValidateSARIFExit  int    `json:"validate_sarif_exit"`
	ComplianceJSONExit int    `json:"compliance_json_exit"`
	OverallExit        int    `json:"overall_exit"`
}

type healthCommandResult struct {
	Exit    int
	Error   string
	Skipped bool
}

func runHealthSubcommand(factory func() *cobra.Command, args []string) healthCommandResult {
	cmd := factory()
	cmd.SetArgs(args)
	err := cmd.Execute()
	if err != nil {
		return healthCommandResult{
			Exit:  ExitCode(err),
			Error: err.Error(),
		}
	}
	return healthCommandResult{}
}

func healthValidateArgs(outputPath, format string, ciMode bool, againstPath string) []string {
	args := []string{"--format", format, "--output", outputPath}
	if ciMode {
		args = append(args, "--strict")
		if againstPath != "" {
			args = append(args, "--fail-on", "new", "--show", "new", "--against", againstPath)
		} else {
			args = append(args, "--fail-on", "any", "--show", "all")
		}
		return args
	}
	args = append(args, "--warn", "--fail-on", "any", "--show", "all")
	if againstPath != "" {
		args = append(args, "--against", againstPath)
	}
	return args
}

func healthComplianceArgs(outputPath, rulesPath string, ciMode bool, againstPath string) []string {
	args := []string{
		"--format", "json",
		"--output", outputPath,
		"--rules", rulesPath,
	}
	if ciMode {
		args = append(args, "--strict")
		if againstPath != "" {
			args = append(args, "--fail-on", "new", "--show", "new", "--against", againstPath)
		} else {
			args = append(args, "--fail-on", "any", "--show", "all")
		}
		return args
	}
	args = append(args, "--warn", "--fail-on", "any", "--show", "all")
	if againstPath != "" {
		args = append(args, "--against", againstPath)
	}
	return args
}

func resolveHealthBaselinePaths(mode, baselineDir string) (string, string, []string) {
	if mode == "none" {
		return "", "", nil
	}
	validatePath := filepath.Join(baselineDir, "validate.json")
	compliancePath := filepath.Join(baselineDir, "compliance.json")
	var missing []string
	if _, err := os.Stat(validatePath); err != nil {
		missing = append(missing, "validate.json")
	}
	if _, err := os.Stat(compliancePath); err != nil {
		missing = append(missing, "compliance.json")
	}
	return validatePath, compliancePath, missing
}

func loadValidateReport(path string) (validateReportJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return validateReportJSON{}, err
	}
	var report validateReportJSON
	if err := json.Unmarshal(data, &report); err != nil {
		return validateReportJSON{}, err
	}
	return report, nil
}

func loadComplianceReport(path string) (compliance.Summary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return compliance.Summary{}, err
	}
	var sum compliance.Summary
	if err := json.Unmarshal(data, &sum); err != nil {
		return compliance.Summary{}, err
	}
	return sum, nil
}

func writeHealthJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func renderHealthSummaryMarkdown(summary healthSummary, runDir string) string {
	var b strings.Builder
	b.WriteString("# docflow health summary\n\n")
	b.WriteString(fmt.Sprintf("- Run ID: `%s`\n", summary.RunID))
	b.WriteString(fmt.Sprintf("- Mode: `%s`\n", summary.Mode))
	b.WriteString(fmt.Sprintf("- Overall exit: `%d`\n", summary.OverallExit))
	b.WriteString(fmt.Sprintf("- Bundle: `%s`\n", runDir))
	b.WriteString("\n## Validate\n")
	b.WriteString(fmt.Sprintf("- exit: `%d` (json=%d, sarif=%d)\n", summary.Validate.Exit, summary.Validate.JSONExit, summary.Validate.SARIFExit))
	b.WriteString(fmt.Sprintf("- new errors/warnings: `%d/%d`\n", summary.Validate.NewErrors, summary.Validate.NewWarnings))
	b.WriteString(fmt.Sprintf("- report: `%s`\n", filepath.Join(runDir, summary.Validate.JSON)))
	b.WriteString(fmt.Sprintf("- sarif: `%s`\n", filepath.Join(runDir, summary.Validate.SARIF)))
	if summary.Validate.Error != "" {
		b.WriteString(fmt.Sprintf("- error: `%s`\n", summary.Validate.Error))
	}
	if summary.Validate.ParseError != "" {
		b.WriteString(fmt.Sprintf("- parse_error: `%s`\n", summary.Validate.ParseError))
	}
	b.WriteString("\n## Compliance\n")
	b.WriteString(fmt.Sprintf("- exit: `%d`\n", summary.Compliance.Exit))
	if summary.Compliance.Skipped {
		b.WriteString("- skipped: `true`\n")
	} else {
		b.WriteString(fmt.Sprintf("- new failed: `%d`\n", summary.Compliance.NewFailed))
	}
	b.WriteString(fmt.Sprintf("- report: `%s`\n", filepath.Join(runDir, summary.Compliance.JSON)))
	if len(summary.Compliance.NewViolations) > 0 {
		keys := make([]string, 0, len(summary.Compliance.NewViolations))
		for k := range summary.Compliance.NewViolations {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString("- new violations:\n")
		for _, k := range keys {
			b.WriteString(fmt.Sprintf("  - `%s`: %d\n", k, summary.Compliance.NewViolations[k]))
		}
	}
	if summary.Compliance.Error != "" {
		b.WriteString(fmt.Sprintf("- error: `%s`\n", summary.Compliance.Error))
	}
	if summary.Compliance.ParseError != "" {
		b.WriteString(fmt.Sprintf("- parse_error: `%s`\n", summary.Compliance.ParseError))
	}
	return b.String()
}

func shouldStopHealth(exitCode int, continueOnError bool) bool {
	if continueOnError {
		return false
	}
	return exitCode == exitCodeUsage || exitCode == exitCodeRuntime
}

func healthOverallExit(codes ...int) int {
	overall := 0
	for _, code := range codes {
		if code > overall {
			overall = code
		}
	}
	return overall
}

func defaultHealthRunID() string {
	return time.Now().UTC().Format("20060102T150405Z")
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func firstHealthErrorByExit(exit int, results ...healthCommandResult) string {
	for _, res := range results {
		if res.Exit == exit && res.Error != "" {
			return res.Error
		}
	}
	if exit == 0 {
		return "none"
	}
	return "unknown error"
}

func resolvedConfigPathForHealth() string {
	if configFile != "" {
		return configFile
	}
	if p := config.FindConfigFile(); p != "" {
		return p
	}
	return "docflow.yaml"
}
