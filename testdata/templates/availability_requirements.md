---
title: Availability Requirements
status: needs_content
---

# Availability Requirements

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Określić wymagania dostępności (SLO/SLA) dla usług, aby kierować architekturą, monitoringiem i kontraktami.

## Zakres i granice
- Obejmuje: cele dostępności i tolerancje (SLO/SLA), definicje mierzenia (MTBF/MTTR/uptime), poziomy krytyczności, zależności, okna serwisowe, raportowanie, wyjątki, wymagania kontraktowe z dostawcami.
- Poza zakresem: szczegółowy design HA/DR (osobne dokumenty), metryki wydajności (oddzielne).

## Wejścia i wyjścia
- Wejścia: priorytety biznesowe, mapa usług i zależności, dane historyczne dostępności, wymagania prawne/kontraktowe, koszty i budżety, okna serwisowe.
- Wyjścia: tabela SLO/SLA per usługa, poziomy krytyczności, definicje pomiaru i wyłączeń, wymagania dla architektury/monitoringu/testów, oczekiwania wobec dostawców, plan raportowania.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: service catalog/CMDB, monitoring/SLI definicje, kontrakty vendorów, polityki change/okna serwisowe; brak – odnotuj.

## Powiązania sekcja↔sekcja
Krytyczność → SLO; zależności → targety; okna serwisowe → obliczanie uptime; raportowanie → wymagania arch/monitoringu.

## Fazy cyklu życia
Discovery → Definicja → Uzgodnienia → Publikacja → Raportowanie/rewizje.

## Struktura sekcji (szkielet)
- Zakres usług i krytyczność.
- Definicje SLO/SLA i SLI (metody pomiaru, wyłączenia, planowane okna).
- Tabela targetów (usługa → SLO/SLA → SLI → okna serwisowe → wyjątki).
- Wymagania architektoniczne (redundancja, HA/DR, capacity).
- Wymagania monitoring/testów (synthetics, alarms, chaos/failover testy).
- Raportowanie i eskalacje.
- Dostawcy i kontrakty (OOS/SLA vendorów).
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- SLI → definicje i źródła metryk.
- Vendor SLA → mapowanie na cele usług.

## Wymagane streszczenia
- Jedna strona: SLO/SLA per usługa + krytyczność.

## Guidance
Cel: jasne cele dostępności. DoR: katalog usług, dane historyczne, wymagania biznesowe/regulacyjne. DoD: SLO/SLA/SLI/tabela, wymagania arch/monitoringu, raportowanie; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Zbierz wymagania i dane; ustal krytyczność i SLO/SLA; zdefiniuj SLI/pomiar; opublikuj; raportuj i rewiduj.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Katalog usług i zależności; [ ] Dane historyczne; [ ] Wymagania biznesowe/regulacyjne; [ ] Okna serwisowe/zmiany.
- DoD: [ ] Tabela SLO/SLA/SLI; [ ] Wymagania arch/monitoringu; [ ] Raportowanie; [ ] Sekcje N/A uzasadnione; metadane aktualne.
