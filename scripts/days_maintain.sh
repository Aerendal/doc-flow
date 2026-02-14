#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
days_maintain.sh <command> [options]

Commands:
  index                 - generuje worklog/days/index.md (deleguje do scripts/days_index.sh)
  lint                  - uruchamia linter days (cmd/worklog-lint-days)
  check [--threshold N] - liczy pliki markdown; exit 1 jeśli > N (default: 130)
  archive --from A --to B [--out PATH] [--dry-run]
                        - pakuje zakres wpisow A..B do tar.gz (domyślnie archive/days_A_B.tar.gz)
  search --pattern REGEX
                        - szuka REGEX w plikach markdown w worklog/days (rg, fallback grep -nR)

Wspólne:
  -h, --help            - pomoc

Uwagi:
- Nie usuwa źródeł; archiwum to kopia. Usuń/przenieś ręcznie jeśli potrzebujesz.
- Lint/index używają Go modules (`-mod=vendor`); wymagany Go 1.25+ w PATH.
USAGE
}

cmd="${1:-}"
if [[ -z "$cmd" || "$cmd" == "-h" || "$cmd" == "--help" ]]; then
  usage
  exit 0
fi
shift

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

go_env() {
  GOFLAGS=${GOFLAGS:-"-mod=vendor"}
  GOCACHE=${GOCACHE:-"/tmp/go-cache"}
  export GOFLAGS GOCACHE
}

case "$cmd" in
  index)
    ./scripts/days_index.sh
    ;;

  lint)
    go_env
    go run ./cmd/worklog-lint-days --root worklog/days
    ;;

  check)
    threshold=130
    while [[ $# -gt 0 ]]; do
      case "$1" in
        --threshold) threshold="$2"; shift 2;;
        *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
      esac
    done
    count=$(find worklog/days -maxdepth 1 -type f -name '*.md' ! -name 'index.md' | wc -l)
    echo "entries count: $count (threshold: $threshold)"
    if (( count > threshold )); then
      exit 1
    fi
    ;;

  archive)
    from=""; to=""; out=""; dry_run=0
    while [[ $# -gt 0 ]]; do
      case "$1" in
        --from) from="$2"; shift 2;;
        --to) to="$2"; shift 2;;
        --out) out="$2"; shift 2;;
        --dry-run) dry_run=1; shift;;
        *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
      esac
    done
    if [[ -z "$from" || -z "$to" ]]; then
      echo "--from i --to są wymagane" >&2; exit 1
    fi
    out=${out:-"archive/days_${from}_${to}.tar.gz"}
    files=()
    for n in $(seq "$from" "$to"); do
      while IFS= read -r f; do
        files+=("$f")
      done < <(find worklog/days -maxdepth 1 -type f -name "*_${n}.md" | LC_ALL=C sort -V)
    done
    if [[ ${#files[@]} -eq 0 ]]; then
      echo "Brak plików w zadanym zakresie" >&2; exit 1
    fi
    echo "Archiwizacja -> $out (${#files[@]} plików)"
    if (( dry_run )); then
      printf '%s\n' "${files[@]}"
      exit 0
    fi
    mkdir -p "$(dirname "$out")"
    tar -czf "$out" "${files[@]}"
    sha256sum "$out" > "${out}.sha256"
    ;;

  search)
    pattern=""
    while [[ $# -gt 0 ]]; do
      case "$1" in
        --pattern) pattern="$2"; shift 2;;
        *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
      esac
    done
    if [[ -z "$pattern" ]]; then
      echo "--pattern wymagany" >&2; exit 1
    fi
    if command -v rg >/dev/null 2>&1; then
      rg --no-heading --glob '*.md' --glob '!index.md' "$pattern" worklog/days
    else
      find worklog/days -maxdepth 1 -type f -name '*.md' ! -name 'index.md' -print0 \
        | xargs -0 -r grep -nH "$pattern" || true
    fi
    ;;

  *)
    echo "Nieznana komenda: $cmd" >&2
    usage
    exit 1
    ;;
esac
