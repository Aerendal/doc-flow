# ALGO_PARAMS — Parametry algorytmów analizy wzorców

Data: 2026-02-09
Etap: day_05
Aktualizacje: day_07 — quality scoring, duplikaty

## 1. Fingerprinting sekcji

| Parametr | Wartość | Uzasadnienie |
|---|---|---|
| Składowe fingerprint | level + text (lowercase, trimmed) | Rozróżnia zarówno strukturę jak i treść sekcji |
| Wyłączenie h1 z fingerprint | TAK | h1 = tytuł dokumentu (unikalny), nie jest częścią wzorca struktury |
| Hash | SHA-256 (24 hex = 12 bajtów) | Wystarczająca unikalność, kompaktowy identyfikator |
| Separator sekcji | `|` (pipe) | Prosty delimiter, nie występuje w nagłówkach |

## 2. N-gramy

| Parametr | Wartość | Uzasadnienie |
|---|---|---|
| n (długość n-gramu) | 2, 3 | Bigramy dla lokalnych wzorców, trigramy dla sekwencji |
| Minimalna częstotliwość | 10% szablonów | Odfiltrowanie rzadkich n-gramów |
| Normalizacja tekstu | lowercase, trim | Spójne porównanie niezależne od wielkości liter |

## 3. Levenshtein / Podobieństwo

| Parametr | Wartość | Uzasadnienie |
|---|---|---|
| Dystans | Levenshtein na sekwencji SectionEntry | Porównanie edit-distance między sekwencjami nagłówków |
| Similarity | 1 - (distance / max(len(a), len(b))) | Znormalizowana miara 0.0–1.0 |
| Próg podobieństwa (merge) | 0.85 | Wzorce z similarity >= 0.85 mogą być kandydatami do mergowania |
| Koszty operacji | insert=1, delete=1, substitute=1 | Standardowe koszty Levenshtein |

## 4. Grupowanie wzorców

| Parametr | Wartość | Uzasadnienie |
|---|---|---|
| Metoda | Exact fingerprint match | Deterministyczna, szybka, brak ambiguity |
| Top N | 50 (domyślnie 10 w CLI) | Pokrycie ~99% szablonów |
| Max examples per group | 5 | Ograniczenie pamięci raportu |

## 5. Wyniki na zbiorze 829 szablonów

| Metryka | Wartość |
|---|---|
| Szablonów | 829 |
| Unikalnych wzorców (h2+) | ~36 |
| Dominujący wzorzec | 718 szablonów (86.6%) — 62 sekcje h2 |
| Top 5 wzorców | pokrycie ~95% |
| Czas analizy | ~0.7s |

## 6. Walidacja sekcji / scoring (day_07)

| Parametr | Wartość |
|---|---|
| Źródło schematu | `section_schema.yaml` |
| Normalizacja nagłówków | lowercase + trim, h1 pomijane |
| Score Structure | 0–60, proporcja braków wymaganych sekcji |
| Score Completeness | 0–30, penalizacja 5 pkt za brak sekcji wymaganej |
| Score Meta | 10 (placeholder) |
| Całkowity | min(100, suma komponentów) |

## 7. Duplikaty szablonów (day_07)

| Parametr | Wartość |
|---|---|
| Reprezentacja | Zbiór tokenów `level:text` dla h2+ |
| Miara | Jaccard similarity |
| Próg domyślny | 0.85 (`docflow find-duplicates --threshold`) |
| Przykład raportu | para ścieżek + sim w [0..1] |

## 8. Rekomendacje szablonów (day_12, MVP)

| Parametr | Wartość |
|---|---|
| Wagi domyślne | doc_type 0.40, lang 0.20, quality 0.25, usage 0.15 |
| Filtry | pomija deprecated |
| Top N | 5 (CLI: --top) |
| Jakość | score z `pkg/quality` (0-100) |
| Użycie | usage count (placeholder, brak persystencji w MVP) |

## 9. Template lifecycle (day_36)

| Parametr | Wartość |
|---|---|
| Stany | draft → active → deprecated → archived |
| Deprecation rule | quality < 50 i usage == 0 przez ≥90 dni |
| Archive rule | manual po migracji (deprecated → archived) |
| Draft → active | manual (po review) |
