---
title: Test Data Preparation Plan
status: needs_content
---

# Test Data Preparation Plan

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zaprojektować przygotowanie danych testowych (functional/perf/security) zapewniając realizm, zgodność z PII i powtarzalność.

## Zakres i granice
- Obejmuje: źródła danych, maskowanie/anonymizacja, generatory syntetyczne, subset sampling, dane do testów wydajności, zgodność prawna (PII/PCI), wersjonowanie zestawów, walidacja jakości danych testowych, dystrybucję i cleanup.
- Poza zakresem: szczegółowe test cases (osobne), produkcyjne dane surowe (bezpośrednie użycie zabronione jeśli PII).

## Wejścia i wyjścia
- Wejścia: wymagania testów, profile obciążeń, polityki PII/PCI, schema danych, narzędzia masking/generacji, zasoby storage.
- Wyjścia: plan zestawów danych, procedury maskowania/generacji, walidacja, wersje i repo danych, checklisty użycia/cleanup, raport zgodności.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: źródła danych, polityki PII/PCI, narzędzia masking/generacji, storage, CI/CD integracje; brak – odnotuj.

## Powiązania sekcja↔sekcja
Wymagania testów → zestawy danych; PII → maskowanie; wersjonowanie → walidacja; cleanup → bezpieczeństwo.

## Fazy cyklu życia
Planowanie → Przygotowanie → Walidacja → Dystrybucja → Użycie → Cleanup → Rewizje.

## Struktura sekcji (szkielet)
- Wymagania testów i profile danych.
- Źródła i zasady pozyskania.
- Maskowanie/anonymizacja i generacja syntetyczna.
- Wersjonowanie i repo danych testowych.
- Walidacja jakości danych (rozklady, referencje, business rules).
- Dystrybucja i dostęp (ACL, audyt).
- Cleanup/retencja i bezpieczeństwo.
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- Maskowanie → techniki i narzędzia.
- Walidacja → testy jakości i raport.

## Wymagane streszczenia
- Tabela zestaw danych → cel → źródło → maskowanie → wersja → właściciel.

## Guidance
Cel: bezpieczne, realistyczne dane testowe. DoR: wymagania testów, polityki PII, źródła, narzędzia. DoD: plan/maskowanie/walidacja/wersje/cleanup; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Zbierz wymagania; wybierz źródła i techniki maskowania; wygeneruj/zweryfikuj dane; udostępnij z ACL; po użyciu wyczyść; wersjonuj.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Wymagania i profile danych; [ ] Polityki PII/PCI; [ ] Źródła; [ ] Narzędzia.
- DoD: [ ] Plan/maskowanie/walidacja/wersje/cleanup; [ ] Sekcje N/A uzasadnione; metadane aktualne.
