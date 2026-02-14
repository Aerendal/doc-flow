---
title: Resilience Testing
status: needs_content
---

# Resilience Testing

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Plan testów odporności (chaos/game days).

## Struktura sekcji (szkielet)
1. Cele i hipotezy.
2. Scenariusze: awarie, degradacje, bezpieczeństwo.
3. Środowisko: prod/pre-prod, zabezpieczenia.
4. Metryki sukcesu: SLO, czas odzysku.
5. Bezpieczeństwo testów: blast radius, rollback.
6. Raport i działania po testach.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Scenariusze zgodne z celami.
- [ ] Środowisko i blast radius uzgodnione.
- [ ] Metryki i plan rollback określone.
- [ ] Raport z wnioskami dostarczony.
