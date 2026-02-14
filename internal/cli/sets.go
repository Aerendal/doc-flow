package cli

import (
	"fmt"

	"docflow/pkg/sets"

	"github.com/spf13/cobra"
)

func setsCmd() *cobra.Command {
	var minCount int
	cmd := &cobra.Command{
		Use:   "template-sets",
		Short: "Pokaż współwystępujące szablony (zestawy)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// placeholder sample data
			usages := []sets.TemplateUsage{
				{SetID: "proj1", Template: "spec.md"},
				{SetID: "proj1", Template: "guide.md"},
				{SetID: "proj2", Template: "spec.md"},
				{SetID: "proj2", Template: "runbook.md"},
				{SetID: "proj3", Template: "spec.md"},
				{SetID: "proj3", Template: "guide.md"},
			}
			co := sets.CoOccurrence(usages, minCount)
			if len(co) == 0 {
				fmt.Println("Brak zestawów powyżej progu.")
				return nil
			}
			fmt.Printf("%-20s %s\n", "TEMPLATE", "CO-OCCURS WITH")
			for t, list := range co {
				fmt.Printf("%-20s %v\n", t, list)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&minCount, "min-count", 2, "minimalna liczba współwystąpień")
	return cmd
}
