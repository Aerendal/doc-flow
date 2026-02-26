# Podsumowanie tygodnia feedbacku - <RRRR-MM-DD> do <RRRR-MM-DD>

Repo: Aerendal/doc-flow
Okno: 7 dni

## 1. Podsumowanie ogólne

- Wszystkich zgłoszeń: <N>
- Issues (bugi/zadania): <N>
- Discussions (pytania/pomysły): <N>
- Zgłoszenia bezpieczeństwa (prywatne): <N lub "brak">
- Potwierdzone bugi: <N>
- Zaakceptowane feature requesty: <N>
- Zamknięte jako duplikat/poza zakresem: <N>

Najważniejsze wnioski (max 5 punktów):
- <wniosek 1>
- <wniosek 2>
- <wniosek 3>

## 2. Zakres i zasady obowiązujące w tygodniu

W zakresie:
- Reprodukowalne defekty (crash, niepoprawny output, zepsuty workflow)
- Problemy z instalacją i pakowaniem
- Regresje determinizmu (niestabilny output, kolejność, timestampy)

Poza zakresem (przykłady):
- Duże redesigny bez acceptance criteria
- Zmiany łamiące wsteczną kompatybilność bez ścieżki opt-in
- Problemy specyficzne dla środowiska bez minimalnego repro

Priorytety:
- P0: bezpieczeństwo, utrata danych, crash, blokada release/CI
- P1: blokuje typowy workflow lub instalację
- P2: ważne, ale jest obejście
- P3: nice-to-have

## 3. Metryki triage

- Procent zgłoszeń z kompletnymi danymi środowiska: <N%>
- Procent zgłoszeń z minimalnym repro: <N%>
- Mediana czasu do pierwszej odpowiedzi: <X godz.>
- Mediana czasu do decyzji triage: <X godz.>

Uwagi:
- <co spowalniało triage>
- <co poprawiło triage>

## 4. Potwierdzone bugi (accepted)

Kolumny: ID | Tytuł | Priorytet | Obszar | Status | Uwagi

- #<ID> | <tytuł> | P? | area/<x> | accepted/in-progress | <krótka uwaga>

## 5. Zaakceptowane feature requesty

### 5.1 #<ID> <Tytuł>

Motywacja:
- <dlaczego>

Proponowane podejście:
- <co zmienia się w CLI/docs/output>

Acceptance criteria:
- [ ] <warunek 1>
- [ ] <warunek 2>

Ryzyka i kompatybilność:
- <ryzyko złamania kompatybilności, ryzyko determinizmu>

Plan opt-in:
- <flaga/konfiguracja jeśli potrzebna>

## 6. Duplikaty i poza zakresem (zamknięte)

Duplikaty:
- #<ID> do #<ID głównego> (powód)

Poza zakresem:
- #<ID> (powód: <krótkie uzasadnienie>)

## 7. Podjęte decyzje

- D1: <co> (uzasadnienie: <dlaczego>)
- D2: <co> (uzasadnienie: <dlaczego>)

## 8. Następne kroki (7-14 dni)

1. <zadanie> (właściciel: <ty/społeczność>, cel: <data>)
2. <zadanie>

## 9. Gdzie można pomóc (help wanted)

- #<ID> (label: meta/help-wanted) - <co zrobić> - <jak testować>

## 10. Linki

- Zapytanie Issues użyte w tygodniu: <link>
- Zapytanie Discussions: <link>
- Baseline wersji/commitu: <commit/tag>
