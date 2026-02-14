#!/usr/bin/env bash
set -euo pipefail

LOG="LOGS/palette_smoke.md"

mkdir -p "$(dirname "$LOG")"

echo "# Palette smoke" > "$LOG"

echo "- palette.sh (list mode)" >> "$LOG"
./scripts/palette.sh <<< "" >> "$LOG" 2>&1 || true

echo "- docflow palette (no selection)" >> "$LOG"
./build/docflow-linux-amd64 palette <<< "" >> "$LOG" 2>&1 || true

echo "- mini-audit before guard (localfold_import)" >> "$LOG"
./scripts/audit.sh --root tmp/localfold_import --log LOGS/audit_smoke.md >> "$LOG" 2>&1 || true

echo "- order_guard (non-strict)" >> "$LOG"
./scripts/order_guard.sh --bundle dist/release/docflow-pr-bundle.zip --health dist/health/localfold_import/health_latest.html --queue dist/health/localfold_import/queue_latest.json --compliance dist/health/localfold_import/compliance_latest.json --log LOGS/order_guard_smoke.md --no-strict >> "$LOG" 2>&1 || true

echo "- artifacts_check --strict" >> "$LOG"
./scripts/artifacts_check.sh --root localfold_import --strict --log LOGS/artifacts_check_smoke.md >> "$LOG" 2>&1 || true

echo "Smoke done" >> "$LOG"
