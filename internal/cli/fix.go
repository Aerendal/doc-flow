package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"docflow/internal/engine"
	engfix "docflow/internal/engine/fix"
	"docflow/internal/validator"
	"docflow/pkg/governance"

	"github.com/spf13/cobra"
)

func fixCmd() *cobra.Command {
	var (
		format          string
		output          string
		apply           bool
		backupDir       string
		scopeCSV        string
		onlyCSV         string
		maxFiles        int
		maxChanges      int
		againstPath     string
		show            string
		statusAware     bool
		governanceRules string
	)

	cmd := &cobra.Command{
		Use:   "fix",
		Short: "Generuj i opcjonalnie aplikuj bezpieczne poprawki",
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "diff" && format != "json" {
				return fmt.Errorf("nieobsługiwany format: %s (dozwolone: diff|json)", format)
			}
			if show != "all" && show != "new" && show != "existing" {
				return fmt.Errorf("nieobsługiwane --show: %s (dozwolone: all|new|existing)", show)
			}
			if againstPath == "" && show != "all" {
				return fmt.Errorf("--show=%s wymaga --against", show)
			}

			scopes, err := engfix.ParseScopesCSV(scopeCSV)
			if err != nil {
				return err
			}

			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			execRes, err := collectExecutionResult("fix", cfg)
			if err != nil {
				return err
			}

			opts := validatorOptions(cfg)
			opts.StatusAware = statusAware
			if governanceRules != "" {
				gr, err := governance.Load(governanceRules)
				if err != nil {
					return fmt.Errorf("nie można wczytać reguł governance z %s: %w", governanceRules, err)
				}
				opts.GovernanceRules = gr
			}

			report, err := validator.ValidateFacts(execRes.Facts, execRes.Index, opts)
			if err != nil {
				return err
			}
			reportJSON := buildValidateJSONReport(report)
			if againstPath != "" {
				reportJSON, err = applyValidateBaseline(reportJSON, againstPath, show)
				if err != nil {
					return err
				}
			}

			issues := make([]engine.Issue, 0, len(reportJSON.Issues))
			for _, issue := range reportJSON.Issues {
				issues = append(issues, toEngineIssue(issue))
			}

			files := make([]engfix.FileInput, 0, len(execRes.Facts))
			for _, fact := range execRes.Facts {
				if fact.Path == "" {
					continue
				}
				files = append(files, engfix.FileInput{
					Path:    fact.Path,
					Content: fact.Content,
				})
			}
			sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })

			plan, err := engfix.BuildPlan(cfg.DocsRoot, issues, files, engfix.Options{
				Scopes:     scopes,
				OnlyCodes:  engfix.ParseOnlyCodesCSV(onlyCSV),
				MaxFiles:   maxFiles,
				MaxChanges: maxChanges,
				Aliases:    cfg.SectionAliases,
			})
			if err != nil {
				return err
			}

			if apply {
				if err := engfix.ApplyPlan(plan, engfix.ApplyOptions{
					Root:      cfg.DocsRoot,
					BackupDir: backupDir,
				}); err != nil {
					return err
				}
			}

			if output == "" {
				output = "-"
			}
			switch format {
			case "json":
				if err := writeFixJSON(plan, output); err != nil {
					return err
				}
			default:
				if err := writeFixDiff(plan, output); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "diff", "format: diff|json")
	cmd.Flags().StringVar(&output, "output", "-", "plik wyjściowy, '-' oznacza stdout")
	cmd.Flags().BoolVar(&apply, "apply", false, "zastosuj plan zmian atomowo (domyślnie dry-run)")
	cmd.Flags().StringVar(&backupDir, "backup-dir", "", "opcjonalny katalog backupu oryginalnych plików przy --apply")
	cmd.Flags().StringVar(&scopeCSV, "scope", "frontmatter,sections", "zakres fixów (CSV): frontmatter,sections")
	cmd.Flags().StringVar(&onlyCSV, "only", "", "filtr code (CSV), np. DOCFLOW.VALIDATE.MISSING_FIELD")
	cmd.Flags().IntVar(&maxFiles, "max-files", 200, "maksymalna liczba plików modyfikowanych w jednym runie")
	cmd.Flags().IntVar(&maxChanges, "max-changes", 1000, "maksymalna liczba zmian w jednym runie")
	cmd.Flags().StringVar(&againstPath, "against", "", "opcjonalny baseline validate JSON")
	cmd.Flags().StringVar(&show, "show", "all", "które findings naprawiać przy --against: all|new|existing")
	cmd.Flags().BoolVar(&statusAware, "status-aware", false, "użyj progressive validation zależnie od statusu dokumentu")
	cmd.Flags().StringVar(&governanceRules, "governance", "", "ścieżka do GOVERNANCE_RULES.yaml (opcjonalnie)")

	cmd.AddCommand(fixGuidedCmd())
	return cmd
}

func toEngineIssue(issue validateIssueJSON) engine.Issue {
	out := engine.Issue{
		Code:     issue.Code,
		Severity: engine.Severity(issue.Level),
		Type:     issue.Type,
		Kind:     engine.IssueKindValidate,
		Path:     issue.Path,
		DocID:    issue.DocID,
		Message:  issue.Message,
		Details:  issue.Details,
	}
	if issue.Line > 0 || issue.Column > 0 {
		out.Location = &engine.Location{
			Line:   issue.Line,
			Column: issue.Column,
		}
	}
	return out
}

func writeFixJSON(plan engfix.Plan, output string) error {
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	if output == "-" {
		_, err = os.Stdout.Write(append(data, '\n'))
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}
	return os.WriteFile(output, data, 0o644)
}

func writeFixDiff(plan engfix.Plan, output string) error {
	diff := engfix.RenderDiff(plan)
	if output == "-" {
		_, err := os.Stdout.Write([]byte(diff))
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}
	return os.WriteFile(output, []byte(diff), 0o644)
}
