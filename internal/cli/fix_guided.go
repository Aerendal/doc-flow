package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func fixGuidedCmd() *cobra.Command {
	var root string
	var logFile string
	cmd := &cobra.Command{
		Use:   "guided",
		Short: "Tryb guided fix (interaktywny / non-destructive)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if root == "" {
				return fmt.Errorf("--root jest wymagane")
			}
			if logFile == "" {
				logFile = "LOGS/GUIDED_FIX_BIN_DAY_232.md"
			}
			guidedCmd := fmt.Sprintf("ROOT=%s LOG=%s ./scripts/guided_fix.sh", root, logFile)
			return runShell(guidedCmd)
		},
	}
	cmd.Flags().StringVar(&root, "root", "", "katalog z docflow.yaml")
	cmd.Flags().StringVar(&logFile, "log", "LOGS/GUIDED_FIX_BIN_DAY_232.md", "plik logu guided fix")
	cmd.MarkFlagRequired("root")
	return cmd
}
