#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(pwd)"
DEMO_DIR="/tmp/docflow-demo"
WORK="$DEMO_DIR/examples"
BACKLOG="$DEMO_DIR/backlog.txt"
CACHE="/tmp/quickstart_cache.json"
QUEUE_JSON="/tmp/quickstart_queue.json"
QUEUE_HTML="/tmp/quickstart_report.html"
LOG="LOGS/quickstart.md"
BIN="$ROOT_DIR/build/docflow-linux-amd64"
HIST="$ROOT_DIR/.docflow/run_history.json"

rm -rf "$DEMO_DIR"
mkdir -p "$DEMO_DIR" LOGS

cp -r examples/simple-api "$DEMO_DIR/examples"
echo "T1 $WORK" > "$BACKLOG"

validate_log=$(/usr/bin/time -f "TIME:%e" bash -lc "cd $WORK && $BIN validate --config docflow.yaml --governance $ROOT_DIR/docs/_meta/GOVERNANCE_RULES.yaml" 2>&1)
compliance_log=$(/usr/bin/time -f "TIME:%e" bash -lc "cd $WORK && $BIN compliance --config docflow.yaml --rules $ROOT_DIR/docs/_meta/GOVERNANCE_RULES.yaml --format text > /tmp/quickstart_compliance.txt" 2>&1)

queue_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$ROOT_DIR/scripts/queue_evaluate.sh --bin $BIN --format json --cache $CACHE --workers 2 $BACKLOG > $QUEUE_JSON" 2>&1)
python3 scripts/queue_report_html.py "$QUEUE_JSON" "$QUEUE_HTML"

ready=$(python3 - <<'PY' "$QUEUE_JSON"
import json,sys
data=json.load(open(sys.argv[1]))
print(data.get("ready",0))
PY
)
blocked=$(python3 - <<'PY' "$QUEUE_JSON"
import json,sys
data=json.load(open(sys.argv[1]))
print(data.get("blocked",0))
PY
)

{
  echo "# Quickstart"
  echo "Demo dir: $DEMO_DIR"
  echo "Backlog: $BACKLOG"
  echo "Queue JSON: $QUEUE_JSON"
  echo "Report HTML: $QUEUE_HTML"
  echo
  echo "## Validate"
  echo '```'
  echo "$validate_log"
  echo '```'
  echo
  echo "## Compliance"
  echo '```'
  echo "$compliance_log"
  echo '```'
  echo
  echo "## Queue"
  echo "ready=$ready blocked=$blocked"
  echo '```'
  echo "$queue_log"
  echo '```'
} > "$LOG"

echo "Quickstart complete. Log: $LOG ; report: $QUEUE_HTML"
