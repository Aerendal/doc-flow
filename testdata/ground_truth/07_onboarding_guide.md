---
doc_id: "onboarding_guide"
title: "Developer Onboarding Guide"
doc_type: "guide"
version: "0.1.0"
status: "draft"
owner: "team-backend"
created: "2026-02-01"
updated: "2026-02-09"
tags:
  - onboarding
  - guide
  - developer
depends_on:
  - "system_architecture"
context_sources:
  - "product_vision_statement"
  - "api_design_specification"
industry: "software"
phase: "operations"
priority: "medium"
language: "pl"
---

# Developer Onboarding Guide

## Cel
Przewodnik dla nowych developerów dołączających do projektu docflow.

## Wymagania wstępne
1. Go >= 1.25 zainstalowane
2. Git skonfigurowany
3. Dostęp do repozytorium

## Kroki
1. Sklonuj repozytorium
2. Uruchom `go mod download`
3. Zbuduj: `go build -o docflow ./cmd/docflow/`
4. Uruchom testy: `go test ./...`
