---
title: User Acceptance Test (UAT) Plan
status: needs_content
---

# User Acceptance Test (UAT) Plan

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zdefiniować plan testów akceptacyjnych użytkownika (UAT) obejmujący zakres, scenariusze, kryteria akceptacji, organizację i raportowanie, aby potwierdzić gotowość rozwiązania do produkcji.

## Zakres i granice
- Obejmuje: cele UAT, zakres funkcjonalny/niefunkcjonalny, kryteria wejścia/wyjścia, scenariusze biznesowe, dane i środowisko, role (biznes/QA/IT), plan defektów, harmonogram, komunikację, ryzyka.
- Poza zakresem: testy jednostkowe/integracyjne/perf (opisane w swoich planach), pełne szkolenia użytkowników (osobne materiały).

## Wejścia i wyjścia
- Wejścia: wymagania biznesowe, przypadki użycia, backlog/us stories, kryteria akceptacji, środowisko UAT, dane testowe, dostęp do systemów, plan release.
- Wyjścia: harmonogram UAT, lista scenariuszy/test cases, matryca ról i RACI, raporty dzienne/postęp, lista defektów, decyzja go/conditional/no-go.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: wymagania, plan release, środowisko i dane UAT, systemy zewnętrzne, polityki bezpieczeństwa danych testowych; brak – odnotuj.

## Powiązania sekcja↔sekcja
Scenariusze → dane → kryteria akceptacji; defekty → decyzja; harmonogram → komunikacja.

## Fazy cyklu życia
Planowanie → Przygotowanie → Wykonanie → Triaging/naprawy → Retest/regresja → Decyzja i zamknięcie.

## Struktura sekcji (szkielet)
- Cele UAT i kryteria wejścia/wyjścia.
- Zakres funkcji i zależności.
- Środowisko i dane (masking, kopie, dostęp).
- Scenariusze/test cases i priorytety.
- Role i RACI; komunikacja (stanówki, eskalacje).
- Defekty: proces zgłoszeń, SLA, severity.
- Harmonogram i kamienie milowe.
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- Dane → zasady anonimizacji/maskowania.
- Kryteria akceptacji → powiązanie z user stories.

## Wymagane streszczenia
- Tabela: scenariusz → owner → data → status → decyzja.

## Guidance
Cel: biznes potwierdza, że produkt spełnia potrzeby. DoR: wymagania, środowisko, dane, role gotowe. DoD: scenariusze/defekty/raporty/dec. go-no-go; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Zdefiniuj cele/zakres/kryteria; przygotuj środowisko i dane; uruchom scenariusze; zarządzaj defektami; raportuj; podejmij decyzję.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Wymagania/kryteria akceptacji; [ ] Środowisko/dane; [ ] Role/RACI i komunikacja.
- DoD: [ ] Scenariusze wykonane; [ ] Defekty i decyzje udokumentowane; [ ] Sekcje N/A uzasadnione; metadane aktualne.

## Definicje robocze
- [Termin 1]
- [Termin 2]
- [Termin 3]

## Przykłady użycia
- [Przykład 1]
- [Przykład 2]

## Ryzyka i ograniczenia
- [Ryzyko 1]
- [Ryzyko 2]

## Decyzje i uzasadnienia
- [Decyzja 1]
- [Decyzja 2]

## Założenia
- [Założenie 1]
- [Założenie 2]

## Otwarte pytania
- [Pytanie 1]
- [Pytanie 2]

## Powiązania z innymi dokumentami
- [Dokument A] — [typ relacji] — [uzasadnienie]
- [Dokument B] — [typ relacji] — [uzasadnienie]

## Powiązania z sekcjami innych dokumentów
- [Dokument X → Sekcja Y] — [powód]
- [Dokument Z → Sekcja W] — [powód]

## Wymagane odwołania do standardów
- [Standard 1]
- [Standard 2]

## Mapa relacji sekcja→sekcja
- [Sekcja A] -> [Sekcja B] : [typ]
- [Sekcja C] -> [Sekcja D] : [typ]

## Mapa relacji dokument→dokument
- [Dokument A] -> [Dokument B] : [typ]
- [Dokument C] -> [Dokument D] : [typ]

## Ścieżki informacji
- [Wejście] → [Źródło] → [Rozwinięcie] → [Wyjście]
- [Wejście] → [Źródło] → [Streszczenie] → [Wyjście]

## Weryfikacja spójności
- [ ] Ścieżki informacji zamknięte
- [ ] Brak sprzecznych relacji
- [ ] Sekcje krytyczne mają źródła

## Lista kontrolna spójności relacji
- [ ] Relacje mają sekcje źródłowe
- [ ] Relacje nie są sprzeczne
- [ ] Cross-doc uzasadnione
- [ ] Rozwinięcia/streszczenia odnotowane

## Artefakty powiązane
- [Artefakt 1]
- [Artefakt 2]

## Ścieżka decyzji
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]

## Użytkownicy i interesariusze
- [Rola] — [potrzeby/odpowiedzialności]
- [Rola] — [potrzeby/odpowiedzialności]

## Ścieżka akceptacji
