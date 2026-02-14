---
title: Data Synchronization Incident
status: needs_content
---

# Data Synchronization Incident

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Opisać postępowanie i raport dla incydentów synchronizacji danych (rozjazdy, opóźnienia, duplikaty, utrata danych) w integracjach/bazach, aby szybko przywrócić spójność, zminimalizować wpływ i zapobiec regresjom.

## Zakres i granice
- Obejmuje: objawy (lag, inconsistency, missing/duplicate rows), systemy źródło↔cel (DB, cache, queue, ETL, CDC, API), detekcja, triage, korekty (replay, recompute, backfill), blokady/rollback, komunikacja, dowody, lessons learned.
- Poza zakresem: incydenty bezpieczeństwa (raportowane osobno), chroniczne problemy architektury (osobne RCA/architecture docs).

## Wejścia i wyjścia
- Wejścia: alerty/metryki (lag, checksum mismatch, DLQ, error rates), logi ETL/CDC/queue, schematy danych, reguły idempotencji/de-dupe, SLA/OLA, listy zależności, dane referencyjne.
- Wyjścia: raport incydentu, plan korekcji (backfill/replay), status zgodności danych, lista defektów/procesów do poprawy, aktualizacje runbooków i konfiguracji monitoringu.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Jeżeli brak danych: wskaż dependency na schematy danych, listę przepływów ETL/CDC, runbook integracji, SLO/SLA, zasady retry/idempotencji; brak – odnotuj.

## Powiązania sekcja↔sekcja
Np. objawy → wybór ścieżki korekty; zakres systemów → plan backfill; reguły idempotencji → procedury replay; lessons learned → aktualizacja monitoringu.

## Fazy cyklu życia
- Detekcja/Triage: identyfikacja symptomów i zakresu.
- Containment/Recovery: wstrzymanie/ograniczenie ruchu, naprawa danych.
- Weryfikacja: potwierdzenie spójności, testy regresji.
- Postmortem: RCA, działania prewencyjne, aktualizacja runbooków/monitoringu.

## Struktura sekcji (szkielet)
- Streszczenie incydentu (systemy, dane, SLA wpływ).
- Timeline i symptomy (metryki, alerty, logi).
- Zakres i wpływ na dane/użytkowników/procesy.
- Przyczyna i czynniki współistniejące (schema drift, backpressure, idempotencja, kolejność zdarzeń).
- Działania naprawcze: replay, backfill, deduplikacja, rekonsyliacja.
- Blokady/guardrails: throttling, stop-the-line, feature flag.
- Weryfikacja spójności: checksumy, próbki, raport porównawczy.
- Lessons learned i plan zapobiegania.
- Aktualizacje monitoringu/runbooków/konfiguracji.

## Wymagane rozwinięcia
- Weryfikacja spójności → metody porównania (checksum, row counts, sampled diff).
- Działania naprawcze → skrypty/komendy/SQL.
- Monitorowanie → konkretne metryki/alerty.

## Wymagane streszczenia
- Tabela: system/zakres/naprawa/status/verif.

## Guidance
- Cel: szybkie przywrócenie spójności przy minimalnym ryzyku degradacji.
- Wejścia: alerty, logi, schematy, SLA, zasady idempotencji.
- Wyjścia: raport, korekty, weryfikacja, lessons learned.
- DoR: zidentyfikowane systemy i schematy, dostęp do logów/monitoringu, zasady replay/backfill.
- DoD: spójność potwierdzona, korekty opisane, lessons learned i aktualizacje runbooków/monitoringu; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Po alercie wypełnij streszczenie/timeline; wybierz ścieżkę korekty.
- Wykonaj backfill/replay/dedupe zgodnie z zasadami; zweryfikuj spójność.
- Uzupełnij lessons learned; zaktualizuj monitoring i runbooki; zamknij DoR/DoD.

## Checklisty jakości (DoR/DoD skrót)
- DoR:
  - [ ] Systemy i schematy zidentyfikowane; logi/metryki dostępne.
  - [ ] Zasady idempotencji/retry/ordering znane; właściciele danych dostępni.
- DoD:
  - [ ] Dane naprawione i zweryfikowane; tabela statusów uzupełniona.
  - [ ] Lekcje wyciągnięte, monitoring/runbooki zaktualizowane; sekcje N/A uzasadnione.
  - [ ] Metadane aktualne; powiązania z innymi rejestrami (incydent, risk) dodane.

## Definicje robocze
- [Termin 1] — [definicja robocza]
- [Termin 2] — [definicja robocza]
- [Termin 3] — [definicja robocza]

## Przykłady użycia
- [Przykład 1 — krótki opis sytuacji i zastosowania]
- [Przykład 2 — krótki opis sytuacji i zastosowania]

## Ryzyka i ograniczenia
- [Ryzyko 1 — wpływ i sposób ograniczenia]
- [Ryzyko 2 — wpływ i sposób ograniczenia]

## Decyzje i uzasadnienia
- [Decyzja 1 — uzasadnienie]
- [Decyzja 2 — uzasadnienie]

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
- [Dokument X → Sekcja Y] — [powód powiązania]
- [Dokument Z → Sekcja W] — [powód powiązania]

## Słownik pojęć w dokumencie
- [Pojęcie 1] — [definicja i źródło]
- [Pojęcie 2] — [definicja i źródło]
- [Pojęcie 3] — [definicja i źródło]

## Wymagane odwołania do standardów
- [Standard 1] — [sekcja/fragment, którego dotyczy]
- [Standard 2] — [sekcja/fragment, którego dotyczy]

## Mapa relacji sekcja→sekcja
- [Sekcja A] -> [Sekcja B] : [typ relacji]
- [Sekcja C] -> [Sekcja D] : [typ relacji]

## Mapa relacji dokument→dokument
- [Dokument A] -> [Dokument B] : [typ relacji]
- [Dokument C] -> [Dokument D] : [typ relacji]

## Ścieżki informacji
- [Wejście] → [Sekcja źródłowa] → [Sekcja rozwinięcia] → [Wyjście]
- [Wejście] → [Sekcja źródłowa] → [Sekcja streszczenia] → [Wyjście]

## Weryfikacja spójności
- [ ] Czy wszystkie ścieżki informacji są zamknięte?
- [ ] Czy istnieją pętle lub sprzeczne relacje?
- [ ] Czy sekcje krytyczne mają wskazane źródła i rozwinięcia?

## Lista kontrolna spójności relacji
- [ ] Czy każda sekcja z relacją ma wskazaną sekcję źródłową?
- [ ] Czy relacje nie tworzą sprzecznych wymagań (np. wzajemne wykluczanie)?
- [ ] Czy relacje cross‑doc mają uzasadnienie i są zgodne z fazą?
- [ ] Czy relacje wymagają rozwinięć lub streszczeń są odnotowane?

## Artefakty powiązane
- [Artefakt 1] — [opis i relacja do dokumentu]
- [Artefakt 2] — [opis i relacja do dokumentu]

## Ścieżka decyzji
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]

## Użytkownicy i interesariusze
- [Rola / interesariusz] — [potrzeby i odpowiedzialności]
- [Rola / interesariusz] — [potrzeby i odpowiedzialności]

## Ścieżka akceptacji
