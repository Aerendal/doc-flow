# Contributing

Dziekujemy za zainteresowanie projektem.

To repo jest utrzymywane bez SLA i bez gwarantowanego wsparcia. Pull Requesty sa mile widziane, ale review i akceptacja nie sa gwarantowane.

## Wymagania

- Go 1.25+
- `git`, `make` (opcjonalny)
- Build i testy dzialaja offline dzieki `vendor/` - nie potrzebujesz sieci dla standardowego buildu

## Srodowisko deweloperskie

```bash
# 1. Sklonuj repo
git clone https://github.com/Aerendal/doc-flow.git
cd doc-flow

# 2. Zbuduj binare
go build -mod=vendor -o build/docflow ./cmd/docflow

# 3. Uruchom testy
go test -mod=vendor ./internal/... ./pkg/... ./tests/...

# 4. Sprawdz kod
go vet -mod=vendor ./...
```

## Zmiany zaleznosci

Jezeli musisz zmienic zaleznosci (go.mod), wymagana jest siec:

```bash
go mod tidy && go mod vendor
```

Zawsze commituj zmiany w `vendor/` razem ze zmiana w `go.mod` / `go.sum`.

## Konwencja commitow

Format: `<typ>: <krotki opis>`

Typy:
- `fix:` naprawa bledu
- `feat:` nowa funkcjonalnosc
- `docs:` tylko dokumentacja
- `refactor:` zmiana wewnetrzna bez zmiany zachowania
- `ci:` zmiany w workflow/skryptach CI
- `chore:` utrzymanie, zaleznosci, konfiguracja

Przyklad: `fix: popraw kolejnosc wyjscia validate --strict`

## Jak proponowac zmiane

1. Sprawdz czy nie ma otwartego Issue lub Discussion na ten temat.
2. Dla bugfixow: upewnij sie ze masz minimalny repro i test regresji.
3. Dla feature'ow: opisz motywacje i acceptance criteria w Issue lub Discussion przed implementacja.
4. Otworz PR z wypelnionym templatem (`.github/PULL_REQUEST_TEMPLATE.md`).

## Wazne zasady

- Zachowaj deterministyczny output CLI - zmiany kolejnosci lub formatu wyjscia musza byc celowe i udokumentowane.
- Nie lamie kontraktu opisanego w `docs/CONTRACT.md` bez uprzedniej dyskusji.
- Testy musza przejsc: `go test -mod=vendor ./...`
- Nie commituj sekretow, tokenow ani prywatnych sciezek.

## Gdzie zadawac pytania

- Pytania o uzycie i pomysly: GitHub Discussions
  https://github.com/Aerendal/doc-flow/discussions
- Bugi: GitHub Issues (uzupelnij szablon)
  https://github.com/Aerendal/doc-flow/issues
- Bezpieczenstwo: tylko przez prywatny kanal
  https://github.com/Aerendal/doc-flow/security/policy

