# Contract-First Engine Architecture

Ten dokument opisuje docelowy kierunek: `contract-first engine + execution layer`.

## Cele
- Jedna kanoniczna reprezentacja wyniku (`Issue`/`Report`) niezależna od:
  - formatu (`text/json/sarif`),
  - trybu wykonania (full/incremental),
  - baseline (`all/new/existing`).
- Renderery nie zawierają logiki reguł.
- Cache/incremental przyspiesza wykonanie, ale nie zmienia semantyki wyniku.

## Niezmienniki kontraktowe
- Exit codes: `0/1/2/3`.
- JSON envelope dla `--format json` przy usage/runtime.
- `schema_version` i kontrakty reportów pozostają stabilne.
- Baseline comparator działa po utworzeniu kanonicznego reportu.

## Pipeline
1. `discover`: wybór dokumentów wejściowych.
2. `extract`: parsowanie i budowa faktów dokumentu.
3. `evaluate`: reguły validate/compliance -> kanoniczny report.
4. `compare`: opcjonalnie baseline (`new/existing/all`).
5. `render`: serializacja `text/json/sarif`.

Warstwa wykonawcza (execution layer) obsługuje IO/cache/równoległość i podaje dane do engine.

## Granice modułów
- `internal/engine`: kanoniczne typy i ocena reguł (bez CLI i bez stdout/stderr).
- `internal/cli`: tylko parsing flag, wywołanie engine, baseline view, render.
- `internal/exec`: sposób wykonania (cache/incremental/concurrency), bez logiki reguł.

## Status wdrożenia (część 1/6)
- Dodane kanoniczne typy: `internal/engine/types.go`.
- Dodane adaptery:
  - `internal/engine/validate.go`
  - `internal/engine/compliance.go`
- CLI validate/compliance używa adapterów engine bez zmiany kontraktu użytkowego.

## Status wdrożenia (część 2/6)
- Dodana warstwa wykonawcza:
  - `internal/exec/cache/sqlite_cache.go` (SQLite cache + schema version + run stats)
  - `internal/exec/runner.go` (discover + incremental + cache hit/miss)
- `validate` i `compliance` korzystają z execution layer do pozyskania facts.
- Dodany build indeksu z facts: `pkg/index/build_from_records.go`.
- Dodana walidacja z facts: `internal/validator/validate_facts.go`.
- Dodany compliance report z facts: `pkg/compliance/reporter.go` (`ReportWithFacts`).
- Flagi CLI dla execution layer: `--cache`, `--no-cache`, `--cache-dir`, `--changed-only`, `--since`.

## Status wdrożenia (część 3/6)
- Dodany kanoniczny moduł grafu:
  - `internal/engine/graph/types.go`
  - `internal/engine/graph/graph.go`
  - `internal/engine/graph/cycles.go`
  - `internal/engine/graph/impact.go`
- Graf jest deterministyczny: sortowanie `nodes[]`, `edges[]`, kanonizacja `cycles[]`.
- Dodana komenda `docflow graph` (`text/json`, `--cycles`, `--include-links`, `--include-templates`).
- `template-impact` korzysta z traversal po kanonicznym grafie (`uses_template`) zamiast osobnej heurystyki.
- Dodane testy:
  - deterministyczność grafu i issues przy różnej kolejności wejścia,
  - stabilność i kanonizacja cykli,
  - impact analysis (forward/reverse).

## Status wdrożenia (część 4/6)
- Comparator baseline v2 oparty o `identity_version=2`:
  - validate identity używa pól strukturalnych (`code`, `path`, `doc_id`, `location`, `details`) zamiast `message`.
  - compliance posiada jawne `identity_version` i metadane baseline z wersją identity.
- Dodane `details` dla typów validate (m.in. `missing_field`, `invalid_yaml`, `missing_expected_dependency`, `governance_violation`, `legacy_section_name`).
- Dodane narzędzie migracji baseline:
  - `docflow baseline migrate --in ... --out ... --kind validate|compliance`
- Dodane testy regresji:
  - wording-only change nie generuje `new` przy baseline v2,
  - multi-issues nie sklejają się przy różnym `details.field`,
  - migracja baseline uzupełnia `identity_version` i `details`.

## Status wdrożenia (część 5/6)
- Dodany `internal/engine/fix`:
  - plan fixów (deterministyczny),
  - unified diff renderer,
  - atomowe apply (`temp + rename`) z kontrolą konfliktu (checksum) i opcjonalnym backupem.
- Dodana komenda `docflow fix`:
  - `--format diff|json`, `--output`, `--apply`, `--backup-dir`,
  - `--scope`, `--only`, `--max-files`, `--max-changes`,
  - opcjonalnie baseline-aware fixing (`--against`, `--show`).
- Zakres bezpiecznych fixów na start:
  - `missing_field(version)` -> uzupełnienie frontmatter,
  - `legacy_section_name` -> migracja nagłówków wg `section_aliases`.
- Dodane testy:
  - deterministyczny diff,
  - apply z zachowaniem BOM/CRLF,
  - wykrywanie konfliktu przy zmianie pliku między planem a apply,
  - test CLI: `--apply` redukuje konkretne błędy validate.

## Status wdrożenia (część 6/6)
- Dodana komenda `docflow health` jako one-command orchestration:
  - tryb lokalny (`warn`) i tryb CI (`--ci`),
  - baseline-aware (`--baseline-mode repo|artifact|none`, `--baseline-dir`),
  - bundle artefaktów: `validate.json`, `validate.sarif`, `compliance.json`, `summary.json`, `summary.md`, `meta.json`.
- `health` agreguje exit codes etapów validate/compliance do jednego `overall_exit` (priorytet `3 > 2 > 1 > 0`).
- Workflow `.github/workflows/docflow.yml` używa `health --ci` i uploaduje SARIF/bundle z katalogu runu.

## Definition of Done (część 1/6)
- Brak zmiany kontraktu wyjścia (`json/sarif`, exit codes, baseline behavior).
- Sortowanie kanonicznych issues realizowane w engine.
- Logika reguł nie jest dodawana do rendererów.
- Testy regresji i Go/No-Go przechodzą.

## Następne kroki
- Rozszerzyć `summary/meta` o granularne czasy etapów i metryki cache hit/miss per etap.
- Rozważyć integrację `fix --format diff` jako opcjonalny artefakt bundle (domyślnie wyłączony).
