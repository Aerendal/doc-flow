## Opis zmiany

<!-- Co robi ten PR? Dlaczego ta zmiana jest potrzebna? -->

## Typ zmiany

<!-- Zaznacz pasujące (usuń niepasujące) -->

- [ ] Bugfix (naprawa błędu, nie zmienia API)
- [ ] Feature (nowa funkcjonalność)
- [ ] Refactor (zmiana wewnętrzna, bez zmiany zachowania)
- [ ] Docs (tylko dokumentacja)
- [ ] CI/build (zmiany w pipeline lub skryptach)
- [ ] Inne: <!-- opisz -->

## Powiązane zgłoszenia

<!-- Closes #XX lub Refs #XX -->

## Determinizm i kompatybilność

- [ ] Output CLI pozostaje stabilny (lub zmiany są udokumentowane i celowe)
- [ ] Brak niedetministycznego zachowania (timestampy, losowość, niestabilna kolejność) chyba że celowo z opt-in flagą
- [ ] Wsteczna kompatybilność zachowana: flagi, formaty konfigów, układy plików

## Checklist

- [ ] Zbudowałem lokalnie: `go build -mod=vendor -o build/docflow ./cmd/docflow`
- [ ] Testy przeszły: `go test -mod=vendor ./internal/... ./pkg/... ./tests/...`
- [ ] `go vet -mod=vendor ./...` bez błędów
- [ ] Zaktualizowałem dokumentację (jeśli zmiana wpływa na zachowanie CLI)
- [ ] Nie wprowadzam zmian łamiących kontrakt z `docs/CONTRACT.md`
- [ ] Nie zawiera sekretów, tokenów ani ścieżek prywatnych

## Jak testować

<!-- Kroki do weryfikacji zmiany przez reviewera -->

1.
2.
3.

## Uwagi dla reviewera

<!-- Cokolwiek co ułatwi review: edge cases, znane ograniczenia, pytania -->
