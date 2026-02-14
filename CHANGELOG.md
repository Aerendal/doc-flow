# Changelog

## v1.0.0-rc2 → GA (2026-02-09)
- Added: testy wydajności 1k/duże pliki i raporty RSS; sanity chaos/security (YAML bomb, symlinks, traversal blokada).
- Fixed: `scan <dir>` respektuje ścieżkę; walidator blokuje path traversal w depends_on/context_sources; governance reporter domyślny status w compliance; cykle wykrywane i raportowane w validate; queue Go/No-Go (text/json, cache, workers).
- Added: observability flags (`--cpu-profile`, `--mem-profile`, `--log-format json`); queue workflow z cache/logami; governance-ready examples; perf 10k (index~0.56s, validate~0.9s, RSS~80MB).
- Tests/QA: e2e (1k + large + 10k), chaos/security, bug bash RC1 (P0=0, P1 naprawiony), smoke RC2, queue smoke.
- Known: README quickstart wciąż czeka na realny URL release; queue violations JSON jeszcze puste; rekomendacje/templates nadal dane demo.

## v0.1.0-mvp (2026-02-09)
- Walidator metadanych (doc_id, context_sources, expected deps, status-aware)
- SectionTree parser + Section metrics (completeness)
- Generator manifestu i prosty generator szablonów z preview
- Rekomendacje szablonów (MVP), zestawy (co-occurrence), plan dzienny, template impact
- Config: promote_context_for, family_rules_path
- CLI: validate, plan daily, recommend (demo), template-sets (demo), templates list/deprecated (demo), template-impact
- Testy `GOFLAGS=-mod=vendor go test ./...` (Go 1.25.7)
