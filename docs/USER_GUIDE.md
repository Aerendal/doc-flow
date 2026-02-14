# docflow — User Guide (MVP)

## 0. Wymagania
- Go 1.25+ (lokalnie) lub gotowy binarek z Releases.
- Brak dostępu do internetu podczas build/test → używaj `GOFLAGS=-mod=vendor`.
- Ustaw `GOCACHE=/tmp/go-cache` jeśli widzisz `permission denied` dla `~/.cache/go-build`.

## 1. Szybki start (5 minut)
```bash
# instalacja (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/jerzy/docflow/main/install.sh | bash

# inicjalizacja projektu (SQLite plan)
./docflow init -n my-docs

# skan istniejących dokumentów
./docflow scan

# walidacja (status-aware + governance)
./docflow validate --status-aware --governance docs/_meta/GOVERNANCE_RULES.yaml

# rekomendacje szablonów (demo)
./docflow recommend --doc-type guide --lang pl

# compliance report (tekst)
./docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml
```

## 2. Kluczowe pojęcia
- **Frontmatter YAML**: `title`, `doc_id`, `doc_type`, `status`, `version`, `depends_on`, `context_sources`, `template_source`, `tags`.
- **Dependencies**: `depends_on` (twarde zależności), `context_sources` (źródła kontekstu).
- **Sections**: nagłówki H2+ budują SectionTree; aliases w `section_aliases` (fuzzy).
- **Templates**: wersjonowane (semver), status (active/deprecated/archived); używane w generatorze i rekomendacjach.
- **Governance**: reguły status/family w `docs/_meta/GOVERNANCE_RULES.yaml`; walidacja `--governance`.
- **Compliance**: raport zbiorczy (`docflow compliance`) — pass rate, violations.

## 3. Typowy workflow
1. **Przygotuj repo**: dodaj `docflow.yaml`, `docs/_meta/GOVERNANCE_RULES.yaml`, `section_aliases`.
2. **Buduj indeks**: `./docflow scan -o .docflow/cache/doc_index.json`.
3. **Waliduj**: `./docflow validate --status-aware --governance docs/_meta/GOVERNANCE_RULES.yaml --strict`.
4. **Migracja sekcji**: `./docflow migrate-sections --apply` (aliasy → kanoniczne).
5. **Analiza zgodności**: `./docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml --format html`.
6. **CI**: uruchamiaj `GOFLAGS=-mod=vendor GOCACHE=/tmp/go-cache go test ./...`.

## 4. Najważniejsze komendy (skrót)
- `scan` — buduje indeks dokumentów.
- `validate` — walidacja metadanych/sekcji (`--status-aware`, `--governance`, `--old-index`, `--auto-bump`).
- `migrate-sections` — zamiana legacy nagłówków wg aliasów.
- `recommend` — demo rekomendacji szablonów.
- `template-impact` — wpływ zmian szablonów na dokumenty.
- `changes --old-index` — różnice między indeksami (checksum/meta/body).
- `compliance` — raport governance (text/json/html).
- `templates list` — status/wersje + metryki content (demo dane).

## 5. Core concepts — przykłady
**Frontmatter (minimum published)**:
```yaml
---
title: "API Onboarding"
doc_id: "api_onboarding"
doc_type: "guide"
status: "published"
version: "v1.2.0"
owner: "docs-team"
depends_on: ["api_overview"]
context_sources: ["product_vision"]
template_source: "templates/guide_v1.md"
---
```

**Sekcje z aliasami**:
- Kanoniczne: `Przegląd`, `Endpoints`
- Legacy: `Overview`, `API Endpoints` → wykrywane jako aliasy, można zautomatyzować migrację.

## 6. Instalacja
- Script: `curl -sSL https://raw.githubusercontent.com/jerzy/docflow/main/install.sh | bash`
- Homebrew (po wydaniu sha256): `brew install --build-from-source ./homebrew/docflow.rb`
- Manual: pobierz `docflow-<os>-<arch>.tar.gz` z GitHub Releases, rozpakuj do PATH.

## 7. FAQ (skrót)
- **Permission denied w ~/.cache/go-build?** Ustaw `GOCACHE=/tmp/go-cache` lub `HOME=/tmp`.
- **Brak internetu na CI?** Użyj `GOFLAGS=-mod=vendor`; wszystkie zależności są vendored.
- **Dlaczego compliance pokazuje fałszywe brakujące pola?** Governance heurystyka szuka `field:` w treści; upewnij się, że frontmatter zawiera pola i jest poprawnie sformatowany.

## 8. Co dalej
- Poprawić governance (parser-based, quality/approvals).
- Dodać smoke-test binarek w release workflow.
- Uzupełnić checksumy w Homebrew formula i release artefaktach.


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

