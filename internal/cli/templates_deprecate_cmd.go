package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func templatesDeprecateCmd() *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "deprecate",
		Short: "Oznacz szablon jako deprecated (demo, in-memory)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("podaj --id (path)")
			}
			found := false
			for i := range sampleTemplatesData {
				if sampleTemplatesData[i].Path == id {
					sampleTemplatesData[i].Status = "deprecated"
					found = true
				}
			}
			if !found {
				return fmt.Errorf("nie znaleziono szablonu %s", id)
			}
			fmt.Printf("Oznaczono %s jako deprecated\n", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "ścieżka szablonu do oznaczenia")
	return cmd
}
