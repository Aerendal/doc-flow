#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
LOG=${LOG:-LOGS/guided_fix_v2_smoke.md}

mkdir -p "$(dirname "$LOG")"

run(){ echo "$ $*" >> "$LOG"; "$@" >> "$LOG" 2>&1; }

{
  echo "# Guided fix v2 smoke"
  echo "Root: $ROOT"
  echo "Timestamp: $(date -Iseconds)"
} > "$LOG"

run ./scripts/guided_fix.sh
run ./scripts/guided_fix_map.sh
run ./scripts/guided_fix_patch.sh
run ./scripts/guided_fix_export.sh --since 20260212 --log LOGS/guided_fix_export.md --zip-name guided_fix_patches_filtered.zip

need=(
  "LOGS/guided_fix.md"
  "LOGS/guided_fix_map.md"
  "LOGS/guided_fix_patch.md"
  "LOGS/guided_fix_export.md"
  "dist/release/guided_fix_patches_filtered.zip"
)
missing=()
for f in "${need[@]}"; do
  [[ -e $f ]] || missing+=("$f")

done

if [[ ${#missing[@]} -gt 0 ]]; then
  {
    echo "Missing:"; printf ' - %s\n' "${missing[@]}"; echo "Result: FAIL";
  } >> "$LOG"
  exit 1
else
  echo "Result: PASS" >> "$LOG"
fi

echo "Smoke log: $LOG"
