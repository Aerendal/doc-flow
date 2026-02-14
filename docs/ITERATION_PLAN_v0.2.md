# Iteration Plan v0.2 (post-RC)

## Zakres (2 tygodnie)
- P0: Podpięcie realnego indeksu/status/usage do recommend/template-sets/templates (usunąć dane demo).
- P0: Flagi CLI `--strict-mode`, `--publish-strict`; migracja sekcji `migrate-sections --dry-run`.
- P1: Zapis usage w generatorze + konsumowanie w rekomendacjach.
- P1: Raport brak przykładów (code/table) w CLI.
- P1: Governance CLI raport, rozbudowa pól/sekcji.

## Kryteria sukcesu
- Recommend/template-sets/templates korzystają z rzeczywistych danych z indeksu; brak placeholderów.
- Usage zwiększane przy generate i widoczne w recommend.
- Migracja sekcji dostępna w CLI (dry-run), strict/publish flagi działają.
- Raport brak przykładów pokazuje doc_id bez code/table.

## Ryzyka
- Brak pełnego `template_source` w dokumentach → impact/recommend może być niepełny.
- Wydajność przy równoległym scan/parse (opcjonalnie) — ryzyko nondeterminismu.

## Capacity / dni
- Dev: 10 dni roboczych (1 os.).

## Kolejność
1) Usunięcie demo danych, podpięcie indeksu do recommend/templates.
2) Flagi CLI strict/publish + migrate-sections (dry-run).
3) Usage store integracja z generate i recommend.
4) Raport brak przykładów + governance CLI.
