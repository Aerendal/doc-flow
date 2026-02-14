# DEPENDENCY MATRIX - PROJEKT DOCFLOW
## Wersja: 1.0 | Data: 2026-02-06

---

## RELATED DOCUMENTS

**Navigation:**
- **[← INDEX](INDEX.md)** - Powrót do głównego indeksu dokumentów
- **[EXTENDED_PLAN.md](EXTENDED_PLAN.md)** - Plan rozszerzony (szczegółowy opis każdego dnia)
- **[RISK_REGISTER.md](RISK_REGISTER.md)** - Ryzyka (mitigation plans dla delays)
- **[VALIDATION_REPORT.md](VALIDATION_REPORT.md)** - Raport walidacji (dependency validation)
- **[FIXES_APPLIED.md](FIXES_APPLIED.md)** - Changelog (L-01: soft dep issue)

**Quick links:**
- Task details: See [EXTENDED_PLAN.md](EXTENDED_PLAN.md) for full description of each day referenced in this matrix
- Risk mitigation: See [RISK_REGISTER.md](RISK_REGISTER.md) for contingency plans when dependencies cause delays
- Validation: See [VALIDATION_REPORT.md - TEST 6](VALIDATION_REPORT.md) for cycle detection & logical integrity checks

---

## WPROWADZENIE

Ten dokument mapuje zależności między wszystkimi etapami projektu docflow (90 dni + pre-work).

**Typy zależności:**
- **HARD (blocking):** Etap B nie może rozpocząć się przed zakończeniem etapu A (np. day_04 wymaga parsera z day_03)
- **SOFT (preferowane):** Etap B korzysta z wyników A, ale może rozpocząć się równolegle jeśli assumptions są known (np. day_05 review preferuje day_01-04 ukończone, ale może zacząć wcześniej)

**Critical Path:** Najdłuższa sekwencja zależności hard - określa minimalny czas projektu.

**Walidacja:**
- Każdy etap musi mieć zdefiniowane wejście (dependencies) i wyjście (outputs)
- Cykle są niedozwolone (A→B→C→A)
- Wszystkie dependencies muszą wskazywać na istniejące etapy

---

## FORMAT TABELI

| Kolumna | Opis |
|---------|------|
| **Day** | ID etapu (PRE-1, day_00, etc.) |
| **Name** | Krótka nazwa etapu |
| **Hard Deps** | Etapy które MUSZĄ być ukończone przed rozpoczęciem (blocking) |
| **Soft Deps** | Etapy których wyniki są użyteczne, ale nie blokują startu |
| **Outputs** | Główne artefakty delivered |
| **Critical** | Czy na critical path (Y/N) |
| **Buffer** | Czy to dzień buforowy (Y/N) |

---

## PRE-WORK DEPENDENCIES (5-7 dni)

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| PRE-1 | Analiza wymagań | - | - | REQUIREMENTS.md, USER_PERSONAS.md, USE_CASES.md | Y | N |
| PRE-2 | Korpus danych | - | PRE-1 | testdata/templates/ (300+), DATA_INVENTORY.md | Y | N |
| PRE-3 | Wybór technologii | - | PRE-1 | TECH_STACK.md, ADR-001..003, ENVIRONMENT.md | Y | N |
| PRE-4 | Architektura | PRE-3 | PRE-1 | ARCHITECTURE.md, DATA_MODEL.md, API_CONTRACTS.md | Y | N |
| PRE-5 | Plan i risk mgmt | PRE-4 | PRE-1, PRE-2 | PROJECT_PLAN.md, RISK_REGISTER.md | Y | N |

**Pre-work critical path:** PRE-1 → PRE-3 → PRE-4 → PRE-5 (4 dni minimum, 5-7 realistic z PRE-2 parallel)

**Walidacja:**
- ✓ No cycles
- ✓ All dependencies exist
- ✓ PRE-2 może być parallel z PRE-3 (both depend tylko na PRE-1)

---

## PHASE 1: FOUNDATION & MVP (day_00-25)

### WEEK 1: Setup & Core Parsing

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_00 | Benchmark + environment | PRE-5 | PRE-2, PRE-3 | Benchmark results, repo structure, DECISIONS.md | Y | N |
| day_01 | Metadata contract (1/2) | day_00 | PRE-1, PRE-4 | DOC_META_SCHEMA.md, DOC_DEPENDENCY_SPEC.md, examples (10) | Y | N |
| day_02 | Config + logger + walker (2/2) | day_01 | PRE-4 | config.go, logger.go, fileutil.go, docflow.yaml, CLI stub | Y | N |
| day_03 | Markdown + YAML parser | day_02 | day_01 | frontmatter.go, markdown.go, DocumentRecord model | Y | N |
| day_04 | Document index | day_03 | day_01 | document_index.go, cache.go, `docflow scan` command | Y | N |
| day_05 | Code review + cleanup | day_00..04 | - | CODE_REVIEW_WEEK1.md, refactored code, CI pipeline | N | Y |

**Week 1 critical path:** day_00 → day_01 → day_02 → day_03 → day_04 (5 dni)

### WEEK 2: Validation & Dependency Graph

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_06 | Metadata validator | day_04 | day_01 | metadata.go, `docflow validate` command | Y | N |
| day_07 | Section parser | day_03 | - | parser.go, SectionTree model | Y | N |
| day_08 | Section schema | day_07 | day_01 | SECTION_SCHEMA.md, section_schema.go, validator.go, schemas (3) | Y | N |
| day_09 | Dependency graph | day_04, day_06 | - | dependency.go, DependencyGraph model, cycle detection | Y | N |
| day_10 | Topo sort + context | day_09 | day_01 | toposort.go, context.go, `docflow graph` command | Y | N |

**Week 2 critical path:** day_06, day_09, day_10 depend na week 1; day_07-08 parallel track
- Longest: day_04 → day_09 → day_10 (+ day_06 parallel)

### WEEK 3: MVP Features

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_11 | Template index | day_08 | PRE-2 | templates/index.go, selector.go | Y | N |
| day_12 | Document generator | day_11 | day_01, day_08 | generator.go, `docflow generate` command | Y | N |
| day_13 | MVP integration test + dogfooding | day_00..12 | - | mvp_pipeline_test.go, MVP_TEST_RESULTS.md, bug fixes | Y | N |
| day_14 | MVP documentation | day_13 | - | README.md, USER_GUIDE.md, CLI_REFERENCE.md | Y | N |
| day_15 | MVP demo + release | day_14 | - | Demo screencast, v0.1.0-mvp release, binaries, CHANGELOG.md | Y | N |

**Week 3 critical path:** day_11 → day_12 → day_13 → day_14 → day_15 (5 dni)

### WEEK 4: Schema Extraction & Quality

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_16 | Pattern extraction | day_07, PRE-2 | - | extractor.go, ALGO_PARAMS.md (patterns) | Y | N |
| day_17 | Schema generation | day_16 | - | schema_generator.go, `docflow analyze-patterns`, generated_schemas/ | Y | N |
| day_18 | Quality metrics | day_08, PRE-2 | - | metrics.go (structure, content, usage scores), ALGO_PARAMS.md (weights) | Y | N |
| day_19 | Quality scoring + validation | day_18 | PRE-2 | scorer.go, `docflow score` command, validation vs ground truth | Y | N |
| day_20 | Buffer + tech debt | day_16..19 | - | Bug fixes, refactoring | N | Y |

**Week 4 critical path:** day_16 → day_17 (pattern track), day_18 → day_19 (quality track) - parallel

### WEEK 5: Testing & Phase 1 Release

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_21 | Test strategy + unit tests | day_00..19 | - | TEST_STRATEGY.md, unit tests (70% coverage) | Y | N |
| day_22 | Integration tests | day_21 | day_13 | tests/integration/, test fixtures (100 docs) | Y | N |
| day_23 | Performance testing | day_22 | day_00 | benchmark_test.go, PERFORMANCE_BASELINE.md | Y | N |
| day_24 | Documentation finalization | day_15 | day_21..23 | Updated README, ARCHITECTURE.md, CONTRIBUTING.md, CHANGELOG | Y | N |
| day_25 | Phase 1 release + retro | day_21..24 | - | v0.1.0 release, PHASE1_RETROSPECTIVE.md, demo | Y | N |

**Week 5 critical path:** day_21 → day_22 → day_23 (parallel z day_24) → day_25 (4 dni)

**Phase 1 total critical path:** day_00 → ... → day_25 = ~21 dni (z 25 available, 4 dni slack z buffers)

**Walidacja Phase 1:**
- ✓ No cycles
- ✓ Critical path identified (21/25 dni)
- ✓ Buffers: day_05, day_20 (2 dni, 8% buffer)
- ✓ All hard deps exist
- ⚠ RISK: day_13 depends na ALL day_00-12 (integration point - jeśli any fails, day_13 blocked)

---

## PHASE 2: INTELLIGENCE & AUTOMATION (day_26-55)

### WEEK 6: Template Recommendation

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_26 | Template metadata expansion | day_11, day_19 | - | Extended metadata (quality, pattern, usage, version, status) | Y | N |
| day_27 | Recommendation scoring | day_26 | - | scorer.go, scoring algorithm, ALGO_PARAMS.md (weights) | Y | N |
| day_28 | Recommender + CLI | day_27 | - | recommender.go, `docflow recommend`, usage tracking | Y | N |
| day_29 | Recommendation evaluation | day_28 | PRE-2 | evaluation_test.go, RECOMMENDER_EVAL.md, tuned weights | Y | N |
| day_30 | Buffer | day_26..29 | - | Bug fixes, tuning | N | Y |

**Week 6 critical path:** day_26 → day_27 → day_28 → day_29 (4 dni)

### WEEK 7: Dependency Rules & Planner

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_31 | Family dependency rules (1/2) | day_01, day_06 | PRE-1 | FAMILY_RULES.yaml (rules dla 5-7 rodzin) | Y | N |
| day_32 | Family validator (2/2) | day_31 | day_09 | family_rules.go, extended `docflow validate --family-rules` | Y | N |
| day_33 | Effort estimation | day_10 | day_19 | effort.go (heuristics), ALGO_PARAMS.md (effort baselines) | Y | N |
| day_34 | Daily planner | day_33, day_10 | day_32 | daily.go, `docflow plan`, DAILY_PLAN.md generation | Y | N |
| day_35 | Buffer | day_31..34 | - | Bug fixes | N | Y |

**Week 7 critical path:** day_31 → day_32, day_33 → day_34 (parallel tracks, longest: 2 dni each)

### WEEK 8: Lifecycle & Versioning

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_36 | Template lifecycle states | day_26, day_19 | - | TEMPLATE_LIFECYCLE.md, lifecycle.go, transition rules | Y | N |
| day_37 | Deprecation + migration | day_36 | - | `docflow templates deprecate/deprecated/suggest-migration` | Y | N |
| day_38 | Template versioning | day_26 | - | versioning.go, multi-version template support | Y | N |
| day_39 | Document versioning | day_04 | - | Document version tracking, change detection | Y | N |
| day_40 | Code review + refactoring | day_36..39 | - | Refactored code, review notes | N | Y |

**Week 8 critical path:** day_36 → day_37, day_38, day_39 parallel (2 dni)

### WEEK 9: Validation Enhancements

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_41 | Section completeness metrics | day_08, day_06 | - | metrics.go (completeness), `docflow stats` command | Y | N |
| day_42 | Progressive validation | day_06, day_08 | - | progressive.go, status-aware validation | Y | N |
| day_43 | Edge case hardening (1/2) | day_00..42 | - | tests/edge_cases/, edge case tests, bug fixes | Y | N |
| day_44 | Fuzzy matching + migration (2/2) | day_08, day_43 | - | fuzzy.go, section aliases, `docflow migrate-sections` | Y | N |
| day_45 | Buffer | day_41..44 | - | Bug fixes | N | Y |

**Week 9 critical path:** day_41, day_42 parallel; day_43 → day_44 (longest: 2 dni)

### WEEK 10: Optimization & Caching

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_46 | Hash tracking | day_04 | - | hash.go, SHA256 per file, invalidation.go | Y | N |
| day_47 | Incremental scan | day_46 | - | Extended `docflow scan --incremental`, benchmark | Y | N |
| day_48 | Template impact analysis | day_26, day_12 | - | impact.go, `docflow template-impact`, auto-update (optional) | Y | N |
| day_49 | Profiling | day_47 | day_23 | PROFILING_RESULTS.md, bottleneck identification | Y | N |
| day_50 | Optimization | day_49 | - | Optimized code, OPTIMIZATION_RESULTS.md, re-benchmark | Y | N |

**Week 10 critical path:** day_46 → day_47, day_49 → day_50 (parallel, longest: day_46→47 + day_49→50 = 4 dni)

### WEEK 11: Analytics & Reporting

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_51 | Analytics engine | day_28, day_19, day_36 | - | aggregator.go, analytics metrics | Y | N |
| day_52 | Analytics reporting | day_51 | - | `docflow analytics`, export.go (HTML, CSV) | Y | N |
| day_53 | Duplicate detection | day_26 | - | detector.go, MinHash+Jaccard, `docflow find-duplicates` | Y | N |
| day_54 | Content hints extraction | day_03, day_26 | - | extractor.go (code blocks, tables, links), content metrics | Y | N |
| day_55 | Phase 2 buffer + retro | day_51..54 | - | Bug fixes, PHASE2_RETROSPECTIVE.md | N | Y |

**Week 11 critical path:** day_51 → day_52, day_53, day_54 parallel (2 dni)

**Phase 2 total critical path:** day_26 → ... → day_55 = ~24 dni (z 30 available, 6 dni slack)

**Walidacja Phase 2:**
- ✓ No cycles
- ✓ Critical path: 24/30 dni
- ✓ Buffers: day_30, 35, 40, 45, 55 (5 dni, 17% buffer)
- ✓ Dependencies na Phase 1: day_26 depends on day_11, day_19 (Phase 1 complete)
- ✓ Most days have parallel work (efficient use of time)

---

## PHASE 3: PRODUCTION READY (day_56-90)

### WEEK 12: Governance & Compliance

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_56 | Governance rules (1/2) | day_06, day_08, day_31 | - | GOVERNANCE_RULES.yaml, required metadata/sections/quality | Y | N |
| day_57 | Governance validator (2/2) | day_56 | - | governance.go, extended `docflow validate --governance` | Y | N |
| day_58 | Compliance reporting (1/2) | day_57 | - | reporter.go, compliance checks | Y | N |
| day_59 | Compliance export (2/2) | day_58 | - | `docflow compliance`, HTML/PDF export | Y | N |
| day_60 | Buffer | day_56..59 | - | Bug fixes | N | Y |

**Week 12 critical path:** day_56 → day_57 → day_58 → day_59 (4 dni)

### WEEK 13: CI/CD & Automation

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_61 | CI pipeline (1/2) | day_21, day_22 | day_05 | .github/workflows/ci.yml, pre-commit hooks | Y | N |
| day_62 | CI finalization (2/2) | day_61 | day_23 | Full CI (build, test, lint, coverage, benchmarks) | Y | N |
| day_63 | Release pipeline (1/2) | day_62 | - | Build automation (3 platforms), packaging | Y | N |
| day_64 | Deployment automation (2/2) | day_63 | - | install.sh, Homebrew formula, release automation | Y | N |
| day_65 | Buffer | day_61..64 | - | CI/CD fixes | N | Y |

**Week 13 critical path:** day_61 → day_62 → day_63 → day_64 (4 dni)

### WEEK 14: Documentation & Examples

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_66 | User documentation (1/2) | day_14, day_25 | day_55 | USER_GUIDE.md (complete), CLI_REFERENCE.md, TROUBLESHOOTING.md | Y | N |
| day_67 | User docs finalization (2/2) | day_66 | - | BEST_PRACTICES.md, complete user docs | Y | N |
| day_68 | Examples & tutorials | day_67 | - | examples/ (3-5 projects), each z README + samples | Y | N |
| day_69 | Video tutorial | day_68 | - | 15-min screencast, published | Y | N |
| day_70 | Buffer | day_66..69 | - | Documentation improvements | N | Y |

**Week 14 critical path:** day_66 → day_67 → day_68 → day_69 (4 dni)

### WEEK 15: Testing & Hardening

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_71 | E2E testing (1/2) | day_22 | day_55 | tests/e2e/, full workflow tests, error scenarios | Y | N |
| day_72 | E2E cross-platform (2/2) | day_71 | day_64 | E2E tests on Linux, macOS, Windows | Y | N |
| day_73 | Security audit | day_00..72 | - | Security checklist, vulnerability scan, SECURITY.md | Y | N |
| day_74 | Chaos testing | day_73 | - | Chaos scenarios, resilience tests, fixes | Y | N |
| day_75 | Buffer | day_71..74 | - | Test fixes, hardening | N | Y |

**Week 15 critical path:** day_71 → day_72, day_73 → day_74 (parallel, 2 dni each)

### WEEK 16: Performance & Scale

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_76 | Scale testing (1/2) | day_23, day_50 | - | Scale tests (1000, 5000 docs), SCALE_TEST_RESULTS.md | Y | N |
| day_77 | Scale analysis (2/2) | day_76 | - | Bottleneck analysis, scale limits identified | Y | N |
| day_78 | Final optimization (1/2) | day_77 | - | Optimizations dla bottlenecks | Y | N |
| day_79 | Optimization validation (2/2) | day_78 | - | Re-benchmark, FINAL_PERFORMANCE.md | Y | N |
| day_80 | Buffer | day_76..79 | - | Performance tuning | N | Y |

**Week 16 critical path:** day_76 → day_77 → day_78 → day_79 (4 dni)

### WEEK 17: Release Preparation

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_81 | Release candidate | day_00..80 | - | v1.0.0-rc1 tag, binaries, internal testing, bug bash | Y | N |
| day_82 | Bug fixes (1/2) | day_81 | - | Critical bug fixes z RC testing | Y | N |
| day_83 | Bug fixes finalization (2/2) | day_82 | - | v1.0.0-rc2 (if needed), retest | Y | N |
| day_84 | Final docs review | day_67, day_83 | - | All docs reviewed, corrections, CHANGELOG finalized, release notes | Y | N |
| day_85 | Buffer | day_81..84 | - | Pre-launch fixes | N | Y |

**Week 17 critical path:** day_81 → day_82 → day_83 → day_84 (4 dni)

### WEEK 18: Launch

| Day | Name | Hard Deps | Soft Deps | Outputs | Critical | Buffer |
|-----|------|-----------|-----------|---------|----------|--------|
| day_86 | Pre-launch checklist | day_84 | day_64 | LAUNCH_CHECKLIST.md verified, all checks passed | Y | N |
| day_87 | v1.0 release | day_86 | - | v1.0.0 tag, GitHub Release, binaries, announcement | Y | N |
| day_88 | Post-launch monitoring | day_87 | - | Monitoring setup, issue responses, hotfix if needed | Y | N |
| day_89 | Project retrospective | day_88 | - | PROJECT_RETROSPECTIVE.md, team debrief | Y | N |
| day_90 | Backlog & roadmap | day_89 | - | BACKLOG.md, ROADMAP.md (v1.1, v1.2, v2.0) | Y | N |

**Week 18 critical path:** day_86 → day_87 → day_88 → day_89 → day_90 (5 dni)

**Phase 3 total critical path:** day_56 → ... → day_90 = ~30 dni (z 35 available, 5 dni slack)

**Walidacja Phase 3:**
- ✓ No cycles
- ✓ Critical path: 30/35 dni
- ✓ Buffers: day_60, 65, 70, 75, 80, 85 (6 dni, 17% buffer)
- ✓ Dependencies na Phase 2: day_56 depends on day_06, 08, 31 (all Phase 1-2)
- ✓ Launch sequence (day_86-90) jest linear (proper gating)

---

## CRITICAL PATH ANALYSIS

### Overall Critical Path (Pre-work → Phase 3)

**Longest sequential path:**
```
PRE-1 (1d) → PRE-3 (2d) → PRE-4 (1d) → PRE-5 (1d) →
day_00 (1d) → day_01 (1d) → day_02 (1d) → day_03 (1d) → day_04 (1d) →
day_09 (1d) → day_10 (1d) → day_11 (1d) → day_12 (1d) → day_13 (1d) → day_14 (1d) → day_15 (1d) →
day_16 (1d) → day_17 (1d) → day_18 (1d) → day_19 (1d) →
day_21 (1d) → day_22 (1d) → day_23 (1d) → day_25 (1d) →
day_26 (1d) → day_27 (1d) → day_28 (1d) → day_29 (1d) →
day_31 (1d) → day_32 (1d) → day_33 (1d) → day_34 (1d) →
day_36 (1d) → day_37 (1d) →
day_41 (1d) → day_42 (1d) →
day_46 (1d) → day_47 (1d) → day_49 (1d) → day_50 (1d) →
day_51 (1d) → day_52 (1d) →
day_56 (1d) → day_57 (1d) → day_58 (1d) → day_59 (1d) →
day_61 (1d) → day_62 (1d) → day_63 (1d) → day_64 (1d) →
day_66 (1d) → day_67 (1d) → day_68 (1d) → day_69 (1d) →
day_71 (1d) → day_72 (1d) → day_73 (1d) → day_74 (1d) →
day_76 (1d) → day_77 (1d) → day_78 (1d) → day_79 (1d) →
day_81 (1d) → day_82 (1d) → day_83 (1d) → day_84 (1d) →
day_86 (1d) → day_87 (1d) → day_88 (1d) → day_89 (1d) → day_90 (1d)
```

**Critical Path Length:** ~75 dni (na 90+5 planned = 95 total)

**Slack (buffer availability):** 20 dni (13 buffer days + 7 dni slack w critical path)

**Utilization:** 79% (75/95)

**Risk assessment:**
- ✓ Reasonable slack (21% buffer)
- ✓ Critical path clearly identified
- ⚠ Długie linear sequences (day_00-15, day_56-90) - mało parallelization opportunities
- ✓ Most 2-day tasks mogą być split jeśli resource available (parallel developers)

---

## PARALLEL WORK OPPORTUNITIES

Następujące etapy MOGĄ być wykonywane równolegle (no hard dependencies między nimi):

### Pre-work:
- PRE-2 || PRE-3 (both depend only on PRE-1)

### Phase 1:
- day_07-08 (section parser) || day_06 (metadata validator) - różne moduły
- day_16-17 (pattern) || day_18-19 (quality) - różne algorithmy

### Phase 2:
- day_31-32 (family rules) || day_33-34 (planner) - różne features
- day_36-37 || day_38 || day_39 - różne lifecycle features
- day_41 || day_42 - różne validation enhancements
- day_46-47 || day_48 - caching vs impact analysis
- day_51-52 || day_53 || day_54 - różne analytics features

### Phase 3:
- day_71-72 || day_73-74 - testing (E2E vs security/chaos)
- (Limited parallelization - większość sequential dla quality gates)

**Note:** Parallelization wymaga:
- Multiple developers (2+)
- Clear module boundaries (no shared code conflicts)
- Coordination (daily standups)

**Estimated speedup z 2 developers:**
- Pre-work: 5→4 dni (20% speedup)
- Phase 1: 25→22 dni (12% speedup)
- Phase 2: 30→25 dni (17% speedup)
- Phase 3: 35→32 dni (9% speedup)
- **Total: 95→83 dni (13% speedup)**

---

## VALIDATION CHECKS

### 1. Dependency Existence Check

Wszystkie hard dependencies wskazują na istniejące etapy:
- ✓ Verified: All referenced days exist w planie
- ✓ No dangling dependencies

### 2. Cycle Detection

Graph zależności jest acyclic (DAG):
- ✓ Verified: No cycles detected
- ✓ Topological order możliwy

### 3. Unreachable Nodes

Wszystkie etapy są osiągalne z pre-work:
- ✓ Verified: All days connected do dependency graph
- ✓ No isolated nodes

### 4. Critical Path Validation

Critical path nie przekracza timeline:
- ✓ Critical path: 75 dni
- ✓ Total timeline: 95 dni (90 + 5 pre-work)
- ✓ Slack: 20 dni (21%)
- ✓ Acceptable

### 5. Buffer Placement Validation

Bufory są umieszczone strategicznie (po zakończeniu major milestones):
- ✓ day_05: Po week 1 (setup)
- ✓ day_20: Po week 4 (schema & quality)
- ✓ day_30, 35, 40, 45, 55: Co ~5-10 dni w Phase 2
- ✓ day_60, 65, 70, 75, 80, 85: Co ~5 dni w Phase 3
- ✓ Total: 13 buffer days (14% of 90 days)

### 6. Phase Boundary Validation

Phases są cleanly separated:
- ✓ Phase 1 (day_00-25): Complete foundation, MVP released
- ✓ Phase 2 (day_26-55): Intelligence features, no dependencies na Phase 3
- ✓ Phase 3 (day_56-90): Production ready, depends on Phase 1-2 complete
- ✓ Clean separation, możliwe early exit po Phase 1 lub 2

### 7. Output Completeness

Wszystkie etapy mają zdefiniowane outputs:
- ✓ Verified: All days list deliverables
- ✓ Outputs traceable (code files, docs, commands, artefacts)

---

## CHANGE IMPACT ANALYSIS

Jeśli etap X jest delayed o N dni, które inne etapy są affected?

### High-Impact Nodes (blocking many others):

**day_04 (Document index):**
- Directly blocks: day_06, day_09
- Transitively blocks: day_10, day_13, day_14, day_15, ..., entire project
- **Impact multiplier: 87 days** (delays całego projektu)
- **Mitigation:** day_04 jest na critical path, ma buffer w day_05 + day_20

**day_13 (MVP integration test):**
- Directly blocks: day_14, day_15
- Transitively blocks: Phase 1 release
- **Impact multiplier: 3 days** (delays MVP release)
- **Mitigation:** Buffers available (day_20 if needed)

**day_25 (Phase 1 complete):**
- Directly blocks: Phase 2 start (day_26)
- Transitively blocks: Entire Phase 2 + 3
- **Impact multiplier: 65 days**
- **Mitigation:** Phase 1 Go/No-Go decision point

**day_81 (Release candidate):**
- Directly blocks: day_82, day_83, day_84
- Transitively blocks: v1.0 release (day_87)
- **Impact multiplier: 10 days**
- **Mitigation:** Buffers day_85, możliwy delay release

### Low-Impact Nodes (few dependencies):

**day_53 (Duplicate detection):**
- Blocks: Nic (standalone feature)
- **Impact multiplier: 0 days**
- **Mitigation:** Can be cut jeśli needed

**day_54 (Content hints):**
- Blocks: Nic
- **Impact multiplier: 0 days**
- **Mitigation:** Can be cut

**day_68-69 (Examples & tutorial):**
- Blocks: Tylko siebie
- **Impact multiplier: 2 days** (affects docs completeness)
- **Mitigation:** Can extend into buffer day_70

---

## DEPENDENCY GRAPH VISUALIZATION

```
PRE-WORK (5d):
PRE-1 ──┬──> PRE-3 ──> PRE-4 ──> PRE-5
        └──> PRE-2

PHASE 1 - WEEK 1 (5d + 1 buffer):
day_00 ──> day_01 ──> day_02 ──> day_03 ──> day_04 ──> [day_05 buffer]

PHASE 1 - WEEK 2 (5d):
                                                    ┌──> day_07 ──> day_08
                                                    │
day_04 ──> day_06 ──────────────────┬──> day_09 ──> day_10
                                     │
                                     └──────────────────────────┘

PHASE 1 - WEEK 3 (5d):
day_08 ──> day_11 ──> day_12 ──┬──> day_13 ──> day_14 ──> day_15
                                │
day_10 ─────────────────────────┘

PHASE 1 - WEEK 4 (4d + 1 buffer):
                        ┌──> day_16 ──> day_17
day_07 ─────────────────┤
                        └──> day_18 ──> day_19 ──> [day_20 buffer]

PHASE 1 - WEEK 5 (5d):
day_19 ──> day_21 ──> day_22 ──> day_23 ──┬──> day_25
                                           │
day_15 ────────────────> day_24 ───────────┘

PHASE 2 - WEEK 6-11 (24d + 6 buffers):
[Similar pattern: linear sequences z parallelization opportunities]

PHASE 3 - WEEK 12-18 (30d + 6 buffers):
[Similar pattern: mostly linear dla quality gates]

LAUNCH:
day_86 ──> day_87 ──> day_88 ──> day_89 ──> day_90
```

---

## MACHINE-READABLE FORMAT (CSV)

```csv
Day,Name,HardDeps,SoftDeps,Outputs,Critical,Buffer
PRE-1,Analiza wymagań,,,"REQUIREMENTS.md,USER_PERSONAS.md,USE_CASES.md",Y,N
PRE-2,Korpus danych,,PRE-1,"testdata/templates/,DATA_INVENTORY.md",Y,N
PRE-3,Wybór technologii,,PRE-1,"TECH_STACK.md,ADR-*,ENVIRONMENT.md",Y,N
PRE-4,Architektura,PRE-3,PRE-1,"ARCHITECTURE.md,DATA_MODEL.md,API_CONTRACTS.md",Y,N
PRE-5,Plan i risk mgmt,PRE-4,"PRE-1,PRE-2","PROJECT_PLAN.md,RISK_REGISTER.md",Y,N
day_00,Benchmark + environment,PRE-5,"PRE-2,PRE-3","Benchmark results,repo structure,DECISIONS.md",Y,N
day_01,Metadata contract,day_00,"PRE-1,PRE-4","DOC_META_SCHEMA.md,DOC_DEPENDENCY_SPEC.md,examples",Y,N
day_02,Config + logger + walker,day_01,PRE-4,"config.go,logger.go,fileutil.go,CLI stub",Y,N
day_03,Markdown + YAML parser,day_02,day_01,"frontmatter.go,markdown.go,DocumentRecord",Y,N
day_04,Document index,day_03,day_01,"document_index.go,cache.go,docflow scan",Y,N
day_05,Code review + cleanup,"day_00,day_01,day_02,day_03,day_04",,"CODE_REVIEW_WEEK1.md,CI pipeline",N,Y
day_06,Metadata validator,day_04,day_01,"metadata.go,docflow validate",Y,N
day_07,Section parser,day_03,,"parser.go,SectionTree",Y,N
day_08,Section schema,day_07,day_01,"SECTION_SCHEMA.md,section_schema.go,schemas",Y,N
day_09,Dependency graph,"day_04,day_06",,"dependency.go,cycle detection",Y,N
day_10,Topo sort + context,day_09,day_01,"toposort.go,context.go,docflow graph",Y,N
day_11,Template index,day_08,PRE-2,"templates/index.go,selector.go",Y,N
day_12,Document generator,day_11,"day_01,day_08","generator.go,docflow generate",Y,N
day_13,MVP integration test,"day_00,day_01,day_02,day_03,day_04,day_06,day_07,day_08,day_09,day_10,day_11,day_12",,"mvp_pipeline_test.go,bug fixes",Y,N
day_14,MVP documentation,day_13,,"README.md,USER_GUIDE.md,CLI_REFERENCE.md",Y,N
day_15,MVP demo + release,day_14,,"Demo,v0.1.0-mvp,binaries,CHANGELOG.md",Y,N
day_16,Pattern extraction,"day_07,PRE-2",,"extractor.go,ALGO_PARAMS.md",Y,N
day_17,Schema generation,day_16,,"schema_generator.go,docflow analyze-patterns",Y,N
day_18,Quality metrics,"day_08,PRE-2",,"metrics.go,ALGO_PARAMS.md",Y,N
day_19,Quality scoring,day_18,PRE-2,"scorer.go,docflow score,validation",Y,N
day_20,Buffer,"day_16,day_17,day_18,day_19",,"Bug fixes,refactoring",N,Y
day_21,Test strategy + unit tests,"day_00..day_19",,"TEST_STRATEGY.md,unit tests",Y,N
day_22,Integration tests,day_21,day_13,"tests/integration/,fixtures",Y,N
day_23,Performance testing,day_22,day_00,"benchmark_test.go,PERFORMANCE_BASELINE.md",Y,N
day_24,Documentation finalization,day_15,"day_21,day_22,day_23","README,ARCHITECTURE.md,CONTRIBUTING.md",Y,N
day_25,Phase 1 release + retro,"day_21,day_22,day_23,day_24",,"v0.1.0,PHASE1_RETROSPECTIVE.md,demo",Y,N
day_26,Template metadata expansion,"day_11,day_19",,"Extended metadata",Y,N
day_27,Recommendation scoring,day_26,,"scorer.go,ALGO_PARAMS.md",Y,N
day_28,Recommender + CLI,day_27,,"recommender.go,docflow recommend",Y,N
day_29,Recommendation evaluation,day_28,PRE-2,"evaluation_test.go,RECOMMENDER_EVAL.md",Y,N
day_30,Buffer,"day_26,day_27,day_28,day_29",,"Bug fixes,tuning",N,Y
day_31,Family rules,"day_01,day_06",PRE-1,"FAMILY_RULES.yaml",Y,N
day_32,Family validator,day_31,day_09,"family_rules.go,docflow validate --family-rules",Y,N
day_33,Effort estimation,day_10,day_19,"effort.go,ALGO_PARAMS.md",Y,N
day_34,Daily planner,"day_33,day_10",day_32,"daily.go,docflow plan",Y,N
day_35,Buffer,"day_31,day_32,day_33,day_34",,Bug fixes,N,Y
day_36,Template lifecycle,"day_26,day_19",,"TEMPLATE_LIFECYCLE.md,lifecycle.go",Y,N
day_37,Deprecation + migration,day_36,,"docflow templates deprecate/deprecated",Y,N
day_38,Template versioning,day_26,,"versioning.go,multi-version support",Y,N
day_39,Document versioning,day_04,,"Version tracking,change detection",Y,N
day_40,Code review,"day_36,day_37,day_38,day_39",,"Refactored code",N,Y
day_41,Section completeness,"day_08,day_06",,"metrics.go,docflow stats",Y,N
day_42,Progressive validation,"day_06,day_08",,"progressive.go,status-aware validation",Y,N
day_43,Edge case hardening,"day_00..day_42",,"tests/edge_cases/,bug fixes",Y,N
day_44,Fuzzy matching,"day_08,day_43",,"fuzzy.go,docflow migrate-sections",Y,N
day_45,Buffer,"day_41,day_42,day_43,day_44",,Bug fixes,N,Y
day_46,Hash tracking,day_04,,"hash.go,invalidation.go",Y,N
day_47,Incremental scan,day_46,,"docflow scan --incremental,benchmark",Y,N
day_48,Template impact,"day_26,day_12",,"impact.go,docflow template-impact",Y,N
day_49,Profiling,day_47,day_23,"PROFILING_RESULTS.md",Y,N
day_50,Optimization,day_49,,"Optimized code,OPTIMIZATION_RESULTS.md",Y,N
day_51,Analytics engine,"day_28,day_19,day_36",,"aggregator.go,metrics",Y,N
day_52,Analytics reporting,day_51,,"docflow analytics,export.go",Y,N
day_53,Duplicate detection,day_26,,"detector.go,docflow find-duplicates",Y,N
day_54,Content hints,"day_03,day_26",,"extractor.go,content metrics",Y,N
day_55,Phase 2 buffer + retro,"day_51,day_52,day_53,day_54",,"Bug fixes,PHASE2_RETROSPECTIVE.md",N,Y
day_56,Governance rules,"day_06,day_08,day_31",,"GOVERNANCE_RULES.yaml",Y,N
day_57,Governance validator,day_56,,"governance.go,docflow validate --governance",Y,N
day_58,Compliance reporting,day_57,,"reporter.go,compliance checks",Y,N
day_59,Compliance export,day_58,,"docflow compliance,HTML/PDF",Y,N
day_60,Buffer,"day_56,day_57,day_58,day_59",,Bug fixes,N,Y
day_61,CI pipeline,"day_21,day_22",day_05,".github/workflows/ci.yml,pre-commit hooks",Y,N
day_62,CI finalization,day_61,day_23,"Full CI pipeline",Y,N
day_63,Release pipeline,day_62,,"Build automation,packaging",Y,N
day_64,Deployment automation,day_63,,"install.sh,release automation",Y,N
day_65,Buffer,"day_61,day_62,day_63,day_64",,CI/CD fixes,N,Y
day_66,User documentation,"day_14,day_25",day_55,"USER_GUIDE.md,CLI_REFERENCE.md,TROUBLESHOOTING.md",Y,N
day_67,User docs finalization,day_66,,"BEST_PRACTICES.md",Y,N
day_68,Examples & tutorials,day_67,,"examples/",Y,N
day_69,Video tutorial,day_68,,15-min screencast,Y,N
day_70,Buffer,"day_66,day_67,day_68,day_69",,Documentation improvements,N,Y
day_71,E2E testing,day_22,day_55,"tests/e2e/",Y,N
day_72,E2E cross-platform,day_71,day_64,"E2E on 3 platforms",Y,N
day_73,Security audit,"day_00..day_72",,"Security checklist,SECURITY.md",Y,N
day_74,Chaos testing,day_73,,"Chaos scenarios,fixes",Y,N
day_75,Buffer,"day_71,day_72,day_73,day_74",,"Test fixes,hardening",N,Y
day_76,Scale testing,"day_23,day_50",,"Scale tests,SCALE_TEST_RESULTS.md",Y,N
day_77,Scale analysis,day_76,,"Bottleneck analysis",Y,N
day_78,Final optimization,day_77,,Optimizations,Y,N
day_79,Optimization validation,day_78,,"Re-benchmark,FINAL_PERFORMANCE.md",Y,N
day_80,Buffer,"day_76,day_77,day_78,day_79",,Performance tuning,N,Y
day_81,Release candidate,"day_00..day_80",,"v1.0.0-rc1,binaries,bug bash",Y,N
day_82,Bug fixes,day_81,,Critical fixes,Y,N
day_83,Bug fixes finalization,day_82,,"v1.0.0-rc2,retest",Y,N
day_84,Final docs review,"day_67,day_83",,"Docs reviewed,CHANGELOG,release notes",Y,N
day_85,Buffer,"day_81,day_82,day_83,day_84",,Pre-launch fixes,N,Y
day_86,Pre-launch checklist,day_84,day_64,"LAUNCH_CHECKLIST.md verified",Y,N
day_87,v1.0 release,day_86,,"v1.0.0,GitHub Release,announcement",Y,N
day_88,Post-launch monitoring,day_87,,"Monitoring,issue responses",Y,N
day_89,Project retrospective,day_88,,PROJECT_RETROSPECTIVE.md,Y,N
day_90,Backlog & roadmap,day_89,,"BACKLOG.md,ROADMAP.md",Y,N
```

---

## USAGE INSTRUCTIONS

### Dla Project Managera:

1. **Tracking progress:**
   - Update CSV z actual completion dates
   - Mark completed days
   - Flag delayed days

2. **Impact analysis:**
   - Jeśli day X delayed: check HardDeps column - które days są blocked?
   - Use Change Impact Analysis section
   - Activate contingency plans z RISK_REGISTER.md

3. **Resource allocation:**
   - Identify parallel work opportunities
   - Assign developers do non-dependent tracks
   - Maximize parallelization

### Dla Developera:

1. **Daily work:**
   - Check HardDeps: czy wszystkie dependencies complete?
   - Jeśli NO: nie zaczynaj, work on parallel task
   - Jeśli YES: proceed

2. **Integration points:**
   - Days z many dependencies (day_13, day_81): traktuj jako integration gates
   - Extra testing needed

### Dla Stakeholdera:

1. **Progress visibility:**
   - Critical path days = najważniejsze (delays directly affect deadline)
   - Buffer days = slack (can absorb delays)

2. **Go/No-Go decisions:**
   - day_25: Phase 1 complete? Proceed do Phase 2?
   - day_55: Phase 2 complete? Proceed do Phase 3?
   - day_86: Launch checklist passed? Release?

---

**Document control:**
- Version: 1.0
- Last updated: 2026-02-06
- Next review: Co 2 tygodnie (wraz z RISK_REGISTER.md)
- Owner: Project Manager & Tech Lead
