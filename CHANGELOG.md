# Dziennik zmian

Wszystkie istotne zmiany w projekcie są dokumentowane tutaj.
Format oparty na [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [1.0.0] — 2026-03-13

Pierwsza stabilna wersja. Kontrakt CLI i format JSON wyjść jest zamrożony dla tej serii.

### Dodano
- **`docflow health`** — orkiestrator CI/CD: validate + compliance + bundle w jednym wywołaniu. Generuje
  `validate.json`, `validate.sarif`, `compliance.json`, `summary.json`, `summary.md`, `meta.json`.
- **`docflow check`** — skrót dla `health` z domyślnymi parametrami (config auto-wykrywany, zero flag).
- **`docflow doctor`** — diagnostyka środowiska: config, rules, baseline, git, vendor, wersja binarna.
- **`docflow fix`** — autofix bezpiecznych klas problemów walidacji (legacy section names, puste sekcje).
- **`docflow baseline`** — zarządzanie baseline'em (save/compare/migrate).
- **`docflow changes`** — raport nowych/usuniętych problemów względem baseline (text/json).
- **`docflow graph`** — wizualizacja grafu zależności dokumentów (DOT/JSON).
- **`docflow compliance`** — reguły governance z GOVERNANCE_RULES.yaml; raport JSON + exit codes.
- **`docflow validate`** — walidacja metadanych z formatami JSON, SARIF 2.1.0, text.
- **`docflow scan`** — pełne skanowanie i indeksowanie katalogu dokumentów.
- **`docflow stats`** — statystyki projektu (rozkład typów, statusów, metryk sekcji).
- **`docflow analyze-patterns`** — analiza wzorców co-occurrence w zależnościach.
- **`docflow find-duplicates`** — wykrywanie duplikatów doc_id w całym drzewie.
- **`docflow recommend`** — rekomendacje szablonów (demo, dane statyczne).
- **`docflow templates`** / **`template-sets`** / **`template-impact`** — zarządzanie szablonami.
- **`docflow migrate-sections`** — migracja legacy nazw sekcji do aktualnego schematu.
- **`docflow init --scaffold`** — tworzy szkielet projektu: `docflow.yaml`, `docs/_meta/GOVERNANCE_RULES.yaml`, przykładowy dokument.
- **Execution layer cache** — deterministyczny cache wyników skanowania; `--no-cache` / `--cache-dir` / `--changed-only` / `--since`.
- **Obserwabilność** — `--cpu-profile`, `--mem-profile`, `--log-format json`, `--log-level`.
- **Kolory ANSI** w terminalu (respektuje `NO_COLOR`, `TERM=dumb`, nie-TTY).
- **Spinner** podczas długich operacji health/scan.
- **Auto-discovery konfiguracji** — `docflow.yaml` wyszukiwany w górę drzewa katalogów (jak git).
- **SARIF 2.1.0** — pełna obsługa w `validate`; properties: `doc_id`, `type`, `identity_version`.
- **GitHub Actions workflow** — gotowy `.github/workflows/docflow.yml` dla CI.
- **Release pipeline** — vendor_guard, cross-platform build, SBOM (CycloneDX), SLSA provenance, cosign signing, checksums.txt.
- **Testy bezpieczeństwa** — chaos: YAML bomb, symlink loop, path traversal (blokada), corrupted cache.
- **Testy wydajności** — 1k dokumentów: index ~100ms, validate ~150ms, RSS <80MB.

### Zmieniono
- **Identity v2** (breaking dla baseline'ów wygenerowanych identity v1) — patrz sekcja [Migracja baseline](#migracja-baseline-identity-v1--v2) poniżej.
- **Go 1.25.8** — zaktualizowano toolchain; naprawiono CVE GO-2026-4602.
- Uprawnienia plików: `0o750` dla katalogów, `0o600` dla plików (poprzednio `0o755`/`0o644`).
- Wszystkie `fmt.Fprintf` zwracają wartości (errcheck); `defer f.Close()` ignoruje błąd explicite.

### Naprawiono
- `scan <dir>` respektuje podaną ścieżkę (poprzednio ignorował argument).
- Path traversal w `depends_on` i `context_sources` jest blokowany w walidatorze.
- Governance reporter używa domyślnego statusu gdy dokument nie ma pola `status`.
- Cykle w grafie zależności są wykrywane i raportowane jako issue `cycle_detected`.
- Queue Go/No-Go działa poprawnie z cache i wieloma workerami.

### Bezpieczeństwo
- Naprawiono ~30 ostrzeżeń gosec (uprawnienia plików, path traversal, subprocess).
- gitleaks: 0 prawdziwych sekretów w codebase (18 false positives w vendor SQLite — hex stałe).
- trivy: 0 CVE w zależnościach projektu.

---

## Migracja baseline: identity v1 → v2

**Czego dotyczy:** baseline'y wygenerowane przez wersje wcześniejsze niż 1.0.0 używają `identity_version: "1"`.
Od v1.0.0 domyślna wersja to `"2"`. Baseline v1 NIE jest kompatybilny z porównywaniem v2.

**Jak sprawdzić, czy masz baseline v1:**
```bash
cat .docflow/baseline/validate.json | grep identity_version
# jeśli brak pola lub "1" — potrzebujesz migracji
```

**Jak zmigrować:**
```bash
# validate baseline
docflow baseline migrate \
  --in  .docflow/baseline/validate.json \
  --out .docflow/baseline/validate.json \
  --kind validate

# compliance baseline
docflow baseline migrate \
  --in  .docflow/baseline/compliance.json \
  --out .docflow/baseline/compliance.json \
  --kind compliance
```

**Co się zmienia w v2:**
- identity każdego issue oparty jest na deterministycznym hash pól `details` (SHA-256 pierwsze 8 bajtów)
  zamiast treści `message`. Oznacza to, że zmiany w tłumaczeniu/literówki w treści komunikatu
  **nie tworzą nowego issue** w baseline — porównanie jest stabilniejsze.
- Pole `location` używa formatu `L<line>:C<column>` zamiast tylko `L<line>`.

**Backup zalecany:**
```bash
cp -r .docflow/baseline .docflow/baseline.v1.bak
```

---

## Wydajność (v1.0.0)

Zmierzone na 8-core laptop, Go 1.25.8, `-mod=vendor`:

| Operacja | 100 plików | 1 000 plików | 10 000 plików |
|----------|-----------|-------------|--------------|
| `scan` (index) | <20ms | ~100ms | ~560ms |
| `validate` | <30ms | ~150ms | ~900ms |
| `health` (pełny bundle) | <100ms | ~300ms | ~1,5s |
| Zużycie RAM (RSS) | <20MB | <80MB | <300MB |

Limity akceptowalne (P99): 1 000 plików <12s, RAM <300MB.
Wydajność >5 000 plików poza zakresem v1.0.0 (planowane w v1.1+).

---

## [v1.0.0-rc2] — 2026-02-09

- Dodano: testy wydajności 1k/duże pliki i raporty RSS; sanity chaos/security (YAML bomb, symlinks, traversal blokada).
- Naprawiono: `scan <dir>` respektuje ścieżkę; walidator blokuje path traversal; governance reporter domyślny status; cykle wykrywane; queue Go/No-Go (text/json, cache, workers).
- Dodano: flagi obserwabilności (`--cpu-profile`, `--mem-profile`, `--log-format json`); queue workflow z cache; governance examples; perf 10k (index~0.56s, validate~0.9s, RSS~80MB).
- Testy/QA: e2e (1k + large + 10k), chaos/security, bug bash RC1 (P0=0, P1 naprawiony), smoke RC2.

## [v0.1.0-mvp] — 2026-02-09

- Walidator metadanych (doc_id, context_sources, expected deps, status-aware).
- SectionTree parser + Section metrics (completeness).
- Generator manifestu i prosty generator szablonów z preview.
- Rekomendacje szablonów (MVP), zestawy (co-occurrence), plan dzienny, template impact.
- Config: `promote_context_for`, `family_rules_path`.
- CLI: validate, plan daily, recommend (demo), template-sets (demo), templates list/deprecated, template-impact.
- Testy `GOFLAGS=-mod=vendor go test ./...` (Go 1.25.7).
