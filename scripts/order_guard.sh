#!/usr/bin/env bash
set -euo pipefail

BUNDLE_ZIP="dist/release/docflow-pr-bundle.zip"
HEALTH_HTML="dist/health/localfold_import/health_latest.html"
QUEUE_JSON="dist/health/localfold_import/queue_latest.json"
COMPLIANCE_JSON="dist/health/localfold_import/compliance_latest.json"
LOG="LOGS/order_guard.md"
STRICT=1

usage(){ echo "usage: $0 [--bundle PATH] [--health HTML] [--queue JSON] [--compliance JSON] [--log FILE] [--no-strict]" >&2; exit 1; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bundle) BUNDLE_ZIP="$2"; shift 2;;
    --health) HEALTH_HTML="$2"; shift 2;;
    --queue) QUEUE_JSON="$2"; shift 2;;
    --compliance) COMPLIANCE_JSON="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    --no-strict) STRICT=0; shift;;
    -h|--help) usage;;
    *) usage;;
  esac
done

mkdir -p "$(dirname "$LOG")"

missing=()
for f in "$BUNDLE_ZIP" "$HEALTH_HTML" "$QUEUE_JSON" "$COMPLIANCE_JSON"; do
  [[ -f "$f" ]] || missing+=("$f")
done

bundle_ts=0
if [[ -f "$BUNDLE_ZIP" ]]; then bundle_ts=$(stat -c %Y "$BUNDLE_ZIP"); fi

health_ts=$(stat -c %Y "$HEALTH_HTML" 2>/dev/null || echo 0)
queue_ts=$(stat -c %Y "$QUEUE_JSON" 2>/dev/null || echo 0)
comp_ts=$(stat -c %Y "$COMPLIANCE_JSON" 2>/dev/null || echo 0)

order_ok=$(( health_ts <= bundle_ts && queue_ts <= bundle_ts && comp_ts <= bundle_ts ))

{
  echo "# Order guard"
  echo "bundle: $BUNDLE_ZIP"
  echo "health: $HEALTH_HTML (ts=$health_ts)"
  echo "queue:  $QUEUE_JSON (ts=$queue_ts)"
  echo "comp:   $COMPLIANCE_JSON (ts=$comp_ts)"
  echo "bundle_ts: $bundle_ts"
  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "Missing:"; printf '%s\n' "${missing[@]}"
  fi
  echo "Order OK: $order_ok"
} > "$LOG"

if [[ ${#missing[@]} -gt 0 ]]; then
  echo "order_guard: missing artifacts (see log)"
  [[ $STRICT -eq 1 ]] && exit 1 || exit 0
fi

if [[ $order_ok -ne 1 ]]; then
  echo "order_guard: bundle newer than upstream artifacts (run audit/health before bundling)"
  [[ $STRICT -eq 1 ]] && exit 1 || exit 0
fi

echo "order_guard: OK"
