#!/usr/bin/env bash
set -euo pipefail
# Build + package + checksums + URL inject (dry-run upload)

usage() {
  cat <<'USAGE'
release_pipeline.sh --linux-url URL [--checksum-url URL] [--no-inject]

Steps:
 1) go build (linux-amd64) to dist/docflow-linux-amd64
 2) smoke: --version/--help/validate --help/compliance --help
 3) package -> dist/docflow-linux-amd64.tar.gz
 4) sha256sum archiwów -> dist/checksums.txt
 5) inject_release_url.sh (unless --no-inject)
 6) print next-step instructions for upload (dry-run)
USAGE
}

LINUX_URL=""
CHECKSUM_URL=""
INJECT=1
while [[ $# -gt 0 ]]; do
  case "$1" in
    --linux-url) LINUX_URL="$2"; shift 2;;
    --checksum-url) CHECKSUM_URL="$2"; shift 2;;
    --no-inject) INJECT=0; shift;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown option: $1" >&2; usage; exit 1;;
  esac
done

if [[ $INJECT -eq 1 && -z "$LINUX_URL" ]]; then
  echo "--linux-url required when inject is enabled" >&2
  usage
  exit 1
fi

repo_root=$(cd "$(dirname "$0")/.." && pwd)
cd "$repo_root"

export GOFLAGS=${GOFLAGS:-"-mod=vendor"}
export GOCACHE=${GOCACHE:-"/tmp/go-cache"}
VERSION=${VERSION:-dev}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}
DATE=${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}
LDFLAGS="-X docflow/internal/buildinfo.Version=${VERSION} -X docflow/internal/buildinfo.Commit=${COMMIT} -X docflow/internal/buildinfo.Date=${DATE}"
mkdir -p dist

echo "[release_pipeline] build..."
go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o dist/docflow-linux-amd64 ./cmd/docflow
chmod +x dist/docflow-linux-amd64

echo "[release_pipeline] smoke..."
./dist/docflow-linux-amd64 --version >/dev/null
./dist/docflow-linux-amd64 --help >/dev/null
./dist/docflow-linux-amd64 validate --help >/dev/null
./dist/docflow-linux-amd64 compliance --help >/dev/null

echo "[release_pipeline] package..."
tar -czf dist/docflow-linux-amd64.tar.gz -C dist docflow-linux-amd64

echo "[release_pipeline] checksums..."
find dist -maxdepth 1 -type f \( -name "*.tar.gz" -o -name "*.zip" \) -print0 \
  | sort -z \
  | xargs -0 sha256sum > dist/checksums.txt

if [[ $INJECT -eq 1 ]]; then
  ./scripts/inject_release_url.sh --linux "$LINUX_URL" ${CHECKSUM_URL:+--checksum "$CHECKSUM_URL"}
fi

echo "[release_pipeline] done"
echo "Artifacts: dist/docflow-linux-amd64, dist/docflow-linux-amd64.tar.gz, dist/checksums.txt"
echo "Next: upload to releases (manual)"
