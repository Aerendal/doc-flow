package cli

import (
	"fmt"

	"docflow/pkg/index"

	"github.com/spf13/cobra"
)

func changesCmd() *cobra.Command {
	var oldIndex string
	cmd := &cobra.Command{
		Use:   "changes",
		Short: "Pokaż zmienione/dodane/usunięte dokumenty względem poprzedniego indeksu",
		RunE: func(cmd *cobra.Command, args []string) error {
			if oldIndex == "" {
				return fmt.Errorf("podaj --old-index (JSON)")
			}
			oldIdx, err := index.LoadJSON(oldIndex)
			if err != nil {
				return err
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			newIdx, err := index.BuildIndex(cfg.DocsRoot, cfg.IgnorePatterns)
			if err != nil {
				return err
			}
			changes := index.DocDiff(oldIdx, newIdx)
			if len(changes) == 0 {
				fmt.Println("Brak zmian.")
				return nil
			}
			fmt.Printf("%-8s %-20s %s\n", "TYPE", "DOC_ID", "PATH")
			for _, c := range changes {
				fmt.Printf("%-8s %-20s %s\n", c.Type, c.DocID, c.Path)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&oldIndex, "old-index", "", "ścieżka do starego indeksu JSON")
	return cmd
}
