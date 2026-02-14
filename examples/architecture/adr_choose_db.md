---
title: "ADR: Choose DB"
doc_id: "adr_choose_db"
doc_type: "architecture"
status: "review"
version: "v1.0.0"
depends_on: ["system_overview"]
owner: "arch_team"
context_sources: ["system_overview"]
---

## Przegląd
Decyzja o doborze bazy danych dla systemu.

## Decyzje architektoniczne
- Wybrano PostgreSQL ze względu na ACID i wsparcie dla JSONB.

## Komponenty
- DB cluster (primary + read replica).
- Migracje zarządzane przez tool X.

## Ryzyka
- Potencjalny bottleneck na primary przy dużym ruchu; potrzebne testy obciążeniowe.
