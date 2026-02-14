#!/usr/bin/env bash
set -euo pipefail

INPUT=""
OUT="/tmp/backlog_pruned.txt"
LOG="LOGS/backlog_prune.md"
APPLY=0
SUGGEST=0

usage() {
  echo "usage: $0 --in backlog.txt [--out file] [--log file] [--apply] [--suggest]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --in) INPUT="$2"; shift 2;;
    --out) OUT="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    --apply) APPLY=1; shift;;
    --suggest) SUGGEST=1; shift;;
    -h|--help) usage;;
    *) usage;;
  esac
done

[ -f "$INPUT" ] || { echo "missing --in file" >&2; exit 1; }
mkdir -p "$(dirname "$OUT")" "$(dirname "$LOG")"

python - "$INPUT" "$OUT" "$LOG" "$APPLY" "$SUGGEST" <<'PY'
import sys, pathlib, datetime
inp, outp, logp, apply_flag, suggest_flag = sys.argv[1:6]
apply_flag = int(apply_flag)
suggest_flag = int(suggest_flag)

lines = []
with open(inp) as f:
    for ln in f:
        ln = ln.strip()
        if not ln:
            continue
        parts = ln.split(maxsplit=1)
        if len(parts) != 2:
            continue
        lines.append((parts[0], parts[1]))

seen_ids = set()
dups = []
missing_path = []
kept = []
for tid, path in lines:
    if tid in seen_ids:
        dups.append((tid, path))
        continue
    seen_ids.add(tid)
    if not pathlib.Path(path).exists():
        missing_path.append((tid, path))
        continue
    kept.append((tid, path))

if apply_flag:
    with open(outp, "w") as f:
        for tid, path in kept:
            f.write(f"{tid} {path}\n")

now = datetime.datetime.now().isoformat(timespec="seconds")
log_lines = [
    "# Backlog prune",
    f"Input: {inp}",
    f"Output: {outp}",
    f"Apply: {bool(apply_flag)}",
    f"Timestamp: {now}",
    "",
    f"Total entries: {len(lines)}",
    f"Kept: {len(kept)}",
    f"Duplicates removed: {len(dups)}",
    f"Missing paths: {len(missing_path)}",
    "",
    "## Duplicates",
]
log_lines += [f"- {tid} {path}" for tid, path in dups] or ["- none"]
log_lines += ["", "## Missing paths"]
log_lines += [f"- {tid} {path}" for tid, path in missing_path] or ["- none"]

if suggest_flag:
    log_lines += ["", "## Hints"]
    if dups:
        for i,(tid,_) in enumerate(dups[:10], start=1):
            log_lines += [f"{i}. Duplikat ID {tid}: zmień na {tid}_2 / {tid}_copy"]
        if len(dups)>10:
            log_lines += [f"+{len(dups)-10} więcej duplikatów..."]
    if missing_path:
        for i,(tid,_) in enumerate(missing_path[:10], start=1):
            log_lines += [f"{i}. Brak ścieżki dla {tid}: ustaw np. ./examples/..."]
        if len(missing_path)>10:
            log_lines += [f"+{len(missing_path)-10} więcej brakujących ścieżek..."]
    if not dups and not missing_path:
        log_lines += ["- Brak podpowiedzi (backlog czysty)"]

pathlib.Path(logp).write_text("\n".join(log_lines), encoding="utf-8")

sys.exit(0 if not missing_path and not dups else 2)
PY

echo "Prune complete. Log: $LOG"
