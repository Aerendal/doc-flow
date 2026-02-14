---
title: Database Incident Response
status: needs_content
---

# Database Incident Response

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Plan reagowania na incydenty bazodanowe (awaria, korupcja, wydajność, bezpieczeństwo): wykrycie, eskalacja, działania, komunikacja, odzyskanie.

## Struktura sekcji (szkielet)
1. Scenariusze: awaria instancji, korupcja danych, spadek wydajności, security breach.
2. Wykrycie i alerty (metryki, logi, thresholdy).
3. Eskalacja i role (DBA, SRE, security, właściciel systemu).
4. Akcje natychmiastowe: izolacja, przełączenie na standby, throttle ruch.
5. Odtwarzanie/naprawa: restore, point‑in‑time, walidacja danych.
6. Komunikacja: statusy, stakeholders, timeline.
7. Post‑incident: RCA, poprawki, testy, dokumentacja.

## Guidance
Cel: skrócone wskazówki do wypełniania szablonów dokumentów (core/satellite).

- Cel dokumentu: 2–3 zdania o decyzjach, ryzykach i wartości dokumentu.
- Zakres i granice: co obejmuje (systemy/procesy/zespoły) i czego nie obejmuje; zaznacz granice odpowiedzialności.
- Wejścia: dane, wymagania, standardy, zależności potrzebne przed startem.
- Wyjścia: artefakty/rezultaty, kto je konsumuje, format (link/plik).
- Zależności dokumentu: wymagane dokumenty lub decyzje; właściciel; wpływ na kolejność prac.
- Powiązania sekcja↔sekcja: które sekcje się rozwijają/streszczają; podaj uzasadnienie.
- Struktura sekcji: utrzymuj układ logiczny; sekcje bez treści oznacz jako N/A z krótkim uzasadnieniem.
- Fazy cyklu życia: zaznacz, w których fazach dokument powstaje/aktualizuje się/archiwizuje; kto odpowiada.
- DoR (Definition of Ready): zakres, wejścia, role, zależności, kryteria akceptacji gotowe.
- DoD (Definition of Done): sekcje uzupełnione lub N/A, powiązania wpisane, checklisty jakości sprawdzone, wersja/data/właściciel, linki/artefakty działają.
- Język: polski; nazwy własne pozostają bez zmian; liczby w nazwach plików usunięte już w szablonach.
- Filozofia: optymalizuj przez rozwój, nie ucinanie — dodawaj, nie kasuj; elementy „satelitarne” zostają.

onitoring.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Alerty i scenariusze zdefiniowane.
- [ ] Role/eskalacja i kontakty DB/SRE.
- [ ] Procedury failover/restore przetestowane.
- [ ] Komunikacja i RCA wykonywane.