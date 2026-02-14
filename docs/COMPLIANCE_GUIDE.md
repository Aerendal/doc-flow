# Compliance Guide (MVP)

## Cel
Szybko sprawdzić zgodność dokumentów z regułami governance (`docs/_meta/GOVERNANCE_RULES.yaml`) i wygenerować raport tekstowy/JSON/HTML.

## Użycie
```bash
./docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml          # raport tekstowy
./docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml --format json --output .docflow/output/compliance.json
./docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml --format html --html .docflow/output/compliance.html
```

## Wyjścia
- **text** (domyślny): podsumowanie + lista niezgodnych dokumentów.
- **json**: pełny Summary (pass/fail, violations_count, docs[]).
- **html**: prosty dashboard (pass rate, violations by type, tabela non-compliant).

## Uwagi / ograniczenia
- Validator governance korzysta z prostych heurystyk (wyszukiwanie `field:` i nagłówków). Dla większej precyzji trzeba podpiąć parser frontmatter/sekcji.
- PDF nieobsługiwany w tej wersji (można wydrukować HTML z przeglądarki).
- Jeśli widzisz `permission denied` dla `~/.cache/go-build`, uruchom z `GOCACHE=/tmp/go-cache`.
