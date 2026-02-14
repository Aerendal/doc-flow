#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"

export GOFLAGS=-mod=vendor
export GOCACHE=${GOCACHE:-/tmp/go-cache}

BIN="./build/docflow-linux-amd64"
BACKLOG=/tmp/backlog_sample.txt

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bin) BIN="$2"; shift 2;;
    --backlog) BACKLOG="$2"; shift 2;;
    -h|--help) echo "usage: scripts/ci_meta.sh [--bin path] [--backlog file]"; exit 0;;
    *) echo "unknown arg $1"; exit 1;;
  esac
done

# 1) regenerate day index
scripts/days_index.sh

# 2) lint days
/usr/local/go/bin/go run ./cmd/worklog-lint-days --root worklog/days

# 3) queue smoke
scripts/queue_evaluate.sh --bin "$BIN" --workers 2 --format json "$BACKLOG"
