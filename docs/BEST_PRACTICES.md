# Best Practices (MVP)

## Struktura repo
- Trzymaj metadane/schematy w `docs/_meta/`.
- Cache i indeks w `.docflow/cache/`; nie commituj cache (opcjonalnie .gitignore).
- Szablony w `templates/` lub `docs/templates/`; aliasy w configu.

## Frontmatter
- Używaj `doc_id` = snake_case, zgodny z nazwą pliku.
- W `depends_on` podawaj doc_id, nie ścieżki.
- `context_sources` dodawaj dla wszystkich statusów (łatwiej o governance pass).
- `version` trzymaj w semver (`v1.2.3`).

## Sekcje
- Kanoniczne nagłówki zgodne z rodziną dokumentu; aliasy tylko tymczasowo.
- Unikaj pustych sekcji w published; sprawdzaj `docflow validate --status-aware`.

## Szablony
- Dodawaj `template_source` w dokumentach generowanych z szablonu (pomaga w template-impact).
- Wersjonuj szablony (v1, v2) i deprecjonuj stare (`templates deprecated`).

## Zależności
- Utrzymuj małe, acykliczne grafy; wykrywaj zmiany `docflow changes --old-index`.
- Dla governance dodawaj required deps per rodzina w rules/family.

## CI/CD
- Używaj `GOFLAGS=-mod=vendor` i `GOCACHE=/tmp/go-cache` na CI.
- Benchmarks jako artefakt; smoke-test binarki przed release.

## Dokumentacja
- Aktualizuj `GOVERNANCE_RULES.yaml` wraz ze zmianami typów dokumentów.
- Przypominaj w README/TROUBLESHOOTING, że build ze źródeł wymaga modułów (sieć lub przygotowany cache).

## Migracje
- `migrate-sections` uruchamiaj na branchu roboczym; review diffów przed merge.
- Wersje dokumentów bumpuj automatycznie (`--auto-bump`) przy zmianach treści/metadanych.
