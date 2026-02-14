---
title: Incident Response Training
status: needs_content
---

# Incident Response Training

## Metadane
- Właściciel: [Security/SRE/Comms]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Plan szkolenia z reagowania na incydenty: zakres, role, scenariusze, ćwiczenia (tabletop/tech), materiały i mierniki. Ma zwiększyć gotowość zespołów, skrócić MTTR i poprawić komunikację.

## Zakres i granice
- Obejmuje: role (Incident Commander, Comms, Scribe, SME), procedury IR, kanały i komunikacja (status page/PR/CS), scenariusze techniczne i komunikacyjne, ćwiczenia (tabletop/live/chaos), materiały (runbooki, checklists, komunikaty), kryteria gotowości, metryki (MTTA/MTTR, coverage, udział), harmonogram i retencję materiałów.
- Poza zakresem: pełny Incident Response Plan (link), plany DR/BCP (linkowane).

## Wejścia i wyjścia
- Wejścia: IR Plan, runbooki, lista ról/on-call, kanały komms, historia incydentów, ryzyka top, wymagania compliance/audytu, kalendarz ćwiczeń.
- Wyjścia: plan szkoleń/ćwiczeń, scenariusze, materiały, harmonogram, lista uczestników, raporty z ćwiczeń, akcje naprawcze, metryki skuteczności.

## Powiązania (meta)
- Key Documents: incident_response_plan, communication_plan_crisis, drp_bcp, security_baseline, system_monitoring_strategy, change_management.
- Key Document Structures: role, scenariusze, harmonogram, materiały, metryki.
- Document Dependencies: on-call roster, runbooki, status page, kanały PR/CS, narzędzia tabletop/chaos.

## Zależności dokumentu
Wymaga IR Plan, runbooków, ról/on-call, kanałów komunikacji, historii incydentów/ryzyk, narzędzi ćwiczeń. Bez tego DoR otwarte.

## Powiązania sekcja↔sekcja
- Scenariusze → Materiały → Ćwiczenia → Raporty → Akcje.
- Metryki → Korekty planu → Kolejny cykl.

## Fazy cyklu życia
- Planowanie: cele, role, scenariusze, harmonogram.
- Przygotowanie: materiały, narzędzia, logistyka.
- Wykonanie ćwiczeń: tabletop/live/chaos.
- Raportowanie: wnioski, akcje, metryki.
- Retrospektywa i poprawa planu.

## Struktura sekcji
1) Cele i kryteria gotowości (MTTA/MTTR, coverage, komunikacja)  
2) Role i odpowiedzialności (IC, Comms, Scribe, SME, Legal, PR, CS)  
3) Scenariusze i typy ćwiczeń (tabletop, tech/chaos, komunikacyjne)  
4) Materiały (runbooki, checklisty, komunikaty, status page szablony)  
5) Harmonogram i uczestnicy (kalendarz, częstotliwość, obowiązkowość)  
6) Metryki i raporty (czas reakcji, błędy, coverage, udział)  
7) Plan akcji po ćwiczeniach i śledzenie  
8) Ryzyka, decyzje, open issues

## Wymagane rozwinięcia
- Lista scenariuszy z celami, wymagane role, kryteria sukcesu.
- Szablony komunikatów/status page, checklisty IC/Scribe/Comms.
- Harmonogram (np. kwartalne tabletop, roczne chaos), metryki i raport.

## Wymagane streszczenia
- Najbliższe ćwiczenia, cele, scenariusze, uczestnicy, metryki do pomiaru.

## Guidance (skrót)
- Różnicuj scenariusze (infra/app/security/comms); ćwicz role i przekazy.
- Mierz MTTA/MTTR i błędy procesu; zbieraj akcje i domykaj je.
- Aktualizuj materiały po każdej lekcji; dbaj o privacy/legal w komunikacji.

## Szybkie powiązania
- linkage_index.jsonl (incident/training)
- incident_response_plan, communication_plan_crisis, drp_bcp, security_baseline, system_monitoring_strategy, change_management

## Jak używać dokumentu
1. Ustal cele/role i scenariusze; przygotuj materiały i harmonogram.
2. Przeprowadź ćwiczenia; zbierz raporty i metryki.
3. Zamknij akcje; zaktualizuj plan i dokument; dodaj do linkage_index.

## Checklisty Definition of Ready (DoR)
- [ ] IR Plan/runbooki dostępne; role/on-call i kanały komunikacji przygotowane.
- [ ] Scenariusze/cele/harmonogram zdefiniowane; narzędzia ćwiczeń gotowe.
- [ ] Struktura sekcji wypełniona/N/A.

## Checklisty Definition of Done (DoD)
- [ ] Ćwiczenia wykonane; raporty i metryki zebrane; akcje przypisane.
- [ ] Materiały/runbooki zaktualizowane; dokument w linkage_index.
- [ ] Wersja/data/właściciel zaktualizowane.

## Definicje robocze
- Tabletop, Chaos, Incident Commander, MTTA, MTTR, Status page.

## Przykłady użycia
- Kwartalny tabletop P1: awaria DB, ćwiczenie IC/Comms/Legal, raport z MTTA/MTTR.
- Chaos ćwiczenie: odcięcie zależności, test runbooków i komunikacji.

## Ryzyka i ograniczenia
- Brak ćwiczeń → słaba reakcja; brak follow-up akcji → brak poprawy; brak privacy/legal → ryzyko komunikacji.

## Decyzje i uzasadnienia
- [Decyzja] Kadencja ćwiczeń i scenariusze — uzasadnienie ryzyk.
- [Decyzja] Metryki sukcesu — uzasadnienie SLO/MTTR.

## Założenia
- Dostęp do narzędzi komunikacji/status; runbooki aktualne; wsparcie leadership.

## Otwarte pytania
- Czy wymagane certyfikacje/audyt dowodów? 
- Jakie scenariusze regulatory (np. breach) włączyć?

## Powiązania z innymi dokumentami
- Incident Response Plan, Crisis Comms, DRP/BCP, Security Baseline, Monitoring Strategy, Change Mgmt.

## Powiązania z sekcjami innych dokumentów
- IR → role/runbooki; Comms → szablony; DRP → scenariusze.

## Słownik pojęć w dokumencie
- Tabletop, Chaos, Incident Commander, MTTA, MTTR, Status page.

## Wymagane odwołania do standardów
- Wewnętrzne polityki IR/DR/BCP, wymagania audytowe (SOC2/ISO) jeśli dotyczy.

## Mapa relacji sekcja→sekcja
- Scenariusze → Ćwiczenia → Raporty → Akcje → Aktualizacja materiałów.

## Mapa relacji dokument→dokument
- IR Training → IR Plan/Comms/DRP → Lessons Learned → Monitoring/Change.

## Ścieżki informacji
- Plan → Ćwiczenia → Raport → Akcje → Aktualizacja → Kolejny cykl.

## Weryfikacja spójności
- [ ] Scenariusze/cele/role zdefiniowane; raporty i akcje domknięte; dokument w linkage_index.

## Lista kontrolna spójności relacji
- [ ] Każde ćwiczenie ma cele, role, materiał, raport, akcje.
- [ ] Każda akcja ma owner/ETA i status; relacje cross‑doc opisane.

## Artefakty powiązane
- Scenariusze, materiały, raporty z ćwiczeń, lista akcji, log komunikacji.

## Ścieżka decyzji
- [Decyzja] → [Uzasadnienie] → [Konsekwencje] → [Właściciel] → [Data].

## Użytkownicy i interesariusze
- Security, SRE/Platform, Comms/PR, Product, Legal, Support.

## Ścieżka akceptacji
- Security/SRE → Comms/PR/Legal → Leadership → Owner sign‑off.

## Kryteria ukończenia
- [ ] Plan i ćwiczenia przeprowadzone; raporty/akcje zaktualizowane; dokument w linkage_index.
- [ ] Wersja/data/właściciel aktualne.

## Metryki jakości
- MTTA/MTTR w ćwiczeniach, udział ról, liczba akcji domkniętych, czas aktualizacji materiałów.
