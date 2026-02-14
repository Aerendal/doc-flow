# Troubleshooting (MVP)

## Build/Test
- `permission denied ~/.cache/go-build`: ustaw `GOCACHE=/tmp/go-cache` lub `HOME=/tmp`.
- `cannot find module providing ... -mod=vendor`: uruchom `go mod download` z dostępem do sieci albo użyj binarki z Releases.
- Nieaktualny cache modułów: wyczyść cache (`go clean -modcache`) i ponów `go mod download`.

## Walidacja
- `missing_frontmatter`: dodaj blok `--- ... ---` na początku pliku.
- `missing_context_sources` w draft: włącz `--status-aware` lub dodaj context_sources.
- `legacy_section_name`: odpal `migrate-sections` lub zaktualizuj nagłówki wg section_aliases.
- Fałszywe `governance` braki pól: upewnij się, że frontmatter zawiera pola i jest poprawnie wcięty; heurystyka szuka `field:` w tekście.

## Compliance
- Raport pokazuje brak sekcji mimo obecności: sprawdź aliasy i poziomy nagłówków (H1 jest pomijane).
- HTML raport pusty: upewnij się, że `docflow compliance --rules ...` wskazuje na prawidłowy plik i indeks obejmuje dokumenty.

## CI/Release
- Artefakty nie budują się: upewnij się, że runner ma dostęp do modułów Go (sieć lub preseedowany `GOMODCACHE`).
- Brak checksumów w Homebrew formula: uzupełnij SHA256 przy pierwszym wydaniu.
- Release changelog pusty: aktualny workflow bierze ostatni commit — dodaj generator changelogów.

## Performance
- Benchmark timeouts: uruchamiaj z `GOCACHE=/tmp/go-cache`; redukuj zakres `-bench`.

## Narzędzia
- `docflow` not found po instalacji: sprawdź PATH (`/usr/local/bin` lub dest w install.sh).
- Pre-commit nie działa: upewnij się, że `./docflow` jest zbudowany i hook zainstalowany (`scripts/install-hooks.sh`).
