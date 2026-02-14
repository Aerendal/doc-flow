# Docflow Gate CI Template

Plik: `templates/ci/docflow-gate.yml`

## Co robi
1) Buduje `docflow` w trybie `-mod=vendor` (Go 1.25).
2) `validate --strict --governance` na `docflow.yaml` w repo.
3) `compliance` → `/tmp/compliance_report.json`.
4) `section_order_lint` na `docs` i `examples` → `.docflow/section_order_lint.json`.
5) `queue_evaluate` na `examples/backlog_ci_template.txt` → `/tmp/queue_report.json`.
6) Uploaduje artefakty jako `docflow-gate-reports`.

## Wymagania w repo
- `docflow.yaml` w root + `docs/_meta/GOVERNANCE_RULES.yaml`.
- Backlog: `examples/backlog_ci_template.txt` (możesz podmienić własny plik).
- Dostęp do Go modules (sieć albo przygotowany cache) i `GOFLAGS=-mod=vendor`.

## Jak użyć
- Skopiuj `templates/ci/docflow-gate.yml` do `.github/workflows/docflow-gate.yml` w swoim repo.
- Dostosuj backlog (ścieżka, taski) i ewentualnie katalogi lintu (`docs`, `examples`).
- Opcjonalnie dodaj matrix Go jeśli potrzebujesz.

## Artefakty
- `compliance_report.json`
- `queue_report.json` (z `cache_status` per task)
- `.docflow/section_order_lint.json`

## Fail conditions
- validate/compliance exit != 0
- lint wykryje missing/out_of_order
- queue zwróci BLOCKED
