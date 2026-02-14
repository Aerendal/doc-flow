---
title: Monitoring Strategy Document
status: needs_content
---

# Monitoring Strategy Document

## Metadane
- Właściciel: [SRE/Platform/Observability]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Ustala strategię monitoringu i obserwowalności dla systemów/aplikacji: co mierzymy, po co, jakim narzędziem, z jakimi progami i ścieżką reakcji. Ma zapewnić spójność KPI/SLO, minimalizować MTTR i blind‑spoty oraz wspierać audyty/regulacje.

## Zakres i granice
- Obejmuje: metryki biznesowe i techniczne, SLO/SLA/SLI, logi, tracery, alerting, dashboards, runbooki, właścicieli, standardy tagów i retencji, budżet kosztów, wymagania compliance (PII, audyt), proces przeglądów i ciągłego doskonalenia.
- Poza zakresem: szczegółowe konfiguracje narzędzi (osobne runbooki), procedury IR (link do incident_response), projekt infrastruktury (link do architecture docs).

## Wejścia i wyjścia
- Wejścia: katalog usług (CMDB), krytyczne ścieżki i zależności, SLO/SLA, mapa KPI biznesowych, standardy tagowania, wymagania compliance/audytu, budżet kosztów observability, narzędzia (Prometheus/Grafana/ELK/APM).
- Wyjścia: standard monitoringu (metryki/logi/traces), lista SLI/SLO z progami, macierz pokrycia (service x signal), standard alertów i eskalacji, dashboardy referencyjne, harmonogram przeglądów, plan optymalizacji kosztów i retencji.

## Powiązania (meta)
- Key Documents: incident_response_runbook, service_level_objectives, observability_architecture, logging_standards, alerting_policy, cost_management_observability.
- Key Document Structures: sygnały (metrics/logs/traces), SLI/SLO, alerting, dashboardy, runbooki, koszt/retencja.
- Document Dependencies: CMDB/usługi, katalog zależności, narzędzia monitoringowe, system ticketowy, on‑call rota, polityki bezpieczeństwa danych.

## Zależności dokumentu
Wymaga: aktualnego CMDB/katalogu usług, zdefiniowanych SLO/SLA, przyjętych standardów tagowania i retencji, dostępów do narzędzi monitoringu/logowania/APM. Braki = DoR otwarte.

## Powiązania sekcja↔sekcja
- SLI/SLO → Alerty → Runbooki → Eskalacje.  
- Katalog usług → Macierz pokrycia → Luki → Plan wdrożeń.  
- Retencja/koszty → Standardy logów/metryk → Budżet/limity.

## Fazy cyklu życia
- Definicja strategii i priorytetów (services tiering).  
- Rollout standardów monitoringu i alertów.  
- Ciągłe przeglądy: coverage, fałszywe alarmy, koszty, audyty.  
- Ewolucja narzędzi/architektury observability.

## Struktura sekcji
1) Kontekst i cele (SLO, MTTR, koszty)  
2) Zakres usług i tiering (Tier0/1/2, krytyczne ścieżki)  
3) Sygnały i standardy (metryki/logi/traces, tagi, sampling, PII)  
4) SLI/SLO/SLA i progi alertów (error rate, latency, dostępność, capacity)  
5) Architektura i narzędzia (collectors, storage, viz, alerting)  
6) Macierz pokrycia (usługa × sygnał × ownership)  
7) Runbooki i eskalacje (on‑call, playbooki, testy alarmów)  
8) Koszty i retencja (budżet, sampling, indeksy, zimne archiwum)  
9) Przeglądy i doskonalenie (AAR, tuning alertów, chaos/DR drills)  
10) Ryzyka, decyzje, zależności

## Wymagane rozwinięcia
- Tabela SLI/SLO per usługa z progami i metodą pomiaru.  
- Standard tagów (service, team, env, region, version, customer/PII flag).  
- Macierz pokrycia i luki wraz z planem wdrożeń i priorytetem.

## Wymagane streszczenia
- Executive snapshot: stan SLO, liczba luk w pokryciu, najdroższe źródła danych, plany redukcji kosztów.  
- Run sheet testu alertów (fire‑drill) z wynikami.

## Guidance (skrót)
- Projektuj „alerty działające na biznes”: mało, głośno, na SLO; reszta jako raporty.  
- Standaryzuj metryki i tagi, inaczej porównywalność spada.  
- Regularnie pruning logów (PII, koszt) i testuj retencję/restaurację.  
- Utrzymuj single source of truth dla SLO i macierzy pokrycia.  
- Automatyzuj testy alertów (synthetics, chaos, replay).

## Szybkie powiązania
- linkage_index.jsonl (monitoring/strategy/document)  
- incident_response_runbook, service_level_objectives, observability_architecture, logging_standards, alerting_policy, cost_management_observability

## Jak używać dokumentu
1. Ustal tiering usług i listę SLI/SLO; dopisz progi i alerty.  
2. Zweryfikuj macierz pokrycia i zaplanuj wdrożenia dla luk.  
3. Ustal retencję/koszt i harmonogram przeglądów; zaktualizuj runbooki i DoR/DoD.

## Checklisty Definition of Ready (DoR)
- [ ] CMDB i zależności usług aktualne.  
- [ ] Wstępne SLI/SLO zdefiniowane dla Tier0/1.  
- [ ] Standard tagów/logów/metryk uzgodniony.  
- [ ] On‑call i kanały eskalacji zdefiniowane.  
- [ ] Dostępy do narzędzi monitoringowych nadane.

## Checklisty Definition of Done (DoD)
- [ ] Sekcje wypełnione lub N/A z uzasadnieniem.  
- [ ] SLI/SLO opublikowane, alerty działają i są przetestowane.  
- [ ] Macierz pokrycia i luki z planem wdrożeń opublikowane.  
- [ ] Koszty/retencja opisane i uzgodnione z FinOps/bezpieczeństwem.  
- [ ] Runbooki i ścieżki eskalacji podlinkowane; status/wersja/data zaktualizowane.

## Definicje robocze
- SLI: mierzalny wskaźnik jakości usługi (np. availability 99.9%, latency p95).  
- SLO: cel dla SLI w okresie (np. 99.9% / 28 dni).  
- Error budget: 1 − SLO; budżet na zmiany/awarie.

## Przykłady użycia
- Zmiana architektury logowania — ocena kosztów i tagów.  
- Nowa usługa Tier1 — nadanie SLI/SLO i alertów.  
- Post‑mortem fałszywych alarmów — tuning progów i reguł.

## Ryzyka i ograniczenia
- Alert fatigue z nadmiarem reguł lub złymi progami.  
- Brak standardu tagów uniemożliwia pivotowanie danych.  
- Niekontrolowane koszty retencji/indeksów.

## Decyzje i uzasadnienia
- Zakres SLO (global vs per region) — zależnie od architektury.  
- Retencja logów/traces — kompromis koszt vs potrzeba audytu/IR.  
- Sampling/aggregation — kompromis dokładność vs koszt.

## Założenia
- Stabilne źródła metryk/logów/traces i kontrola PII.  
- On‑call rota dostępna i aktualna.  
- Narzędzia wspierają etykiety/tagi i multi‑region.

## Otwarte pytania
- Czy wszystkie SLO muszą być customer‑facing czy tylko wewnętrzne?  
- Jakie synthetic tests są wymagane per krytyczna ścieżka?  
- Jakie limity kosztów są akceptowalne per usługa?

## Powiązania z innymi dokumentami
- incident_response_runbook — reakcja na alerty.  
- logging_standards — formaty i PII.  
- cost_management_observability — budżet i optymalizacje.

## Wymagane odwołania do standardów
- ISO 27001 / SOC2 (logowanie, audyt).  
- Wewnętrzne standardy PII/RODO i retencji.
