# WYMAGANIA WYDAJNOŚCIOWE - PROJEKT DOCFLOW
## Wersja: 1.0 | Data: 2026-02-06 | Status: DRAFT

---

## RELATED DOCUMENTS

**Navigation:**
- **[← INDEX](INDEX.md)** - Powrót do głównego indeksu dokumentów
- **[EXTENDED_PLAN.md](EXTENDED_PLAN.md)** - Plan rozszerzony (day_23 baseline, day_76-79 scale testing)
- **[RISK_REGISTER.md](RISK_REGISTER.md)** - Ryzyka (R-F2-002 performance degradation, R-F3-002 scale collapse)
- **[DEPLOYMENT_STRATEGY.md](DEPLOYMENT_STRATEGY.md)** - Strategia deployment (resource constraints)
- **[DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md)** - Zależności (performance testing timeline)

**Quick links:**
- Baseline testing: See [EXTENDED_PLAN.md - day_23](EXTENDED_PLAN.md) for initial benchmarks
- Scale testing: See [EXTENDED_PLAN.md - day_76-77](EXTENDED_PLAN.md) for large-scale tests
- Optimization: See [EXTENDED_PLAN.md - day_78-79](EXTENDED_PLAN.md) for performance tuning
- Risks: See [RISK_REGISTER.md - R-F2-002, R-F3-002](RISK_REGISTER.md) for performance risk mitigation

---

## WPROWADZENIE

Dokument definiuje wymagania wydajnościowe dla narzędzia docflow CLI. Określa target latency, throughput, resource limits oraz failure modes dla wszystkich operacji na różnych skalach projektu.

**Zakres dokumentu:**
- Performance targets: Latency (P50/P95/P99), throughput
- Scale categories: 100/300/1000/5000 files
- Resource limits: Max RAM, max CPU cores
- Failure modes: OOM, timeout handling
- Incremental scan optimizations (day_47)

**Audience:**
- Tech Lead - architecture constraints
- Developers - optimization goals
- QA - performance test criteria
- Users - expected behavior at scale

**Success criteria:**
- All operations meet "Target" benchmarks for 300 files
- All operations meet "Acceptable" benchmarks for 1000 files
- Operations gracefully degrade or fail-fast beyond 5000 files

---

## KATEGORIE SKALI PROJEKTÓW

Klasyfikacja projektów według liczby dokumentów:

| Category | File Count | Use Case | Target Users | Priority |
|----------|------------|----------|--------------|----------|
| **Small** | 10-100 | Single service API docs | Individual developers | P2 |
| **Medium** | 100-300 | Product documentation | Small teams (3-5) | P0 |
| **Large** | 300-1000 | Enterprise knowledge base | Medium teams (10-20) | P0 |
| **Very Large** | 1000-5000 | Multi-product docs | Large orgs (50+) | P1 |
| **Extreme** | 5000+ | Corporate documentation hub | Enterprises | P2 |

**Design focus:** Medium + Large (100-1000 files) = 80% expected users

**Acceptable limitation:** Extreme scale (5000+) may require optimizations (deferred to v1.1+)

---

## OPERACJE I PERFORMANCE TARGETS

### 1. Scan Operation

**Opis:** Skanowanie systemu plików, parsowanie Markdown + YAML frontmatter, budowa indeksu dokumentów.

**Command:** `docflow scan`

**Implementation:** day_04 (initial), day_46-47 (incremental)

---

#### 1.1 Scan - Medium Scale (300 files)

**Target environment:**
- Files: 300 Markdown files
- Avg file size: 5KB
- Total size: ~1.5MB
- File tree depth: 5 levels

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency (P50)** | <2s | <5s | >10s | P0 |
| **Latency (P95)** | <4s | <8s | >15s | P0 |
| **Latency (P99)** | <6s | <12s | >20s | P1 |
| **RAM usage** | <100MB | <300MB | >1GB | P0 |
| **CPU cores** | 1-2 cores | 4 cores | >4 cores | P1 |

**Measurement:** day_23 (performance baseline)

**Benchmark command:**
```bash
time docflow scan --path=/testdata/corpus-300/
# Expected output: "Scanned 300 files in 1.8s"
```

---

#### 1.2 Scan - Large Scale (1000 files)

**Target environment:**
- Files: 1000 Markdown files
- Avg file size: 5KB
- Total size: ~5MB
- File tree depth: 6 levels

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency (P50)** | <8s | <20s | >60s | P0 |
| **Latency (P95)** | <15s | <40s | >120s | P1 |
| **Latency (P99)** | <25s | <60s | >180s | P2 |
| **RAM usage** | <300MB | <1GB | >2GB | P0 |
| **CPU cores** | 2-4 cores | 4 cores | >4 cores | P1 |

**Measurement:** day_76-77 (scale testing)

---

#### 1.3 Scan - Very Large Scale (5000 files)

**Target environment:**
- Files: 5000 Markdown files
- Avg file size: 5KB
- Total size: ~25MB
- File tree depth: 7 levels

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency (P50)** | <60s | <120s | >300s | P1 |
| **Latency (P95)** | <90s | <180s | >600s | P2 |
| **RAM usage** | <1GB | <2GB | >4GB | P0 |
| **CPU cores** | 4 cores | 4 cores | >8 cores | P2 |

**Measurement:** day_76-77 (scale testing)

**Risk:** R-F3-002 (performance collapse at 5000+)

**Mitigation:** If unacceptable, document limitation (max 1000 files recommended) in KNOWN_LIMITATIONS.md

---

#### 1.4 Incremental Scan (day_47 feature)

**Opis:** Scan only changed files (detect via mtime, cache previous results)

**Use case:** Daily workflow, only 10% files changed

**Performance improvement targets:**

| Scale | Full Scan | Incremental Scan (10% changed) | Speedup | Priority |
|-------|-----------|--------------------------------|---------|----------|
| 300 files | 2s | <0.5s | 4x | P1 |
| 1000 files | 10s | <2s | 5x | P0 |
| 5000 files | 60s | <10s | 6x | P1 |

**Implementation:** day_46-47 (incremental parsing)

**Benchmark command:**
```bash
# First scan (full)
docflow scan  # 2s for 300 files

# Modify 30 files (10%)
touch testdata/corpus-300/api/*.md

# Second scan (incremental)
docflow scan  # Expected: <0.5s (4x speedup)
```

**Success criteria:**
- Speedup ≥ 3x for 10% change rate
- Cache invalidation correct (no stale data)

---

### 2. Validate Operation

**Opis:** Walidacja metadanych, sekcji, zależności, quality checks.

**Command:** `docflow validate`

**Implementation:** day_06 (metadata), day_08 (sections), day_42 (progressive), day_56-57 (governance)

---

#### 2.1 Validate - Medium Scale (300 files)

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency (P50)** | <3s | <8s | >20s | P0 |
| **Latency (P95)** | <6s | <15s | >30s | P0 |
| **RAM usage** | <150MB | <500MB | >1GB | P0 |

**Breakdown (cumulative):**
- Load index from cache: <0.5s
- Metadata validation: <1s
- Section validation: <1.5s
- Dependency validation: <0.5s
- Total: ~3.5s

**Measurement:** day_23 (baseline), day_76 (scale testing)

---

#### 2.2 Validate - Large Scale (1000 files)

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency (P50)** | <10s | <30s | >90s | P0 |
| **Latency (P95)** | <20s | <60s | >180s | P1 |
| **RAM usage** | <500MB | <1.5GB | >2GB | P0 |

**Measurement:** day_76-77 (scale testing)

---

### 3. Graph Build Operation

**Opis:** Budowa dependency graph, cycle detection, topological sort.

**Command:** `docflow graph`

**Implementation:** day_09 (graph build), day_10 (toposort)

---

#### 3.1 Graph Build - Medium Scale (100 nodes, 150 edges)

**Target environment:**
- Nodes: 100 documents
- Edges: 150 dependencies
- Avg degree: 1.5
- Max depth: 5 levels

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Build graph** | <0.5s | <1s | >3s | P0 |
| **Cycle detection** | <0.2s | <0.5s | >2s | P0 |
| **Topological sort** | <0.3s | <0.8s | >3s | P0 |
| **Total** | <1s | <2.5s | >8s | P0 |
| **RAM usage** | <50MB | <200MB | >500MB | P0 |

**Algorithm complexity:**
- Build graph: O(V + E) using adjacency list
- Cycle detection: O(V + E) DFS
- Toposort: O(V + E) Kahn's algorithm

**Measurement:** day_23 (baseline)

---

#### 3.2 Graph Build - Large Scale (1000 nodes, 2000 edges)

**Target environment:**
- Nodes: 1000 documents
- Edges: 2000 dependencies
- Avg degree: 2.0
- Max depth: 8 levels

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Build graph** | <3s | <10s | >30s | P0 |
| **Cycle detection** | <2s | <8s | >20s | P0 |
| **Topological sort** | <2s | <8s | >20s | P0 |
| **Total** | <7s | <26s | >70s | P0 |
| **RAM usage** | <300MB | <1GB | >2GB | P0 |

**Measurement:** day_76-77 (scale testing)

**Risk:** R-F3-002 (graph operations may become O(n²) if implemented incorrectly)

**Mitigation:** Use efficient algorithms (adjacency list, not matrix), profiling (day_78-79)

---

### 4. Generate Operation

**Opis:** Generowanie nowego dokumentu z szablonu (template rendering).

**Command:** `docflow generate --template=api-guide`

**Implementation:** day_11 (generator)

---

#### 4.1 Generate - Single Document

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency** | <0.5s | <2s | >5s | P1 |
| **RAM usage** | <50MB | <100MB | >300MB | P2 |

**Breakdown:**
- Load template: <0.1s
- Render (substitute placeholders): <0.2s
- Write to disk: <0.1s
- Total: ~0.4s

**Measurement:** day_23 (baseline)

---

### 5. Recommend Operation

**Opis:** Rekomendacja szablonu na podstawie kontekstu + quality scoring.

**Command:** `docflow recommend --context="API authentication guide"`

**Implementation:** day_12 (recommender), day_27-29 (ML-based scoring)

---

#### 5.1 Recommend - Medium Corpus (300 templates)

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency (P50)** | <1s | <3s | >10s | P0 |
| **Latency (P95)** | <2s | <5s | >15s | P1 |
| **RAM usage** | <100MB | <300MB | >1GB | P1 |

**Breakdown:**
- Load template index: <0.2s
- Compute similarity (TF-IDF): <0.5s
- Rank by quality + similarity: <0.2s
- Return top 5: <0.1s
- Total: ~1s

**Measurement:** day_29 (recommender evaluation)

**Success criteria:** Precision@5 > 0.70 (see EXTENDED_PLAN.md day_29)

---

### 6. Daily Planner Operation

**Opis:** Generowanie planu pracy na dzień (topological order + effort estimation).

**Command:** `docflow plan --max-hours=8`

**Implementation:** day_33-34 (daily planner)

---

#### 6.1 Daily Planner - Medium Project (100 docs)

**Performance targets:**

| Metric | Target | Acceptable | Unacceptable | Priority |
|--------|--------|------------|--------------|----------|
| **Latency** | <5s | <10s | >30s | P1 |
| **RAM usage** | <100MB | <300MB | >500MB | P2 |

**Breakdown:**
- Build graph: <1s
- Compute toposort: <0.5s
- Filter incomplete docs: <0.5s
- Estimate effort: <1s
- Select up to max-hours: <2s
- Total: ~5s

**Measurement:** day_34 (planner testing)

---

## RESOURCE LIMITS

### Acceptable Resource Consumption

| Resource | Target (300 files) | Acceptable (1000 files) | Max Limit (Any scale) |
|----------|--------------------|--------------------------|-----------------------|
| **RAM** | <300MB | <1GB | 2GB (hard limit) |
| **CPU Cores** | 1-2 | 2-4 | 4 (recommended max) |
| **Disk Space (cache)** | <10MB | <50MB | 500MB (auto-cleanup) |
| **Disk I/O** | <1000 IOPS | <5000 IOPS | Unbounded |
| **Network** | 0 (offline tool) | 0 | Optional (version check API) |

**Implementation:**
- RAM limit: day_78-79 (memory optimization, streaming parsers)
- CPU: Configurable in config.yaml (`performance.max_workers: 4`)
- Cache cleanup: day_46-47 (cache invalidation, LRU eviction)

---

### Failure Modes

#### 1. Out of Memory (OOM)

**Trigger:** RAM usage exceeds system limit (e.g., 8GB laptop, docflow uses >6GB)

**Scenarios:**
- Parsing 10000+ files with in-memory index
- Loading entire file content into RAM (large files >100MB)
- Graph build with 100k+ nodes (extreme scale)

**Mitigation:**
1. **Streaming parsers** (day_78-79):
   - Don't load entire file into RAM
   - Parse frontmatter only, skip content if not needed
2. **Chunked processing**:
   - Process files in batches of 1000
   - Clear cache between batches
3. **Fail-fast with clear error**:
   ```
   ERROR: Out of memory while scanning. Try:
     1. Reduce scope: docflow scan --path=./docs/api/ (not entire repo)
     2. Increase --max-file-size limit (currently 10MB)
     3. Use incremental scan: docflow scan --incremental
   ```

**Implementation:** day_78-79 (optimization sprint)

---

#### 2. Timeout (Operation Too Slow)

**Trigger:** Operation exceeds user patience (>2 minutes for interactive command)

**Scenarios:**
- Scan 5000+ files on slow HDD
- Graph build with cycles (infinite loop bug)
- Network timeout (version check API)

**Mitigation:**
1. **Progress indicators** (day_13):
   - Spinner: "Scanning files... (150/300)"
   - ETA: "Estimated 30s remaining"
2. **Configurable timeout** (config.yaml):
   ```yaml
   performance:
     scan_timeout: 300s  # 5 minutes max
   ```
3. **Fail-fast on timeout**:
   ```
   ERROR: Scan timed out after 300s.
   Scanned 3500/5000 files before timeout.
   Try: docflow scan --incremental (faster on subsequent runs)
   ```

**Implementation:** day_43-44 (edge case hardening)

---

#### 3. Disk Space Exhausted

**Trigger:** Cache directory exceeds limit (500MB)

**Mitigation:**
1. **Auto cleanup** (LRU eviction):
   - Keep 100 most recent scans
   - Delete cache older than 30 days
2. **Manual cleanup**:
   ```bash
   docflow cache clear
   # Cleared 450MB of cached data
   ```

**Implementation:** day_46-47 (cache management)

---

## INCREMENTAL SCAN PERFORMANCE (day_47)

### Use Case: Daily Workflow

**Scenario:**
- Developer works on docs project (1000 files)
- Modifies 50-100 files per day (5-10% change rate)
- Runs `docflow scan` multiple times per day

**Without incremental scan:**
- Full scan: 10s (every time)
- Daily overhead: 10s × 10 runs = 100s wasted

**With incremental scan:**
- First scan (cold): 10s
- Subsequent scans (10% changed): 2s
- Daily overhead: 10s + 2s × 9 = 28s (72% time saved)

---

### Implementation Details

**Change detection:**
1. **File mtime (modification time):**
   - Cache: `{file_path: {mtime, parsed_data}}`
   - Compare: `current_mtime > cached_mtime` → re-parse
2. **Content hash (optional, v1.1):**
   - SHA256 hash of file content
   - More accurate (detects mtime-only changes without content change)

**Cache storage:**
- Format: JSON or SQLite (day_46 decision)
- Location: `~/.cache/docflow/index.db`
- Max size: 100MB (auto-evict oldest)

**Cache invalidation:**
- File deleted: Remove from cache
- File moved: Detect via path change, re-parse
- Config changed: Full re-scan (config hash changed)

**Benchmark (day_47):**
```bash
# First scan (cold cache)
docflow scan --path=./testdata/1000-files/
# Output: Scanned 1000 files in 10.2s

# Modify 100 files
touch ./testdata/1000-files/api/*.md  # 100 files

# Second scan (warm cache, 10% changed)
docflow scan --path=./testdata/1000-files/
# Output: Scanned 100 files (900 from cache) in 1.8s
# Speedup: 5.7x
```

**Success criteria:**
- Speedup ≥ 3x for 10% change rate
- Speedup ≥ 5x for 5% change rate
- Cache hit rate ≥ 90%

---

## MAPOWANIE DO EXTENDED_PLAN

| Performance Activity | Day(s) | Deliverable | Metrics |
|----------------------|--------|-------------|---------|
| Performance baseline | day_23 | LOGS/PERFORMANCE_BASELINE.md | 300 files: scan <5s, validate <8s, graph <2.5s |
| Incremental parsing | day_46-47 | Cache implementation, speedup 3-5x | 10% change: <2s scan for 1000 files |
| Scale testing | day_76-77 | LOGS/SCALE_TEST_RESULTS.md | 1000 files, 5000 files tested |
| Final optimization | day_78-79 | LOGS/FINAL_PERFORMANCE.md | Bottlenecks addressed, targets met |

---

## MAPOWANIE DO RISK_REGISTER

| Risk ID | Risk Name | Performance Impact | Mitigation Reference |
|---------|-----------|--------------------|-----------------------|
| R-F2-002 | Performance degradation during Phase 2 | Features slow down core operations | Profile (day_50), optimize (day_78-79) |
| R-F3-002 | Scale testing shows collapse at 5000+ files | Unacceptable latency or OOM | Optimization sprint (day_78-79), document limitations if needed |

---

## PERFORMANCE TESTING STRATEGY

### Benchmark Suites

**Location:** `tests/performance/`

**Suites:**
1. **Unit benchmarks** (day_23):
   - `parser_test.go`: Parse 100 files
   - `graph_test.go`: Build graph 100 nodes
   - `validator_test.go`: Validate 100 files
2. **Integration benchmarks** (day_23):
   - Full pipeline: scan + validate + graph (300 files)
3. **Scale benchmarks** (day_76-77):
   - 1000 files, 5000 files
   - Memory profiling (`go test -memprofile`)
   - CPU profiling (`go test -cpuprofile`)

**CI integration (day_61-62):**
```yaml
# .github/workflows/ci.yml
- name: Performance benchmarks
  run: go test ./tests/performance/... -bench=. -benchmem
  # Fail if regression >20%
```

---

### Profiling Tools

**CPU profiling:**
```bash
go test -cpuprofile=cpu.prof -bench=BenchmarkScan
go tool pprof cpu.prof
# (pprof) top10
# (pprof) list ScanFiles
```

**Memory profiling:**
```bash
go test -memprofile=mem.prof -bench=BenchmarkScan
go tool pprof mem.prof
# (pprof) top10
# (pprof) list ParseMarkdown
```

**Implementation:** day_78-79 (optimization sprint uses profiling to identify bottlenecks)

---

## SUCCESS CRITERIA SUMMARY

### v1.0 Launch Requirements

**Must meet (P0):**
- [x] Scan 300 files: P50 <5s ✓
- [x] Validate 300 files: P50 <8s ✓
- [x] Graph build 100 nodes: <2.5s ✓
- [x] RAM usage 300 files: <300MB ✓
- [x] Incremental scan speedup: ≥3x ✓

**Should meet (P1):**
- [x] Scan 1000 files: P50 <20s
- [x] Validate 1000 files: P50 <30s
- [x] Graph build 1000 nodes: <26s
- [x] RAM usage 1000 files: <1GB

**Nice to have (P2):**
- [ ] Scan 5000 files: P50 <120s (or document limitation)
- [ ] RAM usage 5000 files: <2GB (or fail gracefully)

**Measurement:** day_76-77 (scale testing verification)

**Go/No-Go decision (day_81):**
- If P0 not met: Delay release, optimize (extend Phase 3 +1 week)
- If P1 not met: Release v1.0, add known limitation to docs
- If P2 not met: Acceptable, document limitation

---

## CHANGELOG

### Version 1.0 (2026-02-06)
- Initial performance requirements
- Scale categories: 100/300/1000/5000 files
- Targets defined: Latency (P50/P95/P99), RAM, CPU
- Failure modes: OOM, timeout, disk space
- Incremental scan performance goals (3-5x speedup)
- Mapped to EXTENDED_PLAN: day_23, 46-47, 76-79
- Mapped to RISK_REGISTER: R-F2-002, R-F3-002

---

## NEXT STEPS

### Implementation Timeline:
1. day_23: Establish performance baseline (300 files)
2. day_46-47: Implement incremental scan (verify 3-5x speedup)
3. day_76-77: Scale testing (1000, 5000 files)
4. day_78-79: Optimization sprint (address bottlenecks)
5. day_81: Go/No-Go decision (verify P0 targets met)

### Monitoring (Post-v1.0):
1. Collect real-world performance stats (opt-in telemetry, v1.1)
2. Identify common bottlenecks
3. Prioritize optimizations for v1.1+

---

**END OF PERFORMANCE_REQUIREMENTS.md**
