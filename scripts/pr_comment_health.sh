#!/usr/bin/env bash
set -euo pipefail

ROOT=""
BUNDLE="dist/release/docflow-artifacts.zip"
OUT="LOGS/pr_comment_health.md"

usage(){ echo "usage: $0 --root NAME [--bundle PATH] [--out FILE]" >&2; exit 1; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --bundle) BUNDLE="$2"; shift 2;;
    --out) OUT="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done

[[ -z "$ROOT" ]] && usage

health_dir="dist/health/$ROOT"
health_html="$health_dir/health_latest.html"
health_report="$health_dir/health_report_latest.html"
queue_json="$health_dir/queue_latest.json"
queue_html="$health_dir/queue_latest.html"
compliance_json=$(ls -1t "$health_dir"/compliance_*.json 2>/dev/null | head -n1 || true)

ready="n/a"; blocked="n/a"
if [[ -f "$queue_json" ]]; then
  read -r ready blocked < <(python3 - <<'PY' "$queue_json"
import json,sys
p=sys.argv[1]
d=json.load(open(p))
print(d.get("ready","n/a"), d.get("blocked","n/a"))
PY
)
fi

bundle_sha="brak"; bundle_size="brak"; bundle_status="brak"
if [[ -f "$BUNDLE" ]]; then
  bundle_sha=$(sha256sum "$BUNDLE" | awk '{print $1}')
  bundle_size=$(stat -c%s "$BUNDLE")
  bundle_status="$BUNDLE"
else
  bundle_status="(bundle: brak pliku)"
fi

health_report_entry="$health_report"
[[ -f "$health_report" ]] || health_report_entry="n/a"
[[ -f "$compliance_json" ]] || compliance_json="n/a"

cat > "$OUT" <<EOFMD
# PR comment draft – health
- root: $ROOT
- ready: $ready, blocked: $blocked
- bundle: $bundle_status (SHA: $bundle_sha, size: $bundle_size)
- health: $health_html
- health report: $health_report_entry
- queue: $queue_json (html: $queue_html)
- compliance: $compliance_json
EOFMD

echo "Comment written to $OUT"
