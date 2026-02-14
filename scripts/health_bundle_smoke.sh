#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
LOG=${LOG:-LOGS/health_bundle_smoke.md}
BIN=${BIN:-./build/docflow-linux-amd64}

mkdir -p "$(dirname "$LOG")"

run(){ echo "$ $*" >> "$LOG"; "$@" >> "$LOG" 2>&1; }

{
  echo "# Health bundle smoke"
  echo "Root: $ROOT"
  echo "Bin: $BIN"
  echo "Timestamp: $(date -Iseconds)"
} > "$LOG"

run $BIN health --root "$ROOT" --bundle --report --pr-comment --log LOGS/health_onecmd_smoke.md

# sprawdź artefakty
need=(
  "dist/release/docflow-artifacts.zip"
  "dist/health/$(basename "$ROOT")/health_bundle_report.html"
  "LOGS/pr_comment_health_" "LOGS/health_onecmd_smoke.md"
)
missing=()
for f in "${need[@]}"; do
  if [[ "$f" == LOGS/pr_comment_health_* ]]; then
    latest=$(ls -1t LOGS/pr_comment_health_*.md 2>/dev/null | head -n1 || true)
    [[ -n "$latest" ]] || missing+=("PR_COMMENT")
    continue
  fi
  [[ -e $f ]] || missing+=("$f")

done

if [[ ${#missing[@]} -gt 0 ]]; then
  {
    echo "Missing:"; printf ' - %s\n' "${missing[@]}"; echo "Result: FAIL"
  } >> "$LOG"
  exit 1
else
  echo "Result: PASS" >> "$LOG"
fi

echo "Smoke log: $LOG"
