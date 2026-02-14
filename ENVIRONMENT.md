# ENVIRONMENT — Specyfikacja środowiska docflow

## Język i wersja

- **Język:** Go 1.25.5
- **OS:** Linux (amd64)
- **Środowisko budowania:** NixOS / Replit

## Wymagania runtime

- Go >= 1.25
- System plików z obsługą UTF-8
- Brak wymagań CGO (pure Go)

## Zależności główne

| Pakiet | Wersja | Rola |
|--------|--------|------|
| `github.com/spf13/cobra` | v1.10.2 | Framework CLI |
| `modernc.org/sqlite` | v1.44.3 | SQLite (pure Go, bez CGO) |
| `gopkg.in/yaml.v3` | v3.0.1 | Parser YAML (frontmatter, config) |

## Struktura projektu

```
cmd/docflow/          - Punkt wejścia CLI
pkg/config/           - Loader konfiguracji (docflow.yaml)
internal/
  cli/                - Komendy CLI (cobra)
  db/                 - Warstwa bazy danych SQLite
  logger/             - Logger strukturalny
  model/              - Modele danych
  util/               - Narzędzia (file walker, helpers)
docs/                 - Dokumentacja projektu
  _meta/              - Kontrakty metadanych
  ADR/                - Architecture Decision Records
  days/               - Plan 90-dniowy
testdata/
  templates/          - Szablony dokumentów do przetworzenia
  ground_truth/       - Przykłady referencyjne
tools/
  bench/              - Benchmarki wydajnościowe
  mutations/          - Testy mutacyjne
  scripts/            - Skrypty pomocnicze
LOGS/                 - Dzienniki decyzji
```

## Konfiguracja

- Plik konfiguracji: `docflow.yaml` w katalogu roboczym
- Baza danych: SQLite w `.docflow/docflow.db`
- Cache: `.docflow/cache/`

## Konwencje

- Formatowanie kodu: `gofmt`
- Testy: `go test ./...`
- Build: `go build -o docflow ./cmd/docflow/`
- Linting: `go vet ./...`
