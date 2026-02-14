# Iteration Plan v0.3 (draft) — po RC2

## Cele nadrzędne
- Rozwijamy funkcjonalność (nie tniemy zakresu): cykle, UX, observability, perf 10k, zgodne przykłady.
- Utrzymujemy stabilność (go test ./..., governance-ready examples).

## Priorytety (P1)
1) Wykrywanie cykli — pełniejszy raport (lista cykli, weryfikacja na edge-cases).
2) UX: alias `--dry-run` już dodany, doprecyzować komunikaty brakujących flag (rules/config).
3) Perf 10k — wykonać realny pomiar, zapisać RSS/czasy w LOGS/SCALE_BASELINE_10k.md.
4) Observability — potwierdzić profile CPU/Mem działają, dodać przykład w README (nowy podrozdział).
5) Przykłady governance — utrzymać PASS po ewentualnych zmianach schematów.

## Zadania (propozycja)
- [P1] Uruchomić Perf10k (manual) + zaktualizować log.
- [P1] README: dodać sekcję Observability (cpu/mem profile) + zamienić placeholder URL binarki.
- [P1] Cycle detection: raportować ścieżkę cyklu w validate output (format).
- [P2] CLI UX: lepsze komunikaty gdy brak --rules/--config; dodać krótkie przykłady w help.
- [P2] Packaging: podpisy (opcjonalnie) lub publikacja checksum w release notes.
- [P2] Add tests for multi-cycle graph (disjoint cycles).

## Kryteria ukończenia v0.3
- Perf10k liczby dostępne (czas/RSS) i mieszczą się w budżecie (<500 MB).
- Observability: przykładowe profile CPU/Mem można wygenerować i otworzyć (go tool pprof).
- Validate raportuje cykle z czytelną listą doc_id.
- README/Quickstart bez placeholderów; checksumy opublikowane.
- Wszystkie testy go + compliance examples PASS.
