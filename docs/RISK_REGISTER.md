# RISK REGISTER - PROJEKT DOCFLOW
## Wersja: 1.3 | Data: 2026-02-09 | Status: UPDATED - RC checklist

---

## RELATED DOCUMENTS

**Navigation:**
- **[← INDEX](INDEX.md)** - Powrót do głównego indeksu dokumentów
- **[EXTENDED_PLAN.md](EXTENDED_PLAN.md)** - Plan rozszerzony (90 dni, szczegółowy harmonogram)
- **[DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md)** - Mapa zależności (critical path, bufory)
- **[VALIDATION_REPORT.md](VALIDATION_REPORT.md)** - Raport walidacji (risk coverage 81%)
- **[FIXES_APPLIED.md](FIXES_APPLIED.md)** - Naprawione problemy (v1.0 → v1.1 changes)

**Quick links:**
- Timeline: See [EXTENDED_PLAN.md](EXTENDED_PLAN.md) for day-by-day schedule referenced in triggers
- Buffers: See [DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md) for buffer days (day_05, 20, 30, etc.) referenced in contingencies
- Changes: See [FIXES_APPLIED.md](FIXES_APPLIED.md) for M-01, M-02, M-03 fixes applied in v1.1

---

## LEGENDA

**Prawdopodobieństwo (P):**
- 1 = Bardzo niskie (< 10%)
- 2 = Niskie (10-30%)
- 3 = Średnie (30-50%)
- 4 = Wysokie (50-70%)
- 5 = Bardzo wysokie (> 70%)

**Wpływ (I):**
- 1 = Nieznaczny (opóźnienie < 2 dni, brak wpływu na funkcje)
- 2 = Mały (opóźnienie 2-5 dni, minor feature cut)
- 3 = Średni (opóźnienie 5-10 dni, moderate feature cut)
- 4 = Duży (opóźnienie 10-20 dni, major feature cut)
- 5 = Krytyczny (opóźnienie > 20 dni lub projekt abort)

**Priorytet ryzyka:**
- CRITICAL: P×I ≥ 15
- HIGH: 10 ≤ P×I < 15
- MEDIUM: 6 ≤ P×I < 10
- LOW: P×I < 6

**Status:**
- OPEN: Ryzyko aktywne, wymaga monitoringu
- TRIGGERED: Ryzyko się zmaterializowało, contingency plan aktywny
- MITIGATED: Mitigation skuteczny, ryzyko znacząco zredukowane
- CLOSED: Ryzyko już nie istnieje (minęła faza, feature usunięty, etc.)

---

## RYZYKA FAZY PRE-WORK

### R-PRE-001: Brak dostępu do korpusu szablonów dokumentacji

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-PRE-001 |
| Faza | Pre-work |
| Kategoria | Data Availability |
| Prawdopodobieństwo | 4 (Wysokie) |
| Wpływ | 5 (Krytyczny) |
| Priorytet | **CRITICAL (20)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Po 3 dniach pre-work: brak 300+ plików w testdata/ |
| Aktualny status (2026-02-09) | 829 templates w testdata → GREEN |

**Opis ryzyka:**
Projekt zakłada dostęp do 300+ szablonów dokumentacji (idealnie 10k+) do:
- Ekstrakcji wzorców sekcji (day_16-17)
- Quality scoring training (day_18-19)
- Recommendation engine testing (day_27-29)
Brak korpusu uniemożliwia budowę funkcji intelligence.

**Impact details:**
- Niemożliwe: pattern extraction, quality baseline, recommender evaluation
- Konieczne: pivot na fully synthetic data lub manual curation
- Opóźnienie: +10-20 dni (generowanie synthetic data) lub scope cut (remove recommender)

**Mitigation (proactive):**
1. **Timeline:** Pre-work day 1-2
2. **Actions:**
   - Przeszukaj publicznie dostępne repos (GitHub awesome-* lists, technical writing guides)
   - Scrape 100-200 szablonów z:
     - https://github.com/readme/guides
     - https://github.com/golang/go/wiki
     - https://github.com/microsoft/api-guidelines
   - Kontakt z partnerami: czy mają istniejące repo docs do anonimizacji
3. **Success criteria:** Min. 100 real templates do końca pre-work day 2
4. **Responsible:** Tech Lead + Data Engineer (jeśli dostępny)

**Contingency (reactive):**
1. **Scenario A: Mamy 50-100 templates**
   - Pivot: Użyj LLM (GPT-4) do generowania 200 synthetic templates
   - Process:
     - Zdefiniuj 5 rodzin docs (api, architecture, requirements, guide, runbook)
     - Generate 40 templates per rodzina z różnymi quality levels
     - Manual review 30 templates jako ground truth
   - Timeline: +5 dni (generation + review)
   - Cost: ~$50-100 (LLM API costs)

2. **Scenario B: Mamy < 50 templates**
   - Pivot: Zmniejsz zakres projektu
   - Remove features:
     - Pattern extraction (day_16-17) → manual schemas only
     - Quality scoring automated (day_18-19) → manual quality labels only
     - Recommender (day_27-29) → simple rule-based filter (no scoring)
   - Timeline: -10 dni (features removed), ale ograniczona funkcjonalność
   - Dokumentuj limitation w KNOWN_LIMITATIONS.md

3. **Scenario C: Brak templates w ogóle**
   - **GO/NO-GO DECISION:** Przerwij projekt lub pivot na inny use case
   - Alternatywne use cases:
     - Tool dla single-team docs (nie multi-template)
     - Simple validator only (no intelligence)
   - Timeline: Stop projekt lub restart z nowym scope

**Monitoring:**
- Checkpoint: Pre-work day 2 - count templates in testdata/
- Red flag: < 50 templates → activate contingency
- Green: ≥ 100 templates → continue as planned

**Dependencies:**
- Blocks: day_16-17 (pattern extraction), day_18-19 (quality scoring), day_27-29 (recommender)
- Critical path: YES

---

### R-PRE-002: Wybrana technologia/biblioteka nie wspiera wymaganych funkcji

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-PRE-002 |
| Faza | Pre-work |
| Kategoria | Technical |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (9)** |
| Status | OPEN |
| Aktualny status (2026-02-09) | libs działają na obecnym scope; brak nowych blockerów |
| Owner | Tech Lead |
| Trigger | Day_00-05: Biblioteka nie parsuje edge case lub brak API dla feature X |

**Opis ryzyka:**
Wybrane biblioteki (Markdown parser, YAML parser, graph) mogą nie wspierać:
- Markdown extensions (tables, frontmatter, footnotes)
- Complex YAML (anchors, merge keys)
- Large graphs (10k+ nodes)

**Impact details:**
- Konieczność zamiany biblioteki → przepisanie kodu (2-5 dni)
- Lub workaround (wrapper, preprocessing) → tech debt
- Lub feature cut (nie wspieramy tej funkcjonalności)

**Mitigation (proactive):**
1. **Pre-work day 3: Spike testing**
   - Create 10 edge case test files:
     - Markdown: tables, nested lists, code blocks w różnych językach, frontmatter variants
     - YAML: anchors, merge keys, unicode, very deep nesting (10 levels)
     - Graph: 1000 nodes, cycles, disconnected components
   - Test każdą candidate library na edge cases
   - Document compatibility matrix w TECH_STACK.md
2. **Fail fast:** Jeśli library fails edge case → reject, wybierz alternatywę
3. **Fallback plan:** Zidentyfikuj 2nd choice library dla każdego komponentu

**Contingency (reactive):**
1. **Scenario: Library fails w day_01-10**
   - Immediate switch do 2nd choice library (zidentyfikowanej w pre-work)
   - Przepisz kod (estimated 2-3 dni)
   - Re-test edge cases
   - Update bufory: użyj day_05 (buffer) + day_20 (buffer)

2. **Scenario: Brak dobrej biblioteki (wszystkie fail)**
   - Implement minimal wrapper:
     - Dla Markdown: regex fallback (support tylko headings + frontmatter, no tables)
     - Dla YAML: strict subset (no anchors)
     - Dla Graph: simple adjacency list (no fancy algorithms)
   - Document limitations
   - Timeline: +3 dni (wrapper implementation)

**Monitoring:**
- Pre-work day 3: Edge case test results
- Day_01-10: Track library issues (file bugs, check community activity)
- Red flag: >3 blocker issues w 1 tygodniu → consider switch

**Dependencies:**
- Blocks: Potentially all parsing/graph features
- Critical path: NO (buffers available)

---

### R-PRE-003: Brak team capacity (resource unavailability)

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-PRE-003 |
| Faza | Wszystkie |
| Kategoria | Resources |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 4 (Duży) |
| Priorytet | **HIGH (12)** |
| Status | OPEN |
| Owner | Project Manager |
| Trigger | Developer unavailable >5 dni, lub <50% availability przez >2 tygodnie |

**Opis ryzyka:**
Plan zakłada 1-2 senior developers full-time. Ryzyka:
- Choroba, urlop, turnover
- Konflikt priorytetów (inne projekty)
- Underestimated effort (developer overloaded)

**Impact details:**
- Opóźnienie proporcjonalne do niedostępności
- Jeśli 1 developer i unavailable >10 dni: projekt stalled
- Jeśli capacity <50%: timeline 2x

**Mitigation (proactive):**
1. **Pre-work: Team allocation contract**
   - Dokument: TEAM_ALLOCATION.md
   - Explicit commitment: 1 developer @ 100% FTE, 90 dni
   - Backup: 2nd developer @ 50% FTE (pair programming, code review)
   - Upfront blocker identification: known vacations, conflicting projects
2. **Knowledge sharing:**
   - Pair programming na critical modules (graph, parser)
   - Daily standups (async: Slack update)
   - Code review mandatory (2nd developer reviews all PRs)
3. **Documentation:**
   - Architecture decisions zapisywane (ADR)
   - Code comments dla complex logic
   - Każdy moduł ma README

**Contingency (reactive):**
1. **Scenario: Short unavailability (1-5 dni)**
   - Use buffers: 13 dni buforowych w planie
   - Re-prioritize: focus na critical path
   - Skip nice-to-haves (analytics HTML export, etc.)

2. **Scenario: Medium unavailability (5-15 dni)**
   - Extend timeline: +1 week per 5 days lost
   - Use Phase buffers (każda faza ma buffer na końcu)
   - Re-negotiate scope z stakeholders: defer Phase 3 features

3. **Scenario: Long unavailability (>15 dni) lub turnover**
   - **GO/NO-GO DECISION:**
     - Option A: Pause projekt, resume gdy resource available
     - Option B: Hire replacement (ramp-up: 2-3 tygodnie)
     - Option C: Drastyczny scope cut:
       - Deliver tylko Phase 1 (MVP)
       - Phase 2-3 deferred indefinitely
   - Decision criteria:
     - Business value MVP vs full product
     - Cost of delay vs cost of replacement
     - Stakeholder patience

**Monitoring:**
- Weekly capacity check: actual hours vs planned
- Red flag: <80% capacity 2 tygodnie z rzędu → activate contingency
- Trigger: Developer gives notice (turnover) → immediate GO/NO-GO meeting

**Dependencies:**
- Blocks: Cały projekt
- Critical path: YES

---

## RYZYKA FAZY 1 (FOUNDATION & MVP)

### R-F1-001: Benchmark wydajnościowy pokazuje nieakceptowalne wyniki

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F1-001 |
| Faza | Phase 1 |
| Kategoria | Performance |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (6)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_00 benchmark: scan 300 files > 10s (target: <5s) |

**Opis ryzyka:**
Day_00 benchmark może pokazać, że wybrany język/biblioteka jest zbyt wolny.
Target: 300 plików < 5s (95th percentile)
Acceptable: < 10s
Unacceptable: > 10s

**Impact details:**
- Jeśli 10-20s: Użytkownik experience degraded, ale workable
- Jeśli > 20s: Narzędzie unusable dla dużych projektów
- Konieczność optymalizacji lub zmiany języka

**Mitigation (proactive):**
1. **Pre-work day 3: Wybór szybkich bibliotek**
   - Native parsers (nie regex-based)
   - Streaming parsers dla large files
2. **Day_00: Proper benchmarking methodology**
   - Warm-up runs (ignore 1st run)
   - Median z 5 runs (nie average - reduces outlier impact)
   - Profile z CPU profiler (identify bottlenecks)

**Contingency (reactive):**
1. **Scenario: 5-10s (acceptable, ale nie target)**
   - Defer optimization do day_49-50 (optimization sprint)
   - Document performance jako known limitation
   - Continue z planem

2. **Scenario: 10-20s (degraded)**
   - Immediate optimization (use day_05 buffer):
     - Parallel parsing (goroutines/threads)
     - Lazy loading (parse tylko metadata, not full content)
     - Caching
   - Re-benchmark
   - Jeśli nadal >10s: proceed, note as P1 for Phase 2

3. **Scenario: >20s (unacceptable)**
   - **CRITICAL:** Consider language change
   - Options:
     - Switch do Go (jeśli był Python) - fast compile, good concurrency
     - Optimize inner loops (profiling-guided)
     - Reduce feature scope (parse tylko headings, not full AST)
   - Timeline: +5 dni (re-implementation) lub scope cut
   - GO/NO-GO decision required

**Monitoring:**
- Day_00: Benchmark results
- Day_23: Performance testing (300 files)
- Day_76-77: Scale testing (1000+ files)

**Dependencies:**
- Blocks: Day_23 (performance testing), day_76-77 (scale testing)
- Critical path: NO (można kontynuować z degraded performance, optymalizować później)

---

### R-F1-002: Graf zależności w rzeczywistych danych zawiera cykle

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F1-002 |
| Faza | Phase 1 |
| Kategoria | Data Quality |
| Prawdopodobieństwo | 4 (Wysokie) |
| Wpływ | 2 (Mały) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_09-10: Graph builder wykrywa cycle w >10% dokumentów |

**Opis ryzyka:**
Dokumentacja może zawierać circular dependencies (A→B→C→A).
To nie jest bug w narzędziu, ale problem w danych.

**Impact details:**
- Cycle detection działa (nie blokuje narzędzia)
- Ale: topological sort niemożliwy (no valid order)
- User musi ręcznie naprawić dane (usunąć dependency edge)
- Jeśli >50% docs w cycles: projekt może być unworkable

**Mitigation (proactive):**
1. **Day_01: Spec dependencies clearly**
   - Define: co to jest valid dependency (only forward in layers)
   - Example graph patterns (good vs bad)
   - Educate users w dokumentacji
2. **Day_09-10: Robust cycle detection**
   - Algorytm: strongly connected components (SCC)
   - Raport: każdy cycle z doc_ids i edge path
   - Sugestie: które edge usunąć (heuristic: newest dependency)

**Contingency (reactive):**
1. **Scenario: <10% docs w cycles**
   - Raport cycles do użytkownika
   - Dokumentacja: "How to fix cycles"
   - User manually fixes (remove 1 edge per cycle)
   - Continue

2. **Scenario: 10-50% docs w cycles**
   - Investigate root cause:
     - Czy to bidirectional dependencies? (A needs B for context, B needs A for context)
     - Czy to layers problem? (higher layer depends on lower, but also vice versa)
   - Solutions:
     - Introduce `context_sources` (soft dependency, nie tworzy edge)
     - Relaxed validation: allow cycles w context_sources, not depends_on
   - Timeline: +2 dni (add context_sources semantics)

3. **Scenario: >50% docs w cycles**
   - **CRITICAL:** Data model issue
   - Root cause analysis (1 dzień):
     - Interview stakeholders: dlaczego cycles?
     - Review use cases: czy depends_on ma złą semantykę?
   - Potential fixes:
     - Redefine dependency model: hierarchical layers (no cross-layer)
     - Introduce dependency types (compilation_order vs logical_dependency)
     - Allow cycles, ale provide "best effort" ordering (SCC condensation)
   - Timeline: +5 dni (re-design dependency model, update spec, re-implement)
   - GO/NO-GO decision: jeśli cycles są fundamental, może MVP nie jest viable

**Monitoring:**
- Day_09-10: Cycle detection results
- Day_13: MVP integration test (test na small dataset z known cycles)

**Dependencies:**
- Blocks: Day_13 (daily planner - requires topo order)
- Critical path: PARTIAL (can report cycles, user fixes data)

---

### R-F1-003: MVP integration test fails - core features nie działają end-to-end

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F1-003 |
| Faza | Phase 1 |
| Kategoria | Integration |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 4 (Duży) |
| Priorytet | **HIGH (12)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_13: E2E test fails, pipeline broken |

**Opis ryzyka:**
Poszczególne moduły działają (unit tests pass), ale integracja fails:
- Incompatible data formats między modułami
- State management issues (cache inconsistency)
- CLI command chaining broken

**Impact details:**
- MVP release (day_15) zagrożony
- Konieczność debugowania i refactoringu (3-7 dni)
- Może wymagać re-design interfejsów między modułami

**Mitigation (proactive):**
1. **Day_01: Define API contracts**
   - Document: API_CONTRACTS.md
   - Każdy moduł ma defined input/output types
   - JSON schemas dla cache formats
2. **Day_05, 10: Incremental integration testing**
   - Nie czekaj do day_13
   - Day_05: Test parser→index integration
   - Day_10: Test index→graph integration
3. **Daily builds:**
   - CI pipeline (setup day_05)
   - Build + basic smoke test każdego dnia
   - Fail fast jeśli integration broken

**Contingency (reactive):**
1. **Scenario: Minor issues (1-2 bugs)**
   - Fix bugs w day_14-15 (compress documentation work)
   - Alternatively: use day_20 buffer (dedicated buffer day)
   - Re-test
   - Delay MVP release do day_17 if needed

2. **Scenario: Major issues (5+ bugs, design problems)**
   - Extend Phase 1:
     - Add 1 week (5 dni) for debugging + refactoring
     - Use day_20 buffer (dedicated buffer)
     - MVP release: day_22 (delayed 1 week)
   - Root cause analysis:
     - Contract violations: enforce contracts z validation
     - State issues: introduce state machine or immutability
   - Risk mitigation dla Phase 2: więcej integration tests wcześniej

3. **Scenario: Critical failure (pipeline completely broken)**
   - **GO/NO-GO DECISION:**
     - Stop and re-design architecture (5-10 dni)
     - Lub: pivot do simple tool (no integration, separate commands)
   - Investigate:
     - Czy architecture jest fundamentalnie błędny?
     - Czy assumptions o data flow są niepoprawne?
   - Options:
     - Re-design (use ARCHITECTURE.md revision)
     - Simplify (fewer features, simpler integration)
     - Abort Phase 1, learn lessons

**Monitoring:**
- Daily: CI build status
- Day_05: Parser→Index integration checkpoint
- Day_10: Index→Graph integration checkpoint
- Day_13: Full E2E test

**Dependencies:**
- Blocks: Day_15 (MVP release)
- Critical path: YES

---

### R-F1-004: Dogfooding odkrywa 10+ critical bugs w MVP

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F1-004 |
| Faza | Phase 1 |
| Kategoria | Quality |
| Prawdopodobieństwo | 4 (Wysokie) |
| Wpływ | 2 (Mały) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_13 dogfooding: >10 bugs filed, >3 critical |

**Opis ryzyka:**
Dogfooding (używanie docflow na własnej dokumentacji projektu) może odkryć:
- Edge cases nie pokryte testami
- UX issues (komendy confusing, error messages unclear)
- Performance issues na realnych danych

**Impact details:**
- 10 bugs: normalne, expected
- >3 critical bugs: może opóźnić release
- Critical = blocker dla podstawowego use case (np. scan fails na valid markdown)

**Mitigation (proactive):**
1. **Day_01-12: Comprehensive unit tests**
   - Coverage target: 70%
   - Edge case tests (day_43-44 content moved earlier)
2. **Day_13: Structured dogfooding**
   - Testuj wszystkie komendy systematycznie
   - Document test scenarios przed dogfooding
   - Expected: find 5-10 minor bugs, 0-2 critical

**Contingency (reactive):**
1. **Scenario: 5-10 minor bugs, 0-1 critical**
   - Expected scenario
   - Fix w day_14-15 (compress documentation work to make time)
   - Release on time (day_15)

2. **Scenario: >10 bugs, 2-3 critical**
   - Triage:
     - Fix critical immediately (day_14)
     - Minor bugs: backlog dla v0.1.1
   - Use day_20 buffer (dedicated buffer) jeśli potrzeba więcej czasu
   - Release: day_15 z known minor issues (documented w KNOWN_ISSUES.md)

3. **Scenario: >5 critical bugs**
   - **Delay release:**
     - Extend Phase 1 o 1 tydzień
     - Deep dive: dlaczego tests nie wykryły? (improve test strategy)
     - Fix wszystkie critical
     - Re-dogfood
     - Release: day_22
   - Risk dla Phase 2: może być więcej ukrytych bugs, increase test effort

**Monitoring:**
- Day_13: Bug count after dogfooding
- Severity: Track critical vs minor
- Trend: Czy bugs clustered w jednym module? (indicates quality issue)

**Dependencies:**
- Blocks: Day_15 (MVP release)
- Critical path: PARTIAL (can release z minor bugs)

---

### R-F1-005: MVP release delays - stakeholder rejection lub demo failure

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F1-005 |
| Faza | Phase 1 |
| Kategoria | Stakeholder Management |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (6)** |
| Status | OPEN |
| Owner | Project Manager |
| Trigger | Day_15: Stakeholders reject MVP lub demo fails critically |

**Opis ryzyka:**
MVP demo (day_15) może nie spełnić oczekiwań stakeholders:
- Core features działają, ale UX jest confusing
- Demo technicznie successful, ale business value niejasny
- Stakeholders oczekiwali więcej features
- Live demo failure (environment issues, bugs pod prezencją)

**Impact details:**
- MVP release delayed (minimum day_17, możliwe day_20)
- Morale zespołu obniżone
- Konieczność pivot scope lub re-demo
- Phase 2 opóźnione o 2-5 dni

**Mitigation (proactive):**
1. **Pre-work + day_01: Align expectations**
   - Document: REQUIREMENTS.md z clear MVP scope
   - Stakeholder sign-off przed rozpoczęciem
   - Set realistic expectations: "MVP = core features only, not polished"
2. **Day_13: Internal demo (dry run)**
   - Demo dla team przed stakeholder demo
   - Identify UX issues, fix before day_15
   - Practice presentation, prepare backup plans
3. **Day_14: Demo preparation**
   - Prepare demo script (happy path)
   - Test environment thoroughly
   - Record backup video (jeśli live demo fails)
   - Prepare FAQ (anticipated questions)

**Contingency (reactive):**
1. **Scenario: Minor stakeholder concerns (UX improvements, missing docs)**
   - Acknowledge feedback
   - Quick fixes w day_16-17:
     - Improve error messages
     - Add missing documentation sections
     - Polish CLI output
   - Re-demo (internal) day_17
   - Proceed z Phase 2 (minor delay: 2 dni)

2. **Scenario: Major stakeholder rejection (wrong features, business value unclear)**
   - **Pause and re-align:**
     - Emergency meeting z stakeholders (day_16)
     - Clarify requirements mismatch
     - Options:
       - A) Pivot MVP scope (add critical missing feature, remove less important)
       - B) Better explain value proposition (prepare business case)
       - C) Accept feedback, plan v0.1.1 z improvements
     - Timeline: day_16-18 (re-work), re-demo day_19
     - Use day_20 buffer
     - Phase 2 start: day_21 (delayed 1 week)

3. **Scenario: Demo technical failure (critical bug discovered during demo)**
   - **Immediate recovery:**
     - Switch to backup video recording (if prepared)
     - Acknowledge bug, show mitigation plan
     - Schedule re-demo (day_16-17 after fix)
   - Post-demo:
     - Root cause analysis (why bug not caught in testing?)
     - Fix critical bug (day_16)
     - Additional testing (prevent future demo failures)
     - Re-demo day_17
   - Timeline: +2 dni delay

**Monitoring:**
- Day_14: Demo dry run feedback
- Day_15: Stakeholder satisfaction (poll after demo)
- Red flag: Stakeholders silent lub negative feedback → activate contingency

**Dependencies:**
- Blocks: Phase 2 start (day_26)
- Critical path: PARTIAL (delays affect Phase 2 timeline)

---

## RYZYKA FAZY 2 (INTELLIGENCE & AUTOMATION)

### R-F2-001: Recommender precision@5 < 0.60 (niska jakość rekomendacji)

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F2-001 |
| Faza | Phase 2 |
| Kategoria | Algorithm Quality |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (9)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_29: Evaluation shows precision@5 < 0.60 |

**Opis ryzyka:**
Recommendation engine może mieć niską trafność:
- Scoring algorithm źle dobrany (wagi nieoptymalne)
- Training data za małe lub biased
- Ground truth labels niepoprawne

**Impact details:**
- Precision@5 < 0.60: Users nie ufają recommendations
- Feature ma niską wartość biznesową
- Może wymagać re-design algorytmu (3-5 dni)

**Mitigation (proactive):**
1. **Day_27-28: Research-backed algorithm**
   - Użyj proven techniques (collaborative filtering, content-based)
   - Start z simple baseline (rule-based)
   - Iterate z scoring factors
2. **Day_29: Proper evaluation methodology**
   - Ground truth: 20+ scenarios z expert labels
   - Cross-validation jeśli masz usage data
   - Multiple metrics (precision@K, MRR, NDCG)
3. **Pre-work: High-quality ground truth**
   - Manual labels od domain experts (tech writers)
   - Diverse scenarios (różne rodziny, warstwy, languages)

**Contingency (reactive):**
1. **Scenario: 0.50 < Precision@5 < 0.60**
   - Moderate quality, może być akceptowalne
   - Tuning (day_30 buffer):
     - Adjust wagi w ALGO_PARAMS.md
     - Re-evaluate
     - Iterate 2-3x
   - Jeśli improvement: proceed
   - Jeśli no improvement: accept 0.50-0.60, document limitation

2. **Scenario: Precision@5 < 0.50**
   - **Poor quality**, nie akceptowalne
   - Root cause analysis (1 dzień):
     - Inspect failures: dlaczego wrong templates recommended?
     - Data issue? (labels wrong, templates mislabeled)
     - Algorithm issue? (wagi completely off)
   - Options:
     - A) Re-design algorithm (3 dni):
       - Try different scoring (cosine similarity of sections)
       - Add features (word embeddings, TF-IDF)
     - B) Simplify algorithm (2 dni):
       - Rule-based tylko (exact rodzina+warstwa match)
       - No scoring, just filter
     - C) Cut feature (0 dni):
       - Remove recommender, defer do Phase 3 lub v1.1
   - Decision: Cost/benefit analysis (effort vs value)

3. **Scenario: Cannot evaluate (brak ground truth)**
   - Fallback: Expert review (manual)
   - 5 scenarios, domain expert ranks recommendations
   - Qualitative assessment: "good enough" vs "useless"
   - Proceed z caveat: no quantitative validation

**Monitoring:**
- Day_29: Precision@5 metric
- Day_55: Phase 2 retrospective - user feedback na recommender

**Dependencies:**
- Blocks: Phase 2 feature completeness
- Critical path: NO (can cut feature)

---

### R-F2-002: Pattern extraction nie znajduje spójnych wzorców (chaos w danych)

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F2-002 |
| Faza | Phase 2 |
| Kategoria | Data Quality |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 2 (Mały) |
| Priorytet | **MEDIUM (6)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_16: Analiza pokazuje 0 patterns z frequency >5, wszystkie templates unique |

**Opis ryzyka:**
Korpus szablonów może być zbyt chaotyczny:
- Każdy template unikalny (no common sections)
- Brak standaryzacji w naming (10 variants of "Overview")
- Rozdrobnione: 300 templates, 300 distinct patterns

**Impact details:**
- Pattern extraction bezużyteczny (no reusable patterns)
- Section schemas muszą być ręcznie stworzone
- Opóźnienie: +3 dni (manual schema creation)

**Mitigation (proactive):**
1. **Pre-work day 2: Analyze corpus diversity**
   - Manual inspection: 20 random templates
   - Estimate: czy są common sections? (Overview, Examples, etc.)
   - Jeśli chaos: prepare manual schemas upfront
2. **Day_16: Fuzzy matching dla section names**
   - Normalize titles (lowercase, remove punctuation)
   - Aliases: "Overview" = "Przegląd" = "Summary"
   - Levenshtein distance dla variants

**Contingency (reactive):**
1. **Scenario: Low pattern frequency (best pattern: freq=3-5)**
   - Acceptable, ale not ideal
   - Use top 10 patterns (even jeśli low freq)
   - Supplement z manual schemas (2 dni)
   - Continue

2. **Scenario: No patterns (all unique)**
   - **Pivot: Manual schema creation**
   - Process (3 dni):
     - Domain expert defines 5-7 rodzin
     - Per rodzina: manual schema (required/optional sections)
     - Use examples z corpus jako reference
   - Skip pattern extraction (day_16-17 freed)
   - Use freed time jako buffer

3. **Scenario: Patterns exist, ale wrong (algorytm issue)**
   - Debug algorithm (1 dzień):
     - Adjust similarity threshold (0.80 → 0.70)
     - Change n-gram size (3 → 2)
   - Re-run extraction
   - If still wrong: fallback do manual (scenario 2)

**Monitoring:**
- Pre-work day 2: Corpus analysis
- Day_16: Pattern extraction results (frequency distribution)

**Dependencies:**
- Blocks: Day_17 (schema generation)
- Critical path: NO (manual schemas są workaround)

---

### R-F2-003: Daily planner generuje nierealistyczne plany (effort estimation off)

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F2-003 |
| Faza | Phase 2 |
| Kategoria | Algorithm Quality |
| Prawdopodobieństwo | 4 (Wysokie) |
| Wpływ | 2 (Mały) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_33-34: Plans systematycznie over/underestimate effort (>50% error) |

**Opis ryzyka:**
Effort estimation jest heuristic (brak historical data):
- Może być systematycznie wrong (wszystkie estimates 2x too low)
- Plany generowane są unusable (7h estimated, reality 14h)

**Impact details:**
- Feature ma niską wartość (users ignore planner)
- Wymaga tuning heuristics (1-2 dni)
- Lub wymaga historical data collection (no immediate fix)

**Mitigation (proactive):**
1. **Day_33: Conservative estimates**
   - Baseline effort deliberately higher (avoid underestimation)
   - Add buffers (+20% dla uncertainty)
2. **Day_33: Configurable heuristics**
   - User może override defaults w config:
     ```yaml
     effort_heuristics:
       rodzina:
         api: 3h
         architecture: 5h
     ```
   - Document: "Tune based on your team's velocity"

**Contingency (reactive):**
1. **Scenario: Systematic over-estimation (estimates 2x too high)**
   - Not critical (better than under-estimation)
   - Tune heuristics (reduce baseline by 30%)
   - Re-test (use real data z Phase 1-2 development)
   - Update ALGO_PARAMS.md

2. **Scenario: Systematic under-estimation (estimates 2x too low)**
   - More problematic (plans unachievable)
   - Quick fix (day_35 buffer):
     - Increase baseline by 50%
     - Add +2h buffer per document
     - Re-evaluate
   - Medium-term (Phase 3):
     - Collect actual effort data (time tracking)
     - Build model z historical data
     - Replace heuristics

3. **Scenario: High variance (some right, some 5x off)**
   - Identify patterns: które rodziny są off?
   - Targeted tuning (adjust per rodzina)
   - Jeśli no pattern: too much uncertainty, feature limited value
   - Option: Remove effort estimation, planner tylko orders by dependencies

**Monitoring:**
- Day_34: Test planner na Phase 1 work (retrospective: actual effort vs heuristic)
- Day_55: Phase 2 retrospective - planner usage & feedback

**Dependencies:**
- Blocks: Nic (nice-to-have feature)
- Critical path: NO

---

## RYZYKA FAZY 3 (PRODUCTION READY)

### R-F3-001: Security audit odkrywa critical vulnerability

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-001 |
| Faza | Phase 3 |
| Kategoria | Security |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 4 (Duży) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_73: Security audit finds critical vulnerability (RCE, path traversal, etc.) |

**Opis ryzyka:**
Security issues mogą być discovered late:
- Path traversal (docflow scan ../../etc/passwd)
- YAML bomb (huge nested YAML crashes parser)
- Command injection (w CLI args)
- Dependency vulnerability (library CVE)

**Impact details:**
- Critical vulnerability: Cannot release v1.0 until fixed
- Opóźnienie: 2-10 dni (depending on severity)
- Może wymagać significant refactoring (security by design)

**Mitigation (proactive):**
1. **Day_01-10: Secure coding practices**
   - Input validation (wszystkie CLI args, file paths)
   - Path sanitization (resolve symlinks, check bounds)
   - Resource limits (max file size, max parse depth)
2. **Day_05: Dependency scanning**
   - CI: `npm audit` / `cargo audit` / `go mod verify`
   - Auto-alerts na CVEs
3. **Day_73: Thorough security audit**
   - Checklist: OWASP Top 10
   - Fuzzing (jeśli czas pozwala)
   - Manual code review (security focus)

**Contingency (reactive):**
1. **Scenario: Low/Medium severity issue**
   - Fix w day_74-75 (buffer available)
   - Re-audit
   - Proceed z release

2. **Scenario: High severity (not critical)**
   - Delay release o 3-5 dni
   - Fix vulnerability
   - Add regression test
   - Re-audit
   - Release: day_90-92

3. **Scenario: Critical severity (RCE, data loss)**
   - **Delay release:** +1-2 tygodnie
   - Full security review:
     - Fix immediate issue
     - Review entire codebase dla similar patterns
     - Add comprehensive security tests
     - External security review (jeśli budget allows)
   - GO/NO-GO decision:
     - If fix requires major refactoring: consider v1.0-beta release (limited audience)
     - Full v1.0 deferred until security hardened

**Monitoring:**
- Day_05, 20, 40, 60: Dependency vulnerability scans
- Day_73: Security audit results
- Pre-release: Final vulnerability scan

**Dependencies:**
- Blocks: v1.0 release (day_87)
- Critical path: YES (security blocker)

---

### R-F3-002: Scale testing pokazuje performance collapse na 5000+ docs

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-002 |
| Faza | Phase 3 |
| Kategoria | Performance |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (9)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_76-77: Scan 5000 docs takes >5 minutes lub OOM error |

**Opis ryzyka:**
Performance może degrade non-linearly:
- 300 docs: 5s (OK)
- 1000 docs: 20s (OK)
- 5000 docs: 300s (NOT OK) lub crash
Algorithmic complexity issues (O(n²) graph operations)

**Impact details:**
- Tool unusable dla large projects (5k+ docs jest realistic dla enterprise)
- Need optimization (3-7 dni)
- Lub document scale limitations (max 1000 docs)

**Mitigation (proactive):**
1. **Day_09-10: Use efficient algorithms**
   - Graph: adjacency list (not matrix)
   - Topo sort: O(V+E) Kahn's algorithm
   - Hash maps (not linear search)
2. **Day_46-47: Incremental parsing**
   - Parse tylko changed files
   - Memoization
3. **Day_76: Gradual scale testing**
   - Test 1000, 2000, 3000, 5000 progressively
   - Identify inflection point (gdzie performance drops)

**Contingency (reactive):**
1. **Scenario: Gradual degradation (5000 docs: 2-3 min)**
   - Not ideal, ale workable
   - Optimize w day_78-79 (optimization sprint):
     - Parallel parsing
     - Lazy evaluation
     - Index optimization
   - Target: <1 min dla 5000 docs
   - If achieved: proceed
   - If not: document limitation (recommended max: 2000 docs)

2. **Scenario: Performance collapse (>5 min lub crash)**
   - **Critical optimization needed**
   - Profile (1 dzień):
     - Identify bottleneck (graph build? parsing? validation?)
     - CPU profile, memory profile
   - Fix (2-4 dni):
     - Algorithmic fix (replace O(n²) z O(n log n))
     - Memory fix (streaming, chunking)
     - Incremental computation
   - Re-test
   - Timeline: use day_80 buffer (dedicated buffer) + extend optimization sprint (day_78-79) jeśli needed

3. **Scenario: Fundamental scalability limit (cannot fix w reasonable time)**
   - **Document limitation:**
     - Add to KNOWN_LIMITATIONS.md: "Recommended max: 1000 docs"
     - Add validation: warning jeśli >1000 docs scanned
   - Roadmap v2.0: Re-architecture dla scalability
   - Proceed z release (limitation documented)

**Monitoring:**
- Day_23: Performance baseline (300 docs)
- Day_76-77: Scale testing (1000, 2000, 5000)
- Track: time complexity (linear? quadratic?)

**Dependencies:**
- Blocks: v1.0 release jeśli collapse is severe
- Critical path: PARTIAL (can release z documented limitations)

---

### R-F3-003: User documentation incomplete - users nie mogą onboard

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-003 |
| Faza | Phase 3 |
| Kategoria | Documentation |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (6)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_84: Documentation review finds critical gaps (no quickstart, unclear concepts) |

**Opis ryzyka:**
Dokumentacja może być:
- Niekompletna (brakujące komendy)
- Niejasna (założenia eksperckiej wiedzy)
- Outdated (nie sync z kodem)
Users nie mogą onboard bez help.

**Impact details:**
- Poor adoption (users frustrated, abandon tool)
- Support burden (many questions)
- Delay release do uzupełnienia docs (2-5 dni)

**Mitigation (proactive):**
1. **Day_66-69: Comprehensive documentation sprint**
   - 4 dni dedicated do docs (wystarczająco)
   - User testing: fresh user follows quickstart (identify gaps)
2. **Day_14, 24: Incremental documentation**
   - MVP docs (day_14-15)
   - Phase 1 docs (day_24)
   - Każdy feature dokumentowany when shipped
3. **Examples & tutorials (day_68)**
   - Working examples > długie opisy

**Contingency (reactive):**
1. **Scenario: Minor gaps (1-2 missing sections)**
   - Fix w day_85 buffer
   - Quick additions (1 dzień)
   - Proceed

2. **Scenario: Major gaps (no quickstart, concepts unclear)**
   - Extend documentation sprint (day_85 + extra 2-3 dni)
   - User testing: 2 fresh users try onboarding
   - Iterate based on feedback
   - Delay release: day_89-90

3. **Scenario: Documentation fundamentally confusing**
   - **Re-write needed** (5-7 dni)
   - Hire technical writer (jeśli budget)
   - Lub: Community beta (release v1.0-beta, gather feedback, improve docs)
   - Full v1.0 release: +2 tygodnie

**Monitoring:**
- Day_68: Fresh user testing (example walkthroughs)
- Day_84: Documentation review checklist
- Pre-release: Final user test (cold start)

**Dependencies:**
- Blocks: v1.0 release (poor docs = poor product)
- Critical path: YES (but low probability)

---

### R-F3-004: CI/CD pipeline failures blokują releases

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-004 |
| Faza | Phase 3 |
| Kategoria | Infrastructure |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 2 (Mały) |
| Priorytet | **LOW (4)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_87: Release automation fails, cannot build binaries |

**Opis ryzyka:**
CI/CD może fail:
- Cross-compilation issues (Windows binary not building on Linux CI)
- GitHub Actions quota exceeded
- Signing keys missing
- Flaky tests fail release gate

**Impact details:**
- Delay release o 1-2 dni (debug CI)
- Manual release process (more effort, error-prone)

**Mitigation (proactive):**
1. **Day_61-62: Robust CI pipeline**
   - Test all platforms (Linux, macOS, Windows)
   - Deterministic builds (pinned dependencies)
   - Retry logic dla flaky tests
2. **Day_63-64: Test release process**
   - Dry-run release (v0.9.9-test)
   - Verify binaries work on all platforms
3. **Day_86: Pre-release checklist**
   - Wszystkie CI checks green
   - Test release artifacts ręcznie

**Contingency (reactive):**
1. **Scenario: CI fails, ale known issue**
   - Temporary fix (skip failing test, fix later)
   - Manual build jeśli CI completely broken
   - Timeline: +1 dzień (manual release)

2. **Scenario: CI fails, unknown issue**
   - Debug (use day_88 buffer)
   - Fix CI
   - Re-run release
   - Timeline: +1-2 dni

3. **Scenario: Cannot fix CI w reasonable time**
   - **Manual release:**
     - Build binaries locally (3 platforms)
     - Manual upload do GitHub
     - Manual CHANGELOG generation
   - Timeline: +1 dzień (tedious, ale works)
   - Post-release: fix CI (dla v1.0.1)

**Monitoring:**
- Day_61-62: CI pipeline setup verification
- Day_64: Dry-run release test
- Day_86: Pre-release CI health check

**Dependencies:**
- Blocks: v1.0 release automation
- Critical path: NO (manual release possible)

---

### R-F3-005: Release Candidate exceeds bug threshold - quality below acceptable

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-005 |
| Faza | Phase 3 |
| Kategoria | Quality |
| Prawdopodobieństwo | 3 (Średnie) |
| Wpływ | 3 (Średni) |
| Priorytet | **MEDIUM (9)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Day_81: RC testing finds >5 critical bugs lub >20 total bugs |

**Opis ryzyka:**
Release Candidate (day_81) może mieć więcej bugs niż expected:
- Bug bash (day_81) odkrywa 5+ critical bugs
- Total bug count >20 (indicates quality issues)
- Critical bugs w core features (scan, validate, generate)
- Regression bugs (features broke that worked in Phase 1-2)

**Impact details:**
- Cannot release v1.0 until bugs fixed
- Delay release: minimum 3-5 dni (fix critical bugs)
- Możliwe 10+ dni jeśli regressions wymagają re-design
- Launch date (day_87) opóźniony

**Mitigation (proactive):**
1. **Day_21-23: Comprehensive testing early**
   - Unit tests: 70% coverage
   - Integration tests: all major flows
   - Performance tests: establish baseline
   - Prevent bugs from reaching RC
2. **Day_71-74: E2E + security + chaos testing**
   - Catch bugs przed RC
   - Test on 3 platforms
   - Edge cases, error scenarios
3. **Day_81: Structured bug bash**
   - Clear test scenarios (checklist)
   - Multiple testers (team + stakeholders)
   - Triage immediately (critical vs minor)
   - Set expectations: <3 critical bugs acceptable

**Contingency (reactive):**
1. **Scenario: 3-5 critical bugs, <15 total bugs**
   - Acceptable range (minor delay)
   - Fix critical bugs immediately (day_82-83)
   - Minor bugs: backlog dla v1.0.1
   - Use day_85 buffer
   - Release: day_87-88 (1-2 dni delay)

2. **Scenario: 5-10 critical bugs, 15-30 total bugs**
   - **Quality concerns**
   - Triage (day_82):
     - P0 (blockers): must fix before release
     - P1 (high): fix if time allows
     - P2 (medium/low): backlog
   - Extend bug fixing phase:
     - Fix P0 bugs (day_82-84)
     - Re-test (day_85)
     - Use day_85 buffer + extend 2-3 dni
   - Release: day_89-90 (delayed 1 week)
   - Consider: v1.0-beta release instead (limited audience)

3. **Scenario: >10 critical bugs lub major regressions**
   - **CRITICAL QUALITY ISSUE**
   - **GO/NO-GO DECISION:**
     - Root cause analysis (day_82):
       - Dlaczego bugs nie wykryte wcześniej?
       - Czy to isolated issues czy systemic quality problem?
     - Options:
       - A) Delay v1.0 release o 2-3 tygodnie (fix all critical + improve testing)
       - B) Release v1.0-beta (soft launch, limited users, gather feedback)
       - C) Roll back do last stable (Phase 2 end), stabilize, re-approach Phase 3
     - Decision criteria:
       - Severity of bugs (data loss? crashes? or just UX issues?)
       - Stakeholder tolerance dla delay
       - Business impact of late launch
   - Timeline: +2-3 tygodnie dla full quality fix

**Monitoring:**
- Day_71-74: Bug count from E2E testing (baseline quality)
- Day_81: Bug bash results (triage immediately)
- Day_82: Bug fix velocity (how fast critical bugs resolved?)
- Red flag: >5 critical bugs → activate contingency immediately

**Dependencies:**
- Blocks: v1.0 release (day_87)
- Critical path: YES (quality blocker)

**Notes:**
- This risk is complementary to R-F1-004 (dogfooding bugs in Phase 1)
- R-F1-004 covers MVP testing scope (~50 docs, early features)
- R-F3-005 covers RC testing scope (full product, 1000+ docs, all features)

---

### R-F3-006: UAT reveals critical UX issues - users unable to use tool effectively

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-006 |
| Faza | Phase 3 |
| Kategoria | UX / Product |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 4 (Duży) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | PM |
| Trigger | Day_87: UAT feedback shows ≥3 P0 UX blockers or <70% would recommend |

**Opis ryzyka:**
UAT (day_85-86) może odkryć krytyczne problemy UX:
- Instalacja zbyt trudna (>50% users struggle)
- Core workflow niejasny (users can't complete scenarios without help)
- Error messages confusing (users don't know how to fix issues)
- Missing critical feature ("I can't use this without X")
- Documentation insufficient (users can't self-onboard)

**Impact details:**
- Delay v1.0 launch: 2-5 dni (fix UX issues)
- Scope cut: Remove problematic feature (if low ROI fix)
- Documentation rewrite: 1-2 dni
- Confidence drop: Stakeholders question product readiness
- Potential pivot to v1.0-beta instead of full release

**Mitigation (proactive):**
1. **Continuous dogfooding (day_00-84)**
   - Team uses docflow on real projects
   - Identify UX issues early
   - Fix UX problems before RC
2. **Documentation user testing (day_66-69)**
   - Fresh developer follows USER_GUIDE
   - Iterate based on feedback
   - Ensure docs are clear for new users
3. **Bug bash (day_81)**
   - Internal team tests all workflows
   - Fix obvious UX issues before UAT
   - Prepare RC2 with UX improvements
4. **RC2 quality gate (day_82-83)**
   - Only send RC2 to UAT if bug bash passes (<5 P1 bugs)
   - Verify UX flows work smoothly
   - Review documentation completeness

**Contingency (reactive):**
1. **Scenario A: 1-2 P0 UX blockers (Acceptable)**
   - **Timeline:** day_88 (1 day to fix)
   - **Actions:**
     - Fix blockers (e.g., unclear error message, missing --help text)
     - Update INSTALLATION.md if setup issues
     - Improve error messages (make actionable)
     - Regression test (all platforms)
     - Tag v1.0.0-rc3 (if needed)
   - **Launch:** day_89 (2-day slip from original day_87)
   - **Decision:** Go for launch (minor delay acceptable)

2. **Scenario B: 3-5 P0 UX blockers (Major)**
   - **Timeline:** day_88-90 (3 days to fix + retest)
   - **Actions:**
     - Triage: Fix show-stoppers (P0), defer P1 to v1.0.1
     - Example P0 fixes:
       - Installation fails on macOS ARM → Fix + retest
       - `docflow scan` crashes on Windows WSL2 → Fix
       - Documentation missing critical step → Update docs
       - Error messages unclear → Rewrite + add examples
     - Re-test with 2 UAT participants (smoke test)
     - Tag v1.0.0
   - **Launch:** day_90-92 (3-5 day slip)
   - **Decision:** Go for launch with known limitations (document in KNOWN_LIMITATIONS.md)

3. **Scenario C: >5 P0 blockers or fundamental UX flaw (Critical)**
   - **Timeline:** +1-2 weeks
   - **Actions:**
     - **Root cause analysis (day_88):**
       - Why did we miss this? (requirements gap, testing gap)
       - Is this isolated or systemic issue?
     - **Options:**
       1. **Pivot to v1.0-beta:**
          - Release to limited audience (beta testers only, ~20-50 users)
          - Collect more feedback (2 weeks)
          - Fix critical issues
          - v1.0 full release: +2-4 weeks
       2. **Scope cut:**
          - Remove problematic feature (e.g., recommender if fundamentally broken)
          - v1.0-lite release (core features only: scan, validate, graph)
          - Deferred feature: v1.1 roadmap
          - Timeline: +1 week (remove feature, update docs)
       3. **Delay v1.0:**
          - Full redesign of problematic workflow
          - Re-UAT in 2 weeks
          - v1.0 release: +2-4 weeks
     - **GO/NO-GO DECISION:**
       - Stakeholder meeting (day_88)
       - Assess severity: data loss? unusable? or just confusing?
       - Business impact of delay vs shipping with issues
       - May decide to pause project, reassess product-market fit

**Monitoring:**
- **During UAT (day_85-86):**
  - Monitor survey submissions (real-time)
  - Flag P0 issues immediately (Slack alerts)
  - Daily sync: PM + Tech Lead (triage emerging issues)
- **Red flags:**
  - <50% scenarios completed → workflow too complex
  - Avg ratings <3/5 → major UX issues
  - Multiple "impossible to use" comments → fundamental problem
- **Green signals:**
  - >80% scenarios completed → good UX
  - Avg ratings >4/5 → excellent UX
  - Positive open feedback → product-market fit

**Dependencies:**
- Blocks: v1.0 release (day_87)
- Depends on: day_81 (RC2 ready), day_85-86 (UAT execution)
- Critical path: YES (UAT success required for v1.0 launch)

**Related documents:**
- UAT_PLAN.md (day_85-89 execution plan, scenarios, feedback collection)
- EXTENDED_PLAN.md (day_81 RC, day_85-86 UAT, day_87 launch)
- DEPLOYMENT_STRATEGY.md (v1.0 release process, rollback plan)

**Notes:**
- This risk is NEW (added 2026-02-06 after creating UAT_PLAN.md)
- Complements R-F3-005 (RC quality) - R-F3-005 covers bugs, R-F3-006 covers UX
- UAT is final quality gate before v1.0 launch - critical to project success

---

## RYZYKA CROSS-CUTTING (wszystkie fazy)

### R-CC-001: Scope creep - stakeholders żądają dodatkowych features

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-CC-001 |
| Faza | Wszystkie |
| Kategoria | Scope Management |
| Prawdopodobieństwo | 4 (Wysokie) |
| Wpływ | 3 (Średni) |
| Priorytet | **HIGH (12)** |
| Status | OPEN |
| Owner | Project Manager |
| Trigger | Stakeholder requests feature not w REQUIREMENTS.md |

**Opis ryzyka:**
W trakcie projektu stakeholders mogą prosić o:
- Dodatkowe formaty export (PDF, DOCX)
- Web UI (instead of CLI)
- Integration z zewnętrznymi systemami (Jira, Confluence)
Każda nowa feature: +3-10 dni.

**Impact details:**
- Timeline延长 (każda feature: +5 dni average)
- Zespół demoralized (moving goalposts)
- Original scope nie delivered

**Mitigation (proactive):**
1. **Pre-work: Frozen requirements**
   - Document: REQUIREMENTS.md z all features
   - Stakeholder sign-off przed rozpoczęciem
   - Change control process:
     - Nowa feature = formalna propozycja
     - Impact assessment (effort, timeline)
     - Approval required (project sponsor)
2. **Phased delivery:**
   - v1.0 = core features only (defined w pre-work)
   - v1.1, v1.2 = enhancements (backlog)
   - Clear communication: "this is out of scope dla v1.0"
3. **Demo early & often:**
   - Day_15: MVP demo (set expectations)
   - Day_55: Phase 2 demo
   - Pokazuj progress, zbieraj feedback wcześnie (nie late surprises)

**Contingency (reactive):**
1. **Scenario: Minor feature request (effort <3 dni)**
   - Evaluate:
     - Business value high? Consider adding (use buffer)
     - Business value low? Defer do backlog
   - Decision criteria: P0 dla v1.0? Or can wait dla v1.1?

2. **Scenario: Major feature request (effort 5-10 dni)**
   - **Formal change request:**
     - Document effort estimate
     - Impact: timeline延长 X days
     - Trade-off: które existing features defer?
     - Stakeholder decision: accept delay OR defer new feature
   - Jeśli accept: update timeline, extend buffers
   - Jeśli defer: add do backlog (day_90)

3. **Scenario: Scope creep out of control (>3 major requests)**
   - **Project re-baseline:**
     - Stop development (1-2 dni)
     - Re-negotiate scope:
       - Original v1.0 features
       - Requested features
       - Realistic timeline (original + X weeks)
     - Options:
       - Accept延长 timeline
       - Cut some original features
       - Phase delivery (v1.0 lite → v1.1 full)
     - Formal approval (project sponsor, stakeholders)

**Monitoring:**
- Weekly: Stakeholder check-ins (gather feedback, gauge satisfaction)
- Track: Feature requests count & effort
- Red flag: >2 major requests per sprint → activate formal change control

**Dependencies:**
- Blocks: Timeline adherence
- Critical path: YES (uncontrolled scope = missed deadlines)

---

### R-CC-002: Key dependency (library) deprecated lub critical bug discovered

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-CC-002 |
| Faza | Wszystkie |
| Kategoria | Dependencies |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 4 (Duży) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | Tech Lead |
| Trigger | Library announces deprecation lub critical bug affects docflow |

**Opis ryzyka:**
Third-party libraries mogą:
- Być deprecated (no future updates)
- Have critical bug (data corruption, security)
- Break compatibility (new version breaks our code)

**Impact details:**
- Need migration do innej biblioteki (3-7 dni)
- Lub fork & maintain (ongoing effort)
- Lub workaround (tech debt)

**Mitigation (proactive):**
1. **Pre-work day 3: Choose mature, maintained libraries**
   - Check: GitHub stars, recent commits, active maintainers
   - Prefer: standard libraries > third-party
   - Avoid: abandoned projects (no commits w 1 year)
2. **Dependency pinning:**
   - Pin exact versions (go.mod, Cargo.lock, requirements.txt)
   - Test before upgrading
3. **Monitoring:**
   - GitHub watch (releases, issues)
   - Dependency vulnerability alerts (CI)

**Contingency (reactive):**
1. **Scenario: Library deprecated, ale still works**
   - No immediate action
   - Add do backlog: migrate before library EOL
   - Continue z project (risk accepted)

2. **Scenario: Critical bug, patch available**
   - Upgrade library (1 dzień)
   - Test regression (use buffers)
   - Proceed

3. **Scenario: Critical bug, no patch OR library abandoned**
   - **Migration needed:**
     - Identify replacement library (day 1)
     - Estimate effort (3-7 dni)
     - Options:
       - A) Migrate immediately (use buffers + extend timeline 1 week)
       - B) Workaround (isolate bug, add guard code) + migrate later
       - C) Fork library & patch ourselves (high effort, last resort)
     - Decision criteria:
       - Severity of bug (blocker vs moderate)
       - Time in project (early: migrate now, late: workaround)
   - Update risk register: document migration plan

**Monitoring:**
- Weekly: Check library release notes
- Monthly: Review dependency health (automated tools)

**Dependencies:**
- Blocks: Potentially all features using library
- Critical path: CONDITIONAL (depends on severity)

---

## SUMMARY STATISTICS

**Total ryzyk zidentyfikowanych:** 19

**Breakdown by priority:**
- CRITICAL: 1 (R-PRE-001)
- HIGH: 9 (R-PRE-003, R-F1-001, R-F1-002, R-F1-003, R-F1-004, R-F2-003, R-F3-001, R-F3-006, R-CC-001, R-CC-002)
- MEDIUM: 8 (R-PRE-002, R-F1-005, R-F2-001, R-F2-002, R-F3-002, R-F3-003, R-F3-005)
- LOW: 1 (R-F3-004)

**Breakdown by kategoria:**
- Data Availability: 1
- Technical: 3
- Resources: 1
- Performance: 2
- Data Quality: 2
- Integration: 1
- Quality: 2 (R-F1-004, R-F3-005)
- Algorithm Quality: 3
- Security: 1
- Documentation: 1
- Infrastructure: 1
- Scope Management: 1
- Dependencies: 1
- Stakeholder Management: 1 (R-F1-005)
- UX / Product: 1 (R-F3-006)

**Breakdown by faza:**
- Pre-work: 3
- Phase 1: 6 (R-F1-001, R-F1-002, R-F1-003, R-F1-004, R-F1-005)
- Phase 2: 3
- Phase 3: 6 (R-F3-001, R-F3-002, R-F3-003, R-F3-004, R-F3-005, R-F3-006)
- Cross-cutting: 2

**Critical path ryzyk:** 8 (R-PRE-001, R-PRE-003, R-F1-003, R-F1-005, R-F3-001, R-F3-005, R-CC-001 - bezpośrednio wpływają na deadline)

---

## RISK REVIEW PROCESS

**Frequency:** Co 2 tygodnie (bi-weekly)

**Agenda:**
1. Review status wszystkich OPEN risks (5 min each)
2. Update probability/impact jeśli sytuacja się zmieniła
3. Check triggers - czy któreś się zmaterializowały?
4. Activate contingency plans jeśli needed
5. Identify new risks
6. Close mitigated/obsolete risks

**Ownership:**
- Project Manager: facilitates review
- Tech Lead: technical risks
- Team: contributes observations

**Artifacts:**
- Updated RISK_REGISTER.md (version control)
- Action items (jeśli mitigation needed)

**Escalation:**
- CRITICAL risk triggered → immediate meeting (w 24h)
- HIGH risk triggered → discussion w next daily standup
- MEDIUM/LOW → track w bi-weekly review

---

## RISK RESPONSE DECISION TREE

```
Risk triggered?
│
├─ NO → Continue monitoring (bi-weekly review)
│
└─ YES → Severity?
    │
    ├─ CRITICAL/HIGH
    │   │
    │   ├─ Contingency plan defined?
    │   │   │
    │   │   ├─ YES → Execute contingency
    │   │   │        └─ Monitor effectiveness
    │   │   │             └─ Effective? → Close risk
    │   │   │             └─ Not effective? → Escalate (GO/NO-GO decision)
    │   │   │
    │   │   └─ NO → Emergency meeting (w 24h)
    │   │           └─ Design contingency
    │   │                └─ Execute
    │   │
    │   └─ GO/NO-GO decision needed?
    │       │
    │       ├─ YES → Stakeholder meeting
    │       │        └─ Options: pivot, extend, abort
    │       │
    │       └─ NO → Execute contingency, update timeline
    │
    └─ MEDIUM/LOW
        │
        └─ Use buffers → Fix w scheduled buffer days
            └─ Track resolution
                └─ Close risk when resolved
```

---

## CHANGELOG

### Version 1.2 (2026-02-06)
**Added:**
- R-F3-006: New risk "UAT reveals critical UX issues" (Phase 3, HIGH priority)
  - Addresses UAT testing scenarios in day_85-89
  - Defines 3 contingency scenarios based on P0 blocker count (1-2, 3-5, >5)
  - Includes Go/No-Go decision framework for v1.0 launch

**Updated:**
- Total risk count: 18 → 19
- Phase 3 risks: 5 → 6
- HIGH priority risks: 8 → 9
- Added "UX / Product" category to breakdown

### Version 1.1 (2026-02-06)
**Fixed:**
- M-01: Added buffer day recommendations and risk-buffer mappings
- M-02: Enhanced 13 risk entries with deployment-specific triggers and day_XX timeline references
- M-03: Added missing deployment-related risks (R-F3-004, R-F3-005)
- See [FIXES_APPLIED.md](FIXES_APPLIED.md) for detailed changes

### Version 1.0 (2026-02-06)
**Initial release:**
- 18 risks identified across 4 phases (Pre-work, Phase 1, Phase 2, Phase 3, Cross-cutting)
- Complete risk taxonomy with P×I scoring system
- Mitigation and contingency plans for all risks
- Integration with EXTENDED_PLAN.md and DEPENDENCY_MATRIX.md

---

## ATTACHMENTS

- ALGO_PARAMS.md - parametry algorytmów (progi, wagi) - do stworzenia per feature
- KNOWN_LIMITATIONS.md - documented limitations z risk contingencies - do stworzenia day_31
- TEAM_ALLOCATION.md - resource allocation & commitments - do stworzenia pre-work

---

**Document control:**
- Version: 1.2
- Last updated: 2026-02-06
- Next review: 2026-02-20 (start of Phase 1)
- Owner: Project Manager & Tech Lead
