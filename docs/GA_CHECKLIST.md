# GA Checklist (docflow)

Ten dokument definiuje kryteria dopuszczenia wydania (GA) oraz zasady wersjonowania kontraktu CLI/output.

## 1. Kryteria blokujące (MUST)

### 1.1 Exit codes (stabilne i udokumentowane)
- 0: OK / tryb warn (nieblokujące)
- 1: naruszenia domenowe w trybie strict (validate/compliance)
- 2: błąd użycia CLI (brak wymaganych flag, konflikt flag, zły argument)
- 3: błąd runtime (IO, parse config/rules, internal error)

Wymagane:
- Konflikt `--strict` + `--warn` => exit 2
- Brak wymaganej flagi (np. `--rules`) => exit 2
- Brak pliku rules/config => exit 3

### 1.2 Kontrakt maszynowy JSON
Validate (`--format json`):
- `schema_version`
- `identity_version` (dla comparator v2)
- metryki: `files`, `documents`, `error_count`, `warn_count`
- `issues[]` zawiera minimum:
  - `code`, `level`, `type`, `path`
  - `doc_id` (jeśli dostępny)
  - `message`
  - `line` (jeśli dotyczy)
  - `details` (stabilne pola identity)

Compliance (`--format json`):
- `schema_version`
- `identity_version` (dla comparator v2)
- metryki: `documents`, `passed`, `failed`, `pass_rate`, `violations_count`
- `docs[]` z `path`, `doc_id`, `doc_type`, `status`, `violations[]`
- metadane reguł: `rules_path`, `rules_checksum`
- diagnostyka: `duplicate_doc_ids` (jeśli występują)

### 1.3 JSON error envelope (usage/runtime)
Jeżeli użytkownik wybierze `--format json` i wystąpi błąd użycia lub runtime:
- stdout zawiera JSON envelope:
  - `schema_version`
  - `error.kind` = `usage` lub `runtime`
  - `error.code`, `error.message`
  - `error.details` (jeśli dostępne)
- exit code = 2 dla usage, 3 dla runtime

### 1.4 Baseline/regresje
Dla validate i compliance:
- `--against <baseline.json>` jest wymagane, gdy `--fail-on new` lub `--show new/existing`
- `--fail-on new` blokuje tylko na NOWYCH problemach
- `--show new` filtruje output do samych nowych wyników

Wymagane:
- `--fail-on new` bez `--against` => exit 2 + JSON envelope (dla format json)
- baseline identity v2 nie może zależeć od `message` (wording-only changes nie tworzą „new”)

### 1.5 SARIF (validate)
- `validate --format sarif` generuje SARIF 2.1.0
- `ruleId` mapuje się z `issue.code`
- `level` mapuje się z `issue.level` (error/warn -> error/warning)
- lokalizacja: plik + linia (jeśli dostępne)
- `properties` zawiera minimum: `doc_id`, `type`
- SARIF respektuje `--against/--fail-on/--show`
- deterministyczność: dwa runy na tym samym wejściu => identyczny plik (cmp OK)

### 1.6 CI / Release
- `go test ./...` i `go vet ./...` przechodzi w trybie `-mod=vendor`
- Release publikuje archiwa + `checksums.txt` (basename) i opcjonalnie podpisy (cosign)
- README zawiera “Verify release” (sha256 + opcjonalnie cosign)
- `docflow health --ci` generuje pełny bundle (`validate.json`, `validate.sarif`, `compliance.json`, `summary.json`, `summary.md`, `meta.json`)
- `health --ci` zwraca jeden spójny exit code (`overall_exit`) zgodny z priorytetem `3 > 2 > 1 > 0`

## 2. Kryteria nieblokujące (SHOULD)
- SARIF dla compliance (jeśli będzie miało sens w Twoim modelu lokalizacji)
- JUnit/Markdown summary jako opcjonalne formaty raportu
- Autofix diff-only + apply dla bezpiecznych klas problemów
- Incremental scan/cache
- Provenance attestation i SBOM jako assets (best-effort)

Status:
- Autofix (część 5/6) jest dostępny dla bezpiecznych klas: `missing_field(version)` i `legacy_section_name`, z trybem patch (domyślnie) oraz atomowym `--apply`.

## 3. Wersjonowanie kontraktu (schema_version)

### 3.1 Zasady
- `schema_version` traktuj jako SemVer (np. `1.0`, `1.1`, `2.0`)
- MINOR (kompatybilne wstecz):
  - dodanie nowych pól
  - dodanie nowych kodów issue/violation
  - dodanie nowych formatów output
- MAJOR (breaking):
  - usunięcie/zmiana znaczenia pola
  - zmiana struktury `issues[]` / `docs[]`
  - zmiana definicji tożsamości baseline (identity key)

### 3.2 Zasady kompatybilności
- Parserzy powinni ignorować nieznane pola
- Baseline powinien być używany w obrębie tej samej wersji MAJOR `schema_version`
- Dla migracji identity: `docflow baseline migrate --in ... --out ... --kind validate|compliance`

## 4. Procedura Go/No-Go (przed tagiem GA)
1) compliance bez `--rules` w strict (format json):
   - exit=2 + JSON envelope
2) validate konflikt `--strict --warn` (format json):
   - exit=2 + JSON envelope
3) SARIF deterministyczny:
   - 2 runy => cmp OK
4) Baseline validate “no new”:
   - exit=0, `issues[]` puste przy `--show new`
5) Baseline validate “one new issue”:
   - exit=1, `new_error_count=1`, właściwy `issue.code`
6) Snapshot sanity:
   - entrypoint używa mapowania exit codes
   - release changelog ma poprawny zakres `prev..HEAD`
7) Health bundle sanity:
   - `docflow health --ci` tworzy komplet artefaktów i poprawny `overall_exit` w `summary.json`

## 5. Polityka baseline (repo vs artifact)
- Baseline w repo (`.docflow/baseline/*.json`) aktualizujemy tylko na `main` przez dedykowany PR.
- Baseline nie powinien być automatycznie aktualizowany w PR feature branch (to ukrywa nowy dług).
- Baseline repo jest domyślny: stabilny, offline-friendly i audytowalny w diffie.
- Baseline artifact jest opcjonalny: CI-first, ale ephemeryczny (retencja artifacts, zależność od runów).
- Jeśli używasz trybu artifact, zdefiniuj politykę retencji i źródło „stałego” baseline (np. release asset).
