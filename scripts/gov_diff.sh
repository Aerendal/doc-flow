#!/usr/bin/env bash
set -euo pipefail

SNAPSHOT=".docflow/gov_snapshot.yaml"
LOGFILE="LOGS/gov_diff.md"
UPDATE=0

usage() {
  echo "usage: $0 [--snapshot PATH] [--log PATH] [--update-snapshot]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --snapshot) SNAPSHOT="$2"; shift 2;;
    --log) LOGFILE="$2"; shift 2;;
    --update-snapshot) UPDATE=1; shift;;
    -h|--help) usage;;
    *) usage;;
  esac
done

CURRENT="docs/_meta/GOVERNANCE_RULES.yaml"
[ -f "$CURRENT" ] || { echo "missing $CURRENT" >&2; exit 2; }
mkdir -p "$(dirname "$SNAPSHOT")" LOGS

run_python() {
python - <<'PY' "$CURRENT" "$SNAPSHOT" "$UPDATE"
import sys, json, pathlib
try:
    import yaml
except ImportError as e:
    sys.stderr.write("PyYAML required\n")
    sys.exit(3)

cur_path, snap_path, update_flag = sys.argv[1], sys.argv[2], sys.argv[3]
cur = yaml.safe_load(pathlib.Path(cur_path).read_text()) or {}
if pathlib.Path(snap_path).exists():
    snap = yaml.safe_load(pathlib.Path(snap_path).read_text()) or {}
    snapshot_exists = True
else:
    snap = {}
    snapshot_exists = False

def leaves(obj, prefix=""):
    items=[]
    if isinstance(obj, dict):
        for k in sorted(obj):
            items += leaves(obj[k], f"{prefix}.{k}" if prefix else k)
    elif isinstance(obj, list):
        for i, v in enumerate(obj):
            items += leaves(v, f"{prefix}[{i}]")
        if not obj:
            items.append((prefix, "[]"))
    else:
        items.append((prefix, obj))
    return items

cur_d = dict(leaves(cur))
snap_d = dict(leaves(snap))

added   = sorted(k for k in cur_d if k not in snap_d)
removed = sorted(k for k in snap_d if k not in cur_d)
changed = sorted(k for k in cur_d if k in snap_d and cur_d[k] != snap_d[k])

summary = {
    "snapshot_exists": snapshot_exists,
    "added": added,
    "removed": removed,
    "changed": changed,
}
print(json.dumps(summary))

if update_flag == "1":
    pathlib.Path(snap_path).write_text(pathlib.Path(cur_path).read_text())
PY
}

json_out=$(run_python)
snapshot_exists=$(JSON="$json_out" python - <<'PY'
import os, json
d=json.loads(os.environ["JSON"])
print("1" if d.get("snapshot_exists") else "0")
PY
)

added=$(JSON="$json_out" python - <<'PY'
import os, json
d=json.loads(os.environ["JSON"])
print("\n".join(d["added"]))
PY
)
removed=$(JSON="$json_out" python - <<'PY'
import os, json
d=json.loads(os.environ["JSON"])
print("\n".join(d["removed"]))
PY
)
changed=$(JSON="$json_out" python - <<'PY'
import os, json
d=json.loads(os.environ["JSON"])
print("\n".join(d["changed"]))
PY
)

cat > "$LOGFILE" <<'EOF'
# GOV diff

EOF
{
  echo "Snapshot file: $SNAPSHOT"
  echo "Snapshot existed before run: $([ "$snapshot_exists" = "1" ] && echo yes || echo no)"
  echo
  echo "## Added keys"
  echo "${added:-<none>}"
  echo
  echo "## Removed keys"
  echo "${removed:-<none>}"
  echo
  echo "## Changed keys"
  echo "${changed:-<none>}"
  echo
  if [ "$UPDATE" -eq 1 ]; then
    echo "Snapshot refreshed: yes"
  else
    echo "Snapshot refreshed: no"
  fi
} >> "$LOGFILE"

diff_found=0
if [ -n "$added$removed$changed" ]; then diff_found=1; fi

if [ "$UPDATE" -eq 1 ]; then
  exit 0
elif [ "$snapshot_exists" = "0" ]; then
  echo "snapshot missing (use --update-snapshot to create)" >&2
  exit 1
elif [ $diff_found -eq 1 ]; then
  exit 2
else
  exit 0
fi
