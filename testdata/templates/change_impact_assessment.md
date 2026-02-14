---
title: Change Impact Assessment
status: needs_content
---

# Change Impact Assessment

## Metadane
- Właściciel: [Change Manager / Product Owner / Architect]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Ocenić wpływ planowanej zmiany (produkt, architektura, proces, infrastruktura) na użytkowników, systemy, koszty, ryzyka i zgodność; określić zakres testów, komunikację, plan wdrożenia i kryteria akceptacji, aby zmiana była świadoma i kontrolowana.

## Zakres i granice
- Obejmuje: identyfikację interesariuszy, zależności systemowych, wpływ na SLA/UX/bezpieczeństwo, koszty i zasoby, plan testów i rollback, komunikację, harmonogram, kryteria wejścia/wyjścia, zgodność/regulacje.  
- Poza zakresem: szczegółowa realizacja techniczna (oddzielne plany), budżetowanie portfelowe.

## Wejścia i wyjścia
- Wejścia: opis zmiany (RFC), backlog/ADR, architektura aktualna, lista systemów zależnych, wymagania biznesowe i zgodności, szacowanie ryzyka, plan release.  
- Wyjścia: arkusz oceny wpływu, mapa zależności, matryca ryzyk i kontroli, zakres testów/regresji, plan wdrożenia i rollback, plan komunikacji, checklisty DoR/DoD.

## Powiązania (meta)
- Key Documents: change_management, risk_assessment, release_readiness_statement, rollback_runbook, service_dependency_map, security_assessment, compliance_architecture_review.  
- Key Document Structures: opis zmiany, wpływ, ryzyka/kontrola, testy, wdrożenie/rollback, komunikacja.  
- Document Dependencies: CMDB, monitoring, CI/CD, CAB proces, incident/problem records.

## Zależności dokumentu
Wymaga: aktualnej architektury i CMDB, listy interesariuszy, kryteriów SLA/UX, polityk bezpieczeństwa i zgodności, zasobów na testy, planu rollback. Braki = brak DoR.

## Powiązania sekcja↔sekcja
- Wpływ ↔ Ryzyka/kontrole ↔ Testy.  
- Zależności systemowe ↔ Plan wdrożenia ↔ Rollback.  
- Komunikacja ↔ Harmonogram ↔ Akceptacja interesariuszy.

## Fazy cyklu życia
- Scoping zmiany i identyfikacja wpływu.  
- Analiza ryzyk i kontroli.  
- Plan testów i wdrożenia (z rollbackiem).  
- Decyzja CAB/PO.  
- Wdrożenie, monitorowanie i walidacja.  
- Retrospektywa i aktualizacja procedur.

## Struktura sekcji
1) Opis zmiany i cele  
2) Wpływ na użytkowników, SLA, bezpieczeństwo, zgodność  
3) Zależności i systemy dotknięte  
4) Ryzyka i środki kontrolne  
5) Zakres testów/regresji i dane testowe  
6) Plan wdrożenia i rollback (okno, kroki, kryteria stop)  
7) Komunikacja i akceptacje  
8) Kryteria akceptacji/DoR/DoD  
9) Otwarte pytania

## Wymagane rozwinięcia
- Macierz wpływu: system × wpływ (wys/śr/niski) + właściciel.  
- Plan testów: jakie testy (unit/integration/e2e/perf/sec) i kto wykonuje.  
- Wymagania zgodności i bezpieczeństwa oraz dowody.  
- Plan komunikacji: kogo, kiedy, jak informować; szablony.  
- Plan rollback: kroki, dane do backupu, punkty decyzji stop.

## Wymagane streszczenia
- Executive summary: cel zmiany, wpływ, ryzyko, decyzja.  
- Skrót harmonogramu i okna serwisowego.

## Guidance (skrót)
- Zawsze identyfikuj zależności w CMDB i mapie usług.  
- Włącz bezpieczeństwo/zgodność w ocenie wpływu; nie odkładaj.  
- Ustal jasne kryteria stop/rollback i odpowiedzialnych.  
- Testy dopasuj do wpływu; krytyczne ścieżki muszą mieć regresję.  
- Komunikuj wcześniej; potwierdzaj odbiór przez kluczowych użytkowników.  
- Dokumentuj decyzje CAB/PO wraz z uzasadnieniem.

## Szybkie powiązania
- linkage_index.jsonl (change/impact/assessment)  
- service_dependency_map, rollback_runbook, release_readiness_statement

## Jak używać dokumentu
1. Wypełnij opis zmiany i mapę wpływu.  
2. Oceń ryzyka/kontrole, ustal testy i plan wdrożenia/rollback.  
3. Uzyskaj akceptacje (CAB/PO); poinformuj interesariuszy.  
4. Po wdrożeniu zweryfikuj metryki, zamknij DoD, zaktualizuj linkage_index.

## Checklisty Definition of Ready (DoR)
- [ ] Opis zmiany i cel biznesowy spisane.  
- [ ] Zidentyfikowane systemy zależne i właściciele.  
- [ ] Ocena ryzyka i wymagania zgodności gotowe.  
- [ ] Plan testów/regresji i dane dostępne.  
- [ ] Plan wdrożenia/rollback i komunikacji uzgodniony.

## Checklisty Definition of Done (DoD)
- [ ] Wdrożenie wykonane; testy i monitoring zielone.  
- [ ] Brak otwartych krytycznych incydentów po zmianie.  
- [ ] Dokumentacja wpływu, decyzji i dowodów uzupełniona.  
- [ ] Komunikaty wysłane; interesariusze potwierdzili.  
- [ ] linkage_index/CMDB zaktualizowane; retrospektywa zapisana.

## Definicje robocze
- CAB: Change Advisory Board.  
- Kryteria stop: warunki, przy których wdrożenie jest zatrzymywane.  
- SLA/OLA: parametry usług i wewnętrzne umowy operacyjne.

## Przykłady użycia
- Zmiana schematu API wpływająca na aplikacje mobilne.  
- Migracja bazy danych do nowej wersji.  
- Włączenie nowej polityki bezpieczeństwa haseł.  
- Rebalans usług między regionami chmurowymi.

## Ryzyka i ograniczenia
- Niedoszacowany wpływ → awarie zależnych systemów.  
- Brak testów krytycznych ścieżek → regresje produkcyjne.  
- Słaba komunikacja → incydenty operacyjne i niezadowolenie klientów.  
- Niekompletny rollback → wydłużone przestoje.

## Decyzje i uzasadnienia
- Priorytet i okno serwisowe.  
- Zakres testów/regresji vs czas.  
- Kryteria stop/rollback i właściciele.  
- Akceptacja wyjątków zgodności/bezpieczeństwa.

## Założenia
- CMDB i monitoring są aktualne.  
- Dostępne są środowiska testowe zbliżone do prod.  
- Interesariusze są dostępni do akceptacji i komunikacji.

## Otwarte pytania
- Jakie mierniki sukcesu zmiany (KPI) i horyzont obserwacji?  
- Czy potrzebne jest okno zwrotne (grace period) dla użytkowników?  
- Jak obsłużyć klientów/regulacje w innych regionach?  
- Jak wersjonować dokumentację zmian (linkage_index)?
