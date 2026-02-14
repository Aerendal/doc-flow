#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
OUT=worklog/days/index.md
printf "# Index day_xx\n\n" > "$OUT"
printf "| day | tytuł (pierwsza linia) |\n|---|---|\n" >> "$OUT"
for f in $(ls worklog/days/day_*.md | LC_ALL=C sort -V); do
  day=$(basename "$f" .md)
  title=$(head -n 1 "$f" | sed 's/^# //')
  printf "| [%s](%s) | %s |\n" "$day" "$day.md" "$title" >> "$OUT"
done
