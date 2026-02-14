# Contributing

## Pre-commit hook

Install local hook (Go + docflow validation):
```bash
chmod +x scripts/install-hooks.sh
./scripts/install-hooks.sh
```
Hook runs:
- `go fmt ./...` (GOFLAGS=-mod=vendor)
- `go vet ./...`
- `./build/docflow-linux-amd64 validate --warn --status-aware --governance docs/_meta/GOVERNANCE_RULES.yaml` na staged plikach `.md` (fallback: `./docflow`).

Jeśli nie masz zbudowanego `./docflow`, hook pominie walidację.

## CI
- GitHub Actions workflow `.github/workflows/ci.yml` (test matrix 1.23–1.25, benchmark job).
- Ustaw `GOCACHE=/tmp/go-cache` w lokalnych testach gdy widzisz `permission denied` w `~/.cache/go-build`.

## Style
- Go modules przez `go.mod/go.sum` (`GOFLAGS=-mod=vendor`).
- Build/test ze źródeł na CI wymaga dostępu do modułów (sieć lub preseedowany cache).
