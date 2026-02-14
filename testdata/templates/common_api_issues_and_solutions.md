---
title: Common API Issues and Solutions
status: needs_content
---

# Common API Issues and Solutions

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Lista typowych problemów z API i sposobów ich rozwiązania, żeby skrócić czas diagnozy.

## Typowe problemy / rozwiązania
- **Auth/401/403**: brak/expired token, zły scope → odśwież token, sprawdź scope/role.
- **Rate limit/429**: brak retry/backoff → dodaj exponential backoff, request‑id, optymalizuj batch.
- **Timeouty/latencja**: duże payloady, brak paginacji → włącz paginację, kompresję, zoptymalizuj zapytania.
- **Schema/validation**: breaking change, brak pola → sprawdź wersję, użyj feature flag/compat warstw.
- **Idempotencja**: duplikaty przy retry → użyj idempotency key, zaprojektuj PUT/POST zgodnie z idempotencją.
- **CORS**: brak nagłówków → ustaw `Access-Control-*` zgodnie z polityką.
- **Version mismatch**: stary klient → komunikacja deprecacji, fallback, semver w ścieżce/nagłówku.
- **Error handling**: niespójne kody → ustandaryzuj format błędów (JSON, code/message/request-id).
- **Security**: brak input validation/OWASP → waliduj, limituj, loguj z request-id, skanuj zależności.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Paginacja/kompresja i timeouty ustawione.
- [ ] Rate limiting + retry/backoff po stronie klienta.
- [ ] Idempotency key na operacjach wrażliwych.
- [ ] Spójny format błędów i kody.
- [ ] Proces deprecacji/wersjonowania.
