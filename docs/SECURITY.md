# Security Policy (MVP)

## Reporting vulnerabilities
- Zgłoszenia: otwórz issue typu "Security" lub mail do maintainer (placeholder: security@example.com).
- Dołącz: wersję docflow, OS, minimalny repro.

## Scope
- Binaries `docflow-*` i źródła w tym repo.
- Brak gwarancji dla forków/custom buildów.

## Best practices for users
- Uruchamiaj `docflow` na zaufanych repo; waliduj wejściowe frontmatter/YAML (brak sandboxa).
- Ustaw `GOFLAGS=-mod=vendor`; build ze źródeł uruchamiaj z kontrolowanego źródła kodu i spójnym `vendor/`.
- Waliduj pliki pochodzące z zewnątrz (`docflow validate --strict`).
- Weryfikuj artefakty release wg `docs/SECURITY_VERIFICATION.md` (checksums/cosign/attestations/SBOM).

## Dependency security
- Dependencies kontroluj przez `go.mod/go.sum`; zalecane okresowe `govulncheck ./...` (wymaga dostępu do sieci).

## Supported versions
- MVP (v0.x): brak gwarancji LTS; update do najnowszej wersji.
