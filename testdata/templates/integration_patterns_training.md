---
title: Integration Patterns Training
status: needs_content
---

# Integration Patterns Training

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Przygotować program szkoleniowy z wzorców integracyjnych (API/event/ETL/ESB), aby zespoły stosowały spójne podejście do integracji systemów.

## Zakres i granice
- Obejmuje: patterny (request/response, pub-sub, CQRS, event sourcing, batch/stream, saga/choreography/orchestration, idempotencja, retry, dead-letter), kontrakty API i eventy, bezpieczeństwo, monitoring, testy integracyjne, narzędzia (gateway, ESB, message broker).
- Poza zakresem: szczegółowe projekty pojedynczych integracji (osobne).

## Wejścia i wyjścia
- Wejścia: katalog integracji w organizacji, standardy API/event, narzędzia (gateway/broker/ETL), wymagania bezpieczeństwa/PII, przykłady defektów.
- Wyjścia: sylabus modułów, materiały i laby, checklisty DoR/DoD, scenariusze ćwiczeń (API/event/batch), kryteria oceny, harmonogram sesji.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: standardy integracji, katalog integracji, narzędzia, polityki security/PII, monitoring, testy; brak – odnotuj.

## Powiązania sekcja↔sekcja
Patterny → laby → ocena; bezpieczeństwo → PII; monitoring → testy.

## Fazy cyklu życia
Plan → Materiały/laby → Sesje → Ocena → Utrzymanie programu.

## Struktura sekcji (szkielet)
- Cele i KPI szkolenia.
- Moduły (API patterns, event patterns, batch/ETL, reliability/retry/idempotency, security, monitoring/testing).
- Ćwiczenia i laby (gateway/broker/ETL).
- Narzędzia i środowiska.
- Ocena/certyfikacja i feedback.
- Plan utrzymania i aktualizacji.

## Wymagane rozwinięcia
- Laby → repo i dane.
- Checklisty → DoR/DoD integracji.

## Wymagane streszczenia
- One-pager: moduły, KPI, harmonogram.

## Guidance
Cel: wspólne wzorce integracji. DoR: standardy/narzędzia, katalog integracji, polityki PII. DoD: sylabus/laby/checklisty/ocena; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Dostosuj moduły do grupy; przygotuj laby; przeprowadź sesje; oceń; aktualizuj program.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Standardy integracji i narzędzia; [ ] Katalog integracji; [ ] Polityki security/PII.
- DoD: [ ] Sylabus/laby/checklisty/ocena gotowe; [ ] Sekcje N/A uzasadnione; metadane aktualne.
