#!/usr/bin/env bash
set -euo pipefail
ROOT=${ROOT:-tmp/localfold_import}
LOG=${LOG:-LOGS/guided_fix_patch.md}
SECTION_DIFF=${SECTION_DIFF:-LOGS/SECTION_DIFF_palette_bin.md}
PATCH_DIR=${PATCH_DIR:-LOGS}

ROOT="$ROOT" LOG="$LOG" SECTION_DIFF="$SECTION_DIFF" PATCH_DIR="$PATCH_DIR" ./scripts/guided_fix_context.sh

echo "Patch log written to $LOG"
