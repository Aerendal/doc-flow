---
title: Non-Functional Requirements
status: needs_content
---

# Non-Functional Requirements

## Metadane
- Właściciel: [Architecture/Engineering/QA]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Definiuje mierzalne wymagania niefunkcjonalne (NFR) dla rozwiązania/systemu: wydajność, dostępność, bezpieczeństwo, skalowalność, użyteczność, obserwowalność, zgodność i operowalność. Kieruje decyzjami architektonicznymi, projektowaniem testów i kryteriami akceptacji releasu.

## Zakres i granice
- Obejmuje: SLO/SLA/SLI, opóźnienia/przepustowość, dostępność/HA/DR, odporność, bezpieczeństwo/prywatność, A11y/UX, zgodność (regulacje/branża), obserwowalność (metryki/logi/traces), operowalność (deploy, rollback, monitoring), limity środowiskowe i dane.  
- Poza zakresem: funkcjonalne wymagania biznesowe (FRS/BRD), szczegółowe projekty komponentów (osobne dokumenty).

## Wejścia i wyjścia
- Wejścia: BRD/FRS, ryzyka biznesowe, profile obciążenia, polityki bezpieczeństwa/compliance, dane/PII klasyfikacja, standardy org, budżety (koszt/latency), zależności zewnętrzne.  
- Wyjścia: lista NFR z metrykami i progami, mapowanie do testów (perf/resilience/security/A11y), wymagania architektoniczne (HA/DR/cache/queue), kryteria go/no‑go, aktualizacje do planów monitoringu i runbooków.

## Powiązania (meta)
- Key Documents: architecture_decision_records, performance_test_plan, resilience_testing_plan, security_requirements, accessibility_compliance, observability_architecture, dr_plan.  
- Key Document Structures: SLI/SLO/SLA, wydajność, HA/DR, bezpieczeństwo/prywatność, A11y/UX, obserwowalność, operowalność.  
- Document Dependencies: monitoring stack, CI/CD, load profiles, threat models, data classification, CMDB zależności.

## Zależności dokumentu
Wymaga: znanych scenariuszy biznesowych i profili ruchu, polityk security/privacy, wymagań regulacyjnych, danych o zależnościach zewnętrznych, wstępnej architektury logicznej. Braki = DoR otwarte.

## Powiązania sekcja↔sekcja
- SLI/SLO → Testy wydajności/odporności → Kryteria go/no‑go.  
- Bezpieczeństwo/PII → Kontrole/pen-test → Monitoring/alarmy.  
- HA/DR → Architektura → Testy odtwarzania → Runbooki.

## Fazy cyklu życia
- Definicja NFR (ideation/refinement).  
- Weryfikacja w design review/ADR.  
- Testy i walidacja w CI/CD i przed releasem.  
- Audyty i przeglądy okresowe; aktualizacja progów i SLO.

## Struktura sekcji
1) Kontekst systemu i zakres NFR  
2) SLI/SLO/SLA i budżet błędów  
3) Wydajność i skalowalność (latency, throughput, burst, concurrency)  
4) Dostępność, odporność, HA/DR (RTO/RPO, failover, degradacja)  
5) Bezpieczeństwo i prywatność (authn/z, szyfrowanie, PII, audyt)  
6) Użyteczność i A11y (WCAG/UX)  
7) Obserwowalność i operowalność (metryki/logi/traces, alerty, deploy/rollback)  
8) Zgodność/regulacje (branża, dane, jurysdykcje)  
9) Dane i ograniczenia środowiskowe (rozmiar, retencja, edge cases)  
10) Kryteria akceptacji releasu, ryzyka, decyzje

## Wymagane rozwinięcia
- Tabela SLI/SLO/SLA z progami i metodą pomiaru.  
- Profile obciążenia i scenariusze testów wydajności/odporności.  
- Kontrole bezpieczeństwa (OWASP/ISO/NIST) i wymagania audytu.  
- Plan HA/DR z RTO/RPO i testami.

## Wymagane streszczenia
- Executive snapshot NFR: top SLO, progi, największe ryzyka, luki.  
- Krótka karta HA/DR (RTO/RPO, scenariusze failover).

## Guidance (skrót)
- NFR muszą być mierzalne, testowalne, powiązane z KPI biznesu.  
- Ustal budżet błędów i konsekwencje jego przekroczenia.  
- Testuj pod SLO, nie tylko pod „brak błędów”.  
- Integruj NFR z monitoringiem: te same metryki w testach i prod.  
- Aktualizuj NFR po zmianach architektury lub ruchu.

## Szybkie powiązania
- linkage_index.jsonl (non_functional/requirements)  
- performance_test_plan, resilience_testing_plan, security_requirements, accessibility_compliance, observability_architecture, dr_plan

## Jak używać dokumentu
1. Zdefiniuj SLI/SLO i kluczowe domeny NFR dla systemu.  
2. Ustal progi i metody pomiaru; powiąż je z testami i monitoringiem.  
3. Zweryfikuj w design review, włącz do planu testów i release gate; aktualizuj DoR/DoD.

## Checklisty Definition of Ready (DoR)
- [ ] Scenariusze biznesowe i profile ruchu zebrane.  
- [ ] Polityki security/privacy/regulacje zidentyfikowane.  
- [ ] Wstępne SLI/SLO zdefiniowane; metody pomiaru dostępne.  
- [ ] Zależności zewnętrzne i ograniczenia środowiskowe opisane.  
- [ ] Plan testów niefunkcjonalnych i narzędzia uzgodnione.

## Checklisty Definition of Done (DoD)
- [ ] SLI/SLO/SLA wypełnione, progi zatwierdzone.  
- [ ] Testy wydajności/odporności/bezpieczeństwa/A11y wykonane lub wyjątki zaakceptowane.  
- [ ] HA/DR opisane i przetestowane (failover/fallback).  
- [ ] Monitoring/alerty i logi pod SLO działają; status/wersja/data zaktualizowane.  
- [ ] Linki do ADR/test planów/runbooków dodane; luki NFR odnotowane i przypisane.
