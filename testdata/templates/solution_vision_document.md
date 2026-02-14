---
title: Solution Vision Document
status: needs_content
---

# Solution Vision Document

## Metadane
- Właściciel: [Product/Architecture]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Opisuje wizję rozwiązania: problem/opportunity, wartość dla klientów/biznesu, zakres i kierunek architektury, kryteria sukcesu. Ma spajać zespoły zanim powstaną szczegółowe wymagania i design.

## Zakres i granice
- Obejmuje: problem/opportunity, persony/use cases/JTBD, propozycję wartości, cele/KPI/KR, zakres in/out, kluczowe założenia/ograniczenia, kierunek architektury (high-level), zależności, ryzyka, fazy i kamienie milowe, kryteria sukcesu i miary startowe.
- Poza zakresem: szczegółowy design komponentów, backlog user stories, plan sprintów.

## Wejścia i wyjścia
- Wejścia: insighty z badań/rynku/klientów, dane produktowe/operacyjne, KPI biznesowe, mapy procesów, strategia firmy/produktu, ograniczenia techniczne/prawne/finansowe, polityki security/compliance, dostępne capability/platformy.
- Wyjścia: karta wizji (one-pager), mapa architektury high-level, lista założeń i ryzyk, zakres/fazy i kamienie milowe, kryteria sukcesu/KPI, DoR/DoD dla kolejnych artefaktów (BRD/PRD/architektura), ścieżka komunikacji/alignment.

## Powiązania (meta)
- Key Documents: business_value_proposition, product_strategy_document, technology_strategy, market_analysis, architecture_vision, stakeholder_requirements, risk_register, roadmap, pricing_strategy, go_to_market_strategy.
- Key Document Structures: problem → wartość → cele/KPI → zakres → architektura → fazy → ryzyka/założenia → decyzje.
- Document Dependencies: KPI dashboardy, mapy architektury/CMDB, polityki prawne/security/privacy, budżet/finansowanie, dostępne capability/platformy.

## Zależności dokumentu
Wymaga: zdefiniowanych problemów/okazji i KPI, wyników badań/rynku, głównych ograniczeń, zaangażowania kluczowych interesariuszy. Braki = DoR otwarte.

## Powiązania sekcja↔sekcja
- Problem/okazja → Propozycja wartości → Cele/KPI → Zakres → Architektura high-level → Fazy.
- Ryzyka/założenia → Decyzje/alternatywy → Roadmapa → Kryteria sukcesu.

## Fazy cyklu życia
- Definicja wizji i alignment interesariuszy.
- Komunikacja i decyzje go/no-go dla kolejnych artefaktów (BRD/PRD/Architecture Concept).
- Aktualizacje przy kamieniach milowych (pilot, MVP, GA, pivot).

## Struktura sekcji
1) Problem i kontekst (rynek, użytkownicy, okazja)
2) Persony/use cases/JTBD i propozycja wartości
3) Cele/KPI/KR i kryteria sukcesu
4) Zakres in/out i ograniczenia (techniczne/prawne/finansowe)
5) Architektura high-level (diagram kontekst/komponenty, integracje, dane)
6) Zależności, ryzyka i założenia
7) Fazy i kamienie milowe (pilot/MVP/GA), kryteria go/no-go
8) Decyzje, alternatywy i otwarte pytania

## Wymagane rozwinięcia
- One‑page diagram architektury high-level i mapa interesariuszy.
- Lista założeń i ryzyk z priorytetem oraz właścicielami.
- Kryteria sukcesu, metryki startowe i docelowe; plan pomiaru.
- Komunikacja/alignment plan (fora, cadence, materiały).

## Wymagane streszczenia
- Executive one‑pager: problem, wartość, segmenty/persony, cele/KPI, fazy, ryzyka.
- Karta architektury high-level (diagram + kluczowe decyzje/ograniczenia).

## Guidance (skrót)
- Skup się na problemie/wartości i mierzalnych celach, nie na szczegółach implementacji.
- Określ klarowny in/out of scope; uwzględnij ograniczenia/regulacje.
- Zapewnij alignment: właściciele, fora decyzyjne, plan komunikacji, kryteria go/no-go.
- Dokumentuj założenia/ryzyka i aktualizuj przy kamieniach milowych.

## Szybkie powiązania
- product_strategy_document, business_value_proposition, market_analysis, technology_strategy, architecture_vision, stakeholder_requirements, risk_register, roadmap, pricing_strategy, go_to_market_strategy

## Jak używać dokumentu
1. Zbierz problem/wartość i KPI; opisz zakres.  
2. Dodaj diagram high-level i fazy; uzgodnij z interesariuszami.  
3. Aktualizuj w kamieniach milowych; uzupełnij DoR/DoD i linkage_index.

## Checklisty Definition of Ready (DoR)
- [ ] Problem/opportunity i KPI zdefiniowane.  
- [ ] Persony/scenariusze opisane.  
- [ ] Ograniczenia/założenia zebrane.  
- [ ] Wstępny diagram architektury i fazy.  
- [ ] Interesariusze zidentyfikowani.

## Checklisty Definition of Done (DoD)
- [ ] Wizja opisana; status/wersja/data uzupełnione.  
- [ ] Cele/KPI i zakres/fazy uzgodnione.  
- [ ] Diagram high-level i ryzyka/założenia opublikowane.  
- [ ] Linkage_index zaktualizowany; decyzje zapisane.  
- [ ] Plan przeglądów/aktualizacji ustalony.

## Definicje robocze
- Vision: opis „dlaczego i co”, zanim „jak”.  
- Scope: co robimy i czego nie robimy w tej iteracji.

## Przykłady użycia
- Nowa platforma danych.  
- Replatforming legacy.  
- Wprowadzenie nowej linii produktu.

## Ryzyka i ograniczenia
- Niejasny scope → creep.  
- Cele bez metryk → brak oceny sukcesu.  
- Brak alignment → sprzeczne decyzje.

## Decyzje i uzasadnienia
- Priorytety faz.  
- Architektura docelowa high-level.  
- Kryteria sukcesu.

## Założenia
- Dane i interesariusze dostępni.  
- Budżet/własność zatwierdzone.  
- Strategia firmy stabilna.

## Otwarte pytania
- Jakie są zależności z innymi programami?  
- Jakie ryzyka prawne/regulacyjne?  
- Jakie są limity budżetu/czasu?

## Powiązania z innymi dokumentami
- product_strategy — kierunek biznesu.  
- architecture_vision — kierunek tech.  
- risk_register — ryzyka.

## Wymagane odwołania do standardów
- Wewnętrzne standardy architektoniczne i procesowe.  
- Polityki prawne/compliance jeśli wpływają.
