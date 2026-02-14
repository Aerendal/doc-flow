#!/usr/bin/env bash
set -euo pipefail

# Generate a synthetic backlog for queue smoke tests.
# Usage: gen_backlog_smoke.sh [COUNT] [OUTFILE]
# COUNT   - number of tasks to emit (default: 200)
# OUTFILE - output file path (default: /tmp/backlog_smoke.txt)

COUNT="${1:-200}"
OUT="${2:-/tmp/backlog_smoke.txt}"

paths=(examples/simple-api examples/architecture)
total_paths=${#paths[@]}

mkdir -p "$(dirname "$OUT")"
> "$OUT"

for ((i=1; i<=COUNT; i++)); do
  p="${paths[$(( (i-1) % total_paths ))]}"
  printf "T%03d %s\n" "$i" "$p" >> "$OUT"
done

echo "Generated $COUNT tasks into $OUT"
