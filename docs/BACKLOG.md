# Backlog (post RC)

## P0 (blokery)
- Podpiąć realny indeks/status/usage do `recommend`, `template-sets`, `templates *` (usunąć dane demo)
- Flagi CLI: `--strict-mode`, `--publish-strict`, migracja sekcji `migrate-sections --dry-run`

## P1
- Zapis usage w generatorze + konsumowanie w rekomendacjach
- Raport brakujących przykładów (code/table) w CLI
- Impact: auto-backup/auto-update oparty na template_source
- Governance: rozbudować reguły (sekcje, pola), raport CLI

## P2
- Równoległy scan/parse z deterministycznym porządkiem
- Integracyjne testy pipeline na większym ficie (1000+ md)
- UI HTML dla analytics (quality/usage)

## Quick wins
- Dodać flagi strict/publish do CLI (validate)
- Wpiąć usage_store do generatora (inkrementuj po generate)
- Raport brak przykładów: reuse code_blocks/tables metryk
