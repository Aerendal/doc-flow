#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
QUEUE_JSON=${QUEUE_JSON:-dist/health/$(basename "$ROOT")/queue_latest.json}
LOG=${LOG:-LOGS/guided_fix.md}
PATCH_DIR=${PATCH_DIR:-LOGS}

mkdir -p "$(dirname "$LOG")" "$PATCH_DIR"

echo "# Guided fix" > "$LOG"
echo "Root: $ROOT" >> "$LOG"
echo "Queue: $QUEUE_JSON" >> "$LOG"
echo "Timestamp: $(date -Iseconds)" >> "$LOG"

# zbierz zadania z queue i diff
tasks=()
if [[ -f "$QUEUE_JSON" ]]; then
  while IFS= read -r line; do
    tasks+=("QUEUE|$line")
  done < <(python3 - <<'PY' "$QUEUE_JSON"
import json,sys
q=json.load(open(sys.argv[1]))
tasks=q.get('tasks', [])
for t in tasks:
    status=t.get('status','').lower()
    title=t.get('title') or t.get('id') or t.get('path') or '(no title)'
    if status in ('ready','blocked'):
        print(f"{status}:{title}")
    else:
        print(f"other:{title}")
PY
)
fi
if [[ -f "${SECTION_DIFF:-LOGS/SECTION_DIFF_palette_bin.md}" ]]; then
  while IFS= read -r line; do
    tasks+=("DIFF|${line#'+ '}")
  done < <(grep '^+ ' "${SECTION_DIFF:-LOGS/SECTION_DIFF_palette_bin.md}" || true)
fi

if [[ ${#tasks[@]} -eq 0 ]]; then
  echo "Brak zadań (queue/diff)." >> "$LOG"
  exit 0
fi

select_task() {
  if command -v fzf >/dev/null 2>&1; then
    printf '%s\n' "${tasks[@]}" | fzf --prompt="guided fix > " --with-nth=1.. | head -n1
  else
    nl -w2 -ba <(printf '%s\n' "${tasks[@]}" | cut -d'|' -f2-) >> "$LOG"
    echo "Wybierz numer (Enter=1):" >> "$LOG"
    read -r idx
    [[ -z "$idx" ]] && idx=1
    printf '%s\n' "${tasks[@]}" | sed -n "${idx}p"
  fi
}

selected_task=$(select_task)
source_type=${selected_task%%|*}
task_body=${selected_task#*|}
echo "Selected task: $source_type | $task_body" >> "$LOG"

# heurystyka pliku: pierwszy md
files=("$ROOT"/*.md)
if [[ ${#files[@]} -eq 0 ]]; then
  echo "Brak plików md w $ROOT" >> "$LOG"
  exit 0
fi
selected="${files[0]}"

# generuj patch szkic (non-destructive)
ts=$(date +%Y%m%dT%H%M%S)
patch="$PATCH_DIR/patch_guided_${ts}.patch"
cat > "$patch" <<EOFpatch
*** Begin Patch
*** Update File: $(basename "$selected")
@@
+<!-- TODO guided-fix placeholder: opisz wymagane zmiany sekcji -->
*** End Patch
EOFpatch

echo "Selected file: $selected" >> "$LOG"
echo "Patch: $patch" >> "$LOG"
echo "Status: OK (non-destructive)" >> "$LOG"
