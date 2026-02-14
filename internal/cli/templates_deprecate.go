package cli

import (
	"fmt"

	"docflow/pkg/templates"

	"github.com/spf13/cobra"
)

func templatesDeprecatedReportCmd() *cobra.Command {
	var minQuality int
	var unusedOnly bool
	cmd := &cobra.Command{
		Use:   "deprecated-report",
		Short: "Raport szablonów do deprecjacji (MVP demo)",
		RunE: func(cmd *cobra.Command, args []string) error {
			tpls := sampleTemplates()
			rules := templates.Rules{MinQuality: minQuality, UnusedOnly: unusedOnly}
			rep := templates.EvaluateDeprecation(tpls, rules)
			if len(rep.ToDeprecate) == 0 {
				fmt.Println("Brak szablonów do deprecjacji wg zasad.")
				return nil
			}
			fmt.Printf("%-30s %-8s %-6s %-6s\n", "PATH", "STATUS", "QUAL", "USAGE")
			for _, t := range rep.ToDeprecate {
				fmt.Printf("%-30s %-8s %-6d %-6d\n", t.Path, t.Status, t.Quality, t.Usage)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&minQuality, "min-quality", 50, "próg jakości do deprecjacji")
	cmd.Flags().BoolVar(&unusedOnly, "unused-only", false, "deprecjonuj tylko szablony nieużywane")
	return cmd
}
