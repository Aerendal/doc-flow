---
title: API Reference for Mobile Developers
status: needs_content
---

# API Reference for Mobile Developers

## Metadane
- Właściciel: [API Product Owner / Mobile Lead]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Zapewnić mobilnym deweloperom spójne, stabilne i bezpieczne użycie API: kontrakty, autoryzacja, limity, wersjonowanie, formaty odpowiedzi/błędów oraz przykłady implementacji. Minimalizuje regresje i różnice między platformami iOS/Android.

## Zakres i granice
- Obejmuje: endpointy mobilne, auth (OAuth2/OIDC, PKCE), rate limiting, wersjonowanie, formaty payloadów (JSON/protobuf), obsługę błędów i retry, paginację/filtry, caching, offline/edge cases, telemetry (logs/metrics), bezpieczeństwo (cert pinning, TLS), testy kontraktów.  
- Poza zakresem: pełne wytyczne UI/UX, SDK build system, backend wewnętrzna architektura.

## Wejścia i wyjścia
- Wejścia: specyfikacje OpenAPI/GraphQL, polityka auth, limity, wymagania mobilne (offline, sieć niestabilna), kody błędów, przykładowe flow.  
- Wyjścia: referencja endpointów z przykładami, tablice kodów błędów, zalecenia caching/offline, sample requests w cURL/Kotlin/Swift, checklisty DoR/DoD kontraktu, matryca kompatybilności wersji.

## Powiązania (meta)
- Key Documents: api_versioning_maintenance, authentication_metrics_report, mobile_security_guidelines, error_handling_standards, monitoring_strategy_document, client_release_checklist.  
- Key Document Structures: auth, kontrakty, błędy, wydajność/offline, wersjonowanie, testy.  
- Document Dependencies: API gateway, IdP, rate limiter, monitoring/analytics, schema registry.

## Zależności dokumentu
Wymaga aktualnej specyfikacji API, polityki auth, limitów, listy wspieranych wersji mobilnych SDK, środowisk testowych/stage, oraz dostępnych przykładów kontraktów. Braki = brak DoR.

## Powiązania sekcja↔sekcja
- Auth ↔ Bezpieczeństwo klienta (pinning) ↔ Błędy/retry.  
- Wersjonowanie ↔ Kompatybilność SDK ↔ Migracje.  
- Rate limiting ↔ Retry/backoff ↔ Monitoring.

## Fazy cyklu życia
- Definicja kontraktów i wersji.  
- Implementacja klienta i testy kontraktów.  
- Wydanie i monitorowanie.  
- Deprecacje/migracje i komunikacja.

## Struktura sekcji
1) Auth i bezpieczeństwo klienta (OAuth2/OIDC, PKCE, pinning)  
2) Kontrakty endpointów (metody, parametry, schematy, przykłady)  
3) Błędy i retry (kody, komunikaty, backoff)  
4) Wersjonowanie, zgodność i migracje  
5) Wydajność i offline (cache, batching, timeouts)  
6) Telemetria i monitoring na kliencie  
7) DoR/DoD, ryzyka, pytania

## Wymagane rozwinięcia
- Tabela endpointów z przykładami (cURL/Swift/Kotlin).  
- Matryca kodów błędów i zachowań klienta.  
- Polityka wersjonowania (header/path) i deprecacji.  
- Rekomendacje cache/offline (ETag/If-None-Match, store-and-forward).  
- Konfiguracja pinning TLS i rotacja certów.  
- Testy kontraktów (schemat vs payload) w CI.

## Wymagane streszczenia
- Executive summary: wersja API, wymagany poziom SDK, kluczowe zmiany.  
- Skrót limitów i polityki retry/backoff.

## Guidance (skrót)
- Używaj PKCE i odświeżania tokenów bez przechowywania haseł.  
- Stosuj stałe schematy błędów; nie polegaj na tekstach komunikatów.  
- Dodaj telemetry (request ID, latency) i centralny logging błędów klienta.  
- Wymuszaj TLS 1.2+ i pinning; rotuj piny przed wygaśnięciem.  
- Przy sieci słabej stosuj cache, kolejkę offline i exponential backoff.  
- Trzymaj kompatybilność: nie usuwaj pól, tylko oznaczaj jako deprecated.

## Szybkie powiązania
- linkage_index.jsonl (api/reference/mobile)  
- api_versioning_maintenance, mobile_security_guidelines, error_handling_standards

## Jak używać dokumentu
1. Sprawdź wymagania auth i limity; skonfiguruj klienta.  
2. Implementuj endpointy wg przykładów; dodaj obsługę błędów/retry.  
3. Uruchom testy kontraktów i a11y sieci (offline/slaba sieć).  
4. Publikuj aplikację; monitoruj metryki i aktualizuj przy zmianach API.

## Checklisty Definition of Ready (DoR)
- [ ] Specyfikacja API aktualna i dostępna.  
- [ ] Limity i polityka auth potwierdzone.  
- [ ] Wersje SDK/OS wspierane zdefiniowane.  
- [ ] Strategie błędów/retry/cache ustalone.  
- [ ] Środowiska test/stage i dane testowe gotowe.

## Checklisty Definition of Done (DoD)
- [ ] Endpointy zaimplementowane i pokryte testami kontraktów.  
- [ ] Obsługa błędów/retry zgodna z tabelą.  
- [ ] Pinning TLS i polityka tokenów zaimplementowane.  
- [ ] Monitoring klienta wysyła metryki/logi; brak krytycznych błędów.  
- [ ] Dokumentacja i linkage_index zaktualizowane.

## Definicje robocze
- PKCE: rozszerzenie OAuth2 dla aplikacji publicznych.  
- ETag: nagłówek do walidacji cache.  
- Pinning: weryfikacja certyfikatu serwera przez klienta.

## Przykłady użycia
- Pobieranie feedu produktów z cache i walidacją ETag.  
- Odnawianie tokenu po 401 + retry z backoff.  
- Migracja z v1 na v2 endpointu z polami deprecated.

## Ryzyka i ograniczenia
- Brak pinning → ryzyko MITM.  
- Niespójne błędy → trudna diagnostyka i UX.  
- Agresywne retry → blokady rate limiting.  
- Brak wersjonowania → breaking changes w aplikacjach produkcyjnych.

## Decyzje i uzasadnienia
- Wersjonowanie (path vs header) i okno deprecacji.  
- Strategie cache/offline i limity retry.  
- Zestaw wymaganych nagłówków (request-id, locale, app-version).  
- Format błędów (np. RFC 7807) i logowania klienta.

## Założenia
- Gateway wspiera OIDC/OAuth2 + rate limiting.  
- Monitoring klienta (crash/metrics) jest dostępny.  
- Użytkownicy stosują najnowsze SDK lub mają ścieżkę migracji.

## Otwarte pytania
- Jak długo wspieramy starsze wersje API/SDK?  
- Czy potrzebna jest lokalizacja komunikatów błędów po stronie klienta?  
- Jaki poziom telemetry jest wymagany a co jest opcjonalne?  
- Jak obsługiwać tryb offline przy wymaganym tokenie (grace window)?
