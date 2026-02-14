# docflow

CLI do zarządzania dokumentacją i szablonami (Go, bez CGO). Release Candidate 2 (GA-ready).

## Wymagania
- Go 1.25+ (zainstalowane lokalnie)
- Build ze źródeł (`./cmd/docflow`) jest podstawową ścieżką instalacji
- Tryb offline działa przy aktualnym `vendor/`; gdy `vendor/` nie jest zsynchronizowany, użyj `-mod=mod`

## Quickstart (60s)
```bash
# zbuduj lokalnie binarkę
mkdir -p build
go build -mod=mod -o build/docflow ./cmd/docflow

# szybki audit lokalny (nieblokujący), tworzy bundle w .docflow/out/<run_id>/
./build/docflow health \
  --config examples/simple-api/docflow.yaml \
  --rules docs/_meta/GOVERNANCE_RULES.yaml \
  --baseline-mode repo --baseline-dir .docflow/baseline \
  --bundle-dir .docflow/out

# (opcjonalnie) profiling
./build/docflow scan --cpu-profile /tmp/docflow.cpu examples/simple-api
go tool pprof -top /tmp/docflow.cpu
```

## Maintainer workflows (artefacts/logs)
Wewnętrzne bundlowanie release/PR oraz mapa evidence są opisane w `docs/MAINTAINER_NOTES.md`.

## Quickstart dla Twojego repo (zewnętrznego)
1) Przygotuj pliki wejściowe:
- `docflow.yaml` w katalogu repo (albo wskaż `--config`)
- `docs/_meta/GOVERNANCE_RULES.yaml` (albo wskaż `--rules`)
- baseline w repo: `.docflow/baseline/validate.json` i `.docflow/baseline/compliance.json` (zalecane)

2) Uruchom lokalny audit (nieblokujący) i wygeneruj bundle:
```bash
docflow health \
  --config docflow.yaml \
  --rules docs/_meta/GOVERNANCE_RULES.yaml \
  --baseline-mode repo \
  --baseline-dir .docflow/baseline \
  --bundle-dir .docflow/out
```
Wyniki znajdziesz w `.docflow/out/<run_id>/` (summary + validate/compliance + SARIF).

3) W CI użyj kanonicznego flow (blokuje tylko nowe problemy):
- gotowy workflow: `.github/workflows/docflow.yml`
- lub ręcznie:
```bash
docflow health --ci \
  --config docflow.yaml \
  --rules docs/_meta/GOVERNANCE_RULES.yaml \
  --baseline-mode repo \
  --baseline-dir .docflow/baseline \
  --bundle-dir .docflow/out
```

4) Aktualizacja baseline:
Baseline w repo aktualizuj wyłącznie przez dedykowany PR do `main`.
Szczegóły: `docs/GA_CHECKLIST.md` (Polityka baseline).

Opcjonalnie (advanced): maintainerskie skrypty są opisane w `docs/MAINTAINER_NOTES.md`.

### GitHub Actions (recommended: health --ci)
Minimalny copy/paste dla CI:
```bash
./build/docflow --config docflow.yaml health --ci \
  --rules docs/_meta/GOVERNANCE_RULES.yaml \
  --baseline-mode repo \
  --baseline-dir .docflow/baseline \
  --bundle-dir .docflow/out
```

Po runie powstaje bundle w `.docflow/out/<run_id>/`:
- `validate.sarif`
- `validate.json`
- `compliance.json`
- `summary.json`
- `summary.md`
- `meta.json`

Interpretacja `overall_exit`:
- `0` — brak nowych problemów (PASS)
- `1` — nowe błędy/naruszenia domenowe (FAIL)
- `2` — błąd użycia CLI (flags/args)
- `3` — błąd runtime (IO/config/rules)

Przykład workflow (jedna komenda + upload SARIF i bundle):
```yaml
- name: Docflow health (CI bundle)
  id: health
  run: |
    set +e
    ./build/docflow --config docflow.yaml health --ci \
      --rules docs/_meta/GOVERNANCE_RULES.yaml \
      --baseline-mode repo \
      --baseline-dir .docflow/baseline \
      --bundle-dir .docflow/out
    echo "health_exit=$?" >> "$GITHUB_OUTPUT"
    echo "run_dir=$(ls -1dt .docflow/out/* | head -n 1)" >> "$GITHUB_OUTPUT"
    exit 0

- name: Upload SARIF (validate)
  if: ${{ always() }}
  uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: ${{ steps.health.outputs.run_dir }}/validate.sarif
    category: docflow-validate

- name: Upload health bundle
  if: ${{ always() }}
  uses: actions/upload-artifact@v4
  with:
    name: docflow-bundle
    path: ${{ steps.health.outputs.run_dir }}

- name: Enforce health exit code
  if: ${{ steps.health.outputs.health_exit != '0' }}
  run: exit "${{ steps.health.outputs.health_exit }}"
```
Gotowy workflow w repo: `.github/workflows/docflow.yml`.
Upewnij się, że `docflow.yaml`, `docs/_meta/GOVERNANCE_RULES.yaml` i baseline (`.docflow/baseline/*.json`) są dostępne w runie CI.

## Command palette
- Bash/fzf: `./scripts/palette.sh` (fzf opcjonalny; fallback lista + wybór numeru). Akcje: import, audit, bundle, demo, checks.
- W binarce: `./build/docflow palette` (tekstowe menu bez fzf).

Przykłady `examples/simple-api` i `examples/architecture` są governance-ready (validate/compliance PASS).

## Instalacja ze źródeł
```bash
mkdir -p build

# preferowany fallback (działa gdy vendor jest niespójny)
go build -mod=mod -o build/docflow ./cmd/docflow

# opcjonalnie: build offline, gdy vendor jest zsynchronizowany
GOFLAGS=-mod=vendor go build -o build/docflow ./cmd/docflow
```

## Uruchomienie przykładowe
- Wersja binarki: `./build/docflow --version`
- Skan: `./build/docflow scan -o .docflow/cache/doc_index.json --deterministic`
- Walidacja: `./build/docflow validate --strict`
- Plan dzienny: `./build/docflow plan daily --max 5`
- Rekomendacje (demo): `./build/docflow recommend --doc-type guide --lang pl`
- Impact szablonów: `./build/docflow template-impact --old-index .docflow/cache/doc_index.json`

## Komendy CLI (overview)
- `scan` — buduje indeks dokumentów (metadane, checksum, hints).
- `validate` — sprawdza metadane, doc_id, context_sources, expected deps; flagi `--strict`/`--warn`; opcjonalnie `--old-index FILE --auto-bump --save-index` do śledzenia zmian i podbijania wersji.
- `validate --format json|sarif` oraz `compliance --format json` — stabilne formaty maszynowe do CI.
- `validate|compliance --against BASELINE --fail-on new --show new` — tryb regresji (blokuj tylko nowe problemy).
- `validate --status-aware` — progressive validation zależna od statusu (draft/review/published).
- `plan daily` — kolejność topo + effort heurystyczny; `--max N`.
- `recommend` — demo rekomendacji szablonów (doc_type/lang/quality/usage).
- `template-sets` — demo współwystępowania szablonów.
- `templates list|deprecated|deprecated-report` — demo statusów/wersji/deprecjacji.
- `template-impact --old-index FILE` — wykrywa dokumenty zależne od zmienionych szablonów.
- `stats --schema file.yaml` — metryki kompletności sekcji (grupowanie po doc_type lub status).
- `migrate-sections --apply` — zamienia legacy nazwy sekcji wg `section_aliases` w configu (domyślnie dry-run, alias `--dry-run`).
- `validate --governance docs/_meta/GOVERNANCE_RULES.yaml` — egzekwuje reguły governance (pola i sekcje per status/rodzina).
- `--log-format json`, `--cpu-profile`, `--mem-profile` — globalne flagi obserwowalności.

## Kontrakt i GA (Source of Truth)
- Kontrakt i integracje CI: `docs/CONTRACT.md`
- GA Checklist: `docs/GA_CHECKLIST.md`
- Gotowy workflow: `.github/workflows/docflow.yml`

README opisuje szybkie użycie; szczegóły kontraktu i polityk są utrzymywane w `docs/`.

## Legacy CI (manual flow)
Jeżeli potrzebujesz uruchamiać kroki ręcznie (bez `health --ci`), użyj bezpośrednio komend `validate` i `compliance`.
Kanoniczny sposób integracji CI pozostaje `docflow health --ci`, bo daje jeden spójny `overall_exit` i pełny bundle artefaktów.
Szczegóły parametrów i kontraktu outputów: `docs/CLI_REFERENCE.md` oraz `docs/CONTRACT.md`.

## Instalacja (release)
- Releases są opcjonalne. Jeśli `.../releases/latest/download/...` zwraca `404`, użyj instalacji ze źródeł (sekcja wyżej).
- Ręcznie (gdy release istnieje): pobierz archiwum z Releases (`docflow-<os>-<arch>.tar.gz` lub `.zip`), rozpakuj i uruchom.
- Weryfikacja: `sha256sum -c checksums.txt` (plik `checksums.txt` z tej samej strony release).
- Skrypt instalacyjny: `PREFIX=$HOME/.local ./install.sh` (domyślnie instaluje lokalny build z `./build/docflow`). `--dry-run` aby sprawdzić bez kopiowania.
  - Możesz wskazać własną binarkę: `--from /path/to/docflow`.
  - `--channel latest` używa lokalnych artefaktów `dist/` przygotowanych przez pipeline release.

## Verify release
Release zawiera archiwa (`.tar.gz` dla Linux/macOS, `.zip` dla Windows) oraz:
- `checksums.txt` (SHA-256 wszystkich archiwów)
- opcjonalnie podpisy Cosign: `*.sig` i `*.cert`
- opcjonalnie SBOM: `sbom.cdx.json`
- opcjonalnie provenance attestation (GitHub Attestations)
- pełna checklista: `docs/SECURITY_VERIFICATION.md`

### 1) Verify checksums (baseline)
Linux:
```bash
sha256sum -c checksums.txt
```

macOS (fallback, gdy brak `sha256sum`):
```bash
shasum -a 256 -c checksums.txt
```

Windows PowerShell (single-file):
```powershell
$zip="docflow-windows-amd64.zip"
$expected = (Select-String -Path .\checksums.txt -Pattern $zip).ToString().Split(" ")[0]
$actual = (Get-FileHash .\$zip -Algorithm SHA256).Hash.ToLower()
if ($expected.ToLower() -ne $actual) { throw "SHA256 mismatch" } else { "OK" }
```

### 2) Verify Cosign signatures (optional)
Przykład dla archiwum:
```bash
cosign verify-blob \
  --certificate docflow-linux-amd64.tar.gz.cert \
  --signature   docflow-linux-amd64.tar.gz.sig \
  docflow-linux-amd64.tar.gz
```

I dla `checksums.txt`:
```bash
cosign verify-blob \
  --certificate checksums.txt.cert \
  --signature   checksums.txt.sig \
  checksums.txt
```

## Troubleshooting
- **Brak internetu / vendor**: ustaw `GOFLAGS=-mod=vendor`; Go 1.25+ wymagane. Jeśli CI ma dostęp do sieci tylko do pobrania toolchaina, vendor wystarcza do build/test.
- **`inconsistent vendoring`**: zsynchronizuj vendor `go mod vendor` albo użyj fallback `-mod=mod`.
- **Permission denied w $HOME/.cache/go-build**: ustaw `GOCACHE=/tmp/go-cache` (jak w poleceniu testów) lub `HOME=/tmp` / `XDG_CACHE_HOME=/tmp`.
- **Brak template_source**: `template-impact` zwraca puste wyniki — uzupełnij `template_source` w frontmatter dokumentów generowanych z szablonów.
- **Wysoki czas skanu**: uruchom po raz drugi (checksum cache); rozważ wyłączenie hints jeśli niepotrzebne.

## Roadmapa v0.2 (skrót)
- Podpiąć realny indeks/status/usage do recommend/template-sets/templates.
- Dodać zapis usage w generatorze i konsumować w rekomendacjach.
- Flagi CLI dla StrictMode/PublishStrict; migracja sekcji (dry-run).
- Raport brakujących przykładów w CLI; równoległy scan/parse z zachowaniem deterministyczności.

## Testy
```bash
go test -mod=mod ./internal/... ./pkg/... ./tests/...

# opcjonalnie: offline, gdy vendor jest zsynchronizowany
GOFLAGS=-mod=vendor go test ./internal/... ./pkg/... ./tests/...
```

## Konfiguracja
- Plik `docflow.yaml` (domyślny) — ścieżki, cache, dependency/promote_context_for, family_rules_path.

## Status
- GA (na bazie RC2): queue Go/No-Go (text/json, cache, workers), governance-ready examples, perf 10k chain ~0.56s/0.9s (RSS ~80MB), observability flags.
- Test strategy: `docs/_meta/TEST_STRATEGY.md`.

## Support / Disclaimer
- Oprogramowanie jest dostarczane w modelu "AS IS", bez gwarancji.
- Brak SLA i gwarantowanego wsparcia; Issues/PR mogą pozostać bez odpowiedzi.
- Release'y są publikowane wyłącznie w trybie best-effort (bez gwarancji częstotliwości).
- Używasz na własne ryzyko.
