package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"docflow/internal/util"
	"docflow/pkg/sections"

	"github.com/spf13/cobra"
)

func migrateSectionsCmd() *cobra.Command {
	var apply bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "migrate-sections",
		Short: "Zamień legacy nazwy sekcji na kanoniczne (wg section_aliases)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if len(cfg.SectionAliases) == 0 {
				return fmt.Errorf("brak section_aliases w configu")
			}

			files, err := util.WalkMarkdown(cfg.DocsRoot, cfg.IgnorePatterns)
			if err != nil {
				return err
			}

			// dry-run alias: jeśli podano --dry-run, wymuszamy podgląd (apply=false).
			if dryRun {
				apply = false
			}

			totalHits := 0
			for _, f := range files {
				data, err := os.ReadFile(f.Path)
				if err != nil {
					continue
				}
				newContent, hits := sections.RenameAliases(string(data), cfg.SectionAliases)
				if len(hits) == 0 {
					continue
				}
				totalHits += len(hits)
				rel := filepath.ToSlash(f.RelPath)
				fmt.Printf("%s\n", rel)
				for _, h := range hits {
					fmt.Printf("  %s -> %s\n", h.From, h.To)
				}
				if apply {
					if err := os.WriteFile(f.Path, []byte(newContent), 0o644); err != nil {
						return fmt.Errorf("write %s: %w", rel, err)
					}
				}
			}
			if totalHits == 0 {
				fmt.Println("Brak legacy sekcji do migracji.")
			} else if !apply {
				fmt.Println("\nTryb dry-run (dodaj --apply aby zapisać zmiany).")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&apply, "apply", false, "zapisz zmiany (domyślnie tylko podgląd)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "alias: wymuś podgląd (bez zapisu)")
	return cmd
}
