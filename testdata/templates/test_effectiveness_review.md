---
title: Test Effectiveness Review
status: needs_content
---

# Test Effectiveness Review

## Metadane
- Właściciel: [QA/Engineering]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Ocena skuteczności testów: jak dobrze wykrywają defekty, czy pokrywają ryzyka i ścieżki krytyczne, oraz jakie działania poprawiają jakość. Ma zmniejszyć defect leakage i flakiness oraz zoptymalizować koszty testów.

## Zakres i granice
- Obejmuje: metryki skuteczności (defect leakage, detection effectiveness, coverage risk-based, flake rate), analiza defektów po release, mapowanie testów do ryzyk/ś­cieżek, koszty/test time, jakości danych testowych, stabilność środowisk, rekomendacje (gap fill, automatyzacja, usunięcie flake), plan działania.
- Poza zakresem: szczegółowe przypadki testowe (repo), strategia QA (link).

## Wejścia i wyjścia
- Wejścia: raporty testów, defekty (pre/post-release), SLO/SLA, coverage map, flake log, koszt/test time, dane testowe, środowiska, release data.
- Wyjścia: ocena skuteczności, luki i rekomendacje, plan działań, zaktualizowane KPI/progi, raport dla release/leadership.

## Powiązania (meta)
- Key Documents: qa_strategy, testing_plan_schedule, quality_objectives, performance_metrics, system_monitoring_strategy, incident_postmortems.
- Key Document Structures: metryki, analiza, rekomendacje, plan.
- Document Dependencies: test reports, defect tracker, monitoring, flake tracker, coverage map, data/test env info.

## Zależności dokumentu
Wymaga danych defektów pre/post, raportów testów, mapy coverage/ryzyk, flake log, danych środowisk i release. Bez tego DoR otwarte.

## Powiązania sekcja↔sekcja
- Metryki → Analiza → Rekomendacje → Plan → KPI/progi.
- Defekty po release → Luki coverage/dane/środowisko → Działania.

## Fazy cyklu życia
- Zbieranie danych/metryk.
- Analiza skuteczności i luk.
- Plan działań i zmiana KPI/progów.
- Realizacja i follow-up; kolejny cykl.

## Struktura sekcji
1) Metryki i dane wejściowe (defect leakage, detection rate, flake, coverage, koszt/test time)  
2) Analiza luk (ryzyka/ścieżki, dane, środowiska, test types)  
3) Rekomendacje i priorytety (gap fill, automatyzacja, usuwanie flake, dane)  
4) Plan działań i KPI/progi (owner, ETA, oczekiwany wpływ)  
5) Raport i komunikacja (release/leadership)  
6) Ryzyka, decyzje, open issues

## Wymagane rozwinięcia
- Metryki z ostatnich cykli; mapa defektów do testów; lista flake; koszty/test time.
- Rekomendacje z wpływem; plan działań i KPI/progi.

## Wymagane streszczenia
- Kluczowe metryki, top luki, top rekomendacje/akcje, zmiany KPI/progów.

## Guidance (skrót)
- Mierz defect leakage i flake; linkuj defekty do brakujących/flake testów.
- Naprawiaj dane/test env, gap fill krytyczne ścieżki, usuwaj flake.
- Ustal KPI/progi realistyczne; iteruj cyklicznie po release.

## Szybkie powiązania
- linkage_index.jsonl (qa/test_effectiveness)
- qa_strategy, testing_plan_schedule, quality_objectives, performance_metrics, system_monitoring_strategy, incident_postmortems

## Jak używać dokumentu
1. Zbierz metryki/defekty i coverage; wypełnij sekcje 1–2.
2. Dodaj rekomendacje i plan działań; zaktualizuj KPI/progi.
3. Raportuj do release/leadership; dodaj do linkage_index; powtarzaj cyklicznie.

## Checklisty Definition of Ready (DoR)
- [ ] Dane defektów pre/post, test reports, flake log, coverage map dostępne.
- [ ] Struktura sekcji wypełniona/N/A.

## Checklisty Definition of Done (DoD)
- [ ] Analiza skuteczności i luki opisane; rekomendacje i plan działań gotowe.
- [ ] KPI/progi zaktualizowane; dokument w linkage_index; wersja/data/właściciel aktualne.

## Definicje robocze
- Defect leakage, Detection effectiveness, Flake rate, Coverage, MTTR, SLA/SLO.

## Przykłady użycia
- Po release: wysokie P1 prod → brak testów na ścieżkę X → dodaj testy + dane + alert.
- Flake 10% → refaktoryzacja testów + stabilizacja środowiska.

## Ryzyka i ograniczenia
- Brak danych → brak wniosków; ignorowanie flake → fałszywa wiarygodność; brak follow-up działań.

## Decyzje i uzasadnienia
- [Decyzja] KPI/progi — uzasadnienie ryzyk/SLO; [Decyzja] Priorytety działań — uzasadnienie wpływu.

## Założenia
- Dane metryczne dostępne; zespoły gotowe wdrożyć działania.

## Otwarte pytania
- Jaka kadencja przeglądów (po release vs kwartalnie)? 
- Czy uwzględniać koszty/test time w priorytetyzacji?

## Powiązania z innymi dokumentami
- QA Strategy, Testing Plan, Quality Objectives, Performance Metrics, Monitoring Strategy, Incident Postmortems.

## Powiązania z sekcjami innych dokumentów
- Testing → coverage/testy; Quality Objectives → KPI; Monitoring → obserwacje po release.

## Słownik pojęć w dokumencie
- Defect leakage, Detection effectiveness, Flake rate, Coverage, MTTR, SLA/SLO.

## Wymagane odwołania do standardów
- Polityki QA, SLA/SLO, standardy raportowania jakości.

## Mapa relacji sekcja→sekcja
- Metryki → Luki → Rekomendacje → Plan → Raport → Korekta KPI.

## Mapa relacji dokument→dokument
- Test Effectiveness → QA/Testing/Monitoring/Quality Objectives → Release/IR.

## Ścieżki informacji
- Dane → Analiza → Plan → Raport → Follow-up → Kolejny cykl.

## Weryfikacja spójności
- [ ] Metryki/luki/rekomendacje/plan opisane; dokument w linkage_index.

## Lista kontrolna spójności relacji
- [ ] Każdy wniosek ma dane; każda rekomendacja ma owner/ETA/KPI wpływ.
- [ ] Relacje cross‑doc opisane z uzasadnieniem.

## Artefakty powiązane
- Test reports, defect logs, flake logs, coverage maps, KPI dashboards, action plan.
