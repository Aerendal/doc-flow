# TEST_STRATEGY — docflow

## Framework
- Go standard library `testing`.
- Pokrycie: unit + lekkie integracyjne (mały fixture).
- Uruchomienie: `GOFLAGS=-mod=vendor go test ./...`

## Konwencje
- Pliki testów: `*_test.go`, nazwy funkcji `TestXxx`.
- Testy deterministyczne: sortowanie z tie-break doc_id, brak zależności od czasu/środowiska.
- Seed losowy — niewykorzystywany (brak rng w kodzie).

## Pokrycie minimalne (MVP)
- Parser: frontmatter (valid/invalid), headings, section tree.
- Validator: brak frontmatter, brak pól, duplikat doc_id, missing context (warn/error), missing expected deps, published empty sections.
- Graph: topo sort, cykl, missing nodes.
- Plan daily: topo kolejność, limit.
- Rekomendacje: ranking wg wag.
- Generator/manifest: wczytanie manifestu, generowanie z placeholderami.
- Sets: współwystępowanie.

## Integracja (mini pipeline)
- Fixture: mały katalog 3-5 plików MD (ground_truth) — używany w istniejących testach parser/frontmatter.
- Warunek: testy muszą przechodzić bez dostępu do sieci (vendor).

## Determinizm
- Topo sort i kolejki sortowane alfabetycznie (doc_id/path) po każdym dodaniu.
- Raporty CLI w kolejności deterministycznej.
- Brak losowania / czas w logice biznesowej.

## Coverage cel
- Minimum: wszystkie główne pakiety mają co najmniej 1 test sensowny scenariusz.
- Docelowo: >70% linii w pakietach core (parser, validator, deps, plan, generate, recommend).
