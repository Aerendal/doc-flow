# PLAN UAT (USER ACCEPTANCE TESTING) - PROJEKT DOCFLOW
## Wersja: 1.0 | Data: 2026-02-06 | Status: DRAFT

---

## RELATED DOCUMENTS

**Navigation:**
- **[← INDEX](INDEX.md)** - Powrót do głównego indeksu dokumentów
- **[EXTENDED_PLAN.md](EXTENDED_PLAN.md)** - Plan rozszerzony (day_81 RC, day_85-86 UAT, day_87 v1.0 launch)
- **[RISK_REGISTER.md](RISK_REGISTER.md)** - Ryzyka (R-F3-005 RC quality, nowe R-F3-006 UAT critical UX issues)
- **[DEPLOYMENT_STRATEGY.md](DEPLOYMENT_STRATEGY.md)** - Deployment (RC distribution methods)
- **[DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md)** - Zależności (UAT blocks v1.0 release)

**Quick links:**
- RC build: See [EXTENDED_PLAN.md - day_81](EXTENDED_PLAN.md) for v1.0.0-rc1/rc2
- Bug fixes: See [EXTENDED_PLAN.md - day_82-83](EXTENDED_PLAN.md) for RC bug fixing window
- Launch: See [EXTENDED_PLAN.md - day_87](EXTENDED_PLAN.md) for v1.0 release (depends on UAT success)

---

## WPROWADZENIE

Plan User Acceptance Testing (UAT) definiuje proces walidacji narzędzia docflow przez reprezentatywnych użytkowników końcowych przed wydaniem v1.0. UAT weryfikuje użyteczność, intuicyjność oraz wykrywa krytyczne problemy UX niewidoczne w testach technicznych.

**Cele UAT:**
1. Validate v1.0-rc2 ready for production
2. Identify UX issues (unclear workflows, confusing errors, missing features)
3. Collect structured feedback (quantitative + qualitative)
4. Build confidence in product-market fit

**Zakres:**
- Version tested: v1.0.0-rc2 (post bug-bash fixes)
- Duration: 2 days (day_86-87)
- Participants: 5-10 beta testers (diverse user profiles)
- Format: Guided scenarios + feedback survey

**Success criteria:**
- 80%+ users rate installation/setup as "Easy"
- Average clarity score ≥ 4/5 for documentation
- <3 P0 blockers discovered (critical UX failures)
- 70%+ would recommend tool to colleagues

---

## UAT OBJECTIVES

### Primary Objectives

| Objective | Description | Success Metric | Priority |
|-----------|-------------|----------------|----------|
| **Usability** | Verify tool is intuitive for target users | 80% complete scenarios without help | P0 |
| **Documentation** | Validate docs clarity (installation, usage) | Avg clarity ≥ 4/5 | P0 |
| **Workflows** | Confirm common workflows are smooth | <5 min time-to-value (first scan) | P1 |
| **Error handling** | Ensure errors are clear and actionable | 90% errors understood without docs | P1 |
| **Feature completeness** | Validate MVP scope meets user needs | 70% would recommend | P0 |

---

### Secondary Objectives

| Objective | Description | Success Metric | Priority |
|-----------|-------------|----------------|----------|
| **Performance** | Verify acceptable speed on real projects | No complaints about slowness | P2 |
| **Platform compatibility** | Test on Linux, macOS, Windows WSL2 | All platforms work | P1 |
| **Edge cases** | Identify unexpected use cases | Document edge cases for v1.1 | P2 |

---

## TEST GROUP

### Target Participants (5-10 users)

**Profile distribution:**

| User Type | Count | Rationale | Recruitment Source |
|-----------|-------|-----------|---------------------|
| **Technical Writers** | 3 | Primary users (60% expected audience) | Professional networks, LinkedIn |
| **Software Architects** | 2 | Decision makers for documentation strategy | Past stakeholders, referrals |
| **Developers** | 2 | Secondary users (write API docs, guides) | Open source contributors, GitHub |
| **Product Managers** | 1-2 | Occasional users (requirements docs) | PM communities, Slack groups |
| **QA Engineers** | 1 | Edge case testers (test plan docs) | QA forums |

**Total:** 9-10 participants (target: confirm 5-7 actual participants)

---

### Participant Criteria

**Must have:**
- Experience writing/maintaining technical documentation
- Access to Linux or macOS environment (or Windows WSL2)
- 2+ hours availability during UAT window (day_86-87)
- English proficiency (documentation in EN)

**Nice to have:**
- Experience with CLI tools (git, npm, etc.)
- Markdown knowledge
- Existing docs project (100+ files) to test with

**Exclusions:**
- Internal team members (already dogfooded during development)
- Users with conflicts of interest (competitors)

---

### Recruitment Timeline

| Day | Activity | Owner | Status |
|-----|----------|-------|--------|
| day_75 | Draft recruitment email | PM | Planned |
| day_78 | Send invites (15 candidates) | PM | Planned |
| day_80 | Confirm participants (target: 7) | PM | Planned |
| day_82 | Send UAT package (RC2 + instructions) | PM | Planned |
| day_85 | Remind participants (start UAT) | PM | Planned |
| day_86-87 | UAT execution | Testers | Planned |
| day_88 | Triage feedback | Team | Planned |
| day_88-89 | Fix P0 issues (if any) | Developers | Planned |

---

## UAT SCENARIOS

### Scenario 1: Setup & First Scan

**Objective:** Verify installation and initial project setup.

**User story:** As a technical writer, I want to install docflow and scan my existing docs project.

**Steps:**
1. **Install docflow** (choose method: binary, Homebrew, or go install)
   - Follow INSTALLATION.md
   - Verify: `docflow --version` shows v1.0.0-rc2
2. **Initialize project**
   - Run: `docflow config init`
   - Review generated `~/.config/docflow/config.yaml`
3. **Scan existing docs**
   - Run: `docflow scan --path=./docs/`
   - Expected: "Scanned 150 files in 2.3s"
4. **Review index**
   - Run: `docflow list`
   - Verify: All files listed, metadata parsed

**Success criteria:**
- User completes steps without external help (docs sufficient)
- Time to complete: <10 minutes
- No errors encountered

**Feedback questions:**
- Was installation clear? (1=Confusing, 5=Very clear)
- Was `docflow config init` output understandable?
- Did `docflow scan` behave as expected?

---

### Scenario 2: Validate Documentation

**Objective:** Validate metadata and section structure.

**User story:** As a technical writer, I want to check if my docs comply with team standards.

**Steps:**
1. **Run validation**
   - Run: `docflow validate`
   - Expected output:
     ```
     Validating 150 files...
     ✓ 120 files passed
     ✗ 30 files failed

     Errors:
     - docs/api/auth.md: Missing required field 'status'
     - docs/guides/setup.md: Empty required section 'Prerequisites'
     ```
2. **Fix one error**
   - Edit `docs/api/auth.md`, add `status: draft`
   - Re-run: `docflow validate docs/api/auth.md`
   - Verify: Error resolved
3. **Understand failure reasons**
   - Review validation output
   - Check docs/CLI_REFERENCE.md for validation rules

**Success criteria:**
- User understands validation errors
- User can fix errors without asking for help
- Validation output is actionable

**Feedback questions:**
- Were validation errors clear? (1=Confusing, 5=Very clear)
- Could you fix errors without external help?
- Any errors that were unclear?

---

### Scenario 3: Use Template Recommender

**Objective:** Discover and use template recommendation feature.

**User story:** As a developer, I want to create a new API guide using a recommended template.

**Steps:**
1. **Discover recommender**
   - Run: `docflow recommend --context="REST API authentication guide"`
   - Expected: List of top 5 templates ranked by relevance
2. **Review template**
   - Run: `docflow templates show --id=api-auth-guide`
   - Review template structure
3. **Generate new doc**
   - Run: `docflow generate --template=api-auth-guide --output=docs/api/oauth2.md`
   - Verify: New file created with template structure
4. **Scan and validate**
   - Run: `docflow scan` (incremental)
   - Run: `docflow validate docs/api/oauth2.md`
   - Verify: New doc is valid

**Success criteria:**
- User discovers recommender feature (via `--help` or docs)
- Recommendations are relevant (at least 2/5 useful)
- Generation workflow is smooth

**Feedback questions:**
- Were recommendations relevant? (1=Not at all, 5=Very relevant)
- Was `docflow generate` output clear?
- Would you use this feature in real work? (Yes/No)

---

### Scenario 4: Dependency Graph

**Objective:** Visualize and understand document dependencies.

**User story:** As an architect, I want to see dependencies between architecture docs.

**Steps:**
1. **Build dependency graph**
   - Run: `docflow graph`
   - Expected: Text output of dependency tree or DOT format
2. **Detect cycles** (if any)
   - If cycles exist: Output lists cycle path
   - Example: `Cycle detected: A → B → C → A`
3. **Topological sort**
   - Run: `docflow graph --sort`
   - Expected: Documents in dependency order
4. **Export graph** (optional)
   - Run: `docflow graph --format=dot > deps.dot`
   - Render: `dot -Tpng deps.dot -o deps.png` (if Graphviz installed)
   - Verify: Visual graph readable

**Success criteria:**
- User understands dependency relationships
- Cycle detection output is clear
- Graph output is useful

**Feedback questions:**
- Was dependency graph output understandable? (1=No, 5=Yes)
- Did you find cycles (if any)?
- Would you use this feature regularly? (Yes/No)

---

## FEEDBACK COLLECTION

### Feedback Survey

**Format:** Google Form (or similar)

**Sections:**

#### Section 1: Installation & Setup (Scenario 1)

| Question | Type | Scale/Options |
|----------|------|---------------|
| How easy was installation? | Likert | 1=Very difficult, 5=Very easy |
| Which installation method did you use? | Multiple choice | Binary / Homebrew / go install / Other |
| Did you encounter errors during installation? | Yes/No | + Text (describe error) |
| Time to complete installation | Numeric | Minutes |
| Was INSTALLATION.md helpful? | Likert | 1=Not helpful, 5=Very helpful |

---

#### Section 2: Core Workflows (Scenarios 2-4)

| Question | Type | Scale/Options |
|----------|------|---------------|
| Were validation errors clear? | Likert | 1=Confusing, 5=Very clear |
| Were recommendations relevant? | Likert | 1=Not relevant, 5=Very relevant |
| Was dependency graph understandable? | Likert | 1=No, 5=Yes |
| Which feature was most useful? | Multiple choice | Scan / Validate / Recommend / Graph / Generate / Other |
| Which feature was least useful? | Multiple choice | Same options |

---

#### Section 3: Documentation Quality

| Question | Type | Scale/Options |
|----------|------|---------------|
| Overall documentation clarity | Likert | 1=Very unclear, 5=Very clear |
| Was USER_GUIDE.md helpful? | Likert | 1=Not helpful, 5=Very helpful |
| Was CLI_REFERENCE.md complete? | Likert | 1=Incomplete, 5=Complete |
| Any missing documentation? | Text | Free text |

---

#### Section 4: Overall Satisfaction

| Question | Type | Scale/Options |
|----------|------|---------------|
| Overall satisfaction with docflow | Likert | 1=Very unsatisfied, 5=Very satisfied |
| Would you use this tool in real work? | Yes/No/Maybe | + Text (why/why not) |
| Would you recommend to colleagues? | Likert | 1=Definitely not, 5=Definitely yes |
| Likelihood to adopt (0-10 NPS) | NPS scale | 0=Not likely, 10=Very likely |

---

#### Section 5: Bugs & Issues

| Question | Type | Scale/Options |
|----------|------|---------------|
| Did you encounter bugs? | Yes/No | + Text (describe) |
| Severity of bugs (if any) | Multiple choice | P0 (blocker) / P1 (major) / P2 (minor) |
| Any confusing error messages? | Text | Free text |

---

#### Section 6: Open Feedback

| Question | Type |
|----------|------|
| What did you like most? | Text (free) |
| What frustrated you most? | Text (free) |
| Missing features you expected? | Text (free) |
| Any other feedback? | Text (free) |

---

### Qualitative Feedback (Optional)

**Method:** 15-minute follow-up call with 2-3 participants

**Questions:**
1. Walk me through a real task you'd use docflow for.
2. What would make you abandon the tool?
3. How does it compare to your current workflow?
4. What's the #1 improvement you'd want?

**Recording:** Notes only (no recording unless consent)

---

## SUCCESS CRITERIA

### Quantitative Metrics

| Metric | Target | Acceptable | Unacceptable | Action if Unacceptable |
|--------|--------|------------|--------------|------------------------|
| **Installation ease** (avg rating) | ≥4.5/5 | ≥4.0/5 | <4.0/5 | Improve INSTALLATION.md |
| **Documentation clarity** (avg) | ≥4.5/5 | ≥4.0/5 | <4.0/5 | Rewrite unclear sections |
| **Would recommend** (%) | ≥80% | ≥70% | <70% | Reassess product-market fit |
| **NPS score** | ≥40 | ≥20 | <20 | Major UX issues, delay launch |
| **Scenarios completed** (%) | 100% | ≥80% | <80% | Identify blocking issues |
| **Time to first value** (scan) | <5 min | <10 min | >10 min | Simplify onboarding |

---

### Qualitative Success Criteria

**Must achieve:**
- [x] <3 P0 blockers discovered (critical UX issues)
- [x] No "impossible to use" feedback
- [x] At least 2 participants say "I would use this daily"

**Should achieve:**
- [x] Positive sentiment in open feedback (>70% positive comments)
- [x] No major feature gaps identified (missing MVP features)
- [x] Error messages rated as clear (avg ≥4/5)

**Nice to have:**
- [ ] Participants suggest v1.1 features (engagement)
- [ ] Participants volunteer for future beta testing
- [ ] Organic sharing (participants tell colleagues)

---

## TIMELINE

### UAT Schedule (day_85-89)

| Day | Time | Activity | Owner | Participants |
|-----|------|----------|-------|--------------|
| **day_85 (Mon)** | 09:00 | Send UAT start email + survey link | PM | - |
| **day_85 (Mon)** | 10:00-18:00 | Testers begin scenarios (async) | Testers | 5-7 users |
| **day_86 (Tue)** | 10:00-18:00 | Continue testing (async) | Testers | 5-7 users |
| **day_86 (Tue)** | 18:00 | Survey submission deadline | Testers | - |
| **day_87 (Wed)** | 09:00-12:00 | Triage feedback (team meeting) | Team | Tech Lead, PM, QA |
| **day_87 (Wed)** | 14:00 | Decision: Go/No-Go for v1.0 launch | PM | Stakeholders |
| **day_88 (Thu)** | 10:00-18:00 | Fix P0 issues (if any, ≤3 expected) | Developers | - |
| **day_89 (Fri)** | 10:00-12:00 | Verify fixes, final smoke test | QA | - |
| **day_89 (Fri)** | 14:00 | Final Go/No-Go decision | PM | Stakeholders |

**Note:** day_87 is originally v1.0 launch day in EXTENDED_PLAN. If UAT reveals P0 issues, launch may slip to day_89 (2-day buffer).

---

### UAT Package Contents

**Distributed:** day_82 (after RC2 stabilizes)

**Contents:**
1. **Binary:** `docflow-v1.0.0-rc2-{platform}-{arch}` (per participant platform)
2. **Documentation:**
   - INSTALLATION.md
   - USER_GUIDE.md (quickstart section)
   - CLI_REFERENCE.md
3. **Test project:** Sample docs corpus (50-100 files) with intentional errors
4. **UAT instructions:** `UAT_INSTRUCTIONS.md`
   - Scenario walkthroughs
   - Expected outputs
   - Survey link
5. **Support contact:** Slack channel or email for questions

**Distribution method:** Private GitHub release or Google Drive link

---

## RISK: UAT REVEALS CRITICAL UX ISSUES

### New Risk: R-F3-006

| **Parametr** | **Wartość** |
|--------------|-------------|
| ID | R-F3-006 |
| Faza | Phase 3 (UAT) |
| Kategoria | UX / Product |
| Prawdopodobieństwo | 2 (Niskie) |
| Wpływ | 4 (Duży) |
| Priorytet | **HIGH (8)** |
| Status | OPEN |
| Owner | PM |
| Trigger | day_87: UAT feedback shows ≥3 P0 UX blockers or <70% would recommend |

**Opis ryzyka:**
UAT może odkryć krytyczne problemy UX:
- Instalacja zbyt trudna (>50% users struggle)
- Core workflow niejasny (users can't complete Scenario 2)
- Error messages confusing (users don't know how to fix)
- Missing critical feature ("I can't use this without X")

**Impact details:**
- Delay v1.0 launch: 2-5 dni (fix issues)
- Scope cut: Remove problematic feature (if low ROI)
- Documentation rewrite: 1-2 dni
- Confidence drop: Stakeholders question product readiness

---

### Mitigation (proactive)

**Before UAT (day_00-84):**
1. **Dogfooding** (continuous):
   - Team uses docflow on real projects
   - Identify UX issues early
2. **Documentation user testing** (day_66-69):
   - Fresh developer follows USER_GUIDE
   - Iterate based on feedback
3. **Bug bash** (day_81):
   - Internal team tests all workflows
   - Fix obvious UX issues before UAT
4. **RC2 quality gate** (day_82-83):
   - Only send RC2 to UAT if bug bash passes (<5 P1 bugs)

---

### Contingency (reactive)

**Scenario A: 1-2 P0 blockers (Acceptable)**
- **Timeline:** day_88 (1 day to fix)
- **Actions:**
  - Fix blockers (e.g., unclear error message, missing --help text)
  - Regression test
  - Tag v1.0.0-rc3 (if needed)
- **Launch:** day_89 (2-day slip from original day_87)
- **Decision:** Go for launch (minor delay acceptable)

---

**Scenario B: 3-5 P0 blockers (Major)**
- **Timeline:** day_88-90 (3 days to fix + retest)
- **Actions:**
  - Triage: Fix show-stoppers (P0), defer P1 to v1.0.1
  - Example P0 fixes:
    - Installation fails on macOS ARM → Fix + retest
    - `docflow scan` crashes on Windows WSL2 → Fix
    - Documentation missing critical step → Update docs
  - Re-test with 2 UAT participants (smoke test)
  - Tag v1.0.0
- **Launch:** day_90-92 (3-5 day slip)
- **Decision:** Go for launch with known limitations (document in KNOWN_LIMITATIONS.md)

---

**Scenario C: >5 P0 blockers or fundamental UX flaw (Critical)**
- **Timeline:** +1-2 weeks
- **Actions:**
  - **Root cause analysis:** Why did we miss this? (requirements gap, testing gap)
  - **Options:**
    1. **Pivot to v1.0-beta:**
       - Release to limited audience (beta testers only)
       - Collect more feedback
       - v1.0 full release: +2-4 weeks
    2. **Scope cut:**
       - Remove problematic feature (e.g., recommender if fundamentally broken)
       - v1.0-lite release (core features only)
       - Deferred feature: v1.1 roadmap
    3. **Delay v1.0:**
       - Full redesign of problematic workflow
       - Re-UAT in 2 weeks
       - v1.0 release: +2-4 weeks
- **Decision:** Stakeholder meeting (day_88), assess options
- **GO/NO-GO:** May decide to pause project, reassess product-market fit

---

### Monitoring

**During UAT (day_86-87):**
- Monitor survey submissions (real-time)
- Flag P0 issues immediately (Slack alerts)
- Daily sync: PM + Tech Lead (triage emerging issues)

**Red flags:**
- <50% scenarios completed → workflow too complex
- Avg ratings <3/5 → major UX issues
- Multiple "impossible to use" comments → fundamental problem

**Green signals:**
- >80% scenarios completed → good UX
- Avg ratings >4/5 → excellent UX
- Positive open feedback → product-market fit

---

## MAPOWANIE DO EXTENDED_PLAN

| UAT Activity | Day(s) | Deliverable | Owner |
|--------------|--------|-------------|-------|
| Recruitment invites | day_78 | 15 candidates contacted | PM |
| Confirm participants | day_80 | 5-7 confirmed | PM |
| Send UAT package | day_82 | RC2 + docs + test project | PM |
| UAT execution | day_85-86 | Testers complete scenarios | Testers |
| Survey collection | day_86 (18:00) | All responses submitted | Testers |
| Feedback triage | day_87 (AM) | Issues categorized (P0/P1/P2) | Team |
| Go/No-Go decision | day_87 (PM) | Launch approved or delayed | PM + Stakeholders |
| Fix P0 issues | day_88-89 | Blockers resolved | Developers |
| Final verification | day_89 | Smoke test passed | QA |

---

## MAPOWANIE DO RISK_REGISTER

| Risk ID | Risk Name | UAT Impact | Mitigation Reference |
|---------|-----------|------------|----------------------|
| R-F3-005 | RC quality issues | High bug count delays UAT | Bug bash (day_81), RC2 quality gate |
| R-F3-006 | UAT reveals critical UX issues | Delays v1.0 launch 2-5 days (or more) | Dogfooding, documentation user testing, contingency plan |

---

## DELIVERABLES

### UAT Report (day_87)

**File:** `LOGS/UAT_REPORT.md`

**Contents:**
1. **Executive summary:**
   - Participants: X confirmed, Y completed
   - Success criteria: Met / Partially met / Not met
   - Go/No-Go recommendation
2. **Quantitative results:**
   - Table of all metrics vs targets
   - Charts (if time permits)
3. **Qualitative themes:**
   - Top 3 liked features
   - Top 3 pain points
   - Missing features requested
4. **Issues discovered:**
   - P0 blockers: X (list)
   - P1 major: Y (list)
   - P2 minor: Z (list)
5. **Recommendations:**
   - Fix before v1.0 (P0)
   - Defer to v1.0.1 (P1)
   - Roadmap for v1.1 (P2)

**Owner:** PM (compile report from survey + triage meeting)

**Timeline:** day_87 (12:00 - ready for Go/No-Go decision)

---

## CHANGELOG

### Version 1.0 (2026-02-06)
- Initial UAT plan
- Objectives: Usability, documentation, workflows, error handling
- Test group: 5-10 users (tech writers, architects, developers)
- 4 guided scenarios: Setup, Validate, Recommend, Graph
- Feedback survey: 6 sections, 25+ questions
- Success criteria: 80% easy installation, 70% would recommend, <3 P0 blockers
- Timeline: day_85-89 (UAT execution + triage + fixes)
- New risk: R-F3-006 (UAT critical UX issues)

---

## NEXT STEPS

### Immediate (Before day_78):
1. Finalize recruitment email template
2. Identify 15 candidate testers (3 writers, 2 architects, 2 devs, etc.)
3. Prepare test project corpus (50-100 files with intentional errors)
4. Draft UAT_INSTRUCTIONS.md

### Pre-UAT (day_78-84):
1. Send invites (day_78)
2. Confirm participants (day_80)
3. Distribute UAT package (day_82: RC2 + docs)
4. Setup survey (Google Form)
5. Create Slack channel for support

### During UAT (day_85-86):
1. Monitor survey submissions
2. Respond to participant questions (Slack)
3. Flag critical issues in real-time

### Post-UAT (day_87-89):
1. Compile UAT_REPORT.md
2. Triage issues (P0/P1/P2)
3. Go/No-Go decision
4. Fix P0 blockers (if any)
5. Final smoke test before v1.0 launch

---

**END OF UAT_PLAN.md**
