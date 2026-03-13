package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"docflow/internal/buildinfo"
	"docflow/pkg/config"

	"github.com/spf13/cobra"
)

type doctorCheck struct {
	label  string
	ok     bool
	detail string
}

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnostyka środowiska i konfiguracji projektu",
		Long:  `doctor sprawdza środowisko i konfigurację projektu, wykrywając typowe problemy.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			var checks []doctorCheck

			// 1. Config file
			cfgPath := configFile
			if cfgPath == "" {
				cfgPath = config.FindConfigFile()
			}
			if cfgPath != "" {
				checks = append(checks, doctorCheck{"config (docflow.yaml)", true, cfgPath})
				// Load and validate config
				cfg, err := config.Load(cfgPath)
				if err != nil {
					checks = append(checks, doctorCheck{"config valid", false, err.Error()})
				} else {
					if err := cfg.Validate(); err != nil {
						checks = append(checks, doctorCheck{"config valid", false, err.Error()})
					} else {
						checks = append(checks, doctorCheck{"config valid", true, ""})
					}
					// docs_root exists?
					if _, err := os.Stat(cfg.DocsRoot); err != nil {
						checks = append(checks, doctorCheck{"docs_root istnieje", false, cfg.DocsRoot + " — nie znaleziono"})
					} else {
						checks = append(checks, doctorCheck{"docs_root istnieje", true, cfg.DocsRoot})
					}
				}
			} else {
				checks = append(checks, doctorCheck{"config (docflow.yaml)", false, "nie znaleziono — uruchom z katalogu projektu"})
			}

			// 2. Governance rules
			rulesPath := "docs/_meta/GOVERNANCE_RULES.yaml"
			if _, err := os.Stat(rulesPath); err == nil {
				checks = append(checks, doctorCheck{"rules (GOVERNANCE_RULES.yaml)", true, rulesPath})
			} else {
				checks = append(checks, doctorCheck{"rules (GOVERNANCE_RULES.yaml)", false, rulesPath + " — nie znaleziono (opcjonalne)"})
			}

			// 3. Baseline files
			baselineValidate := filepath.Join(".docflow", "baseline", "validate.json")
			baselineCompliance := filepath.Join(".docflow", "baseline", "compliance.json")
			for _, bp := range []string{baselineValidate, baselineCompliance} {
				if _, err := os.Stat(bp); err == nil {
					checks = append(checks, doctorCheck{"baseline " + filepath.Base(bp), true, bp})
				} else {
					checks = append(checks, doctorCheck{"baseline " + filepath.Base(bp), false, bp + " — brak (uruchom health raz, aby wygenerować)"})
				}
			}

			// 4. Go toolchain
			goVersion := runtime.Version()
			checks = append(checks, doctorCheck{"Go toolchain", true, goVersion})

			// 5. Git available
			if path, err := exec.LookPath("git"); err == nil {
				out, _ := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output() //nolint:gosec
				inRepo := strings.TrimSpace(string(out)) == "true"
				if inRepo {
					checks = append(checks, doctorCheck{"git repo", true, path})
				} else {
					checks = append(checks, doctorCheck{"git repo", false, "bieżący katalog nie jest w repo git"})
				}
			} else {
				checks = append(checks, doctorCheck{"git dostępny", false, "git nie znaleziony w PATH"})
			}

			// 6. vendor/ in sync
			if _, err := os.Stat("vendor"); err == nil {
				if _, err := os.Stat("go.mod"); err == nil {
					// Simple check: vendor/modules.txt exists
					if _, err := os.Stat(filepath.Join("vendor", "modules.txt")); err == nil {
						checks = append(checks, doctorCheck{"vendor/ zsynchronizowany", true, "vendor/modules.txt obecny"})
					} else {
						checks = append(checks, doctorCheck{"vendor/ zsynchronizowany", false, "brak vendor/modules.txt — uruchom: go mod vendor"})
					}
				}
			}

			// 7. Binary version
			checks = append(checks, doctorCheck{"docflow version", true, buildinfo.FullVersion()})

			// Print results
			fmt.Println(bold("docflow doctor"))
			fmt.Println(strings.Repeat("─", 50))

			allOK := true
			for _, c := range checks {
				if c.ok {
					fmt.Printf("  %s  %s", passLabel(), c.label)
					if c.detail != "" {
						fmt.Printf("  %s", cyan(c.detail))
					}
				} else {
					allOK = false
					fmt.Printf("  %s  %s", failLabel() , c.label)
					if c.detail != "" {
						fmt.Printf("  %s", yellow(c.detail))
					}
				}
				fmt.Println()
			}

			fmt.Println(strings.Repeat("─", 50))
			if allOK {
				fmt.Printf("  %s  Środowisko jest gotowe do pracy\n", passLabel())
			} else {
				fmt.Printf("  %s  Znaleziono problemy — sprawdź powyżej\n", warnLabel())
			}

			return nil
		},
	}
}
