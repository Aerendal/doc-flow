#!/usr/bin/env bash
set -euo pipefail

HOOK=".git/hooks/pre-commit"

cat >"$HOOK" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

# Go fmt/vet
if command -v go >/dev/null 2>&1; then
  GOFLAGS=-mod=vendor go fmt ./...
  GOFLAGS=-mod=vendor go vet ./...
fi

# Validate staged Markdown docs
files=$(git diff --cached --name-only -- '*.md')
BIN=""
if [ -x "./build/docflow-linux-amd64" ]; then
  BIN="./build/docflow-linux-amd64"
elif [ -x "./docflow" ]; then
  BIN="./docflow"
fi
if [ -n "$files" ] && [ -n "$BIN" ]; then
  "$BIN" validate --warn --status-aware --governance docs/_meta/GOVERNANCE_RULES.yaml || exit 1
fi
EOF

chmod +x "$HOOK"
echo "Pre-commit hook installed to $HOOK"
