---
title: Deployment Plan
status: needs_content
---

# Deployment Plan

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Zaplanować bezpieczne wdrożenie zmiany/wersji: kroki, zależności, walidacje, ryzyka i ścieżki awaryjne. Ma umożliwić wykonanie deployu powtarzalnie i audytowalnie.

## Zakres i granice
- Obejmuje: zakres releasu, okno i środowiska, kroki techniczne (pre-checks, deploy, migracje), testy po wdrożeniu, plan rollback/abort, komunikację, role on-call.
- Nie obejmuje: strategicznych decyzji releasowych (Release Plan) ani pełnych runbooków operacyjnych (te są referencją).

## Wejścia i wyjścia
- Wejścia: Release Plan, Change Request, lista artefaktów i hashy, zależności systemowe, dane/migracje, okna CAB/blackout, plan monitoringu.
- Wyjścia: wykonany deploy z logiem kroków, wyniki testów smoke/health/func, decyzja Go/No-Go po wdrożeniu, ewentualny rollback, raport z wykonania.

## Struktura sekcji (szkielet)
1. Zakres i wersja release
2. Środowiska, okno wdrożeniowe, zależności
3. Kroki wdrożenia (pre-checks, deploy, migracje, config/feature flags)
4. Testy po wdrożeniu (smoke/health/functional)
5. Plan rollback/backout i kryteria abort
6. Ryzyka i mitigacje
7. Monitoring i SLO/SLA na czas wdrożenia
8. Role, on-call i komunikacja (kanały, odbiorcy, status cadence)
9. Checklist przed/po i log wykonania

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

munikacji.

## Szybkie powiązania
- Dodaj ręcznie 2–3 kluczowe powiązania doc↔doc lub sekcja↔sekcja, korzystając z linkage_index.jsonl / content_links*.json (decyzje, ryzyka, zależności).

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Kroki i testy opisane z ownerami.
- [ ] Rollback/abort kryteria zdefiniowane.
- [ ] Ryzyka i mitigacje spisane.
- [ ] Komunikacja i kanały ustalone.