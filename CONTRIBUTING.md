# Contributing

Dziękuję za zainteresowanie projektem.

To repo jest utrzymywane bez SLA i bez gwarantowanego wsparcia. Pull Requesty są mile widziane, ale review i akceptacja nie są gwarantowane.

## Wymagania

- Go 1.25+
- `git`, `make` (opcjonalny)
- Build i testy działają offline dzięki `vendor/` - nie potrzebujesz sieci dla standardowego buildu

## Środowisko deweloperskie

```bash
# 1. Sklonuj repo
git clone https://github.com/Aerendal/doc-flow.git
cd doc-flow

# 2. Zbuduj binarę
go build -mod=vendor -o build/docflow ./cmd/docflow

# 3. Uruchom testy
go test -mod=vendor ./internal/... ./pkg/... ./tests/...

# 4. Sprawdź kod
go vet -mod=vendor ./...
```

## Zmiany zależności

Jeżeli musisz zmienić zależności (go.mod), wymagana jest sieć:

```bash
go mod tidy && go mod vendor
```

Zawsze commituj zmiany w `vendor/` razem ze zmianą w `go.mod` / `go.sum`.

## Konwencja commitów

Format: `<typ>: <krótki opis>`

Typy:
- `fix:` naprawa błędu
- `feat:` nowa funkcjonalność
- `docs:` tylko dokumentacja
- `refactor:` zmiana wewnętrzna bez zmiany zachowania
- `ci:` zmiany w workflow/skryptach CI
- `chore:` utrzymanie, zależności, konfiguracja

Przykład: `fix: popraw kolejność wyjścia validate --strict`

## Jak proponować zmianę

1. Sprawdź czy nie ma otwartego Issue lub Discussion na ten temat.
2. Dla bugfixów: upewnij się że masz minimalny repro i test regresji.
3. Dla nowych funkcji: opisz motywację i acceptance criteria w Issue lub Discussion przed implementacją.
4. Otwórz PR z wypełnionym templatem (`.github/PULL_REQUEST_TEMPLATE.md`).

## Ważne zasady

- Zachowaj deterministyczny output CLI - zmiany kolejności lub formatu wyjścia muszą być celowe i udokumentowane.
- Nie łam kontraktu opisanego w `docs/CONTRACT.md` bez uprzedniej dyskusji.
- Testy muszą przejść: `go test -mod=vendor ./...`
- Nie commituj sekretów, tokenów ani prywatnych ścieżek.

## Gdzie zadawać pytania

- Pytania o użycie i pomysły: GitHub Discussions
  https://github.com/Aerendal/doc-flow/discussions
- Błędy: GitHub Issues (uzupełnij szablon)
  https://github.com/Aerendal/doc-flow/issues
- Bezpieczeństwo: tylko przez prywatny kanał
  https://github.com/Aerendal/doc-flow/security/policy


