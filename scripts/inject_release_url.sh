#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
inject_release_url.sh --linux URL [--darwin URL] [--windows URL] [--checksum URL]

Podmienia placeholdery URL w README.md i docs/FRONT_OF_BOOK.md.

Argumenty:
  --linux     URL do binarki linux-amd64 (wymagane)
  --darwin    URL do binarki darwin-amd64/arm64 (opcjonalne)
  --windows   URL do binarki windows-amd64 (opcjonalne)
  --checksum  URL do pliku checksums.txt (opcjonalne)

Przykład:
  ./scripts/inject_release_url.sh \
    --linux https://github.com/org/repo/releases/download/v0.4/docflow-linux-amd64 \
    --darwin https://github.com/org/repo/releases/download/v0.4/docflow-darwin-amd64 \
    --windows https://github.com/org/repo/releases/download/v0.4/docflow-windows-amd64.exe \
    --checksum https://github.com/org/repo/releases/download/v0.4/checksums.txt

Idempotentne: wielokrotne uruchomienie z tymi samymi URL nie zmieni plików.
USAGE
}

LINUX=""
DARWIN=""
WINDOWS=""
CHECKSUM=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --linux) LINUX="$2"; shift 2;;
    --darwin) DARWIN="$2"; shift 2;;
    --windows) WINDOWS="$2"; shift 2;;
    --checksum) CHECKSUM="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
  esac
done

if [[ -z "$LINUX" ]]; then
  echo "--linux jest wymagane" >&2
  usage
  exit 1
fi

repo_root="$(cd "$(dirname "$0")/.." && pwd)"

replace_url() {
  local file="$1" placeholder="$2" url="$3"
  if [[ -z "$url" ]]; then return; fi
  if grep -q "$placeholder" "$file"; then
    sed -i "s#${placeholder}#${url}#g" "$file"
  elif grep -q "$url" "$file"; then
    : # już wstawione
  else
    # dodaj informacyjnie na końcu sekcji Quickstart/Download
    echo "" >> "$file"
    echo "Download: ${url}" >> "$file"
  fi
}

replace_url "$repo_root/README.md" "https://example.invalid/docflow-linux-amd64" "$LINUX"
replace_url "$repo_root/docs/FRONT_OF_BOOK.md" "https://example.invalid/docflow-linux-amd64" "$LINUX"

if [[ -n "$DARWIN" ]]; then
  replace_url "$repo_root/README.md" "https://example.invalid/docflow-darwin-amd64" "$DARWIN"
fi
if [[ -n "$WINDOWS" ]]; then
  replace_url "$repo_root/README.md" "https://example.invalid/docflow-windows-amd64.exe" "$WINDOWS"
fi
if [[ -n "$CHECKSUM" ]]; then
  replace_url "$repo_root/README.md" "https://example.invalid/checksums.txt" "$CHECKSUM"
fi

echo "[inject_release_url] updated README.md / FRONT_OF_BOOK.md"
