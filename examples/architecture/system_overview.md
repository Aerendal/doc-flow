---
title: "System Overview"
doc_id: "system_overview"
doc_type: "architecture"
status: "published"
version: "v1.1.0"
owner: "arch_team"
depends_on: []
context_sources: ["product_vision"]
---

## Przegląd
Opis systemu.

## Decyzje architektoniczne
- Modularny podział na API/worker.
- Baza relacyjna z replikacją.

## Komponenty
- API
- Worker
- DB

## Ryzyka
- Single region; ryzyko opóźnień między strefami.
