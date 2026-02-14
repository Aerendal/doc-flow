---
title: Database Cutover Procedure
status: needs_content
---

# Database Cutover Procedure

## Metadane
- Właściciel: [DBA/SRE/Platform]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Opisuje krok‑po‑kroku bezpieczne przełączenie ruchu na nową bazę danych lub nową wersję (cutover/rollback). Minimalizuje ryzyko utraty danych, niedostępności, niespójności i regresji wydajności.

## Zakres i granice
- Obejmuje: typy cutover (switch over/blue‑green/dual‑write/read‑only phase), plan techniczny i komunikację, zamrożenia zmian, walidację danych, smoke/health checks, plan rollback, role i komunikację, kontrolę czasu i warunki go/no‑go.
- Poza zakresem: projekt schematu, migracje danych (osobny runbook), tuning wydajności po cutover (link do performance tuning plan).

## Wejścia i wyjścia
- Wejścia: plan migracji danych, lista aplikacji/klientów, zależności (queues/cache/search/BI), okno serwisowe i SLA, plan komunikacji, checklista ryzyk, backup/restore test, skrypty przełączenia, progi health checks.
- Wyjścia: harmonogram cutover, przypisane role i kanały, wykonane kroki z logiem, wyniki walidacji (pre/post), decyzja go/no‑go, raport z rollbacku (jeśli użyty), aktualizacja CMDB/runbooków.

## Powiązania (meta)
- Key Documents: data_migration_execution, rollback_plan, change_management_request, incident_response_runbook, performance_baseline_report.
- Key Document Structures: przygotowanie, wykonanie, walidacja, rollback, komunikacja.
- Document Dependencies: backup/restore proces, monitoring/alerty, feature flags, klienty aplikacyjne, replikacja/log shipping, time sync.

## Zależności dokumentu
Wymaga: sprawdzonego backupu i testu odtworzenia, zamknięcia otwartych zmian schematu, uzgodnionego okna serwisowego, listy konsumentów DB i ich wersji driverów, planu komunikacji. Brak któregokolwiek = DoR otwarte.

## Powiązania sekcja↔sekcja
- Przygotowanie → Wykonanie → Walidacja → Decyzja go/no‑go.  
- Walidacja → Rollback lub Stabilizacja → Monitoring po cutover.  
- Komunikacja → Akceptacje → Start cutover.

## Fazy cyklu życia
- Planowanie: wybór strategii cutover, okno, testy w staging.  
- Wykonanie: zamrożenie zmian, backup, przełączenie, walidacja, decyzja.  
- Stabilizacja: obserwacja KPI/alertów, performance tuning, usunięcie starych ścieżek.  
- Sunset: wyłączenie starej bazy, archiwizacja, aktualizacja dokumentacji.

## Struktura sekcji
1) Kontekst i zakres cutover (typ, systemy, okno)  
2) Strategia (blue/green, shadow, dual‑write/read, read‑only phase)  
3) Przygotowanie (backup, replikacja, freeze, test w staging, dane referencyjne)  
4) Kroki wykonania (kolejność, polecenia, punkt kontroli czasu)  
5) Walidacja po przełączeniu (integralność, spójność, wydajność, błędy)  
6) Warunki go/no‑go i decyzje  
7) Rollback plan (trigger, kroki, dane, komunikacja)  
8) Komunikacja i RACI (kto, kiedy, kanały)  
9) Monitoring i stabilizacja (KPI, alerty, tuningi)  
10) Dokumentacja i zamknięcie (raport, CMDB, lessons learned)

## Wymagane rozwinięcia
- Listy komend/procedur (SQL, skrypty) dla każdego etapu.  
- Checklista walidacji danych (liczność, checksumy, dane referencyjne, kluczowe tabele).  
- Macierz zgodności klientów (wersje driverów/ORM, retry, timeouty).

## Wymagane streszczenia
- Run sheet na czas okna: sekwencja kroków z timestampami i ownerami.  
- Executive summary po cutover: wynik, KPI, incydenty, decyzje o rollback/stay.

## Guidance (skrót)
- Zawsze miej testowany rollback z danymi produkcyjnymi (maskowanymi) w staging.  
- Zamroź schemat i wdrożenia aplikacji na czas cutover; wyczyść kolejki/cron joby jeśli wymagane.  
- Weryfikuj dane wieloma metodami: counts, checksums, critical rows, biznesowe raporty kontrolne.  
- Ustal pojedynczy punkt decyzji go/no‑go z jasnymi progami.  
- Po cutover monitoruj błędy, latency, blokady, lag replikacji i budżet połączeń.

## Szybkie powiązania
- linkage_index.jsonl (database/cutover/procedure)  
- data_migration_execution, rollback_plan, change_management_request, incident_response_runbook, performance_baseline_report

## Jak używać dokumentu
1. Wypełnij sekcje strategii, przygotowania i komend; uzyskaj akceptację change request.  
2. W oknie serwisowym prowadź cutover z „run sheet” i odhaczaj checklisty.  
3. Po walidacji i decyzji zaktualizuj status DoD, raport i CMDB.

## Checklisty Definition of Ready (DoR)
- [ ] Backup przetestowany (restore + checksum) i dostępny.  
- [ ] Zamrożenie zmian schematu/DDL uzgodnione.  
- [ ] Lista klientów DB i wersji driverów potwierdzona.  
- [ ] Monitoring/alerty i dashboard po cutover skonfigurowane.  
- [ ] Rollback plan z testem w staging zaakceptowany.

## Checklisty Definition of Done (DoD)
- [ ] Cutover wykonany wg run sheet, log z timestampami zapisany.  
- [ ] Walidacje danych i wydajności zaliczone lub zaakceptowane wyjątki.  
- [ ] Decyzja stay/rollback udokumentowana, komunikacja wysłana.  
- [ ] Rollback niewykorzystany lub wykonany z pełnym raportem.  
- [ ] Dokumentacja/CMDB zaktualizowane, lessons learned zapisane.

## Definicje robocze
- Cutover: moment przełączenia ruchu na nowy primary/cluster.  
- Freeze: blokada zmian schematu/deployów na czas okna.  
- Go/No‑Go: formalna decyzja na podstawie checklist i walidacji.

## Przykłady użycia
- Migracja z single‑primary do HA cluster.  
- Przeniesienie na nową wersję bazy lub inny engine (np. MySQL→PostgreSQL) z dual‑write.  
- Cutover regionu w architekturze active‑active (wyłączenie ruchu z jednego regionu).

## Ryzyka i ograniczenia
- Brak spójności danych przy dual‑write bez idempotentnych operacji.  
- Zbyt długie okno read‑only powoduje straty biznesowe.  
- Niedoszacowanie TTL/timeoutów klientów skutkuje falą retry i przeciążeniem.

## Decyzje i uzasadnienia
- Strategia cutover (blue/green vs dual‑write) — zależnie od zgodności schematu i RTO/RPO.  
- Długość okna serwisowego — kompromis między bezpieczeństwem a biznesem.  
- Zakres walidacji — minimalny zestaw blokujący go/no‑go.

## Założenia
- Synchronizacja czasu (NTP) w całym środowisku.  
- Możliwość szybkiego odcięcia ruchu (LB/feature flag).  
- Monitoring i logi dostępne w czasie rzeczywistym.

## Otwarte pytania
- Czy wszystkie batch/cron muszą być zatrzymane, czy wystarczy drenaż kolejek?  
- Jakie są limity PSP/partnerów na przerwy w dostępności?  
- Czy potrzebny jest tymczasowy read‑only fallback dla raportów?

## Powiązania z innymi dokumentami
- change_management_request — formalna akceptacja okna.  
- incident_response_runbook — ścieżka eskalacji, gdy cutover się nie powiedzie.  
- performance_baseline_report — porównanie p99/p999 przed vs po.

## Wymagane odwołania do standardów
- Wewnętrzne standardy ciągłości działania (RPO/RTO).  
- Standardy bezpieczeństwa kopii zapasowych i szyfrowania.
