# Governance model (MVP)

Źródło: day_56

## Zakres
- Status-aware wymagania (draft/review/published).
- Reguły per rodzina (api/architecture/guide).
- Progi jakości i approvals.

## Plik źródłowy
- `docs/_meta/GOVERNANCE_RULES.yaml`

## Struktura YAML
- `statuses`: reguły pola/quality/approvals per status.
- `families`: wymagane sekcje, min_quality_published.
- `approvals`: minimalne role do zatwierdzenia.

## Założenia
- Draft: brak progów jakości, puste sekcje dozwolone.
- Review: 1 approval (owner), min_quality 60.
- Published: 2 approvals (owner + doc_lead), min_quality 75+, brak pustych sekcji, wymagane context_sources.
- Families definiują required_sections i dopuszczalne statusy.

## Kolejne kroki
- day_57: walidator governance (czy pola/sekcje/quality spełniają progi).
