---
title: On-Call Training
status: needs_content
---

# On-Call Training

## Metadane
- Właściciel: [SRE/Support/Operations]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Program szkolenia dyżurujących (on-call): procedury, narzędzia, komunikacja, bezpieczeństwo i wellbeing. Celem jest skrócenie MTTR, poprawa jakości reakcji i zmniejszenie stresu.

## Zakres i granice
- Obejmuje: rolę on-call, kryteria eskalacji, runbooki, monitoring/alerty, narzędzia (pager/chat/ticketing), bezpieczeństwo (dostępy, PII), komunikację (status, klient), ćwiczenia (game day), raportowanie i retro, wellbeing (rotacje, zmiany, handover).  
- Poza zakresem: pełne procedury IR (linkowane), polityki HR.

## Wejścia i wyjścia
- Wejścia: matryca on-call, runbooki, SLO/SLA, narzędzia monitoringu, polityki bezpieczeństwa, harmonogramy, przykładowe incydenty, komunikacja klientów.  
- Wyjścia: sylabus szkoleń, checklisty DoR/DoD, playbooki, plan ćwiczeń, raporty z ewaluacji, poprawki runbooków.

## Powiązania (meta)
- Key Documents: incident_response_runbook, escalation_procedure_design, communication_plan, access_control_policy, observability_plan, game_day_plan.  
- Key Document Structures: rola, narzędzia, alerty, eskalacje, komunikacja, ćwiczenia, wellbeing.  
- Document Dependencies: monitoring/alerting, pager/chat, ticketing, CMDB, runbook repo.

## Zależności dokumentu
Wymaga: aktualnej matrycy on-call, runbooków i SLO, narzędzi monitoringu/pagera, polityk dostępu i komunikacji, harmonogramów rotacji. Braki = DoR otwarte.

## Powiązania sekcja↔sekcja
- Alerty → Eskalacja → Komunikacja → Postmortem.  
- Handover → Jakość reakcji → Wellbeing.  
- Ćwiczenia → Poprawa runbooków → Skrócenie MTTR.

## Fazy cyklu życia
- Przygotowanie materiałów i narzędzi.  
- Szkolenia i ćwiczenia (dry-run, game day).  
- Dyżury i operacje.  
- Retro/ewaluacja i iteracje.

## Struktura sekcji
1) Cel i zakres on-call  
2) Rola i odpowiedzialności, RACI  
3) Narzędzia i dostępy (pager, chat, ticketing, CMDB, logi)  
4) Alerty/SLO/SLA i kryteria eskalacji  
5) Komunikacja (wewnętrzna/zewnętrzna, status page, klient)  
6) Runbooki i playbooki (linki, aktualizacje)  
7) Ćwiczenia i game days (plan, częstotliwość)  
8) Wellbeing i rotacje (handover, zmiany, limit alertów)  
9) Raportowanie i retro (postmortem, poprawki)  
10) Ryzyka, decyzje, otwarte pytania

## Wymagane rozwinięcia
- Plan szkolenia (agenda, laby), lista runbooków krytycznych.  
- Szablon handover i komunikatów status.  
- Harmonogram game day i checklisty.  
- Metryki (MTTA/MTTR, alert fatigue) i progi.

## Wymagane streszczenia
- One‑pager: „jak reagować na alert” + kontakty.  
- Snapshot metryk on-call (MTTR, alerty/osoba, fatigue).

## Guidance (skrót)
- Ucz na realnych narzędziach i danych.  
- Ustal jasne progi eskalacji i komunikacji.  
- Dbaj o wellbeing: rotacje, quiet hours, load balancing.  
- Aktualizuj runbooki po każdym incydencie i ćwiczeniu.  
- Mierz MTTR/MTTA i redukuj alert noise.

## Szybkie powiązania
- linkage_index.jsonl (on_call/training)  
- incident_response_runbook, escalation_procedure_design, communication_plan, observability_plan

## Jak używać dokumentu
1. Przygotuj szkolenie, dostępy i runbooki.  
2. Przeprowadź ćwiczenia; oceń i popraw.  
3. Monitoruj metryki on-call; iteruj proces i linkage_index.

## Checklisty Definition of Ready (DoR)
- [ ] Matryca on-call i runbooki dostępne.  
- [ ] Narzędzia (pager/chat/ticketing/logi) skonfigurowane.  
- [ ] SLO/SLA i kryteria eskalacji znane.  
- [ ] Plan szkolenia/ćwiczeń ustalony.  
- [ ] Polityki dostępu/PII potwierdzone.

## Checklisty Definition of Done (DoD)
- [ ] Szkolenia/ćwiczenia wykonane; wyniki/feedback zapisane; status/wersja/data uzupełnione.  
- [ ] Runbooki zaktualizowane; wnioski z retro wdrożone.  
- [ ] Metryki on-call monitorowane; plany poprawy opisane.  
- [ ] Linkage_index uzupełniony; ryzyka/dec. udokumentowane.  
- [ ] Wellbeing (rotacje/limity) potwierdzone z zespołem.

## Definicje robocze
- MTTA/MTTR: Mean Time To Acknowledge/Resolve.  
- Alert fatigue: przeciążenie liczbą alertów.  
- Game day: ćwiczenie symulujące awarie.

## Przykłady użycia
- Szkolenie nowych SRE przed dołączeniem do rotacji.  
- Game day dla usługi krytycznej.  
- Retro po serii nocnych alertów.

## Ryzyka i ograniczenia
- Brak ćwiczeń → słaba reakcja.  
- Za dużo alertów → burnout/fatigue.  
- Nieaktualne runbooki → dłuższy MTTR.

## Decyzje i uzasadnienia
- Częstotliwość game day.  
- Limity alertów/osoba i rotacje.  
- Zakres uprawnień na dyżurach.

## Założenia
- Narzędzia i logi dostępne.  
- Zespół akceptuje zasady wellbeing.  
- SLO/SLA są zdefiniowane.

## Otwarte pytania
- Jak obsłużyć różne strefy czasowe?  
- Jak mierzyć sukces (MTTR, customer impact)?  
- Czy potrzebny pageduty/inna platforma?

## Powiązania z innymi dokumentami
- escalation_procedure_design — ścieżki eskalacji.  
- communication_plan — komunikaty.  
- observability_plan — alerty.

## Wymagane odwołania do standardów
- Wewnętrzne polityki bezpieczeństwa i dostępu, PII.  
- Standardy SRE/ITIL jeśli przyjęte.
