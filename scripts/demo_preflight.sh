#!/usr/bin/env bash
set -euo pipefail

DOCX="${DOCX:-/home/jerzy/Pobrane/LocalFold-Explorer_Projekt_v2_ASUS_TUF_F15.docx}"
BIN="./build/docflow-linux-amd64"
LOG="LOGS/demo_preflight.md"

missing=()

if [[ ! -f "$DOCX" ]]; then missing+=("docx:$DOCX"); fi
if [[ ! -x "$BIN" ]]; then missing+=("bin:$BIN"); fi
for cmd in python3; do
  command -v "$cmd" >/dev/null 2>&1 || missing+=("cmd:$cmd")
done

warns=()
command -v fzf >/dev/null 2>&1 || warns+=("fzf not installed (optional)")

mkdir -p "$(dirname "$LOG")"
{
  echo "# Demo preflight"
  echo "DOCX: $DOCX"
  echo "BIN: $BIN"
  if [[ ${#missing[@]} -eq 0 ]]; then
    echo "Status: OK"
  else
    echo "Status: MISSING ${#missing[@]}"
    printf '%s\n' "${missing[@]}"
  fi
  if [[ ${#warns[@]} -gt 0 ]]; then
    echo "Warnings:"
    printf '%s\n' "${warns[@]}"
  fi
} > "$LOG"

if [[ ${#missing[@]} -gt 0 ]]; then
  echo "Preflight failed. See $LOG"
  exit 1
fi

echo "Preflight OK. See $LOG"
