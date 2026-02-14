#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
OUT=worklog/days/index.md
printf "# Index wpisow\n\n" > "$OUT"
printf "| wpis | tytul (pierwsza linia) |\n|---|---|\n" >> "$OUT"

mapfile -t files < <(find worklog/days -maxdepth 1 -type f -name '*.md' ! -name 'index.md' | LC_ALL=C sort -V)
for f in "${files[@]}"; do
  entry=$(basename "$f" .md)
  title=$(head -n 1 "$f" | sed 's/^# //')
  printf "| [%s](%s) | %s |\n" "$entry" "$entry.md" "$title" >> "$OUT"
done
