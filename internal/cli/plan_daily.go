package cli

import (
	"fmt"
	"strings"

	"docflow/pkg/index"
	"docflow/pkg/plan"

	"github.com/spf13/cobra"
)

func planDailyCmd() *cobra.Command {
	var max int
	cmd := &cobra.Command{
		Use:   "daily",
		Short: "Wygeneruj dzienny plan dokumentów (topo + effort)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			log := newLogger(cfg)
			log.Info("Plan dzienny: root=%s", cfg.DocsRoot)

			idx, err := index.BuildIndex(cfg.DocsRoot, cfg.IgnorePatterns)
			if err != nil {
				return err
			}

			items, err := plan.PlanDaily(idx, max)
			if err != nil {
				return err
			}

			fmt.Printf("%-20s %-6s %s\n", "DOC_ID", "EFF", "PATH")
			fmt.Println(strings.Repeat("─", 80))
			for _, it := range items {
				path := it.Path
				if len(path) > 60 {
					path = "..." + path[len(path)-57:]
				}
				fmt.Printf("%-20s %-6d %s\n", it.DocID, it.Effort, path)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&max, "max", 5, "maksymalna liczba dokumentów na dzień (0 = wszystkie)")
	return cmd
}
