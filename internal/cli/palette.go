package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func paletteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "palette",
		Short: "Command palette (tekstowa) dla kluczowych akcji",
		RunE: func(cmd *cobra.Command, args []string) error {
			actions := []struct {
				label string
				cmd   string
			}{
				{"Import DOCX (merge)", "./scripts/import_docx.sh --in /home/jerzy/Pobrane/LocalFold-Explorer_Projekt_v2_ASUS_TUF_F15.docx --out-dir tmp/localfold_import --merge-sections --export-patch LOGS/import_docx_palette_bin.patch --status-json LOGS/import_docx_palette_bin.json --section-diff LOGS/SECTION_DIFF_palette_bin.md"},
				{"Audit root tmp/localfold_import", "./scripts/audit.sh --root tmp/localfold_import --log LOGS/AUDIT_PALETTE_BIN.md"},
				{"Bundle release", "./scripts/release_artifacts.sh --log LOGS/RELEASE_ARTIFACT_PALETTE_BIN.md"},
				{"Bundle PR", "./scripts/pr_bundle.sh --log LOGS/PR_BUNDLE_PALETTE_BIN.md"},
				{"Demo full", "./scripts/demo_guided_fix.sh"},
				{"Artifacts check", "./scripts/artifacts_check.sh --root localfold_import --strict --log LOGS/ARTIFACTS_CHECK_PALETTE_BIN.md"},
				{"Order guard", "./scripts/order_guard.sh --bundle dist/release/docflow-pr-bundle.zip --health dist/health/localfold_import/health_latest.html --queue dist/health/localfold_import/queue_latest.json --compliance dist/health/localfold_import/compliance_latest.json --log LOGS/ORDER_GUARD_PALETTE_BIN.md --no-strict"},
				{"Health+bundle+comment", "./build/docflow-linux-amd64 health --root tmp/localfold_import --bundle --report --pr-comment --log LOGS/health_onecmd_palette_bin.md"},
				{"Guided fix", "ROOT=tmp/localfold_import LOG=LOGS/guided_fix_palette_bin.md ./scripts/guided_fix.sh"},
			}
			fmt.Println("Command palette (bez fzf). Wybierz numer:")
			for i, a := range actions {
				fmt.Printf("%d) %s\n", i+1, a.label)
			}
			var idx int
			if _, err := fmt.Scanln(&idx); err != nil {
				return nil
			}
			if idx < 1 || idx > len(actions) {
				fmt.Println("Brak wyboru.")
				return nil
			}
			selected := actions[idx-1]
			fmt.Printf("Wybrano: %s\nCmd: %s\nUruchomić? [Y/n]: ", selected.label, selected.cmd)
			var confirm string
			_, _ = fmt.Scanln(&confirm)
			if confirm != "" && (confirm[0] == 'n' || confirm[0] == 'N') {
				fmt.Println("Przerwano na życzenie użytkownika.")
				return nil
			}
			fmt.Printf("Uruchamiam: %s\n", selected.cmd)
			return executeShell(selected.cmd)
		},
	}
	return cmd
}

func executeShell(command string) error {
	return runBash(command)
}

func runBash(command string) error {
	sh := "/bin/bash"
	if envShell := os.Getenv("SHELL"); envShell != "" {
		sh = envShell
	}
	c := fmt.Sprintf("set -euo pipefail\n%s", command)
	return execCommand(sh, "-c", c)
}
