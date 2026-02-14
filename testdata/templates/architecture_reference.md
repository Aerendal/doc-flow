---
title: Architecture Reference
status: needs_content
---

# Architecture Reference

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Dostarczyć referencyjne artefakty architektoniczne (diagramy, standardy, decyzje) dla danego systemu/portfolio, ułatwiając przeglądy i utrzymanie spójności.

## Zakres i granice
- Obejmuje: kontekst biznesowy, wymagania niefunkcjonalne, diagramy (context/container/component/deployment), decyzje (ADR), standardy bezpieczeństwa/observability/CI-CD, zależności zewnętrzne, wersjonowanie artefaktów, linki do implementacji.
- Poza zakresem: szczegółowy design pojedynczych feature’ów (osobne), backlog produktu.

## Wejścia i wyjścia
- Wejścia: cele systemu, NFR, katalog usług, standardy architektoniczne, decyzje ADR, zależności, SLO.
- Wyjścia: pakiet referencyjny (diagramy, ADR, standardy), mapa zależności, checklisty review, plan przeglądów.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: katalog usług, SLO, ADR, standardy bezpieczeństwa/observability, narzędzia diagramów, repozytoria; brak – odnotuj.

## Powiązania sekcja↔sekcja
NFR → architektura; ADR → standardy; zależności → diagramy; SLO → komponenty krytyczne.

## Fazy cyklu życia
Opracowanie → Publikacja → Przeglądy okresowe → Aktualizacje.

## Struktura sekcji (szkielet)
- Kontekst i NFR.
- Diagramy (C4 lub inne) i wersje.
- Decyzje architektoniczne (ADR) i uzasadnienia.
- Standardy (security, observability, CI/CD, data).
- Zależności zewnętrzne i kontrakty.
- Wersjonowanie artefaktów i repo.
- Plan przeglądów i checklisty.
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- Diagramy → linki/źródła; ADR → format i repo.
- Checklisty review → punkty kontroli.

## Wymagane streszczenia
- One-pager: diagram kontekstowy + top ADR + SLO/NFR.

## Guidance
Cel: jednolite artefakty dla zespołów. DoR: cele/NFR, ADR, standardy, katalog usług. DoD: diagramy/ADR/standardy/zależności/review plan; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Uzupełnij kontekst/NFR, diagramy, ADR, standardy; utrzymuj wersje; używaj checklist w review.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Cele/NFR; [ ] ADR; [ ] Standardy; [ ] Katalog usług.
- DoD: [ ] Diagramy/ADR/standardy/zależności/review plan; [ ] Sekcje N/A uzasadnione; metadane aktualne.
