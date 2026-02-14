---
title: Product Quality Criteria
status: needs_content
---

# Product Quality Criteria

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zdefiniować mierzalne kryteria jakości produktu (funkcjonalne i niefunkcjonalne), aby kierować developmentem, testami i decyzjami go/no-go.

## Zakres i granice
- Obejmuje: kryteria funkcjonalne, użyteczność/UX, wydajność, niezawodność/dostępność, bezpieczeństwo/prywatność, zgodność/regulacje, dostępność (a11y), lokalizacja, obsługę błędów, analitykę i telemetry, maintainability/observability, dokumentację.
- Poza zakresem: szczegółowe test cases (osobne), roadmapa funkcji.

## Wejścia i wyjścia
- Wejścia: wymagania biznesowe, SLO/SLA, regulacje, standardy org (security/a11y), dane z badań UX, metryki prod.
- Wyjścia: lista kryteriów z metrykami/progami, mapping do wymagań/us stories, powiązanie z SLO i testami, szablon oceny release, checklisty QA.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: SLO/SLA, standardy security/a11y/privacy, wymagania regulacyjne, dane UX/telemetry, backlog/test plan; brak – odnotuj.

## Powiązania sekcja↔sekcja
Kryteria → metryki/testy; SLO → wydajność/dostępność; bezpieczeństwo → prywatność; UX → a11y.

## Fazy cyklu życia
Definicja → Review → Aktualizacje per release → Ewaluacja.

## Struktura sekcji (szkielet)
- Kategorie jakości i ich kryteria.
- Metryki/progi i sposób pomiaru.
- Mapowanie do wymagań/us stories i testów.
- Ocena release (go/conditional/no-go) i raport.
- Utrzymanie i rewizje kryteriów.
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- Tabele kryteriów z progami i SLI.
- Szablon oceny release.

## Wymagane streszczenia
- One-pager: top kryteria i progi + link do testów/SLO.

## Guidance
Cel: wspólny język jakości. DoR: SLO, standardy, wymagania, dane UX. DoD: kryteria/metyki/mapowanie/testy/ocena; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Zdefiniuj kryteria/metryki; zmapuj na wymagania/testy; używaj w review release; aktualizuj na podstawie danych prod.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] SLO/SLA; [ ] Standardy security/a11y/privacy; [ ] Wymagania/regulacje; [ ] Dane UX/telemetry.
- DoD: [ ] Kryteria/metryki/mapowanie/ocena opisane; [ ] Sekcje N/A uzasadnione; metadane aktualne.
