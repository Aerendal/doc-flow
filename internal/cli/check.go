package cli

import (
	"github.com/spf13/cobra"
)

// checkCmd is a zero-flag shortcut for "health" that reads all defaults
// from docflow.yaml (or auto-discovered config). Perfect for local dev:
//
//	docflow check       # local mode (non-blocking)
//	docflow check --ci  # CI mode (blocks on new issues)
func checkCmd() *cobra.Command {
	var ciMode bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Skrót dla health z domyślnymi parametrami (config auto-detected)",
		Long: `check uruchamia health z domyślnymi parametrami:
  --baseline-mode repo
  --baseline-dir   .docflow/baseline
  --bundle-dir     .docflow/out
  --rules          docs/_meta/GOVERNANCE_RULES.yaml

Config jest wykrywany automatycznie (docflow.yaml w bieżącym lub nadrzędnym katalogu).
Flaga --ci przełącza w tryb CI (blokuje na nowych problemach).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := resolvedConfigPathForHealth()
			if cfgPath == "" {
				return usageError(
					"DOCFLOW.CLI.CONFIG_NOT_FOUND",
					"nie znaleziono docflow.yaml — uruchom z katalogu projektu lub podaj --config",
					nil,
				)
			}

			// Delegate to healthCmd with pre-set defaults
			healthArgs := []string{
				"--baseline-mode", "repo",
				"--baseline-dir", ".docflow/baseline",
				"--bundle-dir", ".docflow/out",
				"--rules", "docs/_meta/GOVERNANCE_RULES.yaml",
			}
			if ciMode {
				healthArgs = append(healthArgs, "--ci")
			}
			if configFile != "" {
				healthArgs = append(healthArgs, "--config", configFile)
			}

			// Re-use the root command execution
			root := cmd.Root()
			root.SetArgs(append([]string{"health"}, healthArgs...))
			return root.Execute()
		},
	}

	cmd.Flags().BoolVar(&ciMode, "ci", false, "tryb CI: blokuj na nowych problemach")
	return cmd
}
