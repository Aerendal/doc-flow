---
title: Acceptance Criteria
status: needs_content
---

# Acceptance Criteria

## Metadane
- Właściciel: [Product/QA/Engineering Lead]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Ustala jednoznaczne kryteria akceptacji dla user story/feature/zmiany: kiedy uznajemy pracę za ukończoną, zgodną z wymaganiami, bezpieczną i gotową do release. Redukuje ryzyko niedomówień i zwrotów.

## Zakres i granice
- Obejmuje: kryteria funkcjonalne, niefunkcjonalne (wydajność, bezpieczeństwo, dostępność), dane/testy, UX/A11y, zgodność prawna, monitoring/telemetria, regresję/feature flags, kryteria rollout/backout, dowody (testy, logi, screeny).  
- Poza zakresem: pełny plan testów (osobny dokument), decyzje biznesowe o priorytecie backlogu.

## Wejścia i wyjścia
- Wejścia: user story/BRD, definicje Done, wymagania NFR, makiety/UX, dane testowe, polityki bezpieczeństwa/privacy, standardy A11y, dependency list, risk assessment.  
- Wyjścia: lista kryteriów akceptacji, dane testowe i dowody, status spełnienia (checklisty), decyzja go/no‑go dla release, aktualizacja JIRA/ALM i DoD.

## Powiązania (meta)
- Key Documents: qa_strategy_document, testing_vision_statement, non_functional_requirements, security_requirements, accessibility_compliance, release_checklist.  
- Key Document Structures: kryteria funkcjonalne, NFR, dane/testy, bezpieczeństwo/A11y, rollout/backout, dowody.  
- Document Dependencies: system testów (CI/CD), dane testowe, feature flags, monitoring/logi, ticketing.

## Zależności dokumentu
Wymaga: kompletnego opisu user story/feature, zdefiniowanych NFR, dostępnych danych testowych lub sposobu ich pozyskania, ustalonego środowiska i narzędzi testowych, uzgodnionych wymagań bezpieczeństwa/A11y. Braki = DoR otwarte.

## Powiązania sekcja↔sekcja
- Kryteria funkcjonalne + NFR → Testy/Dowody → Decyzja go/no‑go.  
- Dane testowe → Wyniki testów → Akceptacja.  
- Rollout/Backout → Monitoring → Stabilizacja.

## Fazy cyklu życia
- Definicja kryteriów (refinement).  
- Weryfikacja w trakcie implementacji/testów.  
- Akceptacja przed release; ewentualny backout/roll-forward.  
- Retrospektywa i aktualizacja wzorców kryteriów.

## Struktura sekcji
1) Kontekst i zakres user story/feature  
2) Kryteria funkcjonalne (given/when/then, edge cases)  
3) Kryteria niefunkcjonalne (wydajność, bezpieczeństwo, A11y, UX, dane)  
4) Dane i środowiska testowe (źródła, maskowanie, seed)  
5) Dowody i metryki (testy auto/manual, logi, screeny, metryki SLI)  
6) Rollout/feature flags i warunki backout  
7) Otwarte ryzyka/pytania, decyzja go/no‑go

## Wymagane rozwinięcia
- Lista kryteriów w formie Gherkin lub checklist.  
- Progi NFR (np. p95<300 ms, WCAG 2.1 AA, brak high/critical security issues).  
- Plan danych testowych (maskowanie, generowanie) i metody weryfikacji.

## Wymagane streszczenia
- Jednostronicowy snapshot do akceptacji: spełnione/niespełnione, ryzyka, zalecenie go/no‑go.  
- Krótka lista krytycznych edge cases.

## Guidance (skrót)
- Kryteria pisz z perspektywy użytkownika, ale sprawdzaj pokrycie NFR.  
- Wymagaj dowodów: link do testów, zrzuty, logi, metryki.  
- Ustal jasne progi go/no‑go i warunki backout.  
- Dopasuj dane testowe do scenariuszy krytycznych i ryzyk.  
- Synchronizuj z DoR/DoD i ticketami w ALM.

## Szybkie powiązania
- linkage_index.jsonl (acceptance/criteria)  
- qa_strategy_document, testing_vision_statement, non_functional_requirements, release_checklist, security_requirements

## Jak używać dokumentu
1. Zdefiniuj kryteria (funkcjonalne + NFR) wraz z danymi i dowodami.  
2. W trakcie testów odhaczaj spełnienie; zbieraj dowody.  
3. Przed release wypełnij snapshot, podejmij decyzję go/no‑go, zaktualizuj DoD.

## Checklisty Definition of Ready (DoR)
- [ ] User story/feature opisane; scope i zależności znane.  
- [ ] NFR i wymagania bezpieczeństwa/A11y zdefiniowane.  
- [ ] Dane testowe zidentyfikowane, sposób pozyskania uzgodniony.  
- [ ] Środowiska i narzędzia testowe dostępne.  
- [ ] Zdefiniowane warunki backout/rollout (feature flags).

## Checklisty Definition of Done (DoD)
- [ ] Wszystkie kryteria odhaczone; dowody dołączone.  
- [ ] NFR spełnione lub wyjątki zaakceptowane.  
- [ ] Rollout/backout i monitoring zaplanowane; status/wersja/data uzupełnione.  
- [ ] Ticket/ALM zaktualizowany, linkage_index uzupełniony.  
- [ ] Lessons learned/edge cases dopisane do repo wiedzy.

## Definicje robocze
- Acceptance Criteria: warunki uznania user story/feature za ukończone.  
- Backout: procedura powrotu do poprzedniej wersji, gdy kryteria nie są spełnione po deployu.  
- Evidence: obiektywny dowód spełnienia (testy, logi, metryki).

## Przykłady użycia
- Feature płatności: kryteria funkcjonalne + PCI + latency + A11y checkout.  
- API: kryteria kontraktowe (schema), error handling, rate limiting, observability.  
- Mobile: kryteria UX/A11y, wydajność na niskiej klasy urządzeniach.

## Ryzyka i ograniczenia
- Niejasne kryteria → zwroty i opóźnienia.  
- Brak danych testowych → fałszywe akceptacje.  
- Pominięte NFR → regresje wydajności/bezpieczeństwa/A11y.

## Decyzje i uzasadnienia
- Format kryteriów (Given/When/Then vs checklista) zależnie od zespołu.  
- Poziom dowodów wymagany dla produkcji vs testów wewnętrznych.  
- Warunki backout i czas obserwacji po deployu.

## Założenia
- Dostępne środowiska i dane testowe.  
- CI/CD uruchamia testy i zbiera logi/metryki.  
- Zespół ma standard DoR/DoD dla user story.

## Otwarte pytania
- Czy wymagane są dodatkowe testy zgodności (regulatory)?  
- Jakie są minimalne progi dla A11y/bezpieczeństwa?  
- Kto akceptuje wyjątki od NFR?

## Powiązania z innymi dokumentami
- non_functional_requirements — progi i wskaźniki.  
- release_checklist — kroki i dowody.  
- security_requirements — kontrole security.

## Wymagane odwołania do standardów
- WCAG/EN/ADA (A11y), polityki bezpieczeństwa i prywatności.  
- Wewnętrzne standardy jakości i testów.
