---
doc_id: "api_design_specification"
title: "API Design Specification"
doc_type: "specification"
version: "0.1.0"
status: "draft"
owner: "team-platform"
created: "2025-12-01"
updated: "2026-01-20"
tags:
  - api
  - rest
  - design
depends_on:
  - "functional_requirements_specification_frs"
  - "system_architecture"
context_sources:
  - "product_vision_statement"
industry: "software"
phase: "design"
priority: "high"
language: "pl"
---

# API Design Specification

## Cel dokumentu
Opisuje interfejsy CLI i wewnętrzne API modułów docflow.

## Interfejs CLI
- `docflow scan` — skanuj repozytorium
- `docflow validate` — waliduj metadane
- `docflow graph` — buduj graf zależności
- `docflow generate` — generuj dokumenty
