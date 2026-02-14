# STRATEGIA DEPLOYMENT - PROJEKT DOCFLOW
## Wersja: 1.0 | Data: 2026-02-06 | Status: DRAFT

---

## RELATED DOCUMENTS

**Navigation:**
- **[← INDEX](INDEX.md)** - Powrót do głównego indeksu dokumentów
- **[EXTENDED_PLAN.md](EXTENDED_PLAN.md)** - Plan rozszerzony (day_61-64 CI/CD, day_63 cross-compilation)
- **[RISK_REGISTER.md](RISK_REGISTER.md)** - Ryzyka (R-F3-004 CI/CD failures, R-F3-001 security)
- **[PERFORMANCE_REQUIREMENTS.md](PERFORMANCE_REQUIREMENTS.md)** - Wymagania wydajnościowe (deployment constraints)
- **[DEPENDENCY_MATRIX.md](DEPENDENCY_MATRIX.md)** - Zależności (day_61-65 deployment phase)

**Quick links:**
- CI/CD setup: See [EXTENDED_PLAN.md - day_61-62](EXTENDED_PLAN.md) for pipeline implementation
- Cross-compilation: See [EXTENDED_PLAN.md - day_63](EXTENDED_PLAN.md) for multi-platform builds
- Risks: See [RISK_REGISTER.md - R-F3-004](RISK_REGISTER.md) for deployment failure mitigation

---

## WPROWADZENIE

Dokument definiuje strategię wdrożenia narzędzia docflow CLI w środowiskach produkcyjnych. Obejmuje platformy docelowe, metody instalacji, konfigurację, aktualizacje oraz politykę kompatybilności wstecznej.

**Zakres dokumentu:**
- Target platforms: Linux, macOS, Windows
- Installation methods: Binary, package managers, source
- Configuration management: XDG Base Directory compliance
- Update mechanism: Manual + optional version check
- Backwards compatibility: v1.0 → v1.1 stable, v2.0 migration support

**Audience:**
- DevOps Engineers - deployment automation
- End Users - installation procedures
- Tech Lead - architecture decisions
- Release Manager - release process

---

## PLATFORMY DOCELOWE

### Wspierane Platformy (v1.0)

| Platform | Architectures | Priority | Support Level | Min OS Version | Test Coverage |
|----------|---------------|----------|---------------|----------------|---------------|
| **Linux** | x86_64, ARM64 | P0 | Full | Ubuntu 20.04, Debian 11, RHEL 8 | CI automated |
| **macOS** | Intel (x86_64), Apple Silicon (ARM64) | P0 | Full | macOS 11 (Big Sur) | CI automated |
| **Windows** | x86_64 (WSL2 + Native) | P1 | Best effort | Windows 10 21H2, WSL2 | Manual testing |

**Implementacja:** day_63-64

**Decyzje techniczne:**
- **Linux:** Primary target, 70% expected users (tech writers, devops)
- **macOS:** Secondary, 25% expected users (architects, developers)
- **Windows:** Limited support, 5% expected users
  - **WSL2:** Full support (Linux binary works)
  - **Native Windows:** Best-effort via cross-compilation (PowerShell compatibility required)

**Build matrix (CI):**
```yaml
# .github/workflows/build.yml excerpt
matrix:
  os: [ubuntu-22.04, macos-12, macos-14, windows-2022]
  arch: [amd64, arm64]
  exclude:
    - os: windows-2022
      arch: arm64
```

**Related risk:** R-F3-004 (CI/CD pipeline failures)

---

## METODY INSTALACJI

### 1. Binary Tarball (Universal)

**Target:** All platforms
**Priority:** P0
**Implementacja:** day_63-64

**Proces instalacji:**

```bash
# Linux/macOS
curl -LO https://github.com/docflow/docflow/releases/download/v1.0.0/docflow-v1.0.0-linux-amd64.tar.gz
tar -xzf docflow-v1.0.0-linux-amd64.tar.gz
sudo mv docflow /usr/local/bin/
docflow --version
```

**Binary naming convention:**
- `docflow-v{VERSION}-{OS}-{ARCH}.tar.gz`
- Examples:
  - `docflow-v1.0.0-linux-amd64.tar.gz`
  - `docflow-v1.0.0-darwin-arm64.tar.gz`
  - `docflow-v1.0.0-windows-amd64.zip`

**Checksums:**
- Generate SHA256 checksums: `docflow-v1.0.0-checksums.txt`
- Sign with GPG (optional v1.1): `docflow-v1.0.0-checksums.txt.asc`

**Success criteria:**
- Binary size < 20MB per platform (statically linked)
- Download from GitHub Releases
- Checksum verification documented

---

### 2. Package Managers

#### 2a. Debian/Ubuntu (.deb)

**Target:** Ubuntu 20.04+, Debian 11+
**Priority:** P1
**Implementacja:** day_64 (stretch goal)

**Installation:**
```bash
# Add repository (future - v1.1)
curl -fsSL https://apt.docflow.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docflow-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/docflow-archive-keyring.gpg] https://apt.docflow.io stable main" | sudo tee /etc/apt/sources.list.d/docflow.list
sudo apt update
sudo apt install docflow

# v1.0: Manual .deb download
wget https://github.com/docflow/docflow/releases/download/v1.0.0/docflow_1.0.0_amd64.deb
sudo dpkg -i docflow_1.0.0_amd64.deb
```

**Package metadata:**
- Package name: `docflow`
- Maintainer: docflow@example.com
- Depends: `libc6 (>= 2.31)` (if dynamically linked)
- Section: `utils`
- Priority: `optional`

**Installation path:** `/usr/bin/docflow`

---

#### 2b. Red Hat/CentOS (.rpm)

**Target:** RHEL 8+, CentOS Stream, Fedora 35+
**Priority:** P2
**Implementacja:** v1.1 (deferred)

**Reason for deferral:** Lower priority, 15% market share among target users.

---

#### 2c. Homebrew (macOS)

**Target:** macOS 11+
**Priority:** P1
**Implementacja:** day_64

**Installation:**
```bash
brew tap docflow/tap
brew install docflow

# Or direct formula (future)
brew install docflow
```

**Homebrew formula (simplified):**
```ruby
# docflow.rb
class Docflow < Formula
  desc "Documentation workflow automation CLI"
  homepage "https://github.com/docflow/docflow"
  url "https://github.com/docflow/docflow/releases/download/v1.0.0/docflow-v1.0.0-darwin-amd64.tar.gz"
  sha256 "..."
  version "1.0.0"

  def install
    bin.install "docflow"
  end

  test do
    system "#{bin}/docflow", "--version"
  end
end
```

**Success criteria:**
- Formula tested on Intel + Apple Silicon
- Auto-detects architecture
- Formula submitted to homebrew-core (v1.1 goal)

---

#### 2d. Chocolatey (Windows)

**Target:** Windows 10+
**Priority:** P2
**Implementacja:** v1.1 (deferred)

**Reason for deferral:** Limited Windows native users, WSL2 binary sufficient for v1.0.

---

### 3. Go Install (Developer Method)

**Target:** Developers with Go toolchain
**Priority:** P1
**Implementacja:** day_00 (already works)

**Installation:**
```bash
go install github.com/docflow/docflow/cmd/docflow@latest
```

**Advantages:**
- Always latest version
- Auto-compiles for host platform
- No binary download needed

**Disadvantages:**
- Requires Go 1.21+ installed
- Slower (compile time ~30s)

**Documented in:** `README.md`, `docs/INSTALLATION.md`

---

### 4. Install Script (Automated)

**Target:** Linux, macOS
**Priority:** P1
**Implementacja:** day_64

**Installation:**
```bash
curl -fsSL https://install.docflow.io | bash
```

**Script logic (`install.sh`):**
1. Detect OS + architecture (`uname -s`, `uname -m`)
2. Download latest binary from GitHub Releases
3. Verify checksum
4. Extract to `/usr/local/bin/` (or `~/.local/bin/` if no sudo)
5. Verify installation: `docflow --version`
6. Print success message + quickstart

**Safety features:**
- Checksum verification (SHA256)
- Dry-run mode: `curl ... | bash -s -- --dry-run`
- Version pinning: `curl ... | bash -s -- --version=1.0.0`

**Related risk:** R-F3-001 (security - script injection)

**Mitigation:**
- Host script on trusted domain
- Sign releases with GPG (v1.1)
- Document manual installation alternative

---

## KONFIGURACJA

### Configuration File Location (XDG Base Directory Spec)

**Default paths:**

| Platform | Config Path | Cache Path | Data Path |
|----------|-------------|------------|-----------|
| **Linux** | `~/.config/docflow/config.yaml` | `~/.cache/docflow/` | `~/.local/share/docflow/` |
| **macOS** | `~/.config/docflow/config.yaml` | `~/Library/Caches/docflow/` | `~/Library/Application Support/docflow/` |
| **Windows** | `%APPDATA%\docflow\config.yaml` | `%LOCALAPPDATA%\docflow\cache\` | `%APPDATA%\docflow\data\` |

**Environment variable overrides:**
- `DOCFLOW_CONFIG_DIR` - config directory
- `DOCFLOW_CACHE_DIR` - cache directory
- `DOCFLOW_DATA_DIR` - data directory

**Implementacja:** day_02 (config loader), day_63 (platform-specific paths)

---

### Config File Structure

**Default config (`~/.config/docflow/config.yaml`):**

```yaml
# docflow v1.0.0 configuration
version: 1

# Scan settings
scan:
  exclude_patterns:
    - "node_modules/**"
    - ".git/**"
    - "vendor/**"
  max_file_size: 10485760  # 10MB
  follow_symlinks: false

# Validation settings
validation:
  draft:
    allow_empty_sections: true
    allow_missing_context: true
  published:
    allow_empty_sections: false
    require_all_deps: true

# Performance
performance:
  max_workers: 4
  cache_enabled: true
  incremental_scan: true

# Recommendations
recommender:
  max_results: 5
  min_quality_score: 60.0

# Output
output:
  format: "table"  # table | json | yaml
  color: auto      # auto | always | never
  verbose: false
```

**Generation:**
- `docflow config init` - generates default config
- `docflow config validate` - validates config syntax
- `docflow config show` - prints active config (merged: defaults + file + env)

**Implementacja:** day_02, extended day_42 (progressive validation config)

---

### Project-Level Config

**File:** `.docflow.yaml` (project root)

**Use case:** Per-project overrides (team shared config)

**Priority:** Project config > User config > Default config

**Example:**
```yaml
# .docflow.yaml (project-specific)
scan:
  exclude_patterns:
    - "build/**"
    - "dist/**"

validation:
  governance:
    enabled: true
    rules_file: ".docflow/governance.yaml"
```

**Implementacja:** day_02 (config loader priority)

---

## MECHANIZM AKTUALIZACJI

### v1.0: Manual Update

**Method:** User manually downloads + replaces binary

**Check for updates (optional):**
```bash
docflow version --check
# Output:
# Current version: v1.0.0
# Latest version: v1.0.1
# Update available! Download: https://github.com/docflow/docflow/releases/latest
```

**Implementacja:** day_69 (version check API)

**API endpoint:** `https://api.github.com/repos/docflow/docflow/releases/latest`

**Frequency:** Check once per day (cache result in `~/.cache/docflow/version_check.json`)

**Privacy:** No telemetry sent, only HTTP GET request to GitHub API

**Disable:** `docflow config set version_check.enabled false`

---

### v1.1+: Self-Update (Optional)

**Method:** `docflow update` command (future feature)

**Planned implementation:**
1. Check latest release via GitHub API
2. Download binary for current platform
3. Verify checksum
4. Replace current binary (requires sudo if `/usr/local/bin/`)
5. Restart with new version

**Risk mitigation:**
- Backup old binary before replace
- Rollback on failure: `docflow update --rollback`

**Deferred reason:** v1.0 scope too tight, low priority (manual update sufficient)

---

## KOMPATYBILNOŚĆ WSTECZNA

### Versioning Policy (SemVer)

**Format:** `MAJOR.MINOR.PATCH` (e.g., v1.0.0)

**Compatibility guarantees:**

| Version Change | Config Compatibility | CLI Compatibility | Data Format |
|----------------|----------------------|-------------------|-------------|
| **PATCH** (1.0.0 → 1.0.1) | 100% compatible | 100% compatible | 100% compatible |
| **MINOR** (1.0.0 → 1.1.0) | Compatible (new fields optional) | Compatible (new flags optional) | Compatible (extends, not breaks) |
| **MAJOR** (1.0.0 → 2.0.0) | May break (migration tool provided) | May break (deprecation warnings in 1.x) | May break (migration required) |

---

### v1.0 → v1.1 Migration

**Expected changes:**
- New optional config fields (backwards compatible)
- New CLI flags (no breaking changes)
- New metadata fields (old docs still valid)

**Migration:** None required (automatic)

**User action:** Update binary, existing config works

---

### v1.x → v2.0 Migration

**Expected breaking changes (hypothetical):**
- Config schema change (YAML → TOML)
- Metadata format change (frontmatter structure)
- CLI command restructure (`docflow validate` → `docflow check`)

**Migration tool:** `docflow migrate --from=1.x --to=2.0`

**Process:**
1. Backup project
2. Run migration tool
3. Review changes (dry-run mode)
4. Apply migration
5. Test with v2.0

**Timeline:** v2.0 roadmap: +12 months post-v1.0

**Implementacja:** day_90 (roadmap planning)

---

### Deprecation Policy

**Process:**
1. **Announce deprecation** (release notes, warnings in CLI)
2. **Grace period:** Minimum 2 minor versions (e.g., deprecated in 1.2, removed in 1.4)
3. **Warnings:** CLI prints deprecation warnings when using old syntax
4. **Documentation:** Migration guide published before removal

**Example:**
```bash
# v1.2: New command
docflow check --governance

# v1.2: Old command (deprecated, still works)
docflow validate --governance
# WARNING: 'docflow validate --governance' is deprecated. Use 'docflow check --governance'. This will be removed in v1.4.
```

---

## DEPLOYMENT STAGES

### Stage 1: Development Builds (day_00-60)

**Audience:** Internal developers
**Frequency:** Per commit (CI)
**Artifacts:**
- Binary: `docflow-dev-{commit_sha}-{platform}-{arch}`
- Retention: 30 days
- Distribution: CI artifacts, not public

**Purpose:** Daily development, testing, dogfooding

**Implementacja:** day_61-62 (CI pipeline)

---

### Stage 2: Alpha Release (day_55, Phase 2 complete)

**Audience:** Internal + select beta testers (5-10 users)
**Version:** `v0.5.0-alpha`
**Artifacts:**
- Binaries for 3 platforms
- GitHub pre-release tag
- Installation: Manual binary download

**Purpose:**
- Early feedback on Phase 2 features
- Identify integration issues
- Performance testing on real projects

**Timeline:** day_55-60 (1 week alpha testing)

**Success criteria:**
- 5 beta testers onboarded
- 3+ real projects scanned
- 10+ issues reported and triaged

**Related risk:** None (internal only)

---

### Stage 3: Release Candidate (day_81)

**Audience:** Beta testers + public (limited announcement)
**Version:** `v1.0.0-rc1`, `v1.0.0-rc2`
**Artifacts:**
- Binaries for 3 platforms
- Homebrew formula (tap)
- GitHub pre-release tag
- Installation methods: Binary + Homebrew + go install

**Purpose:**
- Final testing before v1.0
- Bug bash (day_81)
- UAT (day_85-86)
- Security audit verification (day_73)

**Timeline:**
- RC1: day_81 (build + internal test)
- RC2: day_82-83 (if bugs found, rebuild)
- Final RC: day_84 (ready for v1.0)

**Success criteria:**
- No P0 bugs
- <3 P1 bugs (fixed before v1.0)
- UAT success (see UAT_PLAN.md)

**Implementacja:** day_81 (RC build)

**Related risk:** R-F3-005 (RC quality issues)

---

### Stage 4: v1.0 Production Release (day_87)

**Audience:** Public
**Version:** `v1.0.0`
**Artifacts:**
- Binaries: Linux (x86_64, ARM64), macOS (Intel, ARM), Windows (x86_64)
- Packages: .deb (Ubuntu/Debian), Homebrew formula
- Checksums: SHA256 + GPG signature (v1.1)
- Documentation: README, USER_GUIDE, CLI_REFERENCE, INSTALLATION
- Examples: 3-5 sample projects

**Distribution channels:**
- GitHub Releases (primary)
- Homebrew tap: `docflow/tap`
- Install script: `https://install.docflow.io`
- go install: `github.com/docflow/docflow/cmd/docflow@v1.0.0`

**Announcement:**
- Blog post (launch announcement)
- Social media (Twitter, LinkedIn, Reddit)
- Mailing list (if exists)
- Hacker News / Product Hunt (optional)

**Success criteria:**
- All deployment methods tested
- Installation guide verified on 3 platforms
- Launch checklist 100% complete (see EXTENDED_PLAN.md day_86)

**Implementacja:** day_87 (launch day)

---

### Stage 5: Post-Launch (day_88-90)

**Monitoring:**
- GitHub: Stars, forks, issues, PRs
- Installation stats: Download counts (GitHub API)
- User feedback: Issues, discussions
- Bugs: Hotfix release if P0/P1 found

**Hotfix policy:**
- P0 (critical, security): Hotfix release within 24h (`v1.0.1`)
- P1 (high, major bug): Hotfix release within 1 week (`v1.0.1` or `v1.0.2`)
- P2 (medium): Batch in next minor release (`v1.1.0`)

**Implementacja:** day_88-90 (monitoring, hotfixes)

**Related document:** WORK_LOG.md (track post-launch issues)

---

## DEPLOYMENT CHECKLIST

### Pre-Release (day_86)

**Build & Artifacts:**
- [ ] All platform binaries built (3 platforms, 4 architectures)
- [ ] Binaries tested on target platforms (manual smoke test)
- [ ] Checksums generated and verified
- [ ] .deb package built and tested (Ubuntu 22.04)
- [ ] Homebrew formula created and tested (macOS Intel + ARM)
- [ ] GitHub Release draft created

**Documentation:**
- [ ] README.md updated (installation methods)
- [ ] INSTALLATION.md complete (all methods documented)
- [ ] USER_GUIDE.md reviewed
- [ ] CLI_REFERENCE.md updated
- [ ] CHANGELOG.md finalized
- [ ] Release notes written

**Testing:**
- [ ] All tests pass (unit + integration + e2e)
- [ ] Performance benchmarks pass (see PERFORMANCE_REQUIREMENTS.md)
- [ ] Security audit complete, no P0 issues (day_73)
- [ ] UAT complete, success criteria met (see UAT_PLAN.md)
- [ ] Examples tested (3-5 projects)

**Infrastructure:**
- [ ] CI/CD pipeline green (all jobs passing)
- [ ] Install script tested (Linux, macOS)
- [ ] GitHub Releases permissions configured
- [ ] Homebrew tap repository ready

**Approvals:**
- [ ] Tech Lead: Code quality approved
- [ ] PM: Release scope approved
- [ ] QA: Testing complete, no blockers
- [ ] Stakeholders: Demo approved

**Source:** EXTENDED_PLAN.md day_86 (LAUNCH_CHECKLIST.md)

---

### Release Day (day_87)

**Actions:**
- [ ] Tag release: `git tag -a v1.0.0 -m "Release v1.0.0"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] Publish GitHub Release (draft → published)
- [ ] Update Homebrew formula (version, checksums)
- [ ] Push Homebrew tap repository
- [ ] Test installation methods (spot check)
- [ ] Announce release (blog, social, mailing list)
- [ ] Monitor initial feedback (first 4 hours)

---

### Post-Release (day_88+)

**Monitoring:**
- [ ] GitHub stars/forks/issues tracked
- [ ] Installation stats reviewed daily (first week)
- [ ] Bugs triaged within 24h
- [ ] P0 bugs → hotfix within 24h
- [ ] P1 bugs → schedule for v1.0.1 or v1.1.0
- [ ] User feedback collected (GitHub Discussions)

**Maintenance:**
- [ ] WORK_LOG.md updated daily (day_88-90)
- [ ] Retrospective conducted (day_89)
- [ ] Roadmap updated (day_90)
- [ ] Backlog prioritized (day_90)

---

## ROLLBACK PLAN

### Scenario: Critical Issue Found Post-Release

**Trigger:** P0 bug discovered within 24h of v1.0.0 release (security, data loss, crash on startup)

**Actions:**
1. **Immediate:**
   - Add warning to GitHub Release page
   - Update README.md with known issue + workaround
   - Create hotfix branch from v1.0.0 tag
2. **Within 4 hours:**
   - Fix issue, add regression test
   - Build hotfix binaries
   - Tag v1.0.1
3. **Within 24 hours:**
   - Release v1.0.1 (hotfix)
   - Update all distribution channels
   - Announce fix (GitHub, social)

**Prevention:**
- Thorough RC testing (day_81-84)
- UAT (day_85-86)
- Security audit (day_73)

**Related risk:** R-F3-005 (RC quality issues)

---

## METRYKI SUKCESU DEPLOYMENT

**Installation success rate:**
- Target: >95% users can install within 5 minutes (binary method)
- Measurement: User feedback, GitHub issues

**Platform coverage:**
- Target: 90% users covered (Linux + macOS)
- Measurement: Download stats per platform

**Update adoption:**
- Target: 50% users update to v1.0.1 within 2 weeks (if released)
- Measurement: Version check API stats (opt-in telemetry, v1.1)

**Time to first value:**
- Target: <10 minutes from install to first scan
- Measurement: Quickstart tutorial completion (user survey)

---

## MAPOWANIE DO EXTENDED_PLAN

| Deployment Activity | Day(s) | Deliverable | Owner |
|---------------------|--------|-------------|-------|
| CI pipeline setup | day_61-62 | `.github/workflows/ci.yml` | Tech Lead |
| Cross-compilation | day_63 | Multi-platform binaries | Tech Lead |
| Deployment automation | day_63-64 | Release pipeline, install script, Homebrew formula | DevOps |
| Examples testing | day_68-69 | 3-5 example projects tested | QA |
| RC build | day_81 | v1.0.0-rc1 binaries | Release Manager |
| Bug fixes | day_82-83 | v1.0.0-rc2 (if needed) | Developers |
| Pre-launch checklist | day_86 | LAUNCH_CHECKLIST.md verified | PM |
| v1.0 release | day_87 | v1.0.0 published, announced | Release Manager |
| Post-launch monitoring | day_88-90 | Hotfixes, WORK_LOG updates | Team |

---

## MAPOWANIE DO RISK_REGISTER

| Risk ID | Risk Name | Deployment Impact | Mitigation Reference |
|---------|-----------|-------------------|----------------------|
| R-F3-001 | Security vulnerability | Cannot release until fixed | Security audit (day_73), GPG signing (v1.1) |
| R-F3-004 | CI/CD pipeline failures | Delays release builds | CI tests (day_61-62), manual fallback |
| R-F3-005 | RC quality issues | Delays v1.0 launch | UAT (day_85-86), bug bash (day_81) |

---

## CHANGELOG

### Version 1.0 (2026-02-06)
- Initial deployment strategy
- Target platforms: Linux, macOS, Windows
- Installation methods: Binary, package managers, go install, install script
- XDG Base Directory compliance
- Manual update mechanism (v1.0)
- Backwards compatibility policy defined
- Deployment stages: Dev → Alpha → RC → v1.0 → Post-launch

---

## NEXT STEPS

### Before v1.0 Launch:
1. ✓ Define deployment strategy (this document)
2. Implement CI pipeline (day_61-62)
3. Setup cross-compilation (day_63)
4. Create install script + Homebrew formula (day_64)
5. Test all installation methods (day_86)
6. Execute launch checklist (day_86-87)

### Post v1.0 (Roadmap):
1. Self-update command (v1.1)
2. APT repository hosting (v1.1)
3. .rpm packages (v1.1)
4. Chocolatey package (v1.2)
5. Docker image (v1.2)
6. GPG signing (v1.1)

---

**END OF DEPLOYMENT_STRATEGY.md**
