# MAINTAINER_NOTES

Ten dokument opisuje utrzymanie repo i wydania `doc-flow` w trybie best-effort.
Nie jest to umowa wsparcia. To są procedury techniczne, aby utrzymać spójność kontraktu i pipeline.

## 1. Polityka projektu (krótko)

- Licencja: MIT (AS IS, bez gwarancji).
- Brak SLA i brak zobowiązań utrzymaniowych (zob. `README.md`, `CONTRIBUTING.md`, `SECURITY.md`).
- Issues/Discussions mogą być wyłączone (preferowane dla "zero presji").
- Model zależności: vendor-first (Model 2): `vendor/` jest trackowany, build/test/release używa `-mod=vendor`.

## 2. Źródła prawdy

- Kontrakt CLI/output (JSON/SARIF/envelope/exit codes): `docs/CONTRACT.md`
- GA / gate / baseline / identity v2: `docs/GA_CHECKLIST.md`
- CI workflow (kanoniczny): `.github/workflows/docflow.yml` (health --ci + bundle + SARIF)
- Release workflow: `.github/workflows/release.yml`
- Build workflow (reusable): `.github/workflows/build.yml`
- CI: `.github/workflows/ci.yml`

## 3. Model zależności: vendor-first (offline build/test)

### 3.1. Zasada

- Wszystkie buildy/testy używają `-mod=vendor`.
- `vendor/` jest trackowany i musi być zsynchronizowany z `go.mod` / `go.sum`.

### 3.2. Jak aktualizować zależności (jedyna poprawna ścieżka)

Wymaga sieci tylko w momencie zmiany zależności.

```bash
make vendor-sync
# równoważne: go mod tidy && go mod vendor && go mod verify
```

### 3.3. Twardy gate "vendor up-to-date"

W CI/build/release istnieje bramka, która failuje, jeśli:

- `go mod vendor` zmienia `vendor/`, `go.mod` lub `go.sum`.

To eliminuje klasę problemów "inconsistent vendoring".

## 4. Lokalny "pre-release audit" (przed tagiem)

Uruchom w katalogu repo:

```bash
git checkout main
git pull --ff-only
git status --porcelain=v1

make vendor-sync
GOFLAGS=-mod=vendor GOCACHE=/tmp/go-cache go test ./...
GOFLAGS=-mod=vendor GOCACHE=/tmp/go-cache go vet ./...
go mod vendor
git diff --exit-code vendor/ go.mod go.sum
```

Opcjonalny sanity-check CLI (smoke):

```bash
./build/docflow --version || true
./build/docflow health --help || true
./build/docflow validate --help || true
./build/docflow compliance --help || true
```

## 5. Wydanie release (tag `v*.*.*`)

### 5.1. Tagowanie

Release workflow uruchamia się na tagach SemVer:

- `v0.1.0`, `v0.1.1`, `v1.0.0` itd.

Przykład:

```bash
git tag -a v0.1.1 -m "v0.1.1"
git push origin v0.1.1
```

### 5.2. Oczekiwane joby w Actions (po tagu)

- `vendor_guard` / gate -> PASS
- `build` (reusable build) -> PASS (artefakty OS/ARCH)
- `sbom` -> PASS (generuje `sbom.cdx.json`, uploaduje artifact `sbom`)
- `release` -> PASS (publikuje release assets + checksums)
- `attest_provenance` -> PASS (best-effort, ale powinno przechodzić)
- `sign_cosign` -> PASS (best-effort, podpisy `.sig/.cert`)

### 5.3. Oczekiwane assets w GitHub Releases

- `docflow-linux-amd64.tar.gz`
- `docflow-linux-arm64.tar.gz`
- `docflow-darwin-amd64.tar.gz`
- `docflow-darwin-arm64.tar.gz`
- `docflow-windows-amd64.zip`
- `checksums.txt` (basename, sha256 dla archiwów)
- Cosign (jeśli włączone): `*.sig`, `*.cert` (dla archiwów i opcjonalnie checksums)
- SBOM: `sbom.cdx.json` (jeśli publikowane jako asset po tagu)
- Attestations: widoczne w GitHub UI (nie zawsze jako asset)

## 6. Test release bez tagu (workflow_dispatch)

### 6.1. Po co

Pozwala przetestować pipeline release (build + sbom + attest + cosign) bez publikowania GitHub Release.

### 6.2. Jak

GitHub -> Actions -> workflow `release` -> Run workflow

- input `ref`: domyślnie `main` albo podaj SHA.

### 6.3. Zasada bezpieczeństwa

W workflow jest guard: krok publikacji release (`softprops/action-gh-release`) wykonuje się tylko dla `refs/tags/v*`.
W `workflow_dispatch`:

- `Create Release` i upload assetów powinny być `skipped`.
- SBOM może być generowany i uploadowany jako artifact.
- Cosign/attestation mogą działać i dawać artifacts.

## 7. Smoke tests binarek (CI)

- Smoke w build workflow jest uruchamiany tylko na natywnej architekturze runnera (native-only).
- Nie używamy emulacji QEMU domyślnie (mniej awarii, szybsze runy).
- Jeśli kiedyś potrzebny smoke dla linux-arm64 na amd64, dodaj osobny job best-effort z QEMU.

## 8. Kanoniczne CI dla repo (health --ci)

W CI (PR/push) kanoniczna ścieżka to:

- `docflow health --ci` generuje bundle `.docflow/out/<run_id>/`
- upload SARIF: `validate.sarif`
- upload bundle jako artifact
- egzekwowanie `overall_exit`

Zasady:

- gating domyślnie blokuje tylko NOWE problemy względem baseline (repo baseline domyślny).
- baseline repo aktualizuje się tylko przez dedykowany PR do `main` (patrz `docs/GA_CHECKLIST.md`).

## 9. Typowe awarie i szybkie naprawy

### 9.1. "inconsistent vendoring"

Objawy:

- Go zwraca "explicitly required in go.mod but not marked as explicit".

Naprawa:

```bash
make vendor-sync
go mod vendor
git diff --exit-code vendor/ go.mod go.sum
git add vendor/ go.mod go.sum
git commit -m "chore: sync vendor"
git push
```

### 9.2. Brak pliku/pakietu w CI, a działa lokalnie

Najczęstsza przyczyna:

- `.gitignore` przypadkowo ignoruje plik w `internal/...` i nie jest trackowany.

Naprawa:

- `git status --ignored` + sprawdź wzorce,
- zmień ignore na "root-only" (np. `/worklog/` zamiast `worklog/`).

### 9.3. SBOM/attestation/cosign

Jeśli SBOM jest best-effort:

- `has_sbom=false` powinno skutkować brakiem download w `release` i brakiem assetu, bez failowania runu.
Jeśli cosign/attestations nie działają:
- sprawdź permissions w workflow (`id-token: write`, `attestations: write`).

## 10. Minimalna checklista przed kolejnym wydaniem

- main: green CI/docflow/build
- `make vendor-sync` przechodzi
- `go test ./...` i `go vet ./...` przechodzą w `-mod=vendor`
- manualny `workflow_dispatch` release działa i nic nie publikuje
- tagowy release publikuje poprawne assets + checksums (+ opcjonalnie cosign/sbom/attestations)
