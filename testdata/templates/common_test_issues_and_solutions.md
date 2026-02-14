---
title: Common Test Issues and Solutions
status: needs_content
---

# Common Test Issues and Solutions

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Typowe problemy w testach (unit/integration/e2e) i sposoby ich eliminacji.

## Typowe problemy / rozwiązania
- **Testy flakujące**: zależne od czasu/sieci → mock time/clock, stub external, retry with cap, izolacja danych.
- **Wolne testy**: zbyt ciężkie E2E → przenieś do integration, równoleglenie, selektywne suite, fixture reuse.
- **Niedeterministyczne dane**: losowość → seedy, deterministic fixtures.
- **Brak pokrycia**: luki → coverage raport, focus na krytyczne ścieżki.
- **Konflikty środowisk**: shared state → hermetyczne środowiska, ephemeral env, containers.
- **Niejasne asercje**: false negative → lepsze matchery, diagnozowalne komunikaty.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Flaky tests zidentyfikowane i naprawione (clock/mock).
- [ ] Czas trwania suite akceptowalny, równoleglenie.
- [ ] Deterministyczne dane/fixtures.
- [ ] Coverage na krytyczne ścieżki monitorowane.
- [ ] Środowiska testowe izolowane/ephemeral.
