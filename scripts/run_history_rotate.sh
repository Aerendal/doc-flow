#!/usr/bin/env bash
set -euo pipefail

HIST=".docflow/run_history.json"
ARCHIVE="LOGS/run_history_archive.json"
LOG="LOGS/run_history_rotate.md"
LIMIT=200
MAX_SIZE=$((1024*1024)) # 1MB

mkdir -p .docflow LOGS
touch "$HIST"

python - <<'PY' "$HIST" "$ARCHIVE" "$LIMIT" "$MAX_SIZE" "$LOG"
import sys, json, pathlib, shutil
hist_path, arch_path, limit, max_size, log_path = sys.argv[1:6]
limit = int(limit); max_size = int(max_size)
hp = pathlib.Path(hist_path)
ap = pathlib.Path(arch_path)
log_lines = []

try:
    data = json.loads(hp.read_text()) if hp.read_text().strip() else []
except Exception:
    shutil.copy2(hp, hp.with_suffix(hp.suffix + ".bak"))
    data = []
    log_lines.append("history corrupted -> reset to empty; backup made")

def rotate(data):
    moved=[]
    if len(data) > limit:
        over = len(data) - limit
        moved = data[:over]
        data = data[over:]
    size = len(json.dumps(data))
    if size > max_size:
        # remove oldest until size ok
        while data and len(json.dumps(data)) > max_size:
            moved.append(data.pop(0))
    return data, moved

new_data, moved = rotate(data)

if moved:
    arch = []
    if ap.exists() and ap.read_text().strip():
        arch = json.loads(ap.read_text())
    arch.extend(moved)
    ap.parent.mkdir(parents=True, exist_ok=True)
    ap.write_text(json.dumps(arch, indent=2))

hp.write_text(json.dumps(new_data, indent=2))

log_lines.append(f"entries_before={len(data)}, entries_after={len(new_data)}, moved={len(moved)}")
log_lines.append(f"archive={arch_path}")

pathlib.Path(log_path).write_text("\n".join(log_lines))
PY

echo "Rotation complete. Log: $LOG"
