package cli

import (
	"fmt"

	"docflow/pkg/recommend"

	"github.com/spf13/cobra"
)

var sampleTemplatesData = []recommend.TemplateInfo{
	{Path: "templates/spec_v1.md", DocType: "specification", Language: "pl", Quality: 90, Usage: 5, Version: "v1", Status: "active", CodeBlocks: 3, Tables: 1, Links: 2, Images: 0},
	{Path: "templates/spec_v2.md", DocType: "specification", Language: "pl", Quality: 80, Usage: 1, Version: "v2", Status: "deprecated", CodeBlocks: 4, Tables: 1, Links: 1, Images: 0},
	{Path: "templates/guide_v1.md", DocType: "guide", Language: "pl", Quality: 85, Usage: 4, Version: "v1", Status: "active", CodeBlocks: 1, Tables: 0, Links: 3, Images: 1},
}

func templatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Lista szablonów (MVP demo)",
	}
	cmd.AddCommand(templatesListCmd())
	cmd.AddCommand(templatesDeprecatedCmd())
	cmd.AddCommand(templatesDeprecatedReportCmd())
	cmd.AddCommand(templatesDeprecateCmd())
	cmd.AddCommand(templatesMigrationCmd())
	return cmd
}

func sampleTemplates() []recommend.TemplateInfo {
	return sampleTemplatesData
}

func templatesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Wyświetl wersje i status szablonów",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%-30s %-14s %-6s %-8s %-10s %-6s %-4s %-4s %-4s %-4s\n", "PATH", "DOC_TYPE", "VER", "STATUS", "QUALITY", "USAGE", "CODE", "TAB", "LINK", "IMG")
			for _, t := range sampleTemplates() {
				fmt.Printf("%-30s %-14s %-6s %-8s %-10d %-6d %-4d %-4d %-4d %-4d\n",
					t.Path, t.DocType, t.Version, t.Status, t.Quality, t.Usage, t.CodeBlocks, t.Tables, t.Links, t.Images)
			}
			return nil
		},
	}
}

func templatesDeprecatedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deprecated",
		Short: "Pokaż szablony deprecated/archived",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%-30s %-6s %-8s\n", "PATH", "VER", "STATUS")
			for _, t := range sampleTemplates() {
				if t.Status == "deprecated" || t.Status == "archived" || t.Deprecated {
					fmt.Printf("%-30s %-6s %-8s\n", t.Path, t.Version, t.Status)
				}
			}
			return nil
		},
	}
}
