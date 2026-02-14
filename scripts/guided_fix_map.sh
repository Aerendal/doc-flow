#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
QUEUE_JSON=${QUEUE_JSON:-dist/health/$(basename "$ROOT")/queue_latest.json}
SECTION_DIFF=${SECTION_DIFF:-LOGS/SECTION_DIFF_palette_bin.md}
LOG=${LOG:-LOGS/guided_fix_map.md}

mkdir -p "$(dirname "$LOG")"

lines=()
if [[ -f "$SECTION_DIFF" ]]; then
  while IFS= read -r l; do
    lines+=("DIFF|${l#'+ '}")
  done < <(grep '^+ ' "$SECTION_DIFF" || true)
fi
if [[ -f "$QUEUE_JSON" ]]; then
  while IFS= read -r l; do
    lines+=("QUEUE|$l")
  done < <(python3 - <<'PY' "$QUEUE_JSON"
import json,sys
q=json.load(open(sys.argv[1]))
for t in q.get('tasks', []):
    title=t.get('title') or t.get('id') or t.get('path') or '(no title)'
    print(title)
PY
)
fi

mapfile -t files < <(ls "$ROOT"/*.md 2>/dev/null || true)

resolve(){
  local needle="$1"; shift
  for f in "${files[@]}"; do
    if grep -iq -- "${needle}" "$f"; then
      echo "$f"; return 0
    fi
  done
  echo "(none)"; return 1
}

{
  echo "# Guided fix map"
  echo "Root: $ROOT"
  echo "Files: ${#files[@]}"
  echo "Timestamp: $(date -Iseconds)"
  echo
  for entry in "${lines[@]}"; do
    src=${entry%%|*}; body=${entry#*|}
    target=$(resolve "$body" || true)
    printf -- "- %s | %s -> %s\n" "$src" "$body" "$target"
  done
} > "$LOG"

echo "Mapping written to $LOG"
