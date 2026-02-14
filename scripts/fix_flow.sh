#!/usr/bin/env bash
set -euo pipefail

ROOT="tmp/localfold_import"
OUT_DIR_SUFFIX="_todo"
OUT_TODO=""
LOG="LOGS/fix_flow.md"
EXPORT_PATCH=0
PATCH_PATH="/tmp/fix_flow_hints.patch"
CLEAN=0
EXPORT_SECTIONS_PATCH=0
SECTIONS_PATCH_PATH="/tmp/fix_flow_sections.patch"

usage() {
  echo "usage: $0 [--root dir] [--out-dir suffix] [--log file] [--export-patch [path]]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --out-dir) OUT_DIR_SUFFIX="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    --export-patch)
      EXPORT_PATCH=1
      if [[ $# -ge 2 && "$2" != --* ]]; then
        PATCH_PATH="$2"; shift 2;
      else
        shift;
      fi;;
    --clean) CLEAN=1; shift;;
    --export-sections-patch)
      EXPORT_SECTIONS_PATCH=1
      if [[ $# -ge 2 && "$2" != --* ]]; then
        SECTIONS_PATCH_PATH="$2"; shift 2;
      else
        shift;
      fi;;
    -h|--help) usage;;
    *) usage;;
  esac
done

mkdir -p "$(dirname "$LOG")"
OUT_DIR="$ROOT/$OUT_DIR_SUFFIX"
mkdir -p "$OUT_DIR"
OUT_TODO="$OUT_DIR/todo_docflow_tasklist.md"
PATCH_TMP=$(mktemp)

# Run import policy lint with hints
./scripts/import_policy_lint.sh --root "$ROOT" --policy docs/_meta/IMPORT_POLICY.yaml --docflow docflow.yaml --log /tmp/fix_flow_import_policy.md --suggest || true

# Run source template check with hints (hints-only)
GOFLAGS=-mod=vendor go run ./cmd/source-template-check docs/_meta/IMPORT_POLICY.yaml docflow.yaml "$ROOT" --suggest > /tmp/fix_flow_source_template.txt || true

cat > "$OUT_TODO" <<'EOF'
---
title: "TODO docflow"
doc_id: todo_docflow_tasklist
doc_type: tasklist
status: draft
owner: "doc_owner"
version: "v0.1.0"
context_sources: ["todo"]
---

# TODO docflow (priority order)

## 1) Sekcje brakujące (kanon) — wstaw poniższe szkielety
EOF

cat /tmp/fix_flow_source_template.txt >> "$OUT_TODO"

cat >> "$OUT_TODO" <<'EOF'

## Wstęp
Lista zadań do uzupełnienia brakujących sekcji i napraw governance.

## Zadania
- Uzupełnij sekcje kanonu w LocalFold (Przegląd, Decyzje architektoniczne, Komponenty, Ryzyka).
- Zweryfikuj owner/version w dokumentach importowanych.
- Powtórz compliance i queue po uzupełnieniach.

## Następne kroki
- Po wykonaniu powyższych, uruchom ponownie lintery i zaktualizuj GO/NO-GO.

## 2) Governance / compliance
- Po uzupełnieniu sekcji uruchom ponownie:
  - ./scripts/import_policy_lint.sh --root tmp/localfold_import --policy docs/_meta/IMPORT_POLICY.yaml --docflow docflow.yaml --log LOGS/import_policy_lint.md --suggest
  - GOFLAGS=-mod=vendor go run ./cmd/source-template-check docs/_meta/IMPORT_POLICY.yaml docflow.yaml tmp/localfold_import --suggest

## 3) Kolejka/backlog
- Zweryfikuj backlog:
  - ./scripts/backlog_prune.sh --in /tmp/backlog_dups.txt --out /tmp/backlog_dups_pruned.txt --log LOGS/backlog_prune.md --suggest

## 4) Potwierdzenie
- Gdy wszystko PASS, usuń ten TODO lub oznacz jako wykonany.
EOF

if [[ "$CLEAN" -eq 1 ]]; then
  cat >> "$OUT_TODO" <<'EOF'

## Status
Brak nowych braków. Dokumenty są kompletne względem bieżących reguł.
EOF
fi

if [[ "$EXPORT_PATCH" -eq 1 ]]; then
  # patch względem ewentualnego poprzedniego TODO (jeśli istnieje)
  ORIG=$(mktemp)
  if [[ -f "$OUT_TODO" ]]; then
    cp "$OUT_TODO" "$ORIG"
  else
    cp /dev/null "$ORIG"
  fi
  diff -u "$ORIG" "$OUT_TODO" > "$PATCH_PATH" || true
  rm -f "$ORIG"
fi

if [[ "$EXPORT_SECTIONS_PATCH" -eq 1 ]]; then
  python3 - <<'PY' "$SECTIONS_PATCH_PATH" "$ROOT" "/tmp/fix_flow_source_template.txt"
import sys, pathlib, re
out_path, root, hints_path = sys.argv[1:4]
hints = pathlib.Path(hints_path)
if not hints.exists() or hints.stat().st_size == 0:
    pathlib.Path(out_path).write_text("", encoding="utf-8")
    sys.exit(0)
content = hints.read_text(encoding="utf-8").splitlines()
patch_lines = []
current_file = None
collect = False
buf = []
for line in content:
    m = re.match(r'^##\s+(.*\.md)', line.strip())
    if m:
        if current_file and buf:
            patch_lines.append(f"--- {current_file}\\n+++ {current_file}\\n@@ -0,0 +1,{len(buf)} @@\\n" + ''.join(buf))
        current_file = m.group(1).strip()
        buf = []
        collect = False
        continue
    if line.strip() == "```markdown":
        collect = True
        continue
    if line.strip() == "```":
        collect = False
        continue
    if collect:
        buf.append("+" + line + "\n")

if current_file and buf:
    patch_lines.append(f"--- {current_file}\\n+++ {current_file}\\n@@ -0,0 +1,{len(buf)} @@\\n" + ''.join(buf))

pathlib.Path(out_path).write_text(''.join(patch_lines), encoding='utf-8')
PY
fi

{
  echo "# Fix flow"
  echo "Root: $ROOT"
  echo "TODO: $OUT_TODO"
  echo "Import policy log: /tmp/fix_flow_import_policy.md"
  echo "Source template hints: /tmp/fix_flow_source_template.txt"
  if [[ "$EXPORT_PATCH" -eq 1 ]]; then
    echo "Patch: $PATCH_PATH"
  fi
} > "$LOG"

echo "Fix flow generated. TODO: $OUT_TODO"
