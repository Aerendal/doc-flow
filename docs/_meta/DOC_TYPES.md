# DOC_TYPES — Typy dokumentów

## Wersja

`v1.0.0`

## Definicja typów

Każdy dokument w repozytorium musi mieć przypisany `doc_type` z poniższej listy. Typ determinuje oczekiwaną strukturę, wymagane sekcje i reguły walidacji.

## Lista typów

### 1. `specification`

Dokument specyfikacji technicznej lub biznesowej.

- **Przeznaczenie:** Definiuje wymagania, architekturę, design, API.
- **Wymagane sekcje:** Cel dokumentu, Zakres, Wymagania, Kryteria akceptacji.
- **Przykłady:** API Design Specification, Functional Requirements, System Architecture.

### 2. `guide`

Przewodnik lub instrukcja.

- **Przeznaczenie:** Opisuje jak coś zrobić krok po kroku.
- **Wymagane sekcje:** Cel, Wymagania wstępne, Kroki, Weryfikacja.
- **Przykłady:** API Integration Guide, Onboarding Guide, Troubleshooting Guide.

### 3. `procedure`

Procedura operacyjna lub procesowa.

- **Przeznaczenie:** Formalny opis procesu z rolami i odpowiedzialnościami.
- **Wymagane sekcje:** Cel, Zakres, Role, Kroki procedury, Eskalacja.
- **Przykłady:** Change Management Procedure, Incident Response Procedure.

### 4. `report`

Raport lub analiza.

- **Przeznaczenie:** Przedstawia wyniki, metryki, wnioski.
- **Wymagane sekcje:** Podsumowanie, Metodyka, Wyniki, Wnioski, Rekomendacje.
- **Przykłady:** Performance Report, Risk Assessment, Data Quality Report.

### 5. `plan`

Plan projektu, wdrożenia lub testów.

- **Przeznaczenie:** Opisuje co, kiedy, kto i jak.
- **Wymagane sekcje:** Cel, Zakres, Harmonogram, Zasoby, Ryzyka.
- **Przykłady:** Implementation Plan, Test Plan, Migration Plan, UAT Plan.

### 6. `checklist`

Lista kontrolna.

- **Przeznaczenie:** Zbiór pozycji do weryfikacji / zatwierdzenia.
- **Wymagane sekcje:** Lista pozycji z checkboxami, Kryteria zaliczenia.
- **Przykłady:** Release Checklist, Security Checklist, Accessibility Checklist.

### 7. `decision`

Rekord decyzji (ADR).

- **Przeznaczenie:** Dokumentuje decyzję architektoniczną lub techniczną.
- **Wymagane sekcje:** Kontekst, Decyzja, Konsekwencje, Status.
- **Przykłady:** ADR-0001, Technology Selection Decision.

### 8. `runbook`

Podręcznik operacyjny.

- **Przeznaczenie:** Instrukcje reagowania na zdarzenia operacyjne.
- **Wymagane sekcje:** Wyzwalacz, Diagnoza, Kroki naprawcze, Eskalacja.
- **Przykłady:** API Monitoring Runbook, Incident Runbook.

### 9. `template`

Szablon dokumentu.

- **Przeznaczenie:** Wzorzec do tworzenia nowych dokumentów.
- **Wymagane sekcje:** Placeholdery, Instrukcje wypełniania.
- **Przykłady:** Szablon retrospektywy, Szablon RFC.

### 10. `policy`

Polityka lub standard.

- **Przeznaczenie:** Definiuje zasady, standardy, reguły compliance.
- **Wymagane sekcje:** Cel, Zakres, Zasady, Egzekwowanie.
- **Przykłady:** Security Policy, Data Retention Policy, SLA Policy.

### 11. `communication`

Dokument komunikacyjny.

- **Przeznaczenie:** Informuje interesariuszy o zmianach, statusie.
- **Wymagane sekcje:** Odbiorcy, Treść, Kanał dystrybucji.
- **Przykłady:** Release Notes, Announcement, Deprecation Notice.

### 12. `unknown`

Typ nierozpoznany.

- **Przeznaczenie:** Dokument bez zdefiniowanego lub rozpoznanego `doc_type`.
- **Akcja:** Loguj ostrzeżenie, zaproponuj klasyfikację.

## Mapowanie automatyczne

Scanner docflow automatycznie proponuje `doc_type` na podstawie:
1. Słów kluczowych w `title` (np. "guide" → `guide`, "checklist" → `checklist`).
2. Struktury dokumentu (np. obecność checkboxów → `checklist`).
3. Wzorców nazw plików.

## Rozszerzalność

Nowe typy można dodawać przez:
1. Dodanie wpisu w tym pliku.
2. Dodanie reguł walidacji w `internal/validator/`.
3. Opcjonalnie: dodanie szablonu w `testdata/templates/`.
