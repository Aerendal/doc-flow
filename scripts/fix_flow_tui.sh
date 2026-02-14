#!/usr/bin/env bash
set -euo pipefail
ROOT="tmp/localfold_import"
TODO=""
PATCH=""
LOG="LOGS/fix_flow_tui.md"
usage(){ echo "usage: $0 [--root dir] [--todo path] [--patch path] [--log file]" >&2; exit 1; }
while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --todo) TODO="$2"; shift 2;;
    --patch) PATCH="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done
mkdir -p "$(dirname "$LOG")"
[[ -z "$TODO" ]] && TODO="$ROOT/_todo/todo_docflow_tasklist.md"
[[ -z "$PATCH" ]] && PATCH="/tmp/fix_flow_hints.patch"
choices=()
[[ -f "$TODO" ]] && choices+=("TODO:$TODO")
[[ -f "$PATCH" && -s "$PATCH" ]] && choices+=("PATCH:$PATCH")
[[ -f /tmp/fix_flow_source_template.txt ]] && choices+=("HINTS:/tmp/fix_flow_source_template.txt")
[[ -f /tmp/fix_flow_import_policy.md ]] && choices+=("POLICY:/tmp/fix_flow_import_policy.md")
if [[ ${#choices[@]} -eq 0 ]]; then
  echo "Brak artefaktów do podglądu" | tee "$LOG"
  exit 0
fi
selected=""
if command -v fzf >/dev/null 2>&1; then
  selected=$(printf '%s\n' "${choices[@]}" | fzf --prompt="fix-flow > ") || true
fi
if [[ -z "$selected" ]]; then
  {
    printf 'Dostępne artefakty:\n'
    printf '%s\n' "${choices[@]}"
    echo "Tip: zainstaluj fzf, np. sudo apt-get install fzf  (brew install fzf na macOS)."
    echo "Fallback: ręcznie otwórz plik z listy powyżej (less/vi)."
  } | tee "$LOG"
  exit 0
fi
path="${selected#*:}"
{
  echo "# Fix flow TUI"
  echo "Root: $ROOT"
  echo "Wybrano: $selected"
} > "$LOG"
less "$path"
