#!/usr/bin/env bash
set -euo pipefail

LOG="LOGS/palette.md"

actions=(
  "Import DOCX (merge) | ./scripts/import_docx.sh --in /home/jerzy/Pobrane/LocalFold-Explorer_Projekt_v2_ASUS_TUF_F15.docx --out-dir tmp/localfold_import --merge-sections --export-patch LOGS/import_docx_palette.patch --status-json LOGS/import_docx_palette.json --section-diff LOGS/SECTION_DIFF_palette.md"
  "Audit root tmp/localfold_import | ./scripts/audit.sh --root tmp/localfold_import --log LOGS/AUDIT_PALETTE.md"
  "Bundle release | ./scripts/release_artifacts.sh --log LOGS/RELEASE_ARTIFACT_PALETTE.md"
  "Bundle PR | ./scripts/pr_bundle.sh --log LOGS/PR_BUNDLE_PALETTE.md"
  "Demo full | ./scripts/demo_guided_fix.sh"
  "Artifacts check | ./scripts/artifacts_check.sh --root localfold_import --strict --log LOGS/ARTIFACTS_CHECK_PALETTE.md"
  "Order guard | ./scripts/order_guard.sh --bundle dist/release/docflow-pr-bundle.zip --health dist/health/localfold_import/health_latest.html --queue dist/health/localfold_import/queue_latest.json --compliance dist/health/localfold_import/compliance_latest.json --log LOGS/ORDER_GUARD_PALETTE.md --no-strict"
  "Health+bundle+comment | ./build/docflow-linux-amd64 health --root tmp/localfold_import --bundle --report --pr-comment --log LOGS/health_onecmd_palette.md"
  "Guided fix | ROOT=tmp/localfold_import LOG=LOGS/guided_fix_palette.md ./scripts/guided_fix.sh"
)

mkdir -p "$(dirname "$LOG")"

select_action() {
  if command -v fzf >/dev/null 2>&1; then
    printf '%s\n' "${actions[@]}" | \
      fzf --prompt="docflow > " \
          --with-nth=1 \
          --preview='awk -F\"|\" \"{print \\\"Cmd: \\\"$2}\" <<< {}' \
          --preview-window=down,5 | \
      cut -d'|' -f2-
  else
    echo "fzf not found; fallback list:" | tee -a "$LOG" >&2
    nl -w2 -ba <(printf '%s\n' "${actions[@]}" | cut -d'|' -f1) | tee -a "$LOG" >&2
    nl -w2 -ba <(printf '%s\n' "${actions[@]}") | \
      sed 's/^[[:space:]]*//' | \
      awk -F'|' '{printf "[%s] %s\n    Cmd: %s\n", $1, $2, $3}' >> "$LOG"
    echo "Wybierz numer (Enter=1):" | tee -a "$LOG" >&2
    read -r idx || true
    [[ -z "$idx" ]] && idx=1
    printf '%s\n' "${actions[@]}" | sed -n "${idx}p" | cut -d'|' -f2-
  fi
}

cmd=$(select_action)

if [[ -z "$cmd" ]]; then
  echo "Brak wyboru." | tee "$LOG"
  exit 0
fi

{
  echo "# Palette run"
  echo "Command: $cmd"
} > "$LOG"

eval "$cmd" >> "$LOG" 2>&1 || true

echo "Palette finished. Log: $LOG"
