---
title: Incident Checklist
status: needs_content
---

# Incident Checklist

## Metadane
- Właściciel: [Incident Manager / SRE / CSIRT]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Zawiera krótką checklistę do użycia w trakcie incydentu, aby zapewnić spójny triage, komunikację, naprawę i zamknięcie z pełną dokumentacją.

## Zakres i granice
- Obejmuje: identyfikację i klasyfikację (sev), przypisanie ról (IM/commander/scribe), komunikację (status page, kanały), diagnostykę, obejścia/rollback, eskalacje, potwierdzenie naprawy, post‑mortem, aktualizację runbooków.  
- Poza zakresem: szczegółowe runbooki techniczne (w innych dokumentach).

## Wejścia i wyjścia
- Wejścia: alert/zgłoszenie, matryca severity, runbooki, dane kontaktowe, SLA/OLA, status page policy.  
- Wyjścia: karta incydentu z timeline, komunikaty, wykonane akcje/rollback, decyzje, post‑mortem, checklisty DoR/DoD.

## Powiązania (meta)
- Key Documents: incident_response_for_customers, communication_failure, rollback_runbook, status_page_policy, incident_pattern_analysis.  
- Key Document Structures: triage, role, komunikacja, działania, zamknięcie.  
- Document Dependencies: monitoring, ticketing, status page, comms tools, CMDB/runbook repo.

## Zależności dokumentu
Wymaga: matrycy severity, listy ról i kontaktów, szablonów komunikatów, runbooków, dostępu do monitoring/ticketing/status page. Brak = brak DoR.

## Powiązania sekcja↔sekcja
- Triage ↔ Role ↔ Komunikacja.  
- Diagnostyka ↔ Działania/rollback ↔ Potwierdzenie.  
- Timeline ↔ Post‑mortem ↔ Aktualizacja runbooków.

## Fazy cyklu życia
- Detekcja i triage.  
- Aktywacja ról i komunikacji.  
- Diagnostyka/działania/rollback.  
- Walidacja i komunikat „resolved”.  
- Post‑mortem i follow‑ups.

## Struktura sekcji
1) Start: severity, role assignment, kanały  
2) Komunikacja (initial/update/resolved) i częstotliwość  
3) Diagnostyka i działania (obejścia, rollback)  
4) Potwierdzenie naprawy i monitoring  
5) Zamknięcie, post‑mortem, DoR/DoD  
6) Lista kontrolna (checkboxy)  
7) Ryzyka, pytania

## Wymagane rozwinięcia
- Szablon karty incydentu (timeline, decyzje).  
- Szablony komunikatów dla P1/P2.  
- Lista eskalacji i kontaktów.  
- Checklista działań (diag, rollback, walidacja).  
- Szablon post‑mortem i właściciele akcji.  
- Linki do runbooków technicznych.

## Wymagane streszczenia
- Executive summary: status, wpływ, ETA.  
- Skrót timeline po zamknięciu.

## Guidance (skrót)
- Najpierw bezpieczeństwo i komunikacja: przydziel role, otwórz kanał.  
- Utrzymuj jeden source of truth (status page/ticket).  
- Dokumentuj timeline na bieżąco; loguj decyzje.  
- Preferuj obejścia i rollback gdy szybciej przywracają SLA.  
- Post‑mortem w ciągu 5 dni; przypisz właścicieli akcji.  
- Aktualizuj linkage_index i runbooki.

## Szybkie powiązania
- linkage_index.jsonl (incident/checklist)  
- incident_response_for_customers, rollback_runbook

## Jak używać dokumentu
1. Otwórz checklistę, określ severity, przypisz role.  
2. Uruchom komunikację, prowadź timeline.  
3. Wykonuj diag/działania, decyzje zapisuj.  
4. Potwierdź przywrócenie, publikuj „resolved”.  
5. Zrób post‑mortem, aktualizuj dokument i linkage_index.

## Checklisty Definition of Ready (DoR)
- [ ] Severity i role zmapowane; kanały komunikacji dostępne.  
- [ ] Runbooki i narzędzia dostępne.  
- [ ] Szablony komunikatów przygotowane.  
- [ ] Dane kontaktowe i eskalacje aktualne.  
- [ ] Wymagania SLA/OLA znane.

## Checklisty Definition of Done (DoD)
- [ ] Komunikat resolved wysłany; status page zaktualizowany.  
- [ ] Timeline i decyzje zapisane; post‑mortem zaplanowane/wykonane.  
- [ ] Akcje follow‑up przypisane; brak otwartych krytycznych ryzyk.  
- [ ] linkage_index/runbooki zaktualizowane.  
- [ ] Dane do raportów SLA/OLA uzupełnione.

## Definicje robocze
- Severity: klasyfikacja wpływu incydentu.  
- Post‑mortem: analiza przyczyn i działań po incydencie.

## Przykłady użycia
- Outage API P1 z komunikacją co 15 min.  
- Degradacja wydajności P2 z rollbackiem release.  
- Incydent bezpieczeństwa przekazany do CSIRT.

## Ryzyka i ograniczenia
- Brak komunikacji → eskalacje klientów.  
- Niedokumentowany timeline → słaby post‑mortem.  
- Brak rollback → dłuższy downtime.  
- Brak aktualnych kontaktów → opóźnienia.

## Decyzje i uzasadnienia
- Częstotliwość update’ów.  
- Kiedy rollback vs hotfix.  
- Kto zatwierdza komunikaty P1/P2.  
- Kadencja post‑mortem i właściciele.

## Założenia
- Monitoring i ticketing działają.  
- Zespół zna role i SLA.  
- Kanały komunikacji są dostępne.

## Otwarte pytania
- Jak długo przechowywać timeline i nagrania?  
- Czy potrzebne są wersje branżowe checklist (np. bezpieczeństwo)?  
- Jak łączyć checklistę z automatyzacją status page?
