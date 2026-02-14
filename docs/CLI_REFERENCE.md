# CLI Reference (MVP)

> Uwaga: część komend (templates/recommend) działa na danych demo; wymagają realnego indeksu w przyszłych iteracjach.

## init
- `docflow init -n <name>` — inicjalizuje bazę SQLite (plan projektu).

## import
- `docflow import -f plan.yaml --clear` — importuje plan (YAML/JSON) do bazy.

## phase / task
- `docflow phase list` — lista faz.
- `docflow phase show 10` — szczegóły dnia 10.
- `docflow task add -p 10 -n "Implement parser"` — dodaje zadanie.

## scan
- `docflow scan -o .docflow/cache/doc_index.json` — buduje indeks dokumentów (metadane, checksum).
- `docflow scan -o .docflow/cache/doc_index.json --deterministic` — zapis bez zmiennego `created_at` (stabilny output do porównań).

## validate
- `docflow validate --strict` — walidacja metadanych; exit code !=0 jeśli błędy.
- `docflow validate --status-aware` — reguły zależne od statusu (draft/review/published).
- `docflow validate --governance docs/_meta/GOVERNANCE_RULES.yaml` — reguły governance (pola/sekcje).
- `docflow validate --old-index .docflow/cache/doc_index.json --auto-bump --save-index` — change tracking i auto-bump wersji.
- `docflow validate --format json --output -` — kontrakt maszynowy JSON.
- `docflow validate --format sarif --output report.sarif` — eksport SARIF 2.1.0.
- `docflow validate --against baseline_validate.json --fail-on new --show new --strict` — gating tylko nowych problemów.

## analyze-patterns
- `docflow analyze-patterns --schema schemas` — skan wzorców sekcji, generuje schematy (demo).

## validate-templates
- `docflow validate-templates --schema section_schema.yaml --dir templates/` — score jakości szablonów vs schema.

## find-duplicates
- `docflow find-duplicates --threshold 0.85` — wyszukuje podobne szablony (demo).

## recommend
- `docflow recommend --doc-type guide --lang pl` — rekomendacje szablonów (demo dane).

## template-sets / templates
- `docflow template-sets` — współwystępowanie szablonów (demo).
- `docflow templates list` — lista szablonów + wersje, metryki content (demo).
- `docflow templates deprecated` — szablony deprecated/archived.
- `docflow template-impact --old-index .docflow/cache/doc_index.json` — dokumenty zależne od zmienionych szablonów.
- `docflow graph --format json --output graph.json` — kanoniczny graf zależności (deterministyczny artefakt source-of-truth).
- `docflow graph --cycles` — szybki podgląd wykrytych cykli.
- `docflow baseline migrate --in baseline_validate_v1.json --out baseline_validate_v2.json --kind validate` — migracja baseline identity v1 → v2.
- `docflow baseline migrate --in baseline_compliance_v1.json --out baseline_compliance_v2.json --kind compliance` — migracja baseline compliance do identity v2.

## migrate-sections
- `docflow migrate-sections --apply` — zamiana legacy nagłówków wg `section_aliases` w configu (domyślnie dry-run).

## changes
- `docflow changes --old-index prev_index.json` — wykrywa added/changed/removed dokumenty.

## plan daily
- `docflow plan daily --max 5` — rekomenduje kolejność pracy (topo + effort heurystyka).

## compliance
- `docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml --format html --html out.html` — raport zgodności (text/json/html).
- `docflow compliance --format json --output - --rules ...` — raport JSON z metadanymi reguł (`rules_path`, `rules_checksum`).
- `docflow compliance --against baseline_compliance.json --fail-on new --show new --strict --rules ...` — gating tylko nowych naruszeń.

## health
- `docflow health` — tryb lokalny (warn), generuje bundle artefaktów w `.docflow/out/<run_id>/`.
- `docflow health --ci --rules docs/_meta/GOVERNANCE_RULES.yaml` — one-command CI: validate (json+sarif) + compliance (json), spójny `overall_exit`.
- `docflow health --ci --baseline-mode repo|artifact|none --baseline-dir .docflow/baseline` — baseline-aware gating (`fail-on new`, gdy baseline dostępny).
- Bundle health zawiera: `validate.json`, `validate.sarif`, `compliance.json`, `summary.json`, `summary.md`, `meta.json`.

## fix
- `docflow fix --format diff --output -` — generuje deterministyczny unified diff (dry-run, bez zapisu).
- `docflow fix --apply --backup-dir .docflow/backup` — stosuje zmiany atomowo (temp+rename), opcjonalnie z backupem.
- `docflow fix --scope frontmatter,sections --only DOCFLOW.VALIDATE.MISSING_FIELD` — ogranicza klasy poprawek / kody.
- `docflow fix --against .docflow/baseline/validate.json --show new` — naprawia wyłącznie nowe findings względem baseline.
- `docflow fix guided --root <dir>` — dotychczasowy guided flow (skryptowy).

## Benchmarks (CI)
- `GOFLAGS=-mod=vendor go test -bench=. -benchmem ./...` — używane w workflow.

## Flagi wspólne
- `--config` (ścieżka do docflow.yaml)
- `--log-level` (debug/info/warn/error)
- `--cache` / `--no-cache` (włącz/wyłącz execution-layer cache SQLite)
- `--cache-dir` (katalog cache execution layer; domyślnie z configu)
- `--changed-only` (parsuj tylko zmienione pliki, ewaluuj pełny zbiór facts)
- `--since <git_ref>` (źródło zmian dla incremental; opcjonalne, fallback gdy brak git)
- `--version` (wersja binarki + commit + data buildu)

## Notatki wydajnościowe
- `docflow graph --include-links` analizuje linki Markdown i na dużych repo może istotnie zwiększyć czas skanu; zalecane użycie z `--cache` i trybem incremental (`--since`, `--changed-only`).
- Domyślnie traversal impact używa bezpiecznych typów krawędzi (`depends_on`, `uses_template`); linki są opt-in.
- `docflow fix` ma bezpieczniki `--max-files` i `--max-changes`; przy dużych repo warto zaczynać od węższego `--scope`.


## Kolejka

Skrypt: `scripts/queue_evaluate.sh` (ocena READY/BLOCKED).

**Parametry**:
- `--bin PATH` — ścieżka do binarki docflow (domyślnie ./build/docflow-linux-amd64)
- `--format text|json` — format wyjścia (domyślnie text)
- `--cache FILE` — plik cache (opcjonalny); unieważnianie: hash docflow.yaml + hash wszystkich *.md w katalogu
- `--log-dir DIR` — logi per zadanie (validate/compliance)
- `--no-cache` — wyłącza cache
- `--workers N` — równoległość (deterministyczne sortowanie inputu)

**Backlog (input)**: plik `task_id path` per linia; `-` czyta ze stdin.

**Wyjście**:
- text: `READY <id> <path>` lub `BLOCKED <id> <path> reason=<...>`
- json: klucze `ready`, `blocked`, `tasks[]` (status, reason, violations[])

**Przykład**:
```bash
scripts/queue_evaluate.sh --bin ./docflow --format json --cache /tmp/qcache.json <<'EOF'
T1 examples/simple-api
T2 examples/architecture
EOF
```

**Cache miss scenariusz**: zmiana pojedynczego MD → miss; zmiana docflow.yaml → miss.
