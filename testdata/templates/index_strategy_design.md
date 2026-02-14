---
title: Index Strategy Design
status: needs_content
---

# Index Strategy Design

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zaprojektować strategię indeksów dla baz danych (SQL/NoSQL), aby zbalansować wydajność zapytań, koszt pamięci i operacje utrzymaniowe.

## Zakres i granice
- Obejmuje: typy indeksów (btree/hash/gin/gist/columnstore), klucze i kolejność, pokrywanie (covering), indeksy częściowe, filtrowane, partition/shard key, maintenance (reindex, autovacuum), wpływ na write/read, strategie dla OLTP/OLAP, monitorowanie i dropping nieużywanych indeksów.
- Poza zakresem: pełne modelowanie danych, optymalizacja sprzętowa.

## Wejścia i wyjścia
- Wejścia: workload i plany zapytań, statystyki użycia indeksów, schema, SLO/latency, limity kosztów/storage, okna utrzymaniowe.
- Wyjścia: plan indeksów (dodane/zmienione/usunięte), rekomendacje per tabela/partycja, harmonogram maintenance, metryki monitoringu, ryzyka write amplification.

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: workload/EXPLAIN, statystyki indeksów, maintenance windows, polityki storage, compliance danych; brak – odnotuj.

## Powiązania sekcja↔sekcja
Workload → wybór indeksów; maintenance → koszty; partition → klucze; monitoring → drop.

## Fazy cyklu życia
Analiza workload → Projekt indeksów → Testy → Rollout → Monitoring/maintenance → Przeglądy okresowe.

## Struktura sekcji (szkielet)
- Zakres tabel/partycji i workload.
- Rekomendowane indeksy (typ, kolumny, warunki, covering, komentarz).
- Wpływ na write/read, storage, autovacuum/reindex.
- Maintenance i harmonogram (okna, limity).
- Monitoring użycia (bloat, scans, hit/miss) i dropping.
- Ryzyka i mitigacje.

## Wymagane rozwinięcia
- Rekomendacje per silnik (Postgres/MySQL/SQL Server/Oracle/NoSQL).
- Skrypty monitoringu i oceny użycia.

## Wymagane streszczenia
- Tabela indeksów: tabela → indeks → cel → koszt → plan rollout.

## Guidance
Cel: optymalne indeksy vs koszt. DoR: workload/EXPLAIN/statystyki, SLO, storage/maintenance limity. DoD: plan indeksów, wpływ/maintenance/monitoring opisane; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Zbierz workload i statystyki; zaprojektuj indeksy; przetestuj; wdrażaj w oknach; monitoruj i sprzątaj.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Workload/EXPLAIN/statystyki; [ ] SLO/koszt; [ ] Okna maintenance.
- DoD: [ ] Plan indeksów + wpływ; [ ] Monitoring/dropping/maintenance; [ ] Sekcje N/A uzasadnione; metadane aktualne.

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
