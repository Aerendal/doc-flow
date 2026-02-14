#!/usr/bin/env bash
set -euo pipefail

DOCX="${1:-/home/jerzy/Pobrane/LocalFold-Explorer_Projekt_v2_ASUS_TUF_F15.docx}"
ROOT="tmp/localfold_import"
LOG="LOGS/demo_guided_fix.md"

mkdir -p "$(dirname "$LOG")" "$ROOT"

echo "# Demo guided fix" > "$LOG"
echo "Docx: $DOCX" >> "$LOG"
echo "Root: $ROOT" >> "$LOG"

# preflight
if ./scripts/demo_preflight.sh >> "$LOG" 2>&1; then
  echo "- preflight: OK" >> "$LOG"
else
  echo "- preflight: FAILED (see log above)" >> "$LOG"
  exit 1
fi

# 1) Import with merge/no-overwrite + patch/status
./scripts/import_docx.sh --in "$DOCX" --out-dir "$ROOT" --merge-sections --export-patch LOGS/import_docx_day200.patch --status-json LOGS/import_docx_day200.json
echo "- import: LOGS/import_docx_day200.json, LOGS/import_docx_day200.patch" >> "$LOG"

# 2) fix_flow idempotent + sections patch
./scripts/fix_flow.sh --root "$ROOT" --log LOGS/fix_flow.md --export-patch LOGS/fix_flow.patch --export-sections-patch LOGS/fix_flow_sections.patch --clean
echo "- fix_flow: LOGS/fix_flow.md, LOGS/fix_flow.patch, LOGS/fix_flow_sections.patch" >> "$LOG"

# 3) Audit/health (one-liner)
./scripts/audit.sh --root "$ROOT" --log LOGS/audit.md
echo "- audit: LOGS/audit.md, dist/health/localfold_import/*latest*" >> "$LOG"

# 4) Release artifacts bundle (health/queue/compliance)
./scripts/release_artifacts.sh --log LOGS/release_artifacts.md
echo "- release bundle: dist/release/docflow-artifacts.zip (log LOGS/release_artifacts.md)" >> "$LOG"

# 5) PR bundle (patch + reports)
./scripts/pr_bundle.sh --log LOGS/pr_bundle.md
echo "- pr bundle: dist/release/docflow-pr-bundle.zip (log LOGS/pr_bundle.md)" >> "$LOG"

echo "Status: DONE" >> "$LOG"
echo "Artifacts:" >> "$LOG"
echo "- import: LOGS/import_docx_day200.json, LOGS/import_docx_day200.patch" >> "$LOG"
echo "- fix_flow: LOGS/fix_flow.md, LOGS/fix_flow.patch, LOGS/fix_flow_sections.patch" >> "$LOG"
echo "- audit: dist/health/localfold_import/health_latest.svg/html, queue_latest.html/json, compliance_latest.json" >> "$LOG"
echo "- release bundle: dist/release/docflow-artifacts.zip" >> "$LOG"
echo "- pr bundle: dist/release/docflow-pr-bundle.zip" >> "$LOG"

echo "Demo guided fix completed. See $LOG for summary."
