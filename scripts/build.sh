#!/usr/bin/env bash
set -euo pipefail

OUT=build
mkdir -p "$OUT"

VERSION=${VERSION:-dev}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}
DATE=${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}
LDFLAGS="-X docflow/internal/buildinfo.Version=${VERSION} -X docflow/internal/buildinfo.Commit=${COMMIT} -X docflow/internal/buildinfo.Date=${DATE}"

build_one() {
  local goos=$1 goarch=$2
  local suffix=""
  if [ "$goos" = "windows" ]; then suffix=".exe"; fi
  echo "Building $goos/$goarch..."
  GOOS=$goos GOARCH=$goarch GOFLAGS=-mod=vendor GOCACHE=/tmp/go-cache go build -ldflags "$LDFLAGS" -o "$OUT/docflow-${goos}-${goarch}${suffix}" ./cmd/docflow
}

build_one linux amd64
build_one linux arm64
build_one darwin amd64
build_one darwin arm64
build_one windows amd64

echo "Artifacts in $OUT/"
