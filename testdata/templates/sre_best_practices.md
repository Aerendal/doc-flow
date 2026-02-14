---
title: SRE Best Practices
status: needs_content
---

# SRE Best Practices

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Zbiór najlepszych praktyk SRE.

## Struktura sekcji (szkielet)
1. SLO/SLI i error budgets.
2. On-call i incident response.
3. Obserwowalność: logi/metyki/trace.
4. Reliability by design: capacity, chaos, rollouts.
5. Toil i automatyzacja.
6. Postmortems i nauki.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] SLO/SLI opisane i stosowane.
- [ ] On-call/runbooki aktualne.
- [ ] Obserwowalność i chaos praktykowane.
- [ ] Toil mierzony i redukowany.
