# Front-of-Book — Docflow GA snapshot (2026-02-09)

## Decyzja release
- GA (na bazie RC2) — patrz `LOGS/RELEASE_DECISION.md`.
- Checksums: `build/checksums.txt` zweryfikowane (`LOGS/CHECKSUM_VERIFY.md`).

## Stan jakości
- Testy: `go test ./...` PASS (`GOFLAGS=-mod=vendor`).
- Perf: `LOGS/SCALE_BASELINE_1k.md`, `LOGS/SCALE_BASELINE_LARGE.md`, `LOGS/SCALE_BASELINE_10k.md` (10k: index ~0.56s, validate ~0.9s, RSS ~80 MB).
- Security/chaos: `LOGS/SECURITY_CHAOS_SANITY.md`.
- Queue Go/No-Go: `scripts/queue_evaluate.sh` (text/json, cache, workers). Ostatni run READY/READY na examples (`LOGS/QUEUE_GO_NO_GO.md`).

## Governance
- Reguły: `docs/_meta/GOVERNANCE_RULES.yaml`.
- Przykłady governance-ready: `examples/simple-api`, `examples/architecture` (validate/compliance PASS).

## Jak uruchomić (3 kroki)
1) `./build/docflow-linux-amd64 validate --config examples/simple-api/docflow.yaml --governance docs/_meta/GOVERNANCE_RULES.yaml`
2) `./build/docflow-linux-amd64 compliance --config examples/simple-api/docflow.yaml --rules docs/_meta/GOVERNANCE_RULES.yaml --format text`
3) Kolejka: `scripts/queue_evaluate.sh --format text examples.backlog` (lub `--format json`, `--workers N`, `--cache cache.json`)

## Linki kluczowe
- Release decision: `LOGS/RELEASE_DECISION.md`
- Perf: `LOGS/SCALE_BASELINE_10k.md`
- Queue: `LOGS/QUEUE_GO_NO_GO.md`
- Observability: `LOGS/OBSERVABILITY_PROFILES.md`

## Otwarte punkty (zaakceptowane na GA)
- Queue JSON nie zawiera violations (plan v0.3).
- Perf grafów fan-out/fan-in niezmierzone (łańcuch 10k OK).

Download: https://github.com/Aerendal/docflow/releases/latest/download/docflow-linux-amd64
