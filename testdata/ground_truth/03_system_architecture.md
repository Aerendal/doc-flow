---
doc_id: "system_architecture"
title: "System Architecture"
doc_type: "specification"
version: "0.2.0"
status: "draft"
owner: "team-backend"
created: "2025-11-15"
updated: "2026-02-05"
tags:
  - architecture
  - design
  - system
depends_on:
  - "functional_requirements_specification_frs"
context_sources:
  - "product_vision_statement"
industry: "software"
phase: "design"
priority: "critical"
language: "pl"
---

# System Architecture

## Cel dokumentu
Opisuje architekturę systemu docflow — komponenty, przepływ danych, decyzje techniczne.

## Komponenty
1. **CLI** — interfejs użytkownika (cobra)
2. **Scanner** — skanowanie plików Markdown
3. **Parser** — parsowanie frontmatter YAML
4. **Validator** — walidacja metadanych
5. **Graph Builder** — budowanie grafu zależności
6. **Generator** — generowanie dokumentów z szablonów
7. **DB** — SQLite do indeksowania
