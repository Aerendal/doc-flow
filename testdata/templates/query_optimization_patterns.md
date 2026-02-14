---
title: Query Optimization Patterns
status: needs_content
---

# Query Optimization Patterns

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zebrać wzorce optymalizacji zapytań (SQL/NoSQL) w różnych silnikach, aby poprawić wydajność i koszty, zachowując poprawność danych.

## Zakres i granice
- Obejmuje: indeksowanie, plan zapytania, partition/sharding, join/aggregation patterns, pushdown, materialized views, caching, statistics, concurrency, anti-patterns, narzędzia profilu/EXPLAIN, różnice między silnikami (Postgres/MySQL/Oracle/SQL Server/BigQuery/Snowflake/NoSQL).
- Poza zakresem: modelowanie danych od zera (osobne), tuning hardware.

## Wejścia i wyjścia
- Wejścia: workload i problemy (wysokie latency, koszt), profile/EXPLAIN, schema, dane o wolumenie i kardynalności, SLO, limity kosztowe.
- Wyjścia: katalog wzorców i przykładów, checklista debug/tuning, zalecenia per silnik, metryki przed/po, ryzyka poprawności.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: schematy, workload/trace, narzędzia EXPLAIN/profil, polityki indeksów/materialized views, limity kosztów; brak – odnotuj.

## Powiązania sekcja↔sekcja
Workload → wzorce; statistics → plany; MV/cache → konsystencja; partition → join.

## Fazy cyklu życia
Zbieranie problemów → Analiza planów → Dobór wzorców → Walidacja → Dokumentacja.

## Struktura sekcji (szkielet)
- Problemy/wymagania i SLO.
- Profilowanie i plany (EXPLAIN/trace).
- Wzorce indeksowania (BTREE/GiST/Gin/covering/partial).
- Wzorce join/aggregation (hash/merge/nested, grouping sets, window, spill).
- Partition/sharding i colocacja.
- Materialized views i cache (TTL, refresh, invalidacja).
- Statistics i cardinality (analyze, histograms, sampling).
- CTE/inline, pushdown i predicate rewrites.
- Concurrency i blokady (isolation, lock avoidance).
- Anti-patterns i checklisty.
- Metryki przed/po, testy regresji i poprawności.

## Wymagane rozwinięcia
- Wzorce per silnik → specyficzne różnice i flagi.
- Testy → skrypty benchmarkowe.

## Wymagane streszczenia
- Tabela wzorzec → użycie → silnik → ryzyko → efekt.

## Guidance
Cel: praktyczne wzorce tuningu. DoR: workload/EXPLAIN/schema; SLO/koszt. DoD: wzorce + przykłady, checklisty, metryki przed/po, sekcje N/A uzasadnione, metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Zbieraj plany; dobierz wzorce; zastosuj i zmierz; dokumentuj; utrzymuj checklistę.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Workload/EXPLAIN/schema; [ ] SLO/koszt; [ ] Polityki indeksów/MV.
- DoD: [ ] Wzorce i przykłady; [ ] Metryki przed/po; [ ] Testy regresji; [ ] Sekcje N/A uzasadnione; metadane aktualne.

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
