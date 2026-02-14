#!/usr/bin/env bash
set -euo pipefail

LOGFILE="LOGS/release_checklist.md"

required=()
warn=()

add_fail() { required+=("$1"); }
add_warn() { warn+=("$1"); }

check_file() {
  local path="$1" label="$2"
  if [ -f "$path" ]; then
    echo "- [PASS] $label ($path)"
  else
    echo "- [FAIL] $label ($path missing)"
    add_fail "$label"
  fi
}

check_readme_url() {
  local label="$1" pattern="$2"
  if grep -q "$pattern" README.md; then
    echo "- [PASS] $label (README)"
  else
    echo "- [FAIL] $label (missing in README)"
    add_fail "$label"
  fi
}

check_checksums_entry() {
  local file="$1" label="$2"
  if [ -f dist/checksums.txt ] && grep -q "$file" dist/checksums.txt; then
    echo "- [PASS] $label in checksums.txt"
  else
    echo "- [FAIL] $label missing in checksums.txt"
    add_fail "$label"
  fi
}

check_perf_logs() {
  local hits=0
for p in LOGS/queue_html.md LOGS/scale_baseline_10k.md LOGS/scale_baseline_1k.md; do
    if [ -f "$p" ]; then hits=1; break; fi
  done
  if [ $hits -eq 1 ]; then
    echo "- [PASS] Perf/queue logs present"
  else
    echo "- [WARN] Perf/queue logs missing"
    add_warn "perf_logs"
  fi
}

mkdir -p LOGS
{
  echo "# Release checklist"
  echo
  check_file "dist/docflow-linux-amd64" "Binary built"
  check_file "dist/checksums.txt" "Checksums file"
  check_checksums_entry "docflow-linux-amd64" "Binary checksum recorded"
  check_readme_url "Release URL present" "releases/latest/download/docflow-linux-amd64"
  check_readme_url "Checksums URL present" "releases/latest/download/checksums.txt"
  check_perf_logs
} > "$LOGFILE"

if [ ${#required[@]} -eq 0 ]; then
  echo "- EXIT: PASS" >> "$LOGFILE"
  exit 0
else
  echo "- EXIT: FAIL (missing ${#required[@]} items)" >> "$LOGFILE"
  exit 1
fi
