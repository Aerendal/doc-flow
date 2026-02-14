package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"docflow/internal/engine"
	enggraph "docflow/internal/engine/graph"

	"github.com/spf13/cobra"
)

func graphCmd() *cobra.Command {
	var format string
	var output string
	var cyclesOnly bool
	var includeLinks bool
	var includeTemplates bool
	var portable bool

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Buduj kanoniczny graf zależności dokumentów",
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "text" && format != "json" {
				return fmt.Errorf("nieobsługiwany format: %s (dozwolone: text|json)", format)
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			execRes, err := collectExecutionResult("graph", cfg)
			if err != nil {
				return fmt.Errorf("nie można przygotować facts execution layer: %w", err)
			}

			g, issues := enggraph.BuildGraph(execRes.Facts, enggraph.BuildGraphOptions{
				IncludeLinks:     includeLinks,
				IncludeTemplates: includeTemplates,
				Deterministic:    true,
				Portable:         portable,
			})

			switch format {
			case "json":
				if output == "" {
					output = "-"
				}
				if err := writeGraphJSON(g, output); err != nil {
					return err
				}
				if output != "-" {
					fmt.Printf("Graf zapisany do %s\n", output)
				}
			default:
				printGraphText(g, issues, cyclesOnly)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "text", "format: text|json")
	cmd.Flags().StringVar(&output, "output", "-", "plik wyjściowy (json), '-' = stdout")
	cmd.Flags().BoolVar(&cyclesOnly, "cycles", false, "wypisz tylko cykle")
	cmd.Flags().BoolVar(&includeLinks, "include-links", false, "uwzględnij linki markdown jako krawędzie links_to (kosztowne na dużych repo; zalecane z cache+incremental)")
	cmd.Flags().BoolVar(&includeTemplates, "include-templates", true, "uwzględnij relacje uses_template")
	cmd.Flags().BoolVar(&portable, "portable", true, "ustaw root='.' w JSON dla przenośnego outputu")
	return cmd
}

func writeGraphJSON(g enggraph.Graph, output string) error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("błąd serializacji grafu: %w", err)
	}
	if output == "" || output == "-" {
		_, err = os.Stdout.Write(append(data, '\n'))
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("błąd utworzenia katalogu dla grafu: %w", err)
	}
	return os.WriteFile(output, data, 0o644)
}

func printGraphText(g enggraph.Graph, issues []engine.Issue, cyclesOnly bool) {
	if cyclesOnly {
		printGraphCycles(g)
		return
	}

	fmt.Printf("schema=%s root=%s\n", g.SchemaVersion, g.Root)
	fmt.Printf("nodes=%d edges=%d cycles=%d max_out=%d max_in=%d\n",
		g.Stats.NodeCount, g.Stats.EdgeCount, g.Stats.CyclesCount, g.Stats.MaxOutDegree, g.Stats.MaxInDegree)
	printGraphCycles(g)

	if len(issues) > 0 {
		fmt.Println("\nGraph issues:")
		for _, is := range issues {
			fmt.Printf("- [%s] %s (%s)\n", is.Code, is.Message, is.Path)
		}
	}
}

func printGraphCycles(g enggraph.Graph) {
	if len(g.Cycles) == 0 {
		fmt.Println("Brak cykli.")
		return
	}
	fmt.Println("Cycles:")
	for _, cyc := range g.Cycles {
		fmt.Printf("- %s\n", strings.Join(cyc, " -> "))
	}
}
