---
doc_id: "incident_runbook"
title: "Incident Response Runbook"
doc_type: "runbook"
version: "0.1.0"
status: "draft"
owner: "team-devops"
created: "2026-02-05"
updated: "2026-02-09"
tags:
  - incident
  - runbook
  - operations
depends_on:
  - "system_architecture"
context_sources:
  - "security_requirements"
  - "release_checklist"
industry: "software"
phase: "operations"
priority: "medium"
language: "pl"
---

# Incident Response Runbook

## Wyzwalacz
Awaria krytyczna w pipeline przetwarzania dokumentów.

## Diagnoza
1. Sprawdź logi: `docflow --log-level debug scan`
2. Sprawdź status cache: `docflow cache status`
3. Zweryfikuj integralność bazy: `docflow db check`

## Kroki naprawcze
1. Wyczyść cache: `docflow cache clear`
2. Przebuduj indeks: `docflow scan --rebuild`
3. Zweryfikuj wynik: `docflow validate`

## Eskalacja
Jeśli problem nie ustąpi — zgłoś issue w repozytorium.
