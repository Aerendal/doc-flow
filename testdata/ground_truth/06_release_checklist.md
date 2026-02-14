---
doc_id: "release_checklist"
title: "Release Checklist"
doc_type: "checklist"
version: "1.0.0"
status: "approved"
owner: "team-devops"
created: "2026-01-01"
updated: "2026-02-09"
tags:
  - release
  - checklist
  - devops
depends_on:
  - "functional_requirements_specification_frs"
  - "security_requirements"
context_sources:
  - "system_architecture"
industry: "software"
phase: "release"
priority: "high"
language: "pl"
---

# Release Checklist

## Przed wydaniem
- [ ] Wszystkie testy przechodzą
- [ ] Dokumentacja zaktualizowana
- [ ] CHANGELOG uzupełniony
- [ ] Wersja podbita (semver)
- [ ] Security review zakończony
- [ ] Performance testy przechodzą

## Po wydaniu
- [ ] Tag w Git
- [ ] Binary opublikowany
- [ ] Release notes wysłane
