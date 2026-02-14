# ROZSZERZONY PLAN PROJEKTU DOCFLOW
## 90 dni, 3 fazy iteracyjne + pre-work

**Źródło:** Szczegółowy 90-dniowy plan operacyjny oparty na koncepcji z [Plan.md](../Plan.md).

---

## RELATED DOCUMENTS

**Navigation:**
- **[← INDEX](INDEX.md)** - Powrót do głównego indeksu dokumentów
- **[RISK_REGISTER.md](RISK_REGISTER.md)** - Ryzyka projektu (18 ryzyk, mitigation plans)
- **[DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md)** - Mapa zależności (critical path, 75 dni)
- **[VALIDATION_REPORT.md](VALIDATION_REPORT.md)** - Raport walidacji spójności
- **[FIXES_APPLIED.md](FIXES_APPLIED.md)** - Changelog naprawionych problemów

**Quick links:**
- Risk management: See [RISK_REGISTER.md](RISK_REGISTER.md) for mitigation plans referenced in this plan
- Dependencies: See [DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md) for complete dependency graph
- Quality: See [VALIDATION_REPORT.md](VALIDATION_REPORT.md) for validation of this plan

---

## PRZEDROSTEK: PRE-WORK (przed day_00)

### Zakres: Przygotowanie fundamentów projektu
### Czas: 5-7 dni roboczych
### Odpowiedzialność: Tech Lead + Product Owner

---

### PRE-1: Analiza wymagań biznesowych (2 dni)

**Cel:** Zdefiniować problem biznesowy i użytkowników docelowych.

**Wejście:**
- Dostęp do interesariuszy (product owner, użytkownicy końcowi)
- Przykłady istniejących problemów z dokumentacją

**Wyjście:**
- `REQUIREMENTS.md` - wymagania biznesowe i funkcjonalne
- `USER_PERSONAS.md` - profile użytkowników (tech writer, architect, developer)
- `USE_CASES.md` - 10-15 głównych scenariuszy użycia
- `SUCCESS_METRICS.md` - KPI projektu (adoption rate, time saved, quality improvement)

**Kryteria ukończenia:**
- Minimum 5 zwalidowanych use cases
- Ustalony primary user persona
- Zdefiniowane metryki sukcesu z progami (np. "50% redukcja czasu tworzenia dokumentu")

**Wyjście awaryjne:**
- Jeśli brak dostępu do użytkowników: użyj desk research + założenia (dokumentuj assumptions)

---

### PRE-2: Pozyskanie korpusu danych (3 dni)

**Cel:** Zdobyć lub wygenerować korpus szablonów dokumentacji.

**Scenariusze:**

**Scenariusz A: Masz dostęp do istniejącego repo**
1. Sklonuj istniejące repo dokumentacji (target: 300+ plików MD)
2. Anonimizuj wrażliwe dane (nazwy projektów, credentials)
3. Waliduj strukturę (min. 50% plików z frontmatter)
4. Sklasyfikuj ręcznie 30-50 przykładów (labels: quality high/medium/low)

**Scenariusz B: Brak istniejącego repo**
1. Zbuduj syntetyczny korpus:
   - Zbierz 100 publicznych szablonów z GitHub (awesome-* repos, tech writing guides)
   - Wygeneruj 200 wariantów z LLM (różne domeny: API, architecture, requirements)
   - Manually curate 50 przykładów jako ground truth
2. Zapisz proces generowania w `DATA_GENERATION.md`
3. Oznacz syntetyczne vs real w metadanych

**Wyjście:**
- `testdata/templates/` - minimum 300 plików Markdown
- `testdata/ground_truth.json` - 30-50 labeled examples
- `DATA_INVENTORY.md` - opis korpusu (źródła, rozkład kategorii, limity)
- `testdata/generator/` - skrypty do powiększania korpusu

**Kryteria ukończenia:**
- Min. 300 plików Markdown z różnorodnością struktur
- Min. 30 labeled examples dla quality scoring
- Zdefiniowany proces powiększania korpusu

**Go/No-Go:**
- Go: Masz min. 300 plików, 30 labels
- No-Go: Zmień zakres na "proof-of-concept with synthetic data only" i dokumentuj limitation

---

### PRE-3: Wybór technologii i bibliotek (1-2 dni)

**Cel:** Podjąć kluczowe decyzje technologiczne przed rozpoczęciem implementacji.

**Decyzje do podjęcia:**

**1. Język (z benchmarku day_00 - ale decyzja wcześniej):**
- Go: szybkość, binary deployment, dobry ekosystem CLI
- Python: prototypowanie, ML libraries, łatwiejsze testowanie wzorców
- Rust: performance, bezpieczeństwo, ale wolniejszy development

**2. Biblioteki Markdown:**
- Go: goldmark (recommended), blackfriday
- Python: mistune, markdown-it-py, python-markdown
- Rust: pulldown-cmark, comrak

**3. Biblioteki YAML:**
- Go: gopkg.in/yaml.v3
- Python: PyYAML, ruamel.yaml
- Rust: serde_yaml

**4. Graf i algorytmy:**
- Go: gonum/graph
- Python: networkx
- Rust: petgraph

**5. CLI framework:**
- Go: cobra + viper
- Python: click + typer
- Rust: clap

**Wyjście:**
- `TECH_STACK.md` - uzasadnione wybory technologiczne
- `ARCHITECTURE_DECISION_RECORDS/` - folder ADR z decyzjami:
  - `ADR-001-language-choice.md`
  - `ADR-002-markdown-parser.md`
  - `ADR-003-dependency-graph-library.md`
- `ENVIRONMENT.md` - wymagania środowiska (OS, wersje, dependencies)

**Kryteria ukończenia:**
- Wszystkie kluczowe biblioteki wybrane
- ADR zapisane i zreviewowane
- Environment setup przetestowany na 2 maszynach

---

### PRE-4: Architektura systemu (1 dzień)

**Cel:** Zaprojektować high-level architekturę przed rozpoczęciem kodowania.

**Wyjście:**
- `ARCHITECTURE.md` - architektura systemu:
  - Component diagram (parser, indexer, validator, graph, generator, recommender)
  - Data flow diagram (od plików MD do raportu/wygenerowanego dokumentu)
  - Moduły i ich interfejsy
  - Strategia cache i invalidation
- `DATA_MODEL.md` - model danych:
  - DocumentRecord schema
  - TemplateRecord schema
  - DependencyGraph structure
  - SectionTree structure
- `API_CONTRACTS.md` - interfejsy między modułami

**Kryteria ukończenia:**
- Diagram architektury zatwierdzony
- Wszystkie główne struktury danych zdefiniowane
- Interfejsy modułów określone (input/output types)

---

### PRE-5: Plan projektu i risk management (1 dzień)

**Cel:** Ustalić plan realizacji i zarządzanie ryzykiem.

**Wyjście:**
- `PROJECT_PLAN.md` - harmonogram 3 faz
- `RISK_REGISTER.md` - rejestr ryzyk:
  - R1: Brak korpusu szablonów → Mitigation: synthetic data
  - R2: Wybrana biblioteka nie wspiera feature X → Mitigation: wrapper pattern
  - R3: Performance issues na dużej skali → Mitigation: early benchmarking
  - R4: Zmiana wymagań → Mitigation: MVP-first approach
- `TEAM_ALLOCATION.md` - alokacja zasobów (kto, kiedy, na czym)
- `COMMUNICATION_PLAN.md` - jak raportować postęp (daily standups? weekly demos?)

**Kryteria ukończenia:**
- Top 10 ryzyk zidentyfikowanych z mitigation plans
- Harmonogram zatwierdzony przez stakeholders
- Team allocation uzgodniony

---

### PRE-WORK GO/NO-GO

**Go:**
- Wszystkie 5 etapów pre-work ukończone
- Minimum 300 plików w korpusie
- Tech stack wybrany i przetestowany
- Architecture zatwierdzony

**No-Go:**
- Brak korpusu danych: pivot na mniejszy zakres lub odłóż projekt
- Brak team capacity: zmniejsz zakres do MVP tylko
- Brak zatwierdzonej architektury: daj sobie dodatkowe 2-3 dni

---

---

## FAZA 1: FOUNDATION & MVP (day_00 - day_25, ~5 tygodni)

### Cel fazy: Zbudować działające minimum viable product
### Milestone: Użytkownik może zeskanować repo, zwalidować metadane, wygenerować brakujący dokument
### Success criteria: MVP działa end-to-end na próbce 50 dokumentów

---

## TYDZIEŃ 1: SETUP & CORE PARSING (day_00-05)

### day_00: Language benchmark + environment setup

**Cel:** Potwierdzić wybór języka benchmarkiem i skonfigurować środowisko.

**Wejście:**
- Pre-work ukończony (tech stack wybrany, ale benchmark potwierdza)
- Korpus 300+ plików z testdata/
- ENVIRONMENT.md

**Wyjście:**
- 3 prototypy POC parsera (Go/Rust/Python) - tylko jeśli benchmark w pre-work był inconclusive
- Benchmark: mediana z 5 runs na 300 plikach
- Decyzja języka zatwierdzona w DECISIONS.md
- Repo Git założone ze strukturą:
  ```
  docflow/
  ├── cmd/          # CLI entrypoint
  ├── pkg/          # public packages
  ├── internal/     # private packages
  ├── testdata/     # test fixtures
  ├── docs/         # project documentation
  │   └── _meta/    # specs
  ├── LOGS/         # decisions, research
  └── .cache/       # runtime cache
  ```
- Environment setup zwalidowany (dependencies installed, builds successfully)

**Kryteria ukończenia:**
- Benchmark pokazuje stabilny czas < 500ms dla 300 plików (95th percentile)
- Repo builds bez błędów
- Struktura katalogów zgodna z PATHS.md

**Testy:**
- Build na czystym środowisku (CI lub clean VM)
- Benchmark powtórzony 3x - variance < 10%

**Czas:** 1 dzień

---

### day_01-02: Metadata contract + Config + Logger + File Walker (2 dni)

**Cel:** Ustalić kontrakt metadanych i zbudować fundamenty narzędzia.

**Wejście:**
- Pre-work (architecture, use cases)
- Wybór języka z day_00
- Korpus danych z analizą 30-50 przykładów

**Wyjście:**

**Dzień 1: Metadata contract**
- `docs/_meta/DOC_META_SCHEMA.md` - pełny schemat metadanych:
  - `doc_id`: string, format: `<family>-<layer>-<slug>`, unique
  - `path`: auto, relative to docs root
  - `rodzina`: enum (api, architecture, requirements, runbook, decision, guide)
  - `warstwa`: int (1-5, C4 model inspired)
  - `status`: enum (draft, review, published, archived, deprecated)
  - `depends_on`: []string, hard dependencies (doc_ids)
  - `context_sources`: []string, soft dependencies (doc_ids)
  - `targets`: []string, output formats or audiences
  - `language`: enum (pl, en), auto-detected
  - `version`: semver (optional)
  - `created`: date (auto from git)
  - `updated`: date (auto from git)
  - `authors`: []string (from git blame or manual)
- `docs/_meta/DOC_DEPENDENCY_SPEC.md` - semantyka zależności:
  - `depends_on`: blocking dependency, required for correctness
  - `context_sources`: informational, provides background
  - Cyclic dependencies handling strategy
- `docs/_meta/DOC_TYPES.md` - typy dokumentów i ich przeznaczenie
- 10 przykładowych dokumentów z poprawnym frontmatter w testdata/
- `LOGS/DECISIONS.md` - zapisane decyzje (format doc_id, reguły zależności)

**Dzień 2: Foundation code**
- `pkg/config/config.go` - loader `docflow.yaml`:
  ```yaml
  docs_root: "./docs"
  cache_dir: "./.cache"
  ignore_patterns:
    - "node_modules"
    - ".git"
    - "_site"
  validation:
    strict_mode: false
    required_fields: [doc_id, rodzina, status]
  ```
- `internal/logger/logger.go` - structured logger (JSON output dla CI)
- `internal/util/fileutil.go` - file walker with ignore patterns
- `cmd/docflow/main.go` - CLI stub z `--help`
- Przykładowy `docflow.yaml` w repo root
- Unit testy dla config loader i file walker

**Kryteria ukończenia:**
- Schemat metadanych jest kompletny bez niejednoznaczności
- 10 przykładowych dokumentów przechodzi ręczną walidację
- `docflow --help` działa
- `docflow` wczytuje config bez błędów
- File walker zwraca listę .md plików z docs/ (ignore działa)

**Testy:**
- Config loader: valid config, invalid YAML, missing required fields
- File walker: 100 plików z 10 ignored, verify count = 90
- End-to-end: `docflow` uruchomiony na testdata/ zwraca listę plików

**Czas:** 2 dni

---

### day_03-04: Markdown + YAML parser + Document Index (2 dni)

**Cel:** Parsować dokumenty i budować indeks.

**Wejście:**
- File walker z day_02
- Metadata schema z day_01
- Korpus 300+ plików

**Wyjście:**

**Dzień 3: Parser**
- `pkg/parser/frontmatter.go` - YAML frontmatter parser:
  - Extract YAML z dokumentu (--- delimiters)
  - Fallback: szukaj bloku `## Metadata` + YAML code block
  - Validate required fields
  - Error reporting (file, line, field)
- `pkg/parser/markdown.go` - Markdown parser:
  - Extract headings (AST traversal)
  - Poziomy h1-h6
  - Anchor extraction (dla heading links)
- `pkg/models/document.go` - DocumentRecord struct
- Unit testy na 20 przykładach (valid, invalid YAML, missing frontmatter)

**Dzień 4: Document Index**
- `pkg/index/document_index.go` - in-memory index:
  - Map: doc_id → DocumentRecord
  - Methods: Add, Get, List, Find
- `pkg/index/cache.go` - JSON persistence:
  - Serialize to `.cache/doc_index.json`
  - Deserialize from cache
  - SHA256 hash per file dla invalidation (zapisz w cache)
- `cmd/docflow/scan.go` - komenda `docflow scan`:
  - Walk docs/
  - Parse każdy .md
  - Build index
  - Save to cache
  - Print summary (files scanned, errors)
- Integration test: scan testdata/ (50 docs), verify index.json

**Kryteria ukończenia:**
- `docflow scan` generuje poprawny indeks JSON
- Indeks zawiera wszystkie pola z schema
- Parser obsługuje oba przypadki (frontmatter + fallback)
- Cache save/load bez utraty danych
- Scan 300 plików < 2s (target), < 5s (acceptable)

**Testy:**
- Parser: 10 valid docs, 5 invalid YAML, 3 missing frontmatter
- Index: add 100 docs, serialize, deserialize, compare
- Cache invalidation: zmień 1 plik, rescan, verify tylko 1 plik reparsed
- Performance: benchmark scan 300 files, measure time

**Czas:** 2 dni

---

### day_05: Code review + tech debt cleanup

**Cel:** Review kodu z tygodnia 1, refactor, fix tech debt.

**Wejście:**
- Kod z day_00-04
- Pull requesty do review

**Wyjście:**
- Code review notes w `LOGS/CODE_REVIEW_WEEK1.md`
- Refactored code (według code review)
- Updated tests
- CI pipeline setup (basic):
  - Build verification
  - Unit tests
  - Lint (golangci-lint / ruff / clippy)

**Kryteria ukończenia:**
- Wszystkie PR z week 1 zreviewowane
- Zero critical issues w code review
- CI pipeline przechodzi green

**Czas:** 1 dzień

---

## TYDZIEŃ 2: VALIDATION & DEPENDENCY GRAPH (day_06-10)

### day_06: Metadata validator

**Cel:** Walidować spójność metadanych.

**Wejście:**
- Document index z day_04
- Metadata schema z day_01

**Wyjście:**
- `pkg/validator/metadata.go` - walidator metadanych:
  - Unikalność doc_id
  - Required fields present
  - Valid enum values (rodzina, status, language)
  - Valid doc_id format (regex)
  - Path consistency (opcjonalne ostrzeżenie)
- `cmd/docflow/validate.go` - komenda `docflow validate`:
  - Load index
  - Run metadata validation
  - Print raport błędów:
    ```
    ERRORS (blocking):
    - docs/api-guide.md: duplicate doc_id 'api-guide-001'
    - docs/arch.md: missing required field 'rodzina'

    WARNINGS (non-blocking):
    - docs/req.md: doc_id 'requirements-l2-auth' doesn't match path
    ```
  - Exit code 1 jeśli errors, 0 jeśli tylko warnings
- Tryby: `--strict` (warnings = errors), `--warn` (ignore errors, report only)
- Testy: 15 test cases (valid, duplicates, missing fields, invalid enums)

**Kryteria ukończenia:**
- `docflow validate` wykrywa wszystkie typy błędów w test cases
- Raport jest czytelny i wskazuje plik + pole + błąd
- Strict mode działa poprawnie

**Testy:**
- Prepare 10 invalid docs: duplicates, missing fields, bad enums
- Run validate, verify all 10 detected
- Verify exit codes

**Czas:** 1 dzień

---

### day_07-08: Section parser + Section schema (2 dni)

**Cel:** Parsować strukturę sekcji i budować drzewo hierarchii.

**Wejście:**
- Markdown parser z day_03
- Przykładowe dokumenty z sekcjami

**Wyjście:**

**Dzień 7: Section parser**
- `pkg/sections/parser.go` - section parser:
  - Traverse Markdown AST
  - Build SectionTree (hierarchia nagłówków)
  - Model:
    ```go
    type Section struct {
        Level    int       // 1-6 (h1-h6)
        Title    string
        Anchor   string    // heading-slug
        Content  string    // text content (bez children)
        Children []*Section
        Parent   *Section
    }
    type SectionTree struct {
        Root     *Section
        Sections []*Section // flat list
    }
    ```
  - Handle edge cases: missing h1, skipped levels (h2 → h4)
- Testy: 10 dokumentów z różnymi strukturami (3-level, missing h2, flat)

**Dzień 8: Section schema**
- `docs/_meta/SECTION_SCHEMA.md` - format schematów sekcji:
  ```yaml
  rodzina: api
  required_sections:
    - title: "Overview"
      level: 2
      required: true
    - title: "Endpoints"
      level: 2
      required: true
      subsections:
        - title: "Authentication"
          level: 3
          required: true
  optional_sections:
    - title: "Examples"
      level: 2
  ```
- `pkg/schema/section_schema.go` - loader schematów
- 2-3 przykładowe schematy dla rodzin (api, architecture, requirements)
- `pkg/sections/validator.go` - walidator sekcji (basic):
  - Check required sections present
  - Report missing sections

**Kryteria ukończenia:**
- Parser buduje poprawną hierarchię dla testów
- Schema loader wczytuje YAML schemas
- Validator wykrywa brakujące required sections

**Testy:**
- Parse 10 docs, verify tree structure ręcznie dla 3
- Load 3 schemas, verify parsed correctly
- Validate 5 docs against schema, 3 pass, 2 fail (missing sections)

**Czas:** 2 dni

---

### day_09-10: Dependency graph + topological sort (2 dni)

**Cel:** Budować graf zależności i sortowanie topologiczne.

**Wejście:**
- Document index z day_04
- Metadata contract z day_01 (depends_on)

**Wyjście:**

**Dzień 9: Dependency graph**
- `pkg/graph/dependency.go` - graf zależności:
  - Build directed graph z depends_on
  - Nodes: doc_id
  - Edges: A depends_on B → edge B→A (B must come before A)
  - Detect cycles (strongly connected components)
  - Detect missing nodes (unresolved dependencies)
- Model:
  ```go
  type DependencyGraph struct {
      Nodes map[string]*DocumentRecord
      Edges map[string][]string  // doc_id → list of dependencies
  }
  ```
- Cycle detection algorithm (DFS)
- Missing dependency detection

**Dzień 10: Topological sort + context sources**
- `pkg/graph/toposort.go` - topological sort:
  - Kahn's algorithm
  - Tie-break: lexicographic sort by doc_id (deterministyczny)
  - Return: ordered list of doc_ids
- `pkg/graph/context.go` - semantyka context_sources:
  - context_sources NIE tworzą edges w grafie
  - Raport: brakujące context_sources (warning, nie error)
  - Opcja config: `promote_context_to_dependency` per rodzina
- `cmd/docflow/graph.go` - komenda `docflow graph`:
  - Print dependency graph (DOT format dla Graphviz)
  - Print topological order
  - Print cycles (jeśli są)
  - Print missing dependencies
- Testy: 20 test cases (valid DAG, cycles, missing deps, ties)

**Kryteria ukończenia:**
- Graf buduje się poprawnie dla valid dependencies
- Topo sort jest deterministyczny (3 runs dają identyczny wynik)
- Cycles są wykrywane
- Missing dependencies są raportowane
- context_sources handling działa (nie blokują, raport oddzielny)

**Testy:**
- Create test graph: 10 nodes, 15 edges, no cycles → verify topo order
- Create test graph with cycle A→B→C→A → verify cycle detected
- Create test graph with missing dep (A depends on X, X not exists) → verify reported
- Run topo sort 10 times, verify identical output

**Czas:** 2 dni

---

## TYDZIEŃ 3: MVP FEATURES (day_11-15)

### day_11-12: Document generator (2 dni)

**Cel:** Generować dokumenty z szablonów.

**Wejście:**
- Section schemas z day_08
- Metadata contract z day_01
- Korpus szablonów z pre-work

**Wyjście:**

**Dzień 11: Template index**
- `pkg/templates/index.go` - indeks szablonów:
  - Scan testdata/templates/
  - Parse metadata templates
  - Store: rodzina, warstwa, language, sections
- `pkg/templates/selector.go` - basic template selection:
  - Filter by rodzina + warstwa + language
  - Return matching templates (nie ma jeszcze scoringu)

**Dzień 12: Generator**
- `pkg/generator/generator.go` - generator dokumentów:
  - Input: doc_id, rodzina, warstwa, language, output_path
  - Steps:
    1. Select template (first match from selector)
    2. Load template content
    3. Replace placeholders: {{doc_id}}, {{title}}, {{date}}
    4. Generate frontmatter z podanych parametrów
    5. Fill required sections z placeholders
    6. Write to output_path
- `cmd/docflow/generate.go` - komenda `docflow generate`:
  - Args: --doc-id, --rodzina, --warstwa, --language, --output
  - Mode: non-interactive (parametry z CLI)
  - Preview: print first 30 lines before write (optional --preview flag)
  - Confirm before write (optional --yes flag to skip)
- Testy: generate 5 documents różnych rodzin

**Kryteria ukończenia:**
- `docflow generate` tworzy poprawny plik MD
- Plik ma frontmatter z wszystkimi required fields
- Plik ma wymagane sekcje (z placeholders)
- Generator działa non-interactively (automatable)

**Testy:**
- Generate doc: rodzina=api, verify output ma Overview, Endpoints sections
- Generate doc z --preview, verify nie zapisuje pliku
- Generate doc z --yes, verify zapisuje bez confirmation

**Czas:** 2 dni

---

### day_13: MVP integration test + dogfooding

**Cel:** Przetestować pełny pipeline MVP end-to-end.

**Wejście:**
- Wszystkie moduły z day_00-12

**Wyjście:**
- `tests/integration/mvp_pipeline_test.go` - test E2E:
  ```
  Test scenario: New documentation project
  1. Create empty docs/ dir
  2. Generate 5 documents (różne rodziny)
  3. Scan docs/
  4. Validate metadata
  5. Validate sections
  6. Build dependency graph
  7. Compute topo order
  8. Verify: no errors, graph is DAG, topo order contains 5 docs
  ```
- Dogfooding: użyj docflow do zarządzania własną dokumentacją projektu
  - Przenies docs/_meta/ files do docs/ i dodaj frontmatter
  - Scan projekt docflow na sobie
  - Validate projekt docflow
- `LOGS/MVP_TEST_RESULTS.md` - wyniki testów i znalezione błędy
- Bug fixes z dogfooding

**Kryteria ukończenia:**
- E2E test przechodzi green
- Dogfooding wykrył min. 3 bugi i są poprawione
- MVP działa na rzeczywistych danych projektu

**Czas:** 1 dzień

---

### day_14-15: MVP documentation + demo (2 dni)

**Cel:** Dokumentacja MVP i przygotowanie demo.

**Wyjście:**

**Dzień 14: Documentation**
- `README.md` - Quick start guide:
  - Installation (binary / from source)
  - 5-minute tutorial
  - Basic commands (scan, validate, generate, graph)
- `docs/USER_GUIDE.md` - Pełny user guide:
  - Metadata schema explained
  - Dependency graph concepts
  - Section schemas
  - Generator usage
- `docs/CLI_REFERENCE.md` - Command reference (auto-generated z --help)
- `docs/TROUBLESHOOTING.md` - Common issues

**Dzień 15: Demo + Release MVP**
- Przygotuj demo (15-min screencast lub live demo):
  - Scenario: "Starting new API documentation project"
  - Show: scan, validate, generate, graph visualization
- Release MVP:
  - Tag v0.1.0-mvp
  - Build binaries (Linux, macOS, Windows)
  - Package release (binary + README + examples)
- `CHANGELOG.md` - Release notes dla MVP

**Kryteria ukończenia:**
- Dokumentacja pozwala nowemu użytkownikowi uruchomić MVP w <30 min
- Demo działa i prezentuje core value
- Release artifacts gotowe

**Czas:** 2 dni

---

## TYDZIEŃ 4: SCHEMA EXTRACTION & QUALITY (day_16-20)

### day_16-17: Section pattern extraction (2 dni)

**Cel:** Wydobyć wzorce sekcji z korpusu szablonów.

**Wejście:**
- Korpus 300+ szablonów z testdata/
- Section parser z day_07

**Wyjście:**

**Dzień 16: Pattern extraction**
- `pkg/patterns/extractor.go` - ekstrakcja wzorców:
  - Parse wszystkie templates
  - Extract section sequences (list of titles + levels)
  - Group podobne sequences:
    - N-gram similarity (n=2,3)
    - Levenshtein distance dla titles
    - Threshold: 0.80 similarity → same group
  - Zapisz parametry algorytmu w `docs/_meta/ALGO_PARAMS.md`:
    ```yaml
    pattern_extraction:
      ngram_size: [2, 3]
      similarity_threshold: 0.80
      levenshtein_max_distance: 3
      min_pattern_frequency: 3  # pattern musi wystąpić min 3x
    ```
- Compute frequency dla każdego patternu

**Dzień 17: Schema generation**
- `pkg/patterns/schema_generator.go` - generowanie schematów:
  - Top 50 najczęstszych patterns
  - Convert pattern → section_schema.yaml
  - Required sections: те, które występują w >80% wzorców w grupie
  - Optional sections: 50-80%
- `cmd/docflow/analyze-patterns.go` - komenda `docflow analyze-patterns`:
  - Print top 50 patterns z częstotliwością
  - Generate schemas dla top 10
  - Save to `generated_schemas/`
- Raport: `LOGS/PATTERN_ANALYSIS.md` z przykładami patterns

**Kryteria ukończenia:**
- `docflow analyze-patterns` wykrywa min. 20 distinct patterns
- Top 10 patterns mają frequency > 5
- Generated schemas są valid YAML

**Testy:**
- Ręcznie sprawdź 5 top patterns - verify grupowanie ma sens
- Validate generated schemas przeciwko sample documents

**Czas:** 2 dni

---

### day_18-19: Quality scoring (2 dni)

**Cel:** Oceniać jakość szablonów i dokumentów.

**Wejście:**
- Ground truth labels z pre-work (30-50 labeled templates)
- Section schemas z day_16-17

**Wyjście:**

**Dzień 18: Quality metrics**
- `pkg/quality/metrics.go` - metryki jakości:
  - Structure score (0-100):
    - Required sections present: +40
    - Correct section order: +20
    - Proper heading levels (no skipped): +20
    - Has examples (code blocks): +20
  - Content score (0-100):
    - Section length (not empty): +30
    - Has tables/lists: +20
    - Has links (internal/external): +20
    - Language quality (word count, no "TODO"): +30
  - Usage score (0-100):
    - Based on usage data (jeśli dostępne)
    - Placeholder: 50 if no data
  - Overall score: weighted average (structure 40%, content 40%, usage 20%)
- Zapisz wagi w `ALGO_PARAMS.md`

**Dzień 19: Scoring + validation**
- `pkg/quality/scorer.go` - scoring engine:
  - Score template/document
  - Return breakdown (structure, content, usage, overall)
- `cmd/docflow/score.go` - komenda `docflow score`:
  - Score wszystkie templates
  - Print raport sorted by overall score
  - Opcja: `--rodzina=api` (filter)
- Validation z ground truth:
  - Compare computed scores z manual labels
  - Compute correlation (Spearman's rho)
  - Target: rho > 0.60 (moderate correlation)
  - If < 0.60: tune weights w ALGO_PARAMS.md
- Testy: score 30 labeled templates, verify correlation

**Kryteria ukończenia:**
- Scoring działa na korpusie
- Correlation z ground truth > 0.60
- Wagi są udokumentowane

**Czas:** 2 dni

---

### day_20: Buffer + tech debt

**Cel:** Buffer na nieprzewidziane problemy + cleanup.

**Wyjście:**
- Bug fixes z week 3-4
- Refactoring z code review
- Dokumentacja update

**Czas:** 1 dzień

---

## TYDZIEŃ 5: TESTING & PHASE 1 RELEASE (day_21-25)

### day_21: Test strategy + unit tests

**Cel:** Ustalić strategię testową i pokryć kod testami.

**Wejście:**
- Kod z day_00-20

**Wyjście:**
- `docs/_meta/TEST_STRATEGY.md` - strategia testowa:
  - Framework (testify, pytest, etc.)
  - Konwencje nazewnictwa (test_*, *_test.go)
  - Coverage target: 70% lines
  - Test categories: unit, integration, e2e
  - Test data strategy (fixtures w testdata/)
- Unit testy dla wszystkich pkg/:
  - parser (frontmatter, markdown, sections)
  - validator (metadata, sections)
  - graph (dependency, toposort)
  - generator
  - quality/metrics
- Coverage raport

**Kryteria ukończenia:**
- Coverage > 70% lines
- Wszystkie critical paths mają testy

**Czas:** 1 dzień

---

### day_22: Integration tests

**Cel:** Testy integracyjne głównych flow.

**Wyjście:**
- `tests/integration/` - testy integracyjne:
  - scan_and_index_test: scan 100 files, verify index
  - validate_pipeline_test: scan + validate metadata + sections
  - graph_pipeline_test: scan + build graph + toposort
  - generate_pipeline_test: generate + scan + validate
  - full_pipeline_test: scan + validate + graph + generate missing
- Test fixtures: 100 dokumentów z różnorodnością (valid, invalid, cycles, etc.)

**Kryteria ukończenia:**
- Wszystkie integration tests pass
- Pipeline działa na 100 docs < 10s

**Czas:** 1 dzień

---

### day_23: Performance testing

**Cel:** Zmierzyć wydajność na większej skali.

**Wejście:**
- Korpus 300+ plików

**Wyjście:**
- `tests/performance/benchmark_test.go` - benchmarki:
  - Scan 300 files: target < 2s, acceptable < 5s
  - Build graph 300 nodes: target < 1s
  - Topo sort 300 nodes: target < 500ms
  - Validate 300 docs: target < 3s
- Raport wydajności: `LOGS/PERFORMANCE_BASELINE.md`
  - Time, memory, CPU per operation
  - Bottlenecks identified (profiling)
  - Optimization recommendations (defer to phase 2)

**Kryteria ukończenia:**
- Benchmarki przechodzą target lub acceptable
- Bottlenecks identified

**Czas:** 1 dzień

---

### day_24: Documentation finalization

**Cel:** Dokończyć dokumentację Phase 1.

**Wyjście:**
- `README.md` updated
- `docs/ARCHITECTURE.md` - actual architecture (as-built)
- `docs/CONTRIBUTING.md` - contribution guide
- `CHANGELOG.md` - updated dla v0.1.0

**Czas:** 1 dzień

---

### day_25: Phase 1 Release + Retrospective

**Cel:** Wydać Phase 1 i przeprowadzić retrospektywę.

**Wyjście:**
- Release v0.1.0:
  - Tag v0.1.0
  - Build binaries (3 platforms)
  - GitHub release z notes
  - Package artifacts
- `LOGS/PHASE1_RETROSPECTIVE.md` - retrospektywa:
  - Co poszło dobrze
  - Co poszło źle
  - Lessons learned
  - Adjustments dla Phase 2
- Demo dla stakeholders (30 min)

**Kryteria ukończenia:**
- v0.1.0 released publicznie
- Retrospektywa przeprowadzona
- Demo pokazany stakeholders

**Czas:** 1 dzień

---

### PHASE 1 GO/NO-GO

**Go:**
- MVP działa end-to-end na 100+ docs
- Wszystkie core features działają (scan, validate, graph, generate)
- Tests coverage > 70%
- Performance acceptable
- Dokumentacja kompletna
- Demo successful

**No-Go:**
- Major bugs w core features → Extend Phase 1 o 1 tydzień
- Performance unacceptable (>10s dla 100 docs) → Investigate bottlenecks, defer some features
- Test coverage < 50% → Extend Phase 1 o 3 dni (testing focus)

**Decision:** If No-Go, assess if project should pivot, pause, or extend.

---

---

## FAZA 2: INTELLIGENCE & AUTOMATION (day_26 - day_55, ~6 tygodni)

### Cel fazy: Dodać inteligencję (rekomendacje, quality) i automatyzację (planner, lifecycle)
### Milestone: System rekomenduje szablony, planuje pracę, śledzi lifecycle
### Success criteria: Recommender precision@5 > 0.70, Daily planner generates useful plans

---

## TYDZIEŃ 6: TEMPLATE RECOMMENDATION (day_26-30)

### day_26: Template metadata expansion

**Cel:** Rozszerzyć metadane szablonów o dodatkowe atrybuty.

**Wejście:**
- Template index z day_11
- Quality scores z day_18-19

**Wyjście:**
- `pkg/templates/metadata.go` - rozszerzone metadane:
  - Quality score (z day_19)
  - Pattern ID (z day_16)
  - Usage count (placeholder: 0)
  - Version (semver)
  - Status (active, deprecated, archived)
  - Last used (timestamp)
  - Examples count (code blocks + tables)
  - Word count
  - Language
- Reindex templates z nowymi polami
- `cmd/docflow/templates.go` - komenda `docflow templates list`:
  - List all templates z metadanymi
  - Filter: --rodzina, --status, --min-quality
  - Sort: --sort-by=quality|usage|updated

**Kryteria ukończenia:**
- Template index ma rozszerzone metadane
- `docflow templates list` działa z filtrami

**Czas:** 1 dzień

---

### day_27-28: Recommendation engine (2 dni)

**Cel:** Zbudować engine rekomendacji szablonów.

**Wejście:**
- Template index z day_26
- Quality scores
- (Opcjonalnie) Usage data

**Wyjście:**

**Dzień 27: Scoring algorithm**
- `pkg/recommender/scorer.go` - algorytm scoringu:
  - Input: rodzina, warstwa, language, (optional) context (existing docs)
  - Scoring factors:
    - Rodzina match: +40 (exact), +20 (compatible)
    - Warstwa match: +20 (exact), +10 (±1 layer)
    - Language match: +20 (exact), +5 (fallback to en)
    - Quality score: +15 (normalized)
    - Usage: +5 (normalized, log scale)
  - Wagi zapisane w `ALGO_PARAMS.md`
  - Output: sorted list of templates z scores + explanations

**Dzień 28: Recommender + CLI**
- `pkg/recommender/recommender.go` - recommender:
  - Recommend(rodzina, warstwa, language) → top K templates
  - Explanations (dlaczego każdy template)
- `cmd/docflow/recommend.go` - komenda `docflow recommend`:
  - Args: --rodzina, --warstwa, --language, --top-k (default 5)
  - Output:
    ```
    Top 5 recommended templates for rodzina=api, warstwa=2, language=pl:

    1. api-guide-v2.md (score: 95/100)
       - Exact rodzina match (+40)
       - Exact warstwa match (+20)
       - Exact language match (+20)
       - High quality (88/100, +13)
       - Usage: 45 times (+5)

    2. api-reference-v1.md (score: 82/100)
       ...
    ```
- Usage tracking:
  - When `docflow generate --template=X` → increment X.usage_count
  - Save to `.cache/template_usage.json`

**Kryteria ukończenia:**
- Recommender zwraca deterministyczne wyniki
- Explanations są czytelne
- Usage tracking działa

**Testy:**
- Create 20 templates z różnymi metadata
- Recommend dla 5 scenarios
- Verify ranking ma sens (ręcznie dla 3 scenarios)

**Czas:** 2 dni

---

### day_29: Recommendation evaluation

**Cel:** Ocenić jakość rekomendacji.

**Wejście:**
- Recommender z day_27-28
- (Idealnie) Ground truth: 10-20 scenarios z expected top templates

**Wyjście:**
- `tests/recommender/evaluation_test.go` - evaluation:
  - 10 test scenarios (rodzina, warstwa, language) z expected top-3
  - Compute metrics:
    - Precision@K (K=1,3,5)
    - Mean Reciprocal Rank (MRR)
  - Target: Precision@5 > 0.60
- `LOGS/RECOMMENDER_EVAL.md` - wyniki evaluation
- Jeśli Precision@5 < 0.60: tune wagi w ALGO_PARAMS.md i retest

**Kryteria ukończenia:**
- Evaluation pokazuje Precision@5 > 0.60 (lub uzasadnij dlaczego nie)

**Czas:** 1 dzień

---

### day_30: Buffer

**Cel:** Buffer na nieprzewidziane problemy z recommendem.

**Czas:** 1 dzień

---

## TYDZIEŃ 7: DEPENDENCY RULES & PLANNER (day_31-35)

### day_31-32: Family dependency rules (2 dni)

**Cel:** Zdefiniować reguły zależności per rodzina dokumentów.

**Wejście:**
- Document types z day_01
- Use cases z pre-work

**Wyjście:**

**Dzień 31: Rules definition**
- `docs/_meta/FAMILY_RULES.yaml` - reguły per rodzina:
  ```yaml
  families:
    - name: architecture
      required_dependencies:
        - rodzina: requirements
          min_count: 1  # must depend on at least 1 requirements doc
      optional_dependencies:
        - rodzina: decision
      expected_context:
        - rodzina: guide

    - name: api
      required_dependencies:
        - rodzina: architecture
          min_count: 1
      expected_sections:
        - Overview
        - Endpoints
        - Authentication
  ```
- Reguły dla 5-7 rodzin

**Dzień 32: Validator integration**
- `pkg/validator/family_rules.go` - walidator reguł:
  - Load FAMILY_RULES.yaml
  - Per dokument: compute expected dependencies
  - Compare expected vs actual (depends_on)
  - Report drift: missing expected dependencies
- `cmd/docflow/validate.go` - extend:
  - Add `--family-rules` flag
  - Print raport drift
- Testy: 10 test docs z różnymi rodzinami

**Kryteria ukończenia:**
- Family rules validator działa
- Raport drift jest czytelny

**Czas:** 2 dni

---

### day_33-34: Daily planner (2 dni)

**Cel:** Stworzyć planer dzienny dla pracy nad dokumentacją.

**Wejście:**
- Dependency graph + toposort z day_09-10
- Quality scores z day_18-19
- Family rules z day_31-32

**Wyjście:**

**Dzień 33: Effort estimation**
- `pkg/planner/effort.go` - estymacja wysiłku:
  - Heuristics (brak historical data):
    - Baseline effort per rodzina:
      - requirements: 2h (simple)
      - api: 3h (moderate)
      - architecture: 5h (complex)
      - guide: 4h (moderate)
    - Modifiers:
      - warstwa +1: +1h
      - Missing template (no match): +2h
      - Complex dependencies (>3): +1h
  - Output: effort_hours per document
  - Zapisz heuristics w `ALGO_PARAMS.md`

**Dzień 34: Planner**
- `pkg/planner/daily.go` - daily planner:
  - Input: max_hours (e.g., 8h), priority (optional)
  - Algorithm:
    1. Get topo order (respects dependencies)
    2. Compute effort dla każdego doc
    3. Greedy select: add docs do planu dopóki sum(effort) < max_hours
    4. Group docs z tej samej rodziny (batch efficiency)
  - Output: ordered list of docs + total effort
- `cmd/docflow/plan.go` - komenda `docflow plan`:
  - Args: --daily (default 8h), --max-hours, --output (MD file)
  - Generate `DAILY_PLAN.md`:
    ```markdown
    # Daily Plan: 2024-02-06

    Total effort: 7.5h / 8h

    ## Documents to complete (5):

    - [ ] requirements-l1-auth (2h) [rodzina: requirements]
    - [ ] architecture-l2-auth (5h) [rodzina: architecture]
    - [ ] api-l3-login (0.5h) [rodzina: api]

    ## Dependencies satisfied:
    - requirements-l1-auth has no dependencies
    - architecture-l2-auth depends on requirements-l1-auth (in plan)
    - api-l3-login depends on architecture-l2-auth (in plan)
    ```
  - Opcja: `--interactive` (zaznaczaj completed docs, update plan)

**Kryteria ukończenia:**
- Planner generuje plan < 10s dla 100 docs
- Plan jest deterministyczny
- Plan respects dependencies (no doc before its deps)

**Testy:**
- Create test graph: 20 docs z dependencies + efforts
- Generate plan --max-hours=8
- Verify: sum(effort) <= 8, topo order preserved

**Czas:** 2 dni

---

### day_35: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 8: LIFECYCLE & VERSIONING (day_36-40)

### day_36-37: Template lifecycle (2 dni)

**Cel:** Zarządzanie cyklem życia szablonów.

**Wejście:**
- Template metadata z day_26
- Quality scores
- Usage data

**Wyjście:**

**Dzień 36: Lifecycle states**
- `docs/_meta/TEMPLATE_LIFECYCLE.md` - lifecycle definition:
  - States: draft → active → deprecated → archived
  - Transitions:
    - draft → active: manual (gdy ready)
    - active → deprecated: quality < 50 AND usage == 0 for 90 days
    - deprecated → archived: manual (safe to remove)
  - Rules zapisane w `ALGO_PARAMS.md`
- `pkg/templates/lifecycle.go` - lifecycle manager:
  - Compute recommended transitions per template
  - Apply transitions (with confirmation)

**Dzień 37: Deprecation + migration**
- `cmd/docflow/templates.go` - extend:
  - `docflow templates deprecate --id=X` - mark as deprecated
  - `docflow templates deprecated` - list deprecated templates
  - `docflow templates suggest-migration` - dla deprecated, suggest active alternative (highest similarity + quality)
- Migration path:
  - Dla każdego deprecated template, find best active replacement (content similarity)
  - List documents używające deprecated template
  - Suggest batch migration

**Kryteria ukończenia:**
- Lifecycle transitions działają
- Deprecation + migration suggestions działają

**Czas:** 2 dni

---

### day_38-39: Versioning (2 dni)

**Cel:** Wersjonowanie szablonów i dokumentów.

**Wejście:**
- Template metadata z day_26

**Wyjście:**

**Dzień 38: Template versioning**
- `pkg/templates/versioning.go` - versioning:
  - Template version w metadanych (semver: v1.0.0)
  - Wiele wersji tego samego template (różne pliki):
    - `api-guide-v1.md`
    - `api-guide-v2.md`
  - Version compatibility rules:
    - Recommender preferuje latest version (if not deprecated)
    - Generowane docs tracą użytą version template

**Dzień 39: Document versioning**
- Document version w metadanych (optional, default v1.0.0)
- Change tracking:
  - Przy każdym `docflow validate`, compare current z cache
  - Detect changes (metadata changed, content changed)
  - Raport: modified documents
  - Opcja: auto-bump version (minor +1 dla content change)

**Kryteria ukończenia:**
- Multiple versions tego samego template współistnieją
- Change tracking wykrywa modyfikacje

**Czas:** 2 dni

---

### day_40: Code review + refactoring

**Czas:** 1 dzień

---

## TYDZIEŃ 9: VALIDATION ENHANCEMENTS (day_41-45)

### day_41: Section completeness metrics

**Cel:** Metryki kompletności sekcji.

**Wejście:**
- Section validator z day_07-08
- Section schemas

**Wyjście:**
- `pkg/sections/metrics.go` - metryki:
  - Per dokument:
    - Required sections present: X / Y
    - Optional sections present: Z / W
    - Completeness score: (X/Y) * 100
    - Empty sections count
  - Per rodzina:
    - Average completeness
    - Most common missing section
- `cmd/docflow/stats.go` - komenda `docflow stats`:
  - Args: --by-rodzina, --by-status
  - Print aggregated metrics

**Kryteria ukończenia:**
- Metrics compute correctly
- Stats command działa

**Czas:** 1 dzień

---

### day_42: Progressive validation (draft vs published)

**Cel:** Różne reguły walidacji dla draft vs published.

**Wejście:**
- Validator z day_06
- Section validator z day_07-08
- Document status z metadata

**Wyjście:**
- `pkg/validator/progressive.go` - progressive validation:
  - Rules:
    - draft: allow empty sections, missing optional deps (warnings only)
    - review: stricter (no empty required sections)
    - published: strictest (no empty sections, all required deps)
  - Config w `docflow.yaml`:
    ```yaml
    validation:
      draft:
        allow_empty_sections: true
        allow_missing_context: true
      published:
        allow_empty_sections: false
        require_all_deps: true
    ```
- Extend `docflow validate` z `--status-aware` flag

**Kryteria ukończenia:**
- Draft docs pass z warnings
- Published docs fail jeśli empty sections

**Testy:**
- Test draft doc z empty section: expect warning
- Test published doc z empty section: expect error

**Czas:** 1 dzień

---

### day_43-44: Edge case hardening (2 dni)

**Cel:** Obsługa edge cases.

**Wejście:**
- Wszystkie validatory

**Wyjście:**

**Dzień 43: Edge cases**
- `tests/edge_cases/` - test cases:
  - Duplicate doc_ids
  - Circular dependencies (A→B→C→A)
  - Missing dependencies (A depends on X, X nie istnieje)
  - Invalid frontmatter (malformed YAML)
  - Empty documents
  - Documents bez headings
  - Documents z duplicate section titles
  - Bardzo długie documents (10k+ lines)
  - Binary files w docs/ (images, PDFs)
- Testy dla każdego edge case
- Fix bugs znalezione przez testy

**Dzień 44: Fuzzy matching + migration**
- `pkg/sections/fuzzy.go` - fuzzy matching dla section names:
  - Handle legacy names (PL/EN variants):
    - "Przegląd" vs "Overview"
    - "Endpoints" vs "API Endpoints"
  - Alias mapping w config:
    ```yaml
    section_aliases:
      "Przegląd": ["Overview", "Podsumowanie"]
      "Endpoints": ["API Endpoints", "API"]
    ```
  - Validator uses fuzzy matching (warning jeśli alias, nie exact match)
- `cmd/docflow/migrate-sections.go` - migrate section names:
  - `docflow migrate-sections --dry-run` - preview changes
  - `docflow migrate-sections --apply` - rewrite files

**Kryteria ukończenia:**
- Wszystkie edge cases mają testy
- Fuzzy matching działa
- Migration tool działa w dry-run

**Czas:** 2 dni

---

### day_45: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 10: OPTIMIZATION & CACHING (day_46-50)

### day_46-47: Incremental parsing (2 dni)

**Cel:** Parsować tylko zmienione pliki.

**Wejście:**
- Cache z day_04
- File hashes

**Wyjście:**

**Dzień 46: Hash tracking**
- `pkg/cache/hash.go` - SHA256 tracking:
  - Per plik: zapisz hash w cache
  - Na scan: compare current hash z cached
  - Parse tylko jeśli hash changed
- `pkg/cache/invalidation.go` - cache invalidation strategy:
  - File changed → invalidate file + dependents (graph traversal)
  - Template changed → invalidate all docs używające tego template

**Dzień 47: Incremental scan**
- `cmd/docflow/scan.go` - extend z `--incremental` flag:
  - Load previous cache
  - Hash wszystkie pliki
  - Parse tylko changed files
  - Merge new + unchanged records
  - Save updated cache
- Benchmark: incremental vs full scan (expect 10x speedup dla 1% changed files)

**Kryteria ukończenia:**
- Incremental scan jest wyraźnie szybszy (>5x dla 10% changes)
- Cache invalidation działa poprawnie

**Testy:**
- Scan 300 files full
- Modify 10 files (3%)
- Scan incremental
- Verify: tylko 10 files reparsed, total time <1s

**Czas:** 2 dni

---

### day_48: Template impact analysis

**Cel:** Śledzić wpływ zmian w szablonach.

**Wejście:**
- Template index z day_26
- Document index z usage tracking

**Wyjście:**
- `pkg/templates/impact.go` - impact analysis:
  - Map: template_id → list of docs używających template
  - Przy zmianie template: compute affected docs
- `cmd/docflow/template-impact.go` - komenda:
  - `docflow template-impact --template=X` - lista affected docs
  - `docflow template-impact --all` - dla wszystkich changed templates
- Auto-update (optional):
  - `docflow template-update --template=X --apply` - regenerate affected docs (z backupem)

**Kryteria ukończenia:**
- Impact analysis działa
- Lista affected docs jest poprawna

**Czas:** 1 dzień

---

### day_49-50: Performance optimization (2 dni)

**Cel:** Optymalizacja bottlenecków.

**Wejście:**
- Performance baseline z day_23
- Profiling data

**Wyjście:**

**Dzień 49: Profiling**
- Profiling (CPU, memory) dla każdej operacji
- Identify top 3 bottlenecks
- `LOGS/PROFILING_RESULTS.md` - raport

**Dzień 50: Optimization**
- Implement optimizations:
  - Typical: parallel parsing, cache warming, lazy loading
- Re-benchmark
- `LOGS/OPTIMIZATION_RESULTS.md` - before/after comparison
- Target: 2x improvement w top bottleneck

**Kryteria ukończenia:**
- Min. 2x improvement w przynajmniej 1 bottleneck

**Czas:** 2 dni

---

## TYDZIEŃ 11: ANALYTICS & REPORTING (day_51-55)

### day_51-52: Analytics dashboard (2 dni)

**Cel:** Raporty analityczne.

**Wejście:**
- Usage tracking z day_28
- Quality scores z day_18-19
- Lifecycle data z day_36-37

**Wyjście:**

**Dzień 51: Analytics engine**
- `pkg/analytics/aggregator.go` - agregacje:
  - Template usage trends (top 10 most used)
  - Template quality distribution
  - Document completeness per rodzina
  - Deprecated templates count
  - Average effort per rodzina
- Metrics:
  - Total templates, total docs
  - % published, % draft
  - Average quality score
  - Coverage: % docs with all required sections

**Dzień 52: Reporting**
- `cmd/docflow/analytics.go` - komenda `docflow analytics`:
  - Print summary report (console)
  - Export HTML report (simple table)
  - Export CSV (dla dalszej analizy)
- `pkg/analytics/export.go` - exporters (HTML, CSV)

**Kryteria ukończenia:**
- Analytics compute correctly
- Export formats działają

**Czas:** 2 dni

---

### day_53: Duplicate detection

**Cel:** Wykrywanie duplikatów szablonów.

**Wejście:**
- Template index z day_26

**Wyjście:**
- `pkg/duplicates/detector.go` - duplicate detection:
  - Algorithm: MinHash + Jaccard similarity
  - Threshold: 0.85 (configurable)
  - Parametry w `ALGO_PARAMS.md`:
    ```yaml
    duplicate_detection:
      algorithm: minhash
      num_hashes: 128
      shingle_size: 3
      jaccard_threshold: 0.85
    ```
- `cmd/docflow/find-duplicates.go` - komenda:
  - `docflow find-duplicates --threshold=0.85`
  - Print grupy duplikatów
  - For each group: suggest master (highest quality)

**Kryteria ukończenia:**
- Duplicate detection działa
- Ręcznie zweryfikuj 3 grupy duplikatów - verify trafność

**Czas:** 1 dzień

---

### day_54: Content hints extraction

**Cel:** Wydobywanie wskazówek treściowych.

**Wejście:**
- Markdown parser z day_03
- Templates

**Wyjście:**
- `pkg/content/extractor.go` - content extraction:
  - Code blocks count
  - Tables count
  - Links count (internal/external)
  - Images count
  - Word count per section
- Add metrics do template metadata
- `cmd/docflow/templates.go` - extend list z content metrics

**Kryteria ukończenia:**
- Content metrics compute correctly

**Czas:** 1 dzień

---

### day_55: Phase 2 buffer + retrospective

**Cel:** Buffer + retrospektywa Phase 2.

**Wyjście:**
- Bug fixes
- `LOGS/PHASE2_RETROSPECTIVE.md`
- Adjustments dla Phase 3

**Czas:** 1 dzień

---

### PHASE 2 GO/NO-GO

**Go:**
- Recommender działa z Precision@5 > 0.60
- Daily planner generuje użyteczne plany
- Lifecycle management działa
- Incremental parsing działa
- Tests coverage maintained > 70%
- Performance acceptable (incremental scan < 2s dla 300 files, 10% changes)

**No-Go:**
- Recommender Precision@5 < 0.50 → Investigate data quality, tune algorithm
- Major performance regression → Extend Phase 2 o 1 tydzień (optimization focus)
- Critical bugs → Extend Phase 2, defer some features to Phase 3

**Decision:** Review stakeholder feedback, decide on Phase 3 scope.

---

---

## FAZA 3: POLISH & PRODUCTION READY (day_56 - day_90, ~7 tygodni)

### Cel fazy: Przygotować system do produkcji
### Milestone: v1.0 release, production-ready
### Success criteria: Zero critical bugs, documentation complete, deployment automated

---

## TYDZIEŃ 12: GOVERNANCE & COMPLIANCE (day_56-60)

### day_56-57: Governance rules (2 dni)

**Cel:** Reguły governance dla szablonów i dokumentów.

**Wyjście:**
- `docs/_meta/GOVERNANCE_RULES.yaml` - reguły:
  - Required metadata fields per status
  - Required sections per rodzina
  - Quality thresholds (min quality score dla published)
  - Review requirements (who approves published docs)
- `pkg/validator/governance.go` - governance validator
- `cmd/docflow/validate.go` - extend z `--governance` flag

**Czas:** 2 dni

---

### day_58-59: Compliance reporting (2 dni)

**Cel:** Raporty compliance.

**Wyjście:**
- `pkg/compliance/reporter.go` - compliance checks:
  - % docs passing governance
  - List non-compliant docs
  - Violations breakdown (by rule)
- `cmd/docflow/compliance.go` - komenda `docflow compliance`
- Export raport (HTML, PDF)

**Czas:** 2 dni

---

### day_60: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 13: CI/CD & AUTOMATION (day_61-65)

### day_61-62: CI pipeline (2 dni)

**Cel:** Automatyzacja CI.

**Wyjście:**
- `.github/workflows/ci.yml` (lub GitLab CI):
  - Build
  - Unit tests
  - Integration tests
  - Lint
  - Coverage report
  - Performance benchmarks (fail if regression >20%)
- Pre-commit hooks:
  - Run `docflow validate` na staged docs

**Czas:** 2 dni

---

### day_63-64: Deployment automation (2 dni)

**Cel:** Automatyzacja deployment.

**Wyjście:**
- Release pipeline:
  - Build binaries (3 platforms: Linux, macOS, Windows)
  - Package artifacts
  - Generate CHANGELOG
  - GitHub Release automation
- Installation script:
  - `install.sh` - curl | bash installer
  - Homebrew formula (dla macOS)
  - APT/YUM package (opcjonalnie)
- Update README z installation instructions

**Czas:** 2 dni

---

### day_65: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 14: DOCUMENTATION & EXAMPLES (day_66-70)

### day_66-67: User documentation (2 dni)

**Cel:** Kompletna dokumentacja użytkownika.

**Wyjście:**
- `docs/USER_GUIDE.md` - complete guide:
  - Getting started (5-min tutorial)
  - Concepts (metadata, dependencies, sections, templates)
  - Workflows (common tasks)
  - Advanced usage
- `docs/CLI_REFERENCE.md` - wszystkie komendy z examples
- `docs/TROUBLESHOOTING.md` - FAQs + common errors
- `docs/BEST_PRACTICES.md` - recommendations

**Czas:** 2 dni

---

### day_68: Examples & tutorials

**Cel:** Przykłady i tutoriale.

**Wyjście:**
- `examples/` - 3-5 przykładowych projektów:
  - `examples/simple-api/` - API documentation project
  - `examples/architecture/` - Architecture docs project
  - `examples/knowledge-base/` - General knowledge base
- Każdy przykład z README + sample docs + expected output

**Czas:** 1 dzień

---

### day_69: Video tutorial

**Cel:** Video walkthrough.

**Wyjście:**
- 15-min screencast:
  - Installation
  - Setup new project
  - Scan, validate, generate
  - Fix validation errors
  - Use recommender
  - Generate daily plan
- Publish na YouTube (opcjonalnie)

**Czas:** 1 dzień

---

### day_70: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 15: TESTING & HARDENING (day_71-75)

### day_71-72: End-to-end testing (2 dni)

**Cel:** Comprehensive E2E tests.

**Wyjście:**
- `tests/e2e/` - E2E test suite:
  - Full workflow: setup project → scan → validate → generate → plan
  - Error scenarios: invalid metadata, cycles, missing deps
  - Performance: 1000 docs (synthetic)
- Run E2E na 3 platformach (Linux, macOS, Windows)

**Czas:** 2 dni

---

### day_73: Security audit

**Cel:** Security review.

**Wyjście:**
- Security checklist:
  - Input validation (CLI args, file paths)
  - Path traversal prevention
  - YAML bomb protection (large/nested YAML)
  - Dependency vulnerabilities scan
- Fixes dla security issues
- `SECURITY.md` - security policy

**Czas:** 1 dzień

---

### day_74: Chaos testing

**Cel:** Test resilience.

**Wyjście:**
- Chaos scenarios:
  - Corrupted cache
  - Huge files (100MB MD file)
  - Deep directory structure (100 levels)
  - Symlink loops
  - Concurrent access (2 processes)
- Fixes dla discovered issues

**Czas:** 1 dzień

---

### day_75: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 16: PERFORMANCE & SCALE (day_76-80)

### day_76-77: Scale testing (2 dni)

**Cel:** Test na dużej skali.

**Wyjście:**
- Scale tests:
  - 1000 documents
  - 5000 documents (synthetic)
  - 10000 documents (if feasible)
- Measure:
  - Scan time
  - Memory usage
  - Graph build time
  - Validation time
- `LOGS/SCALE_TEST_RESULTS.md` - raport
- Identify limits (max practical scale)

**Czas:** 2 dni

---

### day_78-79: Final optimization (2 dni)

**Cel:** Last optimizations.

**Wyjście:**
- Address bottlenecks z scale testing
- Memory optimizations (streaming, lazy loading)
- Re-benchmark
- `LOGS/FINAL_PERFORMANCE.md` - final numbers

**Czas:** 2 dni

---

### day_80: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 17: RELEASE PREPARATION (day_81-85)

### day_81: Release candidate

**Cel:** Zbudować RC.

**Wyjście:**
- Tag v1.0.0-rc1
- Build binaries
- Internal testing (wszystkie komendy na przykładach)
- Bug bash (zespół testuje przez 1 dzień)

**Czas:** 1 dzień

---

### day_82-83: Bug fixes (2 dni)

**Cel:** Fix bugs z RC testing.

**Wyjście:**
- Fix critical bugs
- Tag v1.0.0-rc2 (if needed)
- Retest

**Czas:** 2 dni

---

### day_84: Final documentation review

**Cel:** Review dokumentacji.

**Wyjście:**
- Wszystkie docs reviewed
- Corrections applied
- CHANGELOG finalized
- Release notes written

**Czas:** 1 dzień

---

### day_85: Buffer

**Czas:** 1 dzień

---

## TYDZIEŃ 18: LAUNCH (day_86-90)

### day_86: Pre-launch checklist

**Cel:** Final checklist.

**Wyjście:**
- `LOGS/LAUNCH_CHECKLIST.md`:
  - [ ] All tests pass
  - [ ] Documentation complete
  - [ ] Examples tested
  - [ ] Binaries built (3 platforms)
  - [ ] GitHub release prepared
  - [ ] Installation methods tested
  - [ ] Security audit done
  - [ ] Performance acceptable
  - [ ] No known critical bugs
- Verify checklist

**Czas:** 1 dzień

---

### day_87: v1.0 release

**Cel:** Launch v1.0.

**Wyjście:**
- Tag v1.0.0
- GitHub Release z:
  - Binaries (Linux, macOS, Windows)
  - CHANGELOG
  - Release notes
  - Installation instructions
- Announce (blog post, social media, mailing list)

**Czas:** 1 dzień

---

### day_88: Post-launch monitoring

**Cel:** Monitor adoption.

**Wyjście:**
- Setup monitoring:
  - GitHub stars, forks, issues
  - Installation stats (if telemetry enabled)
- Respond to issues
- Fix urgent bugs (hotfix release if needed)

**Czas:** 1 dzień

---

### day_89: Retrospective

**Cel:** Project retrospective.

**Wyjście:**
- `LOGS/PROJECT_RETROSPECTIVE.md`:
  - What went well
  - What went wrong
  - Lessons learned
  - Recommendations for future
- Team debrief meeting

**Czas:** 1 dzień

---

### day_90: Backlog & roadmap

**Cel:** Plan future iterations.

**Wyjście:**
- `BACKLOG.md` - prioritized backlog:
  - P0: Critical fixes from users
  - P1: High-value features
  - P2: Nice-to-haves
- `ROADMAP.md` - roadmap v1.1, v1.2, v2.0:
  - v1.1 (3 months): Quick wins, polish
  - v1.2 (6 months): Advanced features (AI suggestions, auto-fix)
  - v2.0 (12 months): Major features (web UI, collaboration)
- Next iteration plan (if continuing)

**Czas:** 1 dzień

---

### PHASE 3 GO/NO-GO

**Go:**
- v1.0 launched
- All critical bugs fixed
- Documentation complete
- Tests coverage > 70%
- Performance meets targets
- Zero known security issues

**No-Go (delay launch):**
- Critical bugs unresolved → Extend Phase 3
- Documentation incomplete → Extend 3 dni
- Performance unacceptable → Rollback some features, launch v1.0-lite

**Decision:** If No-Go, re-assess launch date (max +2 weeks extension).

---

---

## SUMMARY: EXTENDED PLAN

**Total duration:** 90 dni robocze (~18 tygodni, ~4.5 miesiąca)

**Fazy:**
- **Pre-work:** 5-7 dni - Fundamenty projektu
- **Phase 1 (Foundation & MVP):** 25 dni - Core features
- **Phase 2 (Intelligence & Automation):** 30 dni - Advanced features
- **Phase 3 (Polish & Production):** 35 dni - Production readiness

**Kluczowe milestones:**
- Pre-work complete: Korpus danych, architektura, tech stack
- day_15: MVP release (v0.1.0)
- day_25: Phase 1 complete
- day_55: Phase 2 complete (recommendations, planner, lifecycle)
- day_90: v1.0 production release

**Bufory:**
- 13 dni buforowych wbudowanych w plan (day_05, 20, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85)
- ~14% czasu na nieprzewidziane problemy

**Team assumptions:**
- 1-2 senior developers full-time
- Code reviews wbudowane w bufory
- Stakeholder reviews przed każdym Phase Go/No-Go

**Risk mitigation:**
- Early MVP (day_15) - fail fast
- Iteracyjne releases (3 fazy)
- Dogfooding (używamy narzędzia na sobie)
- Frequent testing (unit + integration + E2E)
- Go/No-Go gates co fazę

---

## RÓŻNICE OD ORYGINALNEGO PLANU (35 dni)

**Dodane:**
1. **Pre-work (5-7 dni):** Analiza wymagań, pozyskanie danych, architektura
2. **Bufory (13 dni):** Code review, tech debt, nieprzewidziane problemy
3. **MVP milestone (day_15):** Early value delivery
4. **Testing throughout:** Test strategy day_21 (wcześniej day_16)
5. **Rozdzielone równoległe wątki:** day_01-02 zamiast day_01 (2 wątki)
6. **Security audit (day_73):** Production readiness
7. **Chaos testing (day_74):** Resilience
8. **Scale testing (day_76-77):** Performance na dużej skali
9. **Examples & tutorials (day_68-69):** User adoption
10. **Post-launch (day_88-90):** Monitoring, retrospective, roadmap

**Usunięte/Zmodyfikowane:**
- Nie ma "sztucznego" równoległego prowadzenia 2 wątków w 1 dzień
- Funkcje low-ROI odroczone lub usunięte (analytics HTML export uproszczony)
- Research (day_02) skrócony - wykorzystujemy korpus z pre-work

**Czas realokowany:**
- Original: 35 dni
- Extended: 90 dni
- Ratio: 2.57x (zgodne z oceną audytu: 65-80 dni)

**Execution model:**
- Waterfall w każdej fazie, ale fazy są iteracyjne
- Możliwość early exit po Phase 1 (MVP) jeśli projekt nie ma sensu
- Phase 2-3 opcjonalne jeśli MVP wystarczy

---
