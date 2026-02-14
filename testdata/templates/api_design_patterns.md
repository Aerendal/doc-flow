---
title: API Design Patterns
status: needs_content
---

# API Design Patterns

## Metadane
- Właściciel: [Architecture/API Guild]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Katalog wzorców projektowych API (HTTP/gRPC/event) z zasadami użycia, kompromisami i przykładami. Ma ujednolicić projektowanie, poprawić spójność, bezpieczeństwo, wydajność i zgodność z politykami organizacji.

## Zakres i granice
- Obejmuje: style (REST, RPC, event‑driven), wersjonowanie, zasoby i kontrakty, paginację/filtry, błędy, idempotencję, bezpieczeństwo (authz, rate limit, replay), obserwowalność, DTO vs. entity, consistency, streaming/bulk, HATEOAS opcjonalnie.
- Poza zakresem: szczegółowe implementacje per język/framework, specyficzne API produktowe (opis w ADR/Spec).

## Wejścia i wyjścia
- Wejścia: standard organizacji (naming, bezpieczeństwo), wymogi biznesowe/latency, SLO, profile klientów (mobile/web/internal), zgodność (PII/PCI/HIPAA), wzorce istniejących usług.
- Wyjścia: zestaw wzorców z kryteriami wyboru, przykłady kontraktów, sekcja anty‑wzorców, checklista projektowa, odniesienia do security/observability/QA.

## Powiązania (meta)
- Key Documents: api_security_baseline, api_versioning_policy, error_handling_guidelines, observability_standards, change_management.
- Key Document Structures: wzorzec → kiedy użyć → kompromisy → przykład → kontrole.
- Document Dependencies: IAM/authz, logging/trace, schema registry/IDL, SLA/SLO.

## Zależności dokumentu
Wymaga aktualnych polityk bezpieczeństwa, nazewnictwa, limitów i standardów obserwowalności; repozytorium przykładów (OpenAPI/Proto) oraz decyzji wersjonowania. Bez nich DoR pozostaje otwarte.

## Powiązania sekcja↔sekcja
- Kryteria wyboru stylu → Wzorce → Anty‑wzorce → Checklista.
- Bezpieczeństwo → Wzorce authz/rate limit → Testy/observability.
- Kontrakty/wersjonowanie → Change Mgmt → Compatibility policy.

## Fazy cyklu życia
- Analiza: cele biznesowe, latency/SLO, wymagania bezpieczeństwa/compliance.
- Projekt: wybór wzorca, kontrakt, versioning, błędy, idempotencja, limity.
- Implementacja: zgodność ze standardem, testy kontraktowe, observability.
- Wdrożenie: rollout/canary, deprecations/compat.
- Utrzymanie: ewolucja wersji, breaking changes policy, audyt wzorców.

## Struktura sekcji
1) Kryteria wyboru stylu (REST/RPC/event) i kompatybilności  
2) Wzorce zasobów i naming (collection, subresource, actions)  
3) Wersjonowanie (path/header/media type), kompatybilność, deprecations  
4) Paginacja/filtry/sortowanie (limit/offset, cursor, server‑side)  
5) Idempotencja i retry (klucze idempotentne, safe verbs, deduplikacja)  
6) Błędy i kody (problem details, mapowanie wyjątków, korelacja trace)  
7) Bezpieczeństwo (authn/z, rate limit, input validation, payload size, replay)  
8) Kontrakty i schemy (OpenAPI/Proto/AsyncAPI, schema registry)  
9) Format danych (JSON/Proto/CSV/NDJSON), streaming/bulk  
10) Observability (trace/span, log correlation, metrics, audit)  
11) Anty‑wzorce i kiedy NIE używać danego wzorca  
12) Checklista projektowa i linki do przykładów/SDK\n\n## Wymagane rozwinięcia\n- Przykłady kontraktów (OpenAPI/Proto) dla kluczowych wzorców.\n- Mapowanie kodów błędów i format problem+trace id.\n- Reguły idempotencji i retry per typ operacji.\n\n## Wymagane streszczenia\n- Polityka wersjonowania i deprecations.\n- Minimalny zestaw kontroli bezpieczeństwa (authz, rate limit, input validation, size limits).\n\n## Guidance (skrót)\n- Wybierz styl wg konsumenta i charakteru: REST (zasoby), RPC (niskie opóźnienia), event (loose coupling).\n- Zawsze definiuj idempotencję dla operacji modyfikujących; dokumentuj retry/backoff.\n- Ustandaryzuj błędy (problem details) i koreluj je z trace id.\n- Wymuś limity (rate, payload, pagination) i walidację wejścia; unikaj over/under‑fetching.\n- Plan deprecations: komunikaty, okres wsparcia, observability starych wersji.\n\n## Szybkie powiązania\n- linkage_index.jsonl (api/design)\n- api_security_baseline, api_versioning_policy, observability_standards, error_handling_guidelines\n\n## Jak używać dokumentu\n1. Określ styl API i kryteria (konsument, latency, coupling).\n2. Wybierz wzorzec, wypełnij sekcje: kontrakt, wersjonowanie, błędy, idempotencja, bezpieczeństwo.\n3. Dodaj przykład kontraktu i mapę błędów; sprawdź checklistę.\n4. Zaktualizuj DoR/DoD i link do OpenAPI/Proto w repo.\n\n## Checklisty Definition of Ready (DoR)\n- [ ] Cele i SLO/KPI API zdefiniowane; konsument zidentyfikowany.\n- [ ] Polityki bezpieczeństwa/wersjonowania/observability dostępne.\n- [ ] Wstępny wybór stylu i wzorca opisany; alternatywy rozważone.\n- [ ] Struktura sekcji wypełniona/N/A.\n\n## Checklisty Definition of Done (DoD)\n- [ ] Kontrakt (OpenAPI/Proto) gotowy; pola opisane, walidacje podane.\n- [ ] Błędy, idempotencja, retry, limity i bezpieczeństwo opisane.\n- [ ] Observability: trace/log/metrics/audit uzgodnione; przykłady nagłówków.\n- [ ] Deprecations/versioning opisane; linki do przykładów/SDK działają.\n- [ ] Wersja/data/właściciel zaktualizowane; dokument w linkage_index.\n\n## Definicje robocze\n- Idempotencja, Backoff, Retry budget, Cursor pagination, Problem+JSON, gRPC deadline.\n\n## Przykłady użycia\n- Publiczne API katalogu produktów (REST + cursor + problem details).\n- Internal low‑latency API scoringu (gRPC + deadline + idempotent token + rate limit).\n- Event API (AsyncAPI) do integracji partnerów (outbox, deduplikacja, DLQ).\n\n## Ryzyka i ograniczenia\n- Brak spójności wersjonowania → dług techniczny; błędy bez korelacji → trudne debugowanie.\n- Zbyt ogólne endpointy → overfetch/underfetch; brak limitów → DoS/koszty.\n\n## Decyzje i uzasadnienia\n- [Decyzja] Styl i wersjonowanie — uzasadnienie konsument/SLO/kompat.\n- [Decyzja] Limity i polityki retry — uzasadnienie ryzyka/latency.\n\n## Założenia\n- Org ma standard security/observability; repozytorium kontraktów dostępne.\n\n## Otwarte pytania\n- Czy wymagane są kompatybilność wsteczna i jak długo wspieramy stare wersje?\n- Czy konsument wymaga streaming/bulk?\n\n## Powiązania z innymi dokumentami\n- API Security Baseline, Versioning Policy, Observability Standards, Change Mgmt.\n\n## Powiązania z sekcjami innych dokumentów\n- Error Handling Guidelines → Błędy; Security Baseline → authz/limity; Observability → trace/log.\n\n## Słownik pojęć w dokumencie\n- REST, RPC, Event‑driven, HATEOAS, DLQ, Outbox, Backpressure.\n\n## Wymagane odwołania do standardów\n- RFC 9110 (HTTP), gRPC/Proto, JSON:API lub Problem+JSON (jeśli stosowane), AsyncAPI.\n\n## Mapa relacji sekcja→sekcja\n- Kryteria → Wzorce → Kontrakt/Błędy/Security → Observability → Deprecations.\n\n## Mapa relacji dokument→dokument\n- API Design Patterns → API Security/Versioning → Change Mgmt → QA/Testing.\n\n## Ścieżki informacji\n- Wymagania → Wybór wzorca → Kontrakt → Implementacja → Observability → Deprecation.\n\n## Weryfikacja spójności\n- [ ] Wzorzec dobrany do konsumenta/SLO; kontrakt spójny z politykami.\n- [ ] Błędy/idempotencja/retry zgodne z security i observability.\n- [ ] Deprecations opisane i zsynchronizowane z change management.\n\n## Lista kontrolna spójności relacji\n- [ ] Każdy endpoint ma wzorzec, wersję, błędy, limity i observability.\n- [ ] Każdy alert/limit ma uzasadnienie i właściciela.\n- [ ] Relacje cross‑doc opisane z uzasadnieniem.\n\n## Artefakty powiązane\n- OpenAPI/Proto/AsyncAPI, przykładowe payloady, SDK/snippety, testy kontraktowe.\n\n## Ścieżka decyzji\n- [Decyzja] → [Uzasadnienie] → [Konsekwencje] → [Właściciel] → [Data].\n\n## Użytkownicy i interesariusze\n- Architekci, Zespoły produktowe, Security, SRE/Observability, QA, Partnerzy.\n\n## Ścieżka akceptacji\n- Architecture/API Guild → Security → Product/Owner → Release/CAB.\n\n## Kryteria ukończenia\n- [ ] Kontrakt i wzorzec zaakceptowane; linki/artefakty kompletne.\n- [ ] Checklista i powiązania domknięte; dokument w linkage_index.\n\n## Metryki jakości\n- % usług zgodnych ze standardem, liczba/lub severity błędów kontraktu, czas publikacji nowej wersji, liczba breaking changes bez deprecation, pokrycie trace/log/metrics.\n*** End Patch 豪泰။
