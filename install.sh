#!/usr/bin/env bash
set -euo pipefail

PREFIX="${PREFIX:-$HOME/.local}"
BIN="${PREFIX}/bin"
SHARE="${PREFIX}/share"
COMPL="$SHARE/docflow-completions"
SRC_BIN="./build/docflow"
LOG="LOGS/INSTALL_DAY_187.md"
DRY_RUN=0
FROM=""
SHELL_DETECT=""
CHANNEL=""
VERSION=""
CHECKSUMS="dist/checksums.txt"
INSTALL_FZF=0

usage() {
  echo "usage: PREFIX=/desired/prefix ./install.sh [--from url|path] [--channel latest] [--version TAG] [--install-fzf] [--dry-run]" >&2
  echo "default source: ./build/docflow" >&2
  echo "--channel latest expects local dist/ artifacts produced by release tooling" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --from) FROM="$2"; shift 2;;
    --channel) CHANNEL="$2"; shift 2;;
    --version) VERSION="$2"; shift 2;;
    --dry-run) DRY_RUN=1; shift;;
    --install-fzf) INSTALL_FZF=1; shift;;
    -h|--help) usage;;
    *) usage;;
  esac
done

mkdir -p "$BIN" "$COMPL" LOGS

install_bin() {
  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[DRY-RUN] Would install to $BIN/docflow"
    return 0
  fi
  if [[ -n "$FROM" ]]; then
    if [[ "$FROM" =~ ^https?:// ]]; then
      curl -L "$FROM" -o "$BIN/docflow" >/dev/null 2>&1
    else
      cp "$FROM" "$BIN/docflow"
    fi
  else
    cp "$SRC_BIN" "$BIN/docflow"
  fi
  chmod +x "$BIN/docflow"
}

install_compl() {
  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[DRY-RUN] Would copy completions to $COMPL"
    return 0
  fi
  cp completions/docflow.bash "$COMPL/"
  cp completions/docflow.zsh "$COMPL/"
  cp completions/docflow.fish "$COMPL/"
}

install_fzf_opt() {
  if [[ "$INSTALL_FZF" -ne 1 ]]; then
    return 0
  fi
  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[DRY-RUN] Would install fzf (optional)"
    return 0
  fi
  if command -v fzf >/dev/null 2>&1; then
    echo "fzf already installed"
    return 0
  fi
  if command -v apt-get >/dev/null 2>&1; then
    echo "Installing fzf via apt-get"
    sudo apt-get update -y && sudo apt-get install -y fzf || true
  elif command -v brew >/dev/null 2>&1; then
    echo "Installing fzf via brew"
    brew install fzf || true
  else
    echo "[WARN] fzf not installed; install manually (apt-get install fzf / brew install fzf)"
  fi
}

append_rc() {
  local shell_rc="$1" snippet="$2"
  [[ ! -f "$shell_rc" ]] && return 0
  if ! grep -q "$snippet" "$shell_rc"; then
    echo "$snippet" >> "$shell_rc"
  fi
}

verify_checksum() {
  local file="$1"
  local checksums="$2"
  local base
  base=$(basename "$file")
  if [[ ! -f "$checksums" ]]; then
    echo "checksums file not found: $checksums" >&2
    return 1
  fi
  local expected
expected=$(grep -E " (${base}|.*/${base})$" "$checksums" | head -n1 | awk '{print $1}')
if [[ -z "$expected" ]]; then
  echo "no checksum entry for $base in $checksums" >&2
  return 1
fi
  local actual
  actual=$(sha256sum "$file" | awk '{print $1}')
  if [[ "$actual" != "$expected" ]]; then
    echo "checksum mismatch for $base: expected $expected got $actual" >&2
    return 1
  fi
  return 0
}

# channel latest -> use local dist/ artefakty (nie GitHub API)
if [[ "$CHANNEL" == "latest" && -z "$FROM" ]]; then
  FROM="dist/docflow-linux-amd64"
fi

# version tag -> expect dist/releases/<tag>/docflow-linux-amd64 and checksums.txt
if [[ -n "$VERSION" ]]; then
  CHECKSUMS="dist/releases/$VERSION/checksums.txt"
  if [[ -z "$FROM" ]]; then
    FROM="dist/releases/$VERSION/docflow-linux-amd64"
  fi
fi

# verify before install if possible
checksum_result="skipped"
if [[ -f "$FROM" && -f "$CHECKSUMS" ]]; then
  if verify_checksum "$FROM" "$CHECKSUMS"; then
    checksum_result="pass"
  else
    checksum_result="fail"
    if [[ "$DRY_RUN" -eq 0 ]]; then
      echo "Checksum verification failed" >&2
      exit 1
    fi
  fi
fi

if [[ "$DRY_RUN" -eq 0 && -n "$FROM" && ! "$FROM" =~ ^https?:// && ! -f "$FROM" ]]; then
  echo "source binary not found: $FROM" >&2
  echo "build first: go build -mod=mod -o build/docflow ./cmd/docflow" >&2
  echo "or provide an explicit path/url with --from" >&2
  exit 1
fi

if [[ "$DRY_RUN" -eq 0 ]]; then
  install_bin
  install_compl
  install_fzf_opt
fi

SHELL_DETECT="${SHELL##*/}"
case "$SHELL_DETECT" in
  bash) SNIPPET="source $COMPL/docflow.bash" ;;
  zsh)  SNIPPET="source $COMPL/docflow.zsh" ;;
  fish) SNIPPET="source $COMPL/docflow.fish" ;;
  *)    SNIPPET="source $COMPL/docflow.bash" ;;
esac

RC_FILES=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.config/fish/config.fish")
if [[ "$DRY_RUN" -eq 0 ]]; then
  for rc in "${RC_FILES[@]}"; do
    append_rc "$rc" "$SNIPPET" 2>/dev/null || true
  done
fi

{
  echo "# Install"
  echo "Prefix: $PREFIX"
  echo "Bin: $BIN/docflow"
  echo "Completions: $COMPL"
  echo "From: ${FROM:-local build}"
  echo "Dry run: $DRY_RUN"
  echo "Shell detected: $SHELL_DETECT"
  echo "Snippet: $SNIPPET"
  echo "Channel: ${CHANNEL:-manual}"
  echo "Checksum: $checksum_result"
} > "$LOG"

if [[ "$DRY_RUN" -eq 0 ]]; then
  echo "Installed to $BIN/docflow"
  echo "Completions in $COMPL (auto-source attempted; snippet: $SNIPPET)"
else
  echo "[DRY-RUN] Would install to $BIN/docflow"
fi
