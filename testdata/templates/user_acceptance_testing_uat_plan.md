---
title: User Acceptance Testing (UAT) Plan
status: needs_content
---

# User Acceptance Testing (UAT) Plan

## Metadane
- Właściciel: [Product/QA/Business]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Opisuje plan testów akceptacyjnych użytkownika: zakres, role, scenariusze biznesowe, kryteria akceptacji, dane/środowiska, harmonogram, raportowanie i decyzję go/no-go. Ma potwierdzić, że rozwiązanie spełnia wymagania biznesowe i jest gotowe do wdrożenia.

## Zakres i granice
- Obejmuje: scenariusze end-to-end biznesowe, krytyczne ścieżki, kryteria akceptacji, dane i konta UAT, środowisko UAT, role i odpowiedzialności (biznes/QA/dev), harmonogram i komunikację, zarządzanie defektami, kryteria go/conditional/no-go, podpis akceptacji.
- Poza zakresem: testy techniczne/perf/security (linki), pełny plan QA (oddzielny).

## Wejścia i wyjścia
- Wejścia: wymagania biznesowe/user stories, krytyczne ścieżki, definicje done, dane UAT (maskowane/syntetyczne), środowisko UAT, dostęp użytkowników, plan release, polityka defektów/priorów.
- Wyjścia: lista scenariuszy/test cases, plan runów, raporty UAT, lista defektów i status, decyzja akceptacji, podpisy interesariuszy.

## Powiązania (meta)
- Key Documents: qa_strategy, testing_plan_schedule, test_data_preparation, release_plan, change_management, incident_response (na wypadek krytycznych defektów).
- Key Document Structures: scenariusze, dane, role, harmonogram, kryteria, raporty.
- Document Dependencies: środowisko UAT, dane, konta, tracker defektów, komunikacja.

## Zależności dokumentu
Wymaga: zatwierdzonych wymagań/stories, dostępnego środowiska i danych UAT, ról/cont UAT, polityki defektów/priorów, planu release. Bez tego DoR otwarte.

## Powiązania sekcja↔sekcja
- Scenariusze → Dane/środowisko → Runy → Defekty → Decyzja.
- Role/komunikacja → Harmonogram → Raporty → Akceptacja.

## Fazy cyklu życia
- Przygotowanie: scenariusze, dane/konta, role, harmonogram, kryteria.
- Wykonanie: runy, defekty, retesty, aktualizacje statusu.
- Ocena: kryteria go/conditional/no-go, podpisy, plan poprawek.
- Zamknięcie: raport, decyzja, lekcje, aktualizacja release.

## Struktura sekcji
1) Zakres i cele UAT (co w zakresie/poza, krytyczne ścieżki)  
2) Role i odpowiedzialności (biznes, QA, dev, support, decision maker)  
3) Scenariusze i przypadki testowe (mapa do wymagań)  
4) Dane/konta i środowisko UAT (maskowanie/privacy, refresh)  
5) Harmonogram i komunikacja (runy, statusy, kanały)  
6) Kryteria akceptacji/go-no-go i polityka defektów (sev/prio, exit/entry)  
7) Raportowanie i podpisy (raporty dzienne, końcowy, akceptacja)  
8) Ryzyka i plan mitigacji\n\n## Wymagane rozwinięcia\n- Lista scenariuszy mapowanych do wymagań; priorytety; dane/konta per scenariusz.\n- Kryteria go/conditional/no-go (defekty sev/prio, coverage, blokery) i podpisy.\n- Harmonogram runów, status meeting cadence, kanały komunikacji.\n\n## Wymagane streszczenia\n- Zakres i cele, główne scenariusze, harmonogram, kryteria akceptacji, decyzja/go-no-go podpisy.\n\n## Guidance (skrót)\n- Wybierz scenariusze reprezentujące realne ścieżki biznesowe i ryzyka.\n- Zapewnij dane/konta przed startem; blokery środowiska = no-go do UAT start.\n- Jasno zdefiniuj kryteria akceptacji i politykę defektów; decyzje dokumentuj.\n- Raportuj postęp i blokery codziennie; retesty planuj w harmonogramie.\n- Po UAT: raport końcowy, podpisy i lekcje; aktualizuj release/plan QA.\n\n## Szybkie powiązania\n- linkage_index.jsonl (uat/plan)\n- qa_strategy, testing_plan_schedule, test_data_preparation, release_plan, change_management, incident_response\n\n## Jak używać dokumentu\n1. Uzupełnij zakres/scenariusze/dane/role; przygotuj środowisko.\n2. Ustal harmonogram i kryteria go/conditional/no-go; poinformuj interesariuszy.\n3. Podczas UAT zbieraj defekty, raportuj status, planuj retesty; aktualizuj dokument.\n4. Na koniec wykonaj decyzję i podpisy; dodaj lekcje i link do release.\n\n## Checklisty Definition of Ready (DoR)\n- [ ] Wymagania/stories gotowe; scenariusze zmapowane; dane/konta dostępne.\n- [ ] Środowisko UAT gotowe; role/przypisania i polityka defektów ustalone.\n- [ ] Harmonogram wstępny i kryteria go/no-go zdefiniowane.\n- [ ] Struktura sekcji wypełniona/N/A.\n\n## Checklisty Definition of Done (DoD)\n- [ ] UAT wykonane; raporty i defekty zarejestrowane; retesty zaplanowane/wykonane.\n- [ ] Kryteria go/conditional/no-go ocenione; decyzja/podpisy zapisane.\n- [ ] Lekcje/ryzyka zaktualizowane; dokument w linkage_index.\n- [ ] Wersja/data/właściciel zaktualizowane.\n\n## Definicje robocze\n- UAT, Go/Conditional/No-go, Sev/Priority, Entry/Exit criteria.\n\n## Przykłady użycia\n- UAT dla release SaaS: 20 scenariuszy E2E, dane maskowane, decyzja go z jednym warunkowym defektem P3.\n- UAT dla aplikacji mobilnej: testy na P0 urządzeniach, kryteria perf/A11y, podpis product/UX.\n\n## Ryzyka i ograniczenia\n- Brak gotowości środowiska/danych → opóźnienie; niejasne kryteria → spory; brak podpisów → ryzyko wdrożenia.\n\n## Decyzje i uzasadnienia\n- [Decyzja] Zakres scenariuszy i kryteria go/no-go — uzasadnienie ryzyk biznesowych.\n- [Decyzja] Kto podpisuje akceptację — uzasadnienie odpowiedzialności.\n\n## Założenia\n- Środowisko i dane UAT dostępne; interesariusze dostępni w oknie UAT; tracker defektów działa.\n\n## Otwarte pytania\n- Czy UAT obejmuje klienta zewnętrznego? \n- Jakie są minimalne KPI (np. brak P1/P2, max X P3)?\n\n## Powiązania z innymi dokumentami\n- QA Strategy, Testing Plan & Schedule, Test Data Preparation, Release Plan, Change Management, Incident Response.\n\n## Powiązania z sekcjami innych dokumentów\n- Testing Plan → harmonogram/runy; Test Data → dane; Release Plan → decyzje; Incident Response → eskalacje krytyczne.\n\n## Słownik pojęć w dokumencie\n- UAT, Go/Conditional/No-go, Sev/Priority, Entry/Exit criteria.\n\n## Wymagane odwołania do standardów\n- Polityki QA/Release, wymagania regulacyjne klienta jeśli dotyczy.\n\n## Mapa relacji sekcja→sekcja\n- Scenariusze → Dane/Środowisko → Runy → Defekty → Decyzja → Podpisy.\n\n## Mapa relacji dokument→dokument\n- UAT Plan → QA/Testing Plan → Release/Change → Incident Response.\n\n## Ścieżki informacji\n- Wymagania → Scenariusze → UAT run → Defekty → Raport → Decyzja → Podpisy.\n\n## Weryfikacja spójności\n- [ ] Scenariusze pokrywają krytyczne ścieżki; dane/środowisko gotowe.\n- [ ] Kryteria go/no-go jasne; decyzje/podpisy udokumentowane.\n- [ ] Relacje cross‑doc opisane; dokument w linkage_index.\n\n## Lista kontrolna spójności relacji\n- [ ] Każdy scenariusz ma dane, środowisko, właściciela, wynik.\n- [ ] Każda decyzja ma kryteria i podpis; każdy defekt ma priorytet i status.\n- [ ] Relacje cross‑doc opisane z uzasadnieniem.\n\n## Artefakty powiązane\n- Lista scenariuszy/cases, dane/konta, raporty runów, log defektów, komunikaty, decyzje/podpisy.\n\n## Ścieżka decyzji\n- [Decyzja] → [Uzasadnienie] → [Konsekwencje] → [Właściciel] → [Data].\n\n## Użytkownicy i interesariusze\n- Business/Product owners, QA, Dev, UX, Support, Security/Perf (jeśli dotyczy).\n\n## Ścieżka akceptacji\n- Business/Product → QA → (Security/Perf jeśli wymagane) → Release/CAB → Owner sign‑off.\n\n## Kryteria ukończenia\n- [ ] UAT zakończone, decyzja/podpisy udokumentowane, defekty P1/P2 zamknięte lub warunkowo zaakceptowane.\n- [ ] Dokument w linkage_index/checklistach; wersja/data/właściciel aktualne.\n\n## Metryki jakości\n- Pass rate, defekt rate (sev/prio), czas reakcji na defekty, liczba powtórnych runów, dotrzymanie harmonogramu UAT, satysfakcja użytkowników UAT.\n*** End Patch
