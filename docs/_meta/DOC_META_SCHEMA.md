# DOC_META_SCHEMA — Schemat metadanych dokumentów

## Wersja schematu

`v1.0.0`

## Format

Metadane dokumentu są zapisywane jako YAML frontmatter na początku pliku Markdown, otoczone znacznikami `---`.

```yaml
---
doc_id: "api_documentation"
title: "API Documentation"
doc_type: "specification"
version: "0.1.0"
status: "draft"
owner: "team-backend"
created: "2026-02-09"
updated: "2026-02-09"
tags:
  - api
  - backend
depends_on:
  - "functional_requirements_specification_frs"
  - "system_architecture"
context_sources:
  - "product_vision_statement"
industry: "software"
phase: "design"
---
```

## Pola obowiązkowe

| Pole | Typ | Opis | Walidacja |
|------|-----|------|-----------|
| `doc_id` | `string` | Unikalny identyfikator dokumentu. Konwencja: snake_case, bez rozszerzenia. | `^[a-z][a-z0-9_]{2,80}$` |
| `title` | `string` | Tytuł dokumentu czytelny dla człowieka. | Niepusty, max 200 znaków |
| `doc_type` | `string` | Typ dokumentu z `DOC_TYPES.md`. | Wartość z enum w DOC_TYPES |
| `version` | `string` | Wersja dokumentu (semver). | `^\d+\.\d+\.\d+$` |
| `status` | `string` | Status dokumentu w cyklu życia. | Enum: `draft`, `in_review`, `approved`, `published`, `deprecated`, `archived`, `needs_content` |

## Pola opcjonalne

| Pole | Typ | Opis | Domyślna wartość |
|------|-----|------|------------------|
| `owner` | `string` | Osoba lub zespół odpowiedzialny za dokument. | `""` |
| `created` | `string` | Data utworzenia (ISO 8601: `YYYY-MM-DD`). | Data pierwszego skanowania |
| `updated` | `string` | Data ostatniej modyfikacji (ISO 8601). | Data ostatniego skanowania |
| `tags` | `[]string` | Lista tagów klasyfikacyjnych. | `[]` |
| `depends_on` | `[]string` | Lista `doc_id` dokumentów wymaganych. Patrz `DOC_DEPENDENCY_SPEC.md`. | `[]` |
| `context_sources` | `[]string` | Lista `doc_id` dokumentów dostarczających kontekst (niebloking). | `[]` |
| `industry` | `string` | Kod branży / domeny. | `""` |
| `phase` | `string` | Faza projektu, w której dokument jest istotny. | `""` |
| `priority` | `string` | Priorytet dokumentu. | Enum: `critical`, `high`, `medium`, `low`. Domyślnie: `medium` |
| `language` | `string` | Język dokumentu (ISO 639-1). | `"pl"` |
| `template_source` | `string` | Ścieżka do szablonu, z którego wygenerowano dokument. | `""` |

## Reguły walidacji

1. **`doc_id` musi być unikalny** w całym repozytorium.
2. **`depends_on`** może zawierać tylko istniejące `doc_id` (walidacja referencyjna).
3. **`context_sources`** może zawierać tylko istniejące `doc_id`.
4. **Cykl zależności** (`depends_on`) jest zabroniony — graf musi być DAG.
5. **`status`** musi być jedną z dozwolonych wartości.
6. **`doc_type`** musi odpowiadać definicji w `DOC_TYPES.md`.
7. **`version`** musi być poprawnym semver.

## Kompatybilność wsteczna

Dokumenty z minimalnym frontmatter (tylko `title` i `status`) są traktowane jako **częściowo zgodne**. Scanner automatycznie:
- Generuje `doc_id` z nazwy pliku (snake_case).
- Ustawia `version` na `"0.1.0"`.
- Ustawia `doc_type` na `"unknown"`.
- Ustawia puste `depends_on` i `context_sources`.
- Loguje ostrzeżenie o niekompletnych metadanych.

## Przykład minimalny

```yaml
---
doc_id: "release_checklist"
title: "Release Checklist"
doc_type: "checklist"
version: "1.0.0"
status: "approved"
---
```

## Przykład pełny

```yaml
---
doc_id: "api_design_specification"
title: "API Design Specification"
doc_type: "specification"
version: "2.1.0"
status: "published"
owner: "team-platform"
created: "2025-11-01"
updated: "2026-02-09"
tags:
  - api
  - rest
  - design
depends_on:
  - "functional_requirements_specification_frs"
  - "system_architecture"
  - "security_requirements"
context_sources:
  - "product_vision_statement"
  - "market_analysis"
industry: "software"
phase: "design"
priority: "critical"
language: "pl"
template_source: "templates/api_design_specification.md"
---
```
