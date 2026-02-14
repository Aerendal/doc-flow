package cli

import (
	"fmt"

	"docflow/pkg/recommend"

	"github.com/spf13/cobra"
)

func recommendCmd() *cobra.Command {
	var docType string
	var lang string
	var topN int
	cmd := &cobra.Command{
		Use:   "recommend",
		Short: "Rekomenduj szablony na podstawie typu i języka",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Placeholder: demo z kilkoma szablonami; brak indeksu szablonów w tej iteracji
			templates := []recommend.TemplateInfo{
				{Path: "templates/guide.md", DocType: "guide", Language: "pl", Quality: 85, Usage: 3},
				{Path: "templates/spec.md", DocType: "specification", Language: "pl", Quality: 90, Usage: 5},
				{Path: "templates/guide-en.md", DocType: "guide", Language: "en", Quality: 80, Usage: 1},
			}
			params := recommend.Params{
				PreferredDocType: docType,
				PreferredLang:    lang,
				TopN:             topN,
				Weights:          recommend.Weights{DocType: 0.4, Lang: 0.2, Quality: 0.25, Usage: 0.15},
			}
			recs := recommend.Recommend(templates, params)

			if len(recs) == 0 {
				fmt.Println("Brak rekomendacji (brak szablonów spełniających kryteria).")
				return nil
			}

			fmt.Printf("%-40s %-12s %-6s %-6s %-6s %-6s %s\n", "PATH", "DOC_TYPE", "LANG", "QUAL", "USAGE", "SCORE", "REASON")
			for _, r := range recs {
				fmt.Printf("%-40s %-12s %-6s %-6d %-6d %-6.1f %s\n",
					r.Template.Path, r.Template.DocType, r.Template.Language,
					r.Template.Quality, r.Template.Usage, r.Score, r.Reason)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&docType, "doc-type", "", "preferowany doc_type")
	cmd.Flags().StringVar(&lang, "lang", "pl", "preferowany język")
	cmd.Flags().IntVar(&topN, "top", 5, "liczba rekomendacji")
	return cmd
}
