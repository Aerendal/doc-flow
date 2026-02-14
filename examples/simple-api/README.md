---
title: "Example: Simple API Docs"
doc_id: "readme"
doc_type: "guide"
status: "draft"
owner: "api_team"
version: "v1.0.0"
context_sources: ["seed"]
---

## Przegląd
Minimalna dokumentacja API z zależnościami i regułami governance.

## Kroki
1. `./docflow scan -o .docflow/cache/doc_index.json`
2. `./docflow validate --status-aware --governance docs/_meta/GOVERNANCE_RULES.yaml`
3. `./docflow compliance --rules docs/_meta/GOVERNANCE_RULES.yaml --format text`
4. `./docflow changes --old-index .docflow/cache/doc_index.json` (po modyfikacji)

## FAQ
- **Co zawiera przykład?** `endpoints.md` (depends_on: auth), `authentication.md` (auth flows).
- **Po co?** Pokazuje zależności i walidację governance.
