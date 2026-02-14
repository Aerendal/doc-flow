package cli

import (
	"fmt"
	"sort"

	enggraph "docflow/internal/engine/graph"
	"docflow/pkg/index"

	"github.com/spf13/cobra"
)

func templateImpactCmd() *cobra.Command {
	var oldIndex string
	cmd := &cobra.Command{
		Use:   "template-impact",
		Short: "Pokaż dokumenty zależne od zmienionych szablonów",
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
			execRes, err := collectExecutionResult("template-impact", cfg)
			if err != nil {
				return err
			}
			newIdx := execRes.Index

			changed := index.ChangedTemplates(oldIdx, newIdx)
			if len(changed) == 0 {
				fmt.Println("Brak zmian w szablonach.")
				return nil
			}

			g, _ := enggraph.BuildGraph(execRes.Facts, enggraph.BuildGraphOptions{
				IncludeTemplates: true,
				Deterministic:    true,
				Portable:         true,
			})

			templateIDs := make([]string, 0, len(changed))
			for tmplID := range changed {
				templateIDs = append(templateIDs, tmplID)
			}
			sort.Strings(templateIDs)

			fmt.Printf("Zmienione szablony: %d\n", len(changed))
			fmt.Printf("%-20s %s\n", "TEMPLATE", "AFFECTED DOCS")
			for _, tmpl := range templateIDs {
				startIDs := enggraph.NodeIDsForDocID(g, tmpl)
				docs := enggraph.ReachableDocPaths(g, startIDs, enggraph.ImpactOptions{
					MaxDepth:     1,
					IncludeStart: false,
					EdgeKinds: map[enggraph.EdgeKind]bool{
						enggraph.EdgeUsesTemplate: true,
					},
				})
				fmt.Printf("%-20s %v\n", tmpl, docs)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&oldIndex, "old-index", "", "ścieżka do starego indeksu JSON")
	return cmd
}
