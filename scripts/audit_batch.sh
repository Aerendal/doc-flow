#!/usr/bin/env bash
set -euo pipefail

LOG="LOGS/audit_batch.md"
ROOTS=()

usage(){ echo "usage: $0 --roots \"dir1 dir2 ...\" [--log FILE]" >&2; exit 1; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --roots) IFS=' ' read -r -a ROOTS <<< "$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done

[[ ${#ROOTS[@]} -eq 0 ]] && usage

mkdir -p "$(dirname "$LOG")"

results=()
status=0

for r in "${ROOTS[@]}"; do
  out_log="LOGS/AUDIT_BATCH_${r//\//_}.md"
  echo "[audit_batch] auditing $r (log: $out_log)"
  if ./scripts/audit.sh --root "$r" --log "$out_log"; then
    ready=$(python3 - "$out_log" <<'PY'
import sys,re
text=open(sys.argv[1]).read()
m=re.search(r"Ready:\s*(\d+)", text)
print(m.group(1) if m else "0")
PY)
    blocked=$(python3 - "$out_log" <<'PY'
import sys,re
text=open(sys.argv[1]).read()
m=re.search(r"Blocked:\s*(\d+)", text)
print(m.group(1) if m else "0")
PY)
    results+=("$r\tREADY\tready=$ready blocked=$blocked\t$out_log")
    if [[ "$blocked" != "0" ]]; then status=2; fi
  else
    results+=("$r\tFAIL\tblocked>0 or error\t$out_log")
    status=2
  fi
done

{
  echo "# Audit batch"
  echo "Roots: ${ROOTS[*]}"
  echo "Results:"
  printf '%s\n' "${results[@]}"
} > "$LOG"

cat "$LOG"
exit $status
