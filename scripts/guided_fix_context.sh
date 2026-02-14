#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
QUEUE_JSON=${QUEUE_JSON:-dist/health/$(basename "$ROOT")/queue_latest.json}
SECTION_DIFF=${SECTION_DIFF:-LOGS/SECTION_DIFF_palette_bin.md}
LOG=${LOG:-LOGS/guided_fix_context.md}
PATCH_DIR=${PATCH_DIR:-LOGS}

mkdir -p "$(dirname "$LOG")" "$PATCH_DIR"

echo "# Guided fix (context)" > "$LOG"
echo "Root: $ROOT" >> "$LOG"
echo "Queue: $QUEUE_JSON" >> "$LOG"
echo "Section diff: $SECTION_DIFF" >> "$LOG"
echo "Timestamp: $(date -Iseconds)" >> "$LOG"

# heurystyka: wybierz pierwszą linię dodaną z SECTION_DIFF
context_line=$(grep '^\+' "$SECTION_DIFF" | head -n1 | sed 's/^+//') || true
file_candidate=$(printf "%s\n" "$ROOT"/*.md | head -n1)
context_source="diff"

# fallback: pierwsza niepusta linia po frontmatterze, jeśli diff nie daje kontekstu
if [[ -z "$context_line" && -n "$file_candidate" && -f "$file_candidate" ]]; then
  context_line=$(python3 - <<'PY' "$file_candidate"
import sys, pathlib
text=pathlib.Path(sys.argv[1]).read_text(errors='ignore').splitlines()
front=False; found=False
for line in text:
    if line.strip() == '---':
        front = not front
        continue
    if front:
        continue
    if line.strip():
        print(line.strip())
        found=True
        break
if not found:
    print('')
PY
)
  context_source="file-fallback"
  echo "Kontekst (fallback z pliku): $context_line" >> "$LOG"
fi

if [[ -z "$context_line" ]]; then
  context_source="none"
  echo "Brak kontekstu z section diff i pliku; patch będzie z pustą adnotacją" >> "$LOG"
fi

snippet="(brak)"
if [[ -n "$file_candidate" && -f "$file_candidate" ]]; then
  snippet=$(python3 - <<'PY' "$file_candidate" "$context_line"
import sys, pathlib
path=sys.argv[1]; needle=sys.argv[2]
text=pathlib.Path(path).read_text(errors='ignore').splitlines()
if needle:
    for i,line in enumerate(text):
        if needle.lower() in line.lower():
            start=max(0,i-2); end=min(len(text), i+3)
            for j in range(start,end):
                print(f"{j+1}: {text[j]}")
            break
    else:
        for j,line in enumerate(text[:5]):
            print(f"{j+1}: {line}")
else:
    for j,line in enumerate(text[:5]):
        print(f"{j+1}: {line}")
PY
)
fi

ts=$(date +%Y%m%dT%H%M%S)
patch="$PATCH_DIR/patch_guided_${ts}.patch"
cat > "$patch" <<EOFpatch
*** Begin Patch
*** Update File: $(basename "$file_candidate")
@@
+<!-- TODO guided-fix context: $context_line -->
*** End Patch
EOFpatch

echo "Selected file: $file_candidate" >> "$LOG"
echo "Context source: $context_source" >> "$LOG"
echo "Context line: $context_line" >> "$LOG"
echo "Snippet:" >> "$LOG"
echo "$snippet" >> "$LOG"
echo "Patch: $patch" >> "$LOG"
echo "Status: OK (context heuristic)" >> "$LOG"
