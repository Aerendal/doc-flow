#!/usr/bin/env bash
set -euo pipefail

ROOT="docs"
STRICT=0
LOG="LOGS/artifacts_check.md"

usage(){ echo "usage: $0 [--root dir] [--strict] [--log file]" >&2; exit 1; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --strict) STRICT=1; shift;;
    --log) LOG="$2"; shift 2;;
    -h|--help) usage;;
    *) usage;;
  esac
done

mkdir -p "$(dirname "$LOG")"

checks=(
  "dist/health/health_latest.svg"
  "dist/health/health_latest.html"
  "dist/health/${ROOT}/health_latest.svg"
  "dist/health/${ROOT}/queue_latest.json"
  "dist/health/${ROOT}/queue_latest.html"
  "dist/health/${ROOT}/compliance_latest.json"
  "dist/release/docflow-artifacts.zip"
  "dist/release/docflow-pr-bundle.zip"
  "LOGS/audit.md"
  "LOGS/demo_guided_fix.md"
)

missing=()
for f in "${checks[@]}"; do
  if [[ ! -f "$f" ]]; then
    missing+=("$f")
  fi
done

{
  echo "# Artifacts check"
  echo "Root: $ROOT"
  echo "Strict: $STRICT"
  if [[ ${#missing[@]} -eq 0 ]]; then
    echo "Status: OK (all required artifacts present)"
  else
    echo "Status: MISSING ${#missing[@]}"
    printf '%s\n' "${missing[@]}"
  fi
} > "$LOG"

if [[ ${#missing[@]} -gt 0 && $STRICT -eq 1 ]]; then
  exit 1
fi

echo "Artifacts check complete. See $LOG"
