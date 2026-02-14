---
title: Automation Strategy
status: needs_content
---

# Automation Strategy

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Określić strategię automatyzacji (IT/operacje/biznes/RPA), aby zwiększyć efektywność, jakość i zgodność, przy kontroli kosztów i ryzyk.

## Zakres i granice
- Obejmuje: katalog procesów do automatyzacji, kryteria wyboru (wartość/ryzyko/kompleksowość), techniki (script/IaC/CI-CD/RPA/API), architekturę i standardy, governance (waivery/exception), bezpieczeństwo/PII, mierzenie ROI/KPI, plan wdrożeń i kompetencji.
- Poza zakresem: szczegółowe buildy botów/skryptów (osobne), pełny plan transformacji organizacyjnej.

## Wejścia i wyjścia
- Wejścia: procesy biznesowe/IT, pain points, metryki koszt/czas/błędy, polityki bezpieczeństwa, budżet, ograniczenia prawne, dostępne narzędzia.
- Wyjścia: roadmap automatyzacji, standardy i guardrails, RACI i role (CoE), backlog inicjatyw z priorytetami, KPI (koszt/czas/defekty), plan narzędzi i szkoleń.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: katalog procesów, narzędzia (CI/CD/RPA/IaC), polityki security/PII, budżet, governance, teamy/CoE; brak – odnotuj.

## Powiązania sekcja↔sekcja
Kryteria → backlog; standardy → governance; KPI → roadmap; bezpieczeństwo → narzędzia.

## Fazy cyklu życia
Discovery → Roadmap → Pilotaż → Skalowanie → Utrzymanie/rewizje.

## Struktura sekcji (szkielet)
- Cele i KPI automatyzacji.
- Kryteria wyboru procesów i scoring.
- Techniki i narzędzia (script/API/IaC/CI-CD/RPA/ML).
- Architektura i standardy (naming, repo, logi, audyt).
- Governance i wyjątki (waivery, risk assessment).
- Bezpieczeństwo/PII i compliance.
- Backlog i roadmapa (priorytety, ETA).
- Plan kompetencji i szkoleń.
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- Scoring → wzór i progi.
- Standardy → szablony repo/CI/CD.

## Wymagane streszczenia
- One-pager: top inicjatywy, KPI, timeline.

## Guidance
Cel: skoordynowana automatyzacja. DoR: katalog procesów, polityki, budżet, narzędzia. DoD: kryteria/roadmap/standardy/governance/KPI; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Oceń procesy, ustal scoring i backlog; wybierz techniki i standardy; zaplanuj roadmapę; zarządzaj wyjątkami; mierz KPI.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Katalog procesów; [ ] Polityki security/PII; [ ] Budżet i narzędzia; [ ] Zespoły/CoE.
- DoD: [ ] Scoring/backlog/roadmap; [ ] Standardy/governance/KPI; [ ] Sekcje N/A uzasadnione; metadane aktualne.
