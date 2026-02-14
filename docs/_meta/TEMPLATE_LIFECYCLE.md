# Template Lifecycle

## Stany
- `draft` — nowy/niegotowy szablon
- `active` — używany szablon
- `deprecated` — do zastąpienia, nie używać w nowych dokumentach
- `archived` — wyłączony z użycia

## Reguły przejść
- draft → active: manual (po review)
- active → deprecated: quality < 50 **i** usage == 0 przez >= 90 dni (lub ręcznie)
- deprecated → archived: manual (po migracji)
- active → archived: manual (wyjątek, jeśli zastąpiony i nieużywany)

## Parametry
- Quality threshold: 50
- Usage window: 90 dni (brak użycia)

## Działanie
- Lifecycle manager rekomenduje przejścia; wykonanie przejścia wymaga potwierdzenia.
