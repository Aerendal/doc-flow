# Maintainer Notes

Ten dokument zawiera wewnętrzne workflow maintainerskie (bundle artefaktów, evidence map, logi operacyjne), które celowo nie są eksponowane w głównym quickstarcie README.

## Release Artifacts Bundle

```bash
./scripts/release_artifacts.sh
```

Skrypt tworzy archiwum release artefaktów (zip) i zapisuje log/checksumy w `LOGS/`.

## PR Bundle

```bash
./scripts/pr_bundle.sh
```

Skrypt tworzy bundle PR (patch + raporty pomocnicze) i zapisuje log w `LOGS/`.

## Evidence Map

Spis logów i artefaktów operacyjnych jest utrzymywany w `LOGS/EVIDENCE_MAP_*.md`.
Aktualizacja mapy odbywa się skryptem:

```bash
./scripts/evidence_map.sh
```
