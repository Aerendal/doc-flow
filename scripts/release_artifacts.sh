#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="dist/release"
ZIP_NAME="docflow-artifacts.zip"
LOG="LOGS/release_artifacts.md"

HEALTH_SVG="dist/health/health_latest.svg"
HEALTH_HTML="dist/health/health_latest.html"
QUEUE_JSON="/tmp/queue_report.json"
QUEUE_HTML="/tmp/queue_report.html"
COMPLIANCE_JSON="/tmp/compliance_report.json"
SECTION_DIFF="$(ls -1t LOGS/SECTION_DIFF_*.md 2>/dev/null | head -n1 || true)"
HEALTH_PER_ROOT=("$(find dist/health -mindepth 2 -maxdepth 2 -name health_report_latest.html 2>/dev/null | tr '\n' ' ')")
STAGE="$(mktemp -d)"
trap 'rm -rf "${STAGE}"' EXIT

usage() {
  echo "usage: $0 [--log FILE] [--zip-name NAME]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --log) LOG="$2"; shift 2;;
    --zip-name) ZIP_NAME="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done

missing=()
for f in "$HEALTH_SVG" "$HEALTH_HTML" "$QUEUE_JSON" "$QUEUE_HTML" "$COMPLIANCE_JSON"; do
  [[ -f "$f" ]] || missing+=("$f")
done
if [[ -z "$SECTION_DIFF" || ! -f "$SECTION_DIFF" ]]; then
  SECTION_DIFF=""
fi

shopt -s nullglob
per_root_copies=()
for f in ${HEALTH_PER_ROOT[@]}; do
  root="$(basename "$(dirname "$f")")"
  cp "$f" "$STAGE/health_report_latest_${root}.html"
  per_root_copies+=("$STAGE/health_report_latest_${root}.html")
done
shopt -u nullglob

if [[ ${#missing[@]} -gt 0 ]]; then
  echo "Missing artifacts:" >&2
  printf ' - %s\n' "${missing[@]}" >&2
  exit 1
fi

mkdir -p "$OUT_DIR" "$(dirname "$LOG")"
ZIP_PATH="$OUT_DIR/$ZIP_NAME"

rm -f "$ZIP_PATH"
zip -j "$ZIP_PATH" \
  "$HEALTH_SVG" \
  "$HEALTH_HTML" \
  "$QUEUE_JSON" \
  "$QUEUE_HTML" \
  "$COMPLIANCE_JSON" \
  ${per_root_copies:+${per_root_copies[@]}} \
  ${SECTION_DIFF:+$SECTION_DIFF} >/tmp/release_zip.log 2>&1

SHA=$(sha256sum "$ZIP_PATH" | awk '{print $1}')
SIZE=$(stat -c%s "$ZIP_PATH")

{
  echo "# Release artifact bundle"
  echo "Zip: $ZIP_PATH"
  echo "SHA256: $SHA"
  echo "Size: $SIZE bytes"
  echo "Included:"
  echo "- $HEALTH_SVG"
  echo "- $HEALTH_HTML"
  if [[ ${#per_root_copies[@]} -gt 0 ]]; then
    for c in "${per_root_copies[@]}"; do
      echo "- $(basename "$c")"
    done
  else
    echo "- (per-root health_report_latest) brak — pominięto"
  fi
  echo "- $QUEUE_JSON"
  echo "- $QUEUE_HTML"
  echo "- $COMPLIANCE_JSON"
  if [[ -n "$SECTION_DIFF" ]]; then
    echo "- $SECTION_DIFF"
  else
    echo "- (section diff) brak — pominięto"
  fi
} > "$LOG"

echo "Bundle created: $ZIP_PATH (SHA256 $SHA, size $SIZE)"
