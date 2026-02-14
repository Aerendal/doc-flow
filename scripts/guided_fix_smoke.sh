#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
LOG=${LOG:-LOGS/guided_fix_smoke.md}
BIN=${BIN:-./build/docflow-linux-amd64}

mkdir -p "$(dirname "$LOG")"

run(){ echo "$ $*" >> "$LOG"; "$@" >> "$LOG" 2>&1; }

{
  echo "# Guided fix smoke"
  echo "Root: $ROOT"
  echo "Bin: $BIN"
  echo "Timestamp: $(date -Iseconds)"
} > "$LOG"

# guided_fix basic
run ./scripts/guided_fix.sh
# guided_fix context
run ./scripts/guided_fix_context.sh
# export
run ./scripts/guided_fix_export.sh

need=(
  "LOGS/guided_fix.md"
  "LOGS/guided_fix_context.md"
  "LOGS/guided_fix_export.md"
  "dist/release/guided_fix_patches.zip"
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
