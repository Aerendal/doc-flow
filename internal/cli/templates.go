package cli

import (
	"fmt"
	"os"
	"strings"

	"docflow/internal/util"
	"docflow/pkg/pattern"
	"docflow/pkg/quality"
	"docflow/pkg/schema"

	"github.com/spf13/cobra"
)

func validateTemplatesCmd() *cobra.Command {
	var schemaPath string
	var dir string
	cmd := &cobra.Command{
		Use:   "validate-templates",
		Short: "Waliduj szablony względem schematu sekcji i oblicz score",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			root := cfg.DocsRoot
			if dir != "" {
				root = dir
			}

			var sch *schema.SectionSchema
			if schemaPath != "" {
				sch, err = schema.LoadSectionSchema(schemaPath)
				if err != nil {
					return err
				}
			}

			files, err := util.WalkMarkdown(root, cfg.IgnorePatterns)
			if err != nil {
				return err
			}

			fmt.Printf("%-50s %-8s %-8s %s\n", "SZABLON", "SCORE", "MISS", "MISSING_SECTIONS")
			fmt.Println(strings.Repeat("─", 100))

			for _, f := range files {
				data, err := os.ReadFile(f.Path)
				if err != nil {
					continue
				}
				s := quality.ScoreContent(string(data), sch)
				missingStr := "-"
				if len(s.Missing) > 0 {
					missingStr = strings.Join(s.Missing, ", ")
				}
				path := f.RelPath
				if len(path) > 48 {
					path = "..." + path[len(path)-45:]
				}
				fmt.Printf("%-50s %-8d %-8d %s\n", path, s.Total, len(s.Missing), missingStr)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&schemaPath, "schema", "", "ścieżka do section_schema.yaml")
	cmd.Flags().StringVar(&dir, "dir", "", "katalog z szablonami (domyślnie docs_root)")
	return cmd
}

func findDuplicatesCmd() *cobra.Command {
	var threshold float64
	cmd := &cobra.Command{
		Use:   "find-duplicates",
		Short: "Wykryj podobne szablony na podstawie struktury sekcji",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			patterns, err := pattern.ScanPatterns(cfg.DocsRoot, cfg.IgnorePatterns)
			if err != nil {
				return err
			}

			pairs := pattern.FindSimilar(patterns, threshold)
			if len(pairs) == 0 {
				fmt.Println("Brak duplikatów powyżej progu.")
				return nil
			}

			fmt.Printf("%-40s %-40s %s\n", "SZABLON A", "SZABLON B", "SIM")
			fmt.Println(strings.Repeat("─", 100))
			for _, p := range pairs {
				fmt.Printf("%-40s %-40s %.2f\n", p.APath, p.BPath, p.Score)
			}
			return nil
		},
	}
	cmd.Flags().Float64Var(&threshold, "threshold", 0.85, "próg podobieństwa (0-1)")
	return cmd
}
