#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(pwd)"
LOG="LOGS/healthcheck.md"
TS="$(date +%Y%m%dT%H%M%S)"
HTML="/tmp/health_report.html"
BADGE="/tmp/health_badge.svg"
BACKLOG="examples/backlog_ci_template.txt"
BIN="$ROOT_DIR/build/docflow-linux-amd64"
CACHE="/tmp/health_cache.json"
QUEUE_JSON="/tmp/health_queue.json"
OUTDIR="dist/health"
KEEP=5

mkdir -p LOGS

validate_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$BIN validate --config examples/simple-api/docflow.yaml --governance $ROOT_DIR/docs/_meta/GOVERNANCE_RULES.yaml" 2>&1 || true)
compliance_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$BIN compliance --config examples/simple-api/docflow.yaml --rules $ROOT_DIR/docs/_meta/GOVERNANCE_RULES.yaml --format text > /tmp/health_compliance.txt" 2>&1 || true)

prune_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$ROOT_DIR/scripts/backlog_prune.sh --in $BACKLOG --out /tmp/health_backlog.txt --log /tmp/health_prune.log --apply" 2>&1 || true)

queue_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$ROOT_DIR/scripts/queue_evaluate.sh --bin $BIN --format json --cache $CACHE --workers 2 /tmp/health_backlog.txt > $QUEUE_JSON" 2>&1 || true)
python3 "$ROOT_DIR/scripts/queue_report_html.py" "$QUEUE_JSON" "$HTML"

ready=$(python3 - <<'PY' "$QUEUE_JSON"
import json,sys
d=json.load(open(sys.argv[1]))
print(d.get("ready",0))
PY
)
blocked=$(python3 - <<'PY' "$QUEUE_JSON"
import json,sys
d=json.load(open(sys.argv[1]))
print(d.get("blocked",0))
PY
)

python - <<'PY' "$BADGE" "$blocked"
import sys, pathlib
badge_path, blocked = sys.argv[1], int(sys.argv[2])
color = "#4caf50" if blocked == 0 else "#e53935"
text = "PASS" if blocked == 0 else "FAIL"
svg = f"""<svg xmlns="http://www.w3.org/2000/svg" width="110" height="20">
<rect width="50" height="20" fill="#555"/>
<rect x="50" width="60" height="20" fill="{color}"/>
<text x="25" y="14" fill="#fff" font-size="11" text-anchor="middle">health</text>
<text x="80" y="14" fill="#fff" font-size="11" text-anchor="middle">{text}</text>
</svg>"""
pathlib.Path(badge_path).write_text(svg)
PY

{
  echo "# Healthcheck"
  echo "Validate:"
  echo '```'
  echo "$validate_log"
  echo '```'
  echo
  echo "Compliance:"
  echo '```'
  echo "$compliance_log"
  echo '```'
  echo
  echo "Backlog prune:"
  echo '```'
  echo "$prune_log"
  echo '```'
  echo
  echo "Queue (ready=$ready blocked=$blocked):"
  echo '```'
  echo "$queue_log"
  echo '```'
  echo
  echo "Reports:"
  echo "- Queue JSON: $QUEUE_JSON"
  echo "- Queue HTML: $HTML"
  echo "- Badge SVG: $BADGE"
} > "$LOG"

# ensure outdir
mkdir -p "$OUTDIR"

# timestamped copies + latest symlinks
BADGE_TS="$OUTDIR/health_${TS}.svg"
HTML_TS="$OUTDIR/health_${TS}.html"
cp "$BADGE" "$BADGE_TS"
cp "$HTML" "$HTML_TS"
ln -sf "$(basename "$BADGE_TS")" "$OUTDIR/health_latest.svg"
ln -sf "$(basename "$HTML_TS")" "$OUTDIR/health_latest.html"

# rotate keep $KEEP newest
ls -1t "$OUTDIR"/health_*.svg 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f
ls -1t "$OUTDIR"/health_*.html 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f

# copy artifacts to LOGS (latest)
cp "$BADGE" LOGS/health_badge.svg
cp "$HTML" LOGS/health_report.html

python - <<'PY' "$LOG" "$ready" "$blocked" ".docflow/run_history.json"
import sys, json, pathlib, datetime
log, ready, blocked, hist = sys.argv[1:5]
entry = {
    "ts": datetime.datetime.now().isoformat(timespec="seconds"),
    "kind": "health",
    "exit": 0 if int(blocked)==0 else 2,
    "artifact": log,
    "duration": "",
    "ready": int(ready),
    "blocked": int(blocked)
}
hpath = pathlib.Path(hist)
data = json.loads(hpath.read_text()) if hpath.exists() and hpath.read_text().strip() else []
hpath.parent.mkdir(parents=True, exist_ok=True)
data.append(entry)
hpath.write_text(json.dumps(data, indent=2))
PY

if [[ "$blocked" -gt 0 ]]; then
  exit 2
fi
