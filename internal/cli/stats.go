package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"docflow/pkg/index"
	"docflow/pkg/schema"
	"docflow/pkg/sections"

	"github.com/spf13/cobra"
)

func statsCmd() *cobra.Command {
	var schemaPath string
	var byStatus bool
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Pokaż metryki kompletności sekcji",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			var sch *schema.SectionSchema
			if schemaPath != "" {
				sch, err = schema.LoadSectionSchema(schemaPath)
				if err != nil {
					return err
				}
			}

			idx, err := index.BuildIndex(cfg.DocsRoot, cfg.IgnorePatterns)
			if err != nil {
				return err
			}

			type agg struct {
				count    int
				sumComp  float64
				sumEmpty int
			}
			groupBy := func(rec *index.DocumentRecord) string {
				if byStatus {
					return rec.Status
				}
				return rec.DocType
			}

			stats := make(map[string]*agg)

			for _, rec := range idx.All() {
				path := filepath.Join(cfg.DocsRoot, rec.Path)
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				res := sections.ComputeCompleteness(string(data), sch, nil)
				key := groupBy(rec)
				a := stats[key]
				if a == nil {
					a = &agg{}
					stats[key] = a
				}
				a.count++
				a.sumComp += res.Completeness
				a.sumEmpty += res.Empty
			}

			fmt.Printf("%-20s %-6s %-8s %-10s\n", "GROUP", "DOCS", "EMPTY", "AVG_COMP")
			fmt.Println(strings.Repeat("─", 50))
			for k, v := range stats {
				avgComp := 0.0
				if v.count > 0 {
					avgComp = (v.sumComp / float64(v.count)) * 100
				}
				fmt.Printf("%-20s %-6d %-8.1f %-9.1f%%\n", k, v.count, float64(v.sumEmpty)/float64(max(v.count, 1)), avgComp)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&schemaPath, "schema", "", "plik YAML sekcji (opcjonalnie)")
	cmd.Flags().BoolVar(&byStatus, "by-status", false, "grupuj po statusie zamiast doc_type")
	return cmd
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
