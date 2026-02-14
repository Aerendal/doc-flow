---
title: Post-Deployment Support Plan
status: needs_content
---

# Post-Deployment Support Plan

## Metadane
- Właściciel: [Support / SRE / Product]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Opisuje wsparcie tuż po wdrożeniu (stabilizacja/hypercare): monitoring, obsługa incydentów, komunikacja, eskalacje, runbooki, kryteria zakończenia hypercare i przekazania do BAU.

## Zakres i granice
- Obejmuje: okres hypercare, zakres monitoringu, obsługę incydentów/problemów, kanały i SLA, eskalacje/on-call, runbooki krytycznych ścieżek, kryteria exit, raportowanie, komunikację do klientów i wewnętrzną.
- Poza zakresem: długoterminowy support BAU (opisany w Solution Support Plan); pełne playbooki IR (osobny dokument IR).

## Wejścia i wyjścia
- Wejścia: plan wdrożenia, znane ryzyka/defekty, runbooki, monitoring/alerty, release notes, lista właścicieli komponentów, SLA/OLA, kontakt do klientów.
- Wyjścia: plan hypercare (czas, role, on-call), zakres monitoringu i alertów, drabinka eskalacji, runbooki i checklisty walidacyjne, raportowanie (status, incydenty, defekty), kryteria exit i przekazanie do BAU.

## Powiązania (meta)
- Key Documents: solution_support_plan, incident_response_playbook, change_management_process, release_plan, rollback_plan, risk_register.
- Dependencies: monitoring/alerting, właściciele komponentów, znane defekty, komunikacja, SLA.
- RACI: Product/Release Owner, SRE, Support, Engineering, Security/Compliance, Comms.

## Zależności dokumentu
- Upstream: release plan, znane ryzyka/defekty, monitoring/alerting, SLA.
- Downstream: raporty hypercare, action items, przekazanie do BAU, retrospektywa.
- Zewnętrzne: klienci/partnerzy (komunikacja), dostawcy (jeśli zależności 3rd party).

## Powiązania sekcja↔sekcja
- Monitoring/alerty → obsługa incydentów → komunikacja → raportowanie → kryteria exit.
- Runbooki → walidacje powdrożeniowe → decyzje rollback/continue.

## Fazy cyklu życia
- Przygotowanie: zakres hypercare, monitoring/alerty, runbooki, role/on-call, komunikacja.
- Wykonanie: aktywny hypercare, incydenty, raporty statusu.
- Zakończenie: spełnienie kryteriów exit, przekazanie BAU, retrospektywa.

## Struktura sekcji (szkielet)
1) Streszczenie i cel hypercare (czas trwania, KPI: incydenty, MTTR, błędy krytyczne)
2) Zakres hypercare (komponenty, funkcje krytyczne, klientów/regiony objęte)
3) Monitoring i alerty (metryki, progi, dashboardy, SLO)
4) Role/on-call i eskalacje (drabinka, kontakty, SLA/OLA)
5) Runbooki i checklisty walidacji powdrożeniowej (sanity checks, smoke, dane)
6) Obsługa incydentów/problemów (proces, komunikacja wew./zewn., warunki rollback)
7) Raportowanie i komunikacja (status daily, klient, stakeholderzy)
8) Kryteria exit hypercare i przekazanie do BAU (SLO/KPI, defekty otwarte, dokumentacja)
9) Ryzyka i założenia; decyzje (ADR) i otwarte pytania

## Wymagane rozwinięcia
- Lista dashboardów/alertów i progów; drabinka eskalacji; runbooki sanity/smoke/rollback; checklisty komunikacji; kryteria exit.
- Raport template (status, incydenty, defekty, działania, ryzyka).

## Wymagane streszczenia
- Executive summary: czas hypercare, KPI, incydenty/defekty, decyzje rollback/continue, data exit.
- One-pager: kontakty/on-call, SLO, kryteria exit, linki do runbooków.

## Guidance (skrót)
- DoR: release plan, monitoring/alerty gotowe; runbooki i rollback; on-call i kontakty; komunikacja przygotowana; SLO/exit criteria uzgodnione.
- DoD: hypercare przeprowadzony; incydenty/defekty zarejestrowane; kryteria exit spełnione; przekazanie BAU i retrospektywa; metadane aktualne; dokument w linkage_index.
- Spójność: monitoring pokrywa krytyczne ścieżki; eskalacje/on-call działają; decyzje rollback mają warunki i ownerów; exit criteria są mierzalne.

## Szybkie powiązania
- solution_support_plan, incident_response_playbook, change_management_process, release_plan, rollback_plan, risk_register

## Checklisty Definition of Ready (DoR)
- [ ] Monitoring/alerty i runbooki gotowe; on-call i kontakty; SLO/exit criteria uzgodnione.
- [ ] Komunikacja (wew./zewn.) przygotowana; release/rollback plan dostępny.

## Checklisty Definition of Done (DoD)
- [ ] Hypercare zakończony; kryteria exit spełnione; incydenty/defekty zamknięte lub przekazane; przekazanie BAU; retrospektywa.
- [ ] Metadane aktualne; dokument w linkage_index.

## Artefakty powiązane
- Dashboardy/alerty, escalation ladder, runbooki sanity/smoke/rollback, raporty hypercare, checklisty exit, ADR log.

## Weryfikacja spójności
- [ ] Krytyczne ścieżki mają monitoring/alerty; runbooki powdrożeniowe wykonane.
- [ ] Incydenty/defekty obsłużone; komunikacja prowadzona; decyzje rollback/continue udokumentowane.
- [ ] Kryteria exit spełnione; przekazanie BAU i retrospektywa wykonane.
