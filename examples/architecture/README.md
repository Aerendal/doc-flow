---
title: "Example: Architecture Docs"
doc_id: "readme"
doc_type: "architecture"
status: "draft"
owner: "arch_team"
version: "v1.0.0"
context_sources: ["architecture_playbook"]
---

## Przegląd
Przykładowy zestaw ADR + system overview z zależnościami.

## Decyzje architektoniczne
- system_overview.md opisuje kontekst.
- adr_choose_db.md dokumentuje wybór bazy.

## Komponenty
- API, Worker, DB.

## Ryzyka
- Wybór bazy wpływa na skalowanie; konieczny monitoring replikacji.
