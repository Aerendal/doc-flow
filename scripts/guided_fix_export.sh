#!/usr/bin/env bash
set -euo pipefail
OUT_DIR=${OUT_DIR:-dist/release}
ZIP_NAME=${ZIP_NAME:-guided_fix_patches.zip}
LOG=${LOG:-LOGS/guided_fix_export.md}
PATCH_DIR=${PATCH_DIR:-LOGS}
SINCE=""

mkdir -p "$OUT_DIR" "$(dirname "$LOG")"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --since) SINCE="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    --zip-name) ZIP_NAME="$2"; shift 2;;
    --patch-dir) PATCH_DIR="$2"; shift 2;;
    -h|--help)
      echo "usage: $0 [--since YYYYMMDD] [--log FILE] [--zip-name NAME] [--patch-dir DIR]" >&2
      exit 0
      ;;
    *) echo "unknown arg: $1" >&2; exit 1;;
  esac
done

patches=("$PATCH_DIR"/patch_guided_*.patch)
if [[ ! -e ${patches[0]} ]]; then
  {
    echo "# Guided fix export"; echo "Brak patchy w $PATCH_DIR"; echo "Result: SKIP";
  } > "$LOG"
  exit 0
fi

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
selected=()
for p in "${patches[@]}"; do
  base=$(basename "$p")
  if [[ -n "$SINCE" ]]; then
    stamp=$(echo "$base" | sed -n 's/patch_guided_\([0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]\).*/\1/p')
    if [[ -z "$stamp" ]]; then
      continue
    fi
    if [[ "$stamp" < "$SINCE" ]]; then
      continue
    fi
  fi
  selected+=("$p")
done

if [[ ${#selected[@]} -eq 0 ]]; then
  {
    echo "# Guided fix export"
    echo "Brak patchy po filtrze"
    echo "Filter --since: ${SINCE:-n/a}"
    echo "Result: SKIP"
  } > "$LOG"
  exit 0
fi

cp "${selected[@]}" "$TMP"/
ZIP_PATH="$OUT_DIR/$ZIP_NAME"
rm -f "$ZIP_PATH"
zip -j "$ZIP_PATH" "$TMP"/*.patch >/dev/null
SHA=$(sha256sum "$ZIP_PATH" | awk '{print $1}')
SIZE=$(stat -c%s "$ZIP_PATH")

{
  echo "# Guided fix export"
  echo "Patches total: ${#patches[@]}"
  echo "Patches selected: ${#selected[@]}"
  echo "Filter --since: ${SINCE:-n/a}"
  echo "Zip: $ZIP_PATH"
  echo "SHA256: $SHA"
  echo "Size: $SIZE"
  for p in "${selected[@]}"; do
    echo "- $(basename "$p")"
  done
} > "$LOG"

echo "Export written: $ZIP_PATH"
