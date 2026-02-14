---
title: Common Issues & Solutions
status: needs_content
---

# Common Issues & Solutions

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Uniwersalna lista najczęstszych problemów (cross‑domain) i szybkich rozwiązań/diagnostyki.

## Sekcje przykładowe
- **Wydajność**: wolne odpowiedzi → profil, cache, paginacja, indeksy.
- **Stabilność**: crashe → logi/crash reports, reprodukcja, feature flags.
- **Sieć**: 5xx/timeout → retry/backoff, circuit breaker, limit concurrency.
- **Bezpieczeństwo**: brak MFA/secrets → włącz MFA, secret store, rotacja.
- **Dane**: brak spójności → idempotencja, walidacja, monitoring DQ.
- **Deploy**: rollback → blue-green/canary, health checks, automatyczny rollback.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Logi/metryki/trace dostępne dla diagnozy.
- [ ] Retry/backoff/circuit breaker gdzie trzeba.
- [ ] Cache/paginacja dla ciężkich zapytań.
- [ ] MFA i secret store włączone.
- [ ] Monitoring DQ i alerty na błędy.
