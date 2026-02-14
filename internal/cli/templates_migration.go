package cli

import (
	"fmt"

	"docflow/pkg/templates"

	"github.com/spf13/cobra"
)

func templatesMigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest-migration",
		Short: "Zaproponuj migracje dla szablonów deprecated (demo)",
		RunE: func(cmd *cobra.Command, args []string) error {
			deprecated := []int{}
			for i, t := range sampleTemplatesData {
				if t.Status == "deprecated" {
					deprecated = append(deprecated, i)
				}
			}
			if len(deprecated) == 0 {
				fmt.Println("Brak szablonów deprecated.")
				return nil
			}
			for _, idx := range deprecated {
				dep := sampleTemplatesData[idx]
				best := templates.SuggestMigration(dep, sampleTemplatesData)
				if best == nil {
					fmt.Printf("%s: brak migracji (brak aktywnych kandydatów)\n", dep.Path)
					continue
				}
				fmt.Printf("%s -> %s (quality %d, usage %d)\n", dep.Path, best.Path, best.Quality, best.Usage)
			}
			return nil
		},
	}
	return cmd
}
