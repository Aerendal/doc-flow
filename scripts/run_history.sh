#!/usr/bin/env bash
set -euo pipefail

HIST=".docflow/run_history.json"
LOG="LOGS/run_history.md"
ACTION="${1:-list}"  # list|add|open|tui
ENTRY="${2:-}"
TARGET="${3:-}"

mkdir -p .docflow LOGS
touch "$HIST"

add_entry() {
  local kind="$1" exitc="$2" artifact="$3" duration="$4"
  python - "$HIST" "$kind" "$exitc" "$artifact" "$duration" <<'PY'
import sys, json, pathlib, datetime
hist_path, kind, exitc, artifact, duration = sys.argv[1:6]
path = pathlib.Path(hist_path)
data = json.loads(path.read_text()) if path.read_text().strip() else []
data.append({
    "ts": datetime.datetime.now().isoformat(timespec="seconds"),
    "kind": kind,
    "exit": int(exitc),
    "artifact": artifact,
    "duration": duration
})
path.write_text(json.dumps(data, indent=2))
PY
}

list_entries() {
  python - "$HIST" <<'PY'
import sys, json, pathlib
path=pathlib.Path(sys.argv[1])
data = json.loads(path.read_text()) if path.read_text().strip() else []
for i, d in enumerate(reversed(data[-50:])):
    print(f"{len(data)-i}: {d['ts']} kind={d['kind']} exit={d['exit']} artifact={d['artifact']} duration={d.get('duration','')}")
PY
}

open_entry() {
  local idx="$1"
  python - "$HIST" "$idx" <<'PY'
import sys, json, pathlib, os, subprocess
path=pathlib.Path(sys.argv[1])
idx=int(sys.argv[2])
data = json.loads(path.read_text()) if path.read_text().strip() else []
if idx<=0 or idx>len(data):
    sys.exit("bad index")
entry=data[idx-1]
art=entry.get("artifact","")
if not art:
    sys.exit("no artifact")
if not pathlib.Path(art).exists():
    sys.exit(f"artifact not found: {art}")
os.execvp("less", ["less", "-R", art] if art.endswith(".md") or art.endswith(".txt") else ["less", "-R", art])
PY
}

tui() {
  if ! command -v fzf >/dev/null 2>&1; then
    echo "fzf not found; fallback to list/open (run run_history.sh open <idx>)" >&2
    echo "install fzf e.g.: sudo apt-get install fzf   # lub brew install fzf" >&2
    list_entries
    exit 0
  fi
  python - "$HIST" <<'PY'
import sys, json, subprocess, pathlib, os
path=pathlib.Path(sys.argv[1])
data = json.loads(path.read_text()) if path.read_text().strip() else []
items=[]
for i,d in enumerate(reversed(data[-200:])):
    items.append(f"{len(data)-i}\t{d['ts']}\t{d['kind']}\t{d['exit']}\t{d.get('artifact','')}")
fzf = subprocess.run(["fzf","--with-nth=1..","--prompt","history> "], input="\n".join(items), text=True, capture_output=True)
if fzf.returncode!=0 or not fzf.stdout.strip():
    sys.exit(0)
line=fzf.stdout.strip().split("\t")[0]
idx=int(line)
entry=data[idx-1]
art=entry.get("artifact","")
if not art:
    sys.exit(0)
if not pathlib.Path(art).exists():
    sys.exit(f"artifact not found: {art}")
os.execvp("less", ["less","-R", art] if art.endswith(".md") or art.endswith(".txt") else ["less","-R", art])
PY
}

case "$ACTION" in
  add) add_entry "${2:-kind}" "${3:-0}" "${4:-}" "${5:-}";;
  list) list_entries;;
  open) open_entry "$ENTRY";;
  tui) tui;;
  *) echo "usage: run_history.sh [list|add|open <idx>]"; exit 1;;
esac

echo "# Run history" > "$LOG"
list_entries >> "$LOG"
