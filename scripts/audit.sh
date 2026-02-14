#!/usr/bin/env bash
set -euo pipefail

ROOT=""
LOG="LOGS/audit.md"
BIN="./build/docflow-linux-amd64"
RULES="docs/_meta/GOVERNANCE_RULES.yaml"
TS="$(date +%Y%m%dT%H%M%S)"
KEEP=5

usage() {
  echo "usage: $0 --root DIR [--log FILE]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done

[[ -z "$ROOT" ]] && usage
if [[ ! -f "$ROOT/docflow.yaml" ]]; then
  echo "Brak docflow.yaml w $ROOT" >&2
  exit 1
fi
if [[ ! -x "$BIN" ]]; then
  echo "Brak binarki $BIN (zbuduj: ./scripts/build.sh)" >&2
  exit 1
fi

mkdir -p "$(dirname "$LOG")"

ROOT_NAME=$(basename "$ROOT")
OUT_DIR="dist/health/$ROOT_NAME"
mkdir -p "$OUT_DIR"

VALIDATE_LOG="$OUT_DIR/validate.log"
COMPLIANCE_JSON="$OUT_DIR/compliance_${TS}.json"
QUEUE_JSON="$OUT_DIR/queue_${TS}.json"
QUEUE_HTML="$OUT_DIR/queue_${TS}.html"
BADGE="$OUT_DIR/health_${TS}.svg"
BADGE_LATEST="$OUT_DIR/health_latest.svg"
QUEUE_HTML_LATEST="$OUT_DIR/queue_latest.html"
HEALTH_HTML="$OUT_DIR/health_${TS}.html"
HEALTH_HTML_LATEST="$OUT_DIR/health_latest.html"
HEALTH_REPORT="$OUT_DIR/health_report_${TS}.html"
HEALTH_REPORT_LATEST="$OUT_DIR/health_report_latest.html"

TS=$(date -Iseconds)
EXIT_CODE=0

validate_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$BIN validate --config $ROOT/docflow.yaml --governance $RULES" 2>&1 || true)
echo "$validate_log" > "$VALIDATE_LOG"

compliance_log=$(/usr/bin/time -f "TIME:%e" bash -lc "$BIN compliance --config $ROOT/docflow.yaml --rules $RULES --format json --output $COMPLIANCE_JSON" 2>&1 || true)

backlog=$(mktemp)
printf "T1 %s\n" "$ROOT" > "$backlog"
CACHE="/tmp/audit_cache.json"
queue_log=$(/usr/bin/time -f "TIME:%e" bash -lc "./scripts/queue_evaluate.sh --format json --bin $BIN --cache $CACHE $backlog > $QUEUE_JSON" 2>&1 || true)
python3 ./scripts/queue_report_html.py "$QUEUE_JSON" "$QUEUE_HTML"
cp "$QUEUE_HTML" "$HEALTH_HTML" 2>/dev/null || true

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

python3 - <<'PY' "$BADGE" "$blocked"
import sys, pathlib
badge_path, blocked = sys.argv[1], int(sys.argv[2])
color = "#4caf50" if blocked == 0 else "#e53935"
text = "PASS" if blocked == 0 else "FAIL"
svg = f"""<svg xmlns="http://www.w3.org/2000/svg" width="110" height="20">
<rect width="50" height="20" fill="#555"/>
<rect x="50" width="60" height="20" fill="{color}"/>
<text x="25" y="14" fill="#fff" font-size="11" text-anchor="middle">audit</text>
<text x="80" y="14" fill="#fff" font-size="11" text-anchor="middle">{text}</text>
</svg>"""
pathlib.Path(badge_path).write_text(svg)
PY

# symlinks latest
ln -sf "$(basename "$BADGE")" "$BADGE_LATEST"
ln -sf "$(basename "$QUEUE_HTML")" "$QUEUE_HTML_LATEST"
ln -sf "$(basename "$HEALTH_HTML")" "$HEALTH_HTML_LATEST"
if [[ -f /tmp/health_report.html ]]; then
  cp /tmp/health_report.html "$HEALTH_REPORT"
  ln -sf "$(basename "$HEALTH_REPORT")" "$HEALTH_REPORT_LATEST"
fi

# rotate keep latest $KEEP
ls -1t "$OUT_DIR"/health_*.svg 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f
ls -1t "$OUT_DIR"/queue_*.html 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f
ls -1t "$OUT_DIR"/queue_*.json 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f
ls -1t "$OUT_DIR"/compliance_*.json 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f
ls -1t "$OUT_DIR"/health_*.html 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f
ls -1t "$OUT_DIR"/health_report_*.html 2>/dev/null | tail -n +$((KEEP+1)) | xargs -r rm -f

{
  echo "# Audit"
  echo "Timestamp: $TS"
  echo "Root: $ROOT"
  echo "Ready: $ready"
  echo "Blocked: $blocked"
  echo
  echo "Validate log (snippet):"
  echo '```'
  echo "$validate_log"
  echo '```'
  echo
  echo "Compliance log (stderr):"
  echo '```'
  echo "$compliance_log"
  echo '```'
  echo
  echo "Queue log:"
  echo '```'
  echo "$queue_log"
  echo '```'
  echo
  echo "Artifacts:"
  echo "- Compliance JSON: $COMPLIANCE_JSON"
  echo "- Queue JSON: $QUEUE_JSON"
  echo "- Queue HTML: $QUEUE_HTML"
  echo "- Badge: $BADGE"
} > "$LOG"

rm -f "$backlog"

if [[ "$blocked" -gt 0 ]]; then
  EXIT_CODE=2
fi

exit "$EXIT_CODE"
