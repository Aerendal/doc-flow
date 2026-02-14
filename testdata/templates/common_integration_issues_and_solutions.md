---
title: Common Integration Issues and Solutions
status: needs_content
---

# Common Integration Issues and Solutions

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Najczęstsze problemy przy integracjach systemów i jak je szybko diagnozować/naprawiać.

## Typowe problemy / rozwiązania
- **Brak idempotencji/duplikaty**: retry bez key → wprowadź idempotency key, upsert/merge.
- **Ordering i consistency**: eventy poza kolejnością → klucze partycjonowania, sequence, dedupe okna.
- **Schema drift**: zmiany pól → schema registry, versioning, kompatybilność wsteczna, kontraktowe testy.
- **Timeouty/latencja**: duże payloady/brak paginacji → paginacja, kompresja, asynchroniczne kolejki.
- **Rate limits/429**: brak backoff → exponential backoff, jitter, limity po stronie klienta.
- **Auth/keys**: wygasłe klucze → rotacja, monitoring expiry, mTLS.
- **Mapping błędów**: różne kody → ustalony format błędów, tłumaczenie na własne kody.
- **Data quality**: null/format → walidacja po obu stronach, kontrakty, monit na anomaly.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Idempotencja + retry/backoff.
- [ ] Paginacja/kompresja i limity.
- [ ] Schema registry/versioning i testy kontraktowe.
- [ ] Standaryzowany format błędów.
- [ ] Monitoring lag/error i expiry kluczy.
