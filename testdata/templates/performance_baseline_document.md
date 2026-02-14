---
title: Performance Baseline Document
status: needs_content
---

# Performance Baseline Document

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Ustalić bazową wydajność systemu/usługi (baseline), aby mierzyć regresje i skuteczność optymalizacji.

## Struktura sekcji (szkielet)
1. Zakres: komponenty, ścieżki krytyczne, środowisko, wersja.
2. Metryki bazowe: latency p50/p95/p99, throughput, error rate, zasoby (CPU/RAM/IO), koszty.
3. Metodologia pomiaru: narzędzia, scenariusze, dane testowe, powtarzalność.
4. Wyniki i wnioski: tabele/wykresy, bottlenecks, limity.
5. Progi i SLO: ustalone cele, budżety wydajności, alerty.
6. Utrzymanie: kiedy aktualizować baseline (release/infra zmiana), wersjonowanie, repo wyników.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Zakres i metodologia pomiaru opisane; dane powtarzalne.
- [ ] Metryki bazowe zebrane i wizualizowane.
- [ ] SLO/alerty ustalone na podstawie baseline.
- [ ] Wersjonowanie i harmonogram aktualizacji baseline opisane.
