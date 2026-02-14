#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
backlog_validate.sh --file PATH [--root DIR] [--strict]

Sprawdza backlog (linie: task_id path):
- duplikaty task_id (ERROR)
- istnienie pliku docflow.yaml w path (ERROR w --strict, WARN w trybie domyślnym)
- ścieżki wychodzące poza root (ERROR w --strict, WARN domyślnie)

Zwraca 0, jeśli brak błędów.
USAGE
}

FILE=""
ROOT="."
STRICT=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --file) FILE="$2"; shift 2;;
    --root) ROOT="$2"; shift 2;;
    --strict) STRICT=1; shift;;
    -h|--help) usage; exit 0;;
    *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
  esac
done

if [[ -z "$FILE" ]]; then
  usage; exit 1
fi

python3 - <<'PY' "$FILE" "$ROOT" "$STRICT"
import sys, pathlib
from collections import Counter
file, root, strict = sys.argv[1], pathlib.Path(sys.argv[2]), int(sys.argv[3])
errors = []
warns = []
ids = []
lines = []
for lineno, line in enumerate(pathlib.Path(file).read_text().splitlines(), 1):
    line = line.strip()
    if not line or line.startswith("#"):
        continue
    parts = line.split()
    if len(parts) < 2:
        warns.append(f"{file}:{lineno} brak path/task_id")
        continue
    tid, path = parts[0], parts[1]
    ids.append(tid)
    lines.append((lineno, tid, path))

dup = [tid for tid, c in Counter(ids).items() if c > 1]
for d in dup:
    errors.append(f"duplikat task_id: {d}")

for lineno, tid, p in lines:
    pth = pathlib.Path(p)
    if not pth.is_absolute():
        pth = (root / pth).resolve()
    if strict and root not in pth.parents and pth != root:
        errors.append(f"{file}:{lineno} path poza root: {pth}")
    docflow = pth / "docflow.yaml"
    if not docflow.exists():
        msg = f"{file}:{lineno} brak docflow.yaml w {pth}"
        if strict:
            errors.append(msg)
        else:
            warns.append(msg)

for w in warns:
    print(f"WARN: {w}")
for e in errors:
    print(f"ERROR: {e}")

if errors:
    sys.exit(1)
sys.exit(0)
PY
