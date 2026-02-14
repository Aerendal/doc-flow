#!/usr/bin/env bash
set -euo pipefail

FROM=200
TO=214
OUT="LOGS/evidence_map.md"

usage(){ echo "usage: $0 [--from N] [--to N] [--out FILE]" >&2; exit 1; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --from) FROM="$2"; shift 2;;
    --to) TO="$2"; shift 2;;
    --out) OUT="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done

mkdir -p "$(dirname "$OUT")"

rows=()
for n in $(seq "$FROM" "$TO"); do
  day=day_${n}
  log_candidates=$(ls -1 LOGS/*${n}*.md 2>/dev/null || true)
  bundle_candidates=$(ls -1 dist/release/*${n}*.zip 2>/dev/null || true)
  if [[ -z "$log_candidates" && -z "$bundle_candidates" ]]; then
    rows+=("|${n}|brak logów/artefaktów|")
    continue
  fi
  entry="|${n}|"
  entry+="${log_candidates:-brak logów}"; entry+="; "
  entry+="${bundle_candidates:-brak bundli}"; entry+="|"
  rows+=("$entry")
done

{
  echo "# Evidence map (days ${FROM}–${TO})"
  echo
  echo "| Day | Logi / Artefakty |"
  echo "|-----|-------------------|"
  printf '%s\n' "${rows[@]}"
} > "$OUT"

echo "Evidence map written to $OUT"
