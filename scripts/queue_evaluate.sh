#!/usr/bin/env bash
set -euo pipefail

BIN="./build/docflow-linux-amd64"
FORMAT="text" # text|json
INPUT=""
CACHE_FILE=""
LOG_DIR=""
NO_CACHE=0
WORKERS=1
HIST=".docflow/run_history.json"

usage() {
  echo "usage: $0 [--bin path] [--format text|json] [--cache file] [--log-dir dir] [--no-cache] [--workers N] <backlog_file_or_->" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bin) BIN="$2"; shift 2;;
    --format) FORMAT="$2"; shift 2;;
    --cache) CACHE_FILE="$2"; shift 2;;
    --log-dir) LOG_DIR="$2"; shift 2;;
    --no-cache) NO_CACHE=1; shift;;
    --workers) WORKERS="$2"; shift 2;;
    --history) HIST="$2"; shift 2;;
    -h|--help) usage;;
    -*) usage;;
    *) INPUT="$1"; shift;;
  esac
done

[ -z "${INPUT:-}" ] && usage
[ -n "$LOG_DIR" ] && mkdir -p "$LOG_DIR"
mkdir -p "$(dirname "$HIST")"

sha() { sha256sum "$1" 2>/dev/null | awk '{print $1}'; }
read_lines() { if [ "$INPUT" = "-" ]; then cat; else cat "$INPUT"; fi; }
hash_md_files() {
  local dir="$1"
  python - <<'PY' "$dir"
import sys, pathlib, hashlib
root = pathlib.Path(sys.argv[1])
paths = sorted(p for p in root.glob("*.md") if p.is_file())
h = hashlib.sha256()
for p in paths:
    h.update(p.name.encode())
    h.update(b":")
    h.update(hashlib.sha256(p.read_bytes()).digest())
digest = h.hexdigest()
print(digest)
PY
}

# JSON helpers
_dict_get() { python - <<'PY' "$1" "$2"
import json,sys
cache=json.loads(sys.argv[1]) if sys.argv[1] else {}
print(cache.get(sys.argv[2],""))
PY
}
_dict_set() { python - <<'PY' "$1" "$2" "$3"
import json,sys
cache=json.loads(sys.argv[1]) if sys.argv[1] else {}
cache[sys.argv[2]]=json.loads(sys.argv[3]) if sys.argv[3] not in ("", None) else sys.argv[3]
print(json.dumps(cache))
PY
}

cache_json=""
if [ -n "$CACHE_FILE" ] && [ -f "$CACHE_FILE" ]; then cache_json="$(cat "$CACHE_FILE")"; fi

run_task() {
  local id="$1" path="$2" outfile="$3" tmpdir="$4"
  local cache_key="$id:$path"
  local hash_cfg="$(sha "$path/docflow.yaml")"
  local hash_files="$(hash_md_files "$path")"
  local cached_reason="" cached_hash="" cached_vios=""
  local cached_fhash=""
  local cache_status="miss"

  if [ $NO_CACHE -eq 0 ] && [ -n "$cache_json" ]; then
    cached_reason="$(_dict_get "$cache_json" "$cache_key:reason")"
    cached_hash="$(_dict_get "$cache_json" "$cache_key:hash")"
    cached_fhash="$(_dict_get "$cache_json" "$cache_key:fhash")"
    cached_vios="$(_dict_get "$cache_json" "$cache_key:vios")"
  fi

  local reason="$cached_reason"
  local vios_json="$cached_vios"

  if [ $NO_CACHE -eq 1 ]; then
    cache_status="nocache"
  elif [ -n "$cached_fhash" ] && [ "$cached_fhash" = "$hash_files" ]; then
    cache_status="hit"
  else
    cache_status="miss"
  fi

  if [ "$cache_status" != "hit" ]; then
    v_log=${LOG_DIR:+"$LOG_DIR/$id-validate.log"}; v_log=${v_log:-/tmp/queue_validate.log}
    c_log=${LOG_DIR:+"$LOG_DIR/$id-compliance.log"}; c_log=${c_log:-/tmp/queue_compliance.log}
    comp_json="$tmpdir/$id-compliance.json"

    reason=""
    vios_json="[]"

    # Compliance JSON to file
    "$BIN" compliance --config "$path/docflow.yaml" --rules docs/_meta/GOVERNANCE_RULES.yaml --format json --output "$comp_json" >/dev/null 2>"$c_log" || true
    if [ -s "$comp_json" ]; then
      failed=$(python - <<'PY' "$comp_json"
import json,sys
with open(sys.argv[1]) as f:
    data=json.load(f)
print(data.get('failed',0))
PY
)
      if [ "$failed" != "0" ]; then
        reason="compliance_fail"
        vios_json=$(python - <<'PY' "$comp_json"
import json,sys
with open(sys.argv[1]) as f:
    data=json.load(f)
viol=[]
for d in data.get('docs',[]):
    for v in d.get('violations',[]) or []:
        viol.append(str(v))
print(json.dumps(viol))
PY
)
      fi
    fi

    if [ -z "$reason" ]; then
      # fallback validate exit code
      if ! "$BIN" validate --config "$path/docflow.yaml" --governance docs/_meta/GOVERNANCE_RULES.yaml >"$v_log" 2>&1; then
        reason="validate_fail"
        vios_json='["validate_fail"]'
      fi
    fi

    if [ $NO_CACHE -eq 0 ]; then
      cache_json="$(_dict_set "$cache_json" "$cache_key:reason" "\"$reason\"")"
      cache_json="$(_dict_set "$cache_json" "$cache_key:fhash" "\"$hash_files\"")"
      cache_json="$(_dict_set "$cache_json" "$cache_key:vios" "${vios_json:-[]}")"
    fi
  fi

  if [ -z "$reason" ]; then
    echo "READY $id $path" >"$outfile"
    echo "[]" >"$outfile.vios"
  else
    echo "BLOCKED $id $path reason=$reason" >"$outfile"
    echo "${vios_json:-[]}" >"$outfile.vios"
  fi
  echo "$hash_files" >"$outfile.hash"
  echo "$cache_status" >"$outfile.cstat"
}

mapfile -t tasks < <(read_lines | awk 'NF>=2{print $1" "$2}' | LC_ALL=C sort)

pids=()
tmpdir=$(mktemp -d)
sem_fifo=$(mktemp -u); mkfifo "$sem_fifo"
exec 3<>"$sem_fifo"; rm "$sem_fifo"
for ((i=0;i<WORKERS;i++)); do echo >&3; done

for idx in "${!tasks[@]}"; do
  IFS=' ' read -r tid tpath <<< "${tasks[$idx]}"
  read -u 3
  {
    run_task "$tid" "$tpath" "$tmpdir/$idx.out" "$tmpdir"
    echo >&3
  } &
  pids+=($!)
done
wait "${pids[@]}"
exec 3>&-

status=0
json_items=()
ready_count=0
blocked_count=0

for idx in "${!tasks[@]}"; do
  line=$(cat "$tmpdir/$idx.out")
  vios=$(cat "$tmpdir/$idx.out.vios")
  hfile=$(cat "$tmpdir/$idx.out.hash")
  cstat=$(cat "$tmpdir/$idx.out.cstat")
  IFS=' ' read -r state tid tpath reason <<< "$line"
  if [ "$state" = "READY" ]; then
    ready_count=$((ready_count+1))
    [ "$FORMAT" = "text" ] && echo "READY   $tid $tpath cache=$cstat"
    json_items+=("{\"id\":\"$tid\",\"path\":\"$tpath\",\"status\":\"READY\",\"cache_status\":\"$cstat\",\"violations\":[]}")
  else
    blocked_count=$((blocked_count+1))
    [ "$FORMAT" = "text" ] && echo "BLOCKED $tid $tpath reason=$reason cache=$cstat"
    json_items+=("{\"id\":\"$tid\",\"path\":\"$tpath\",\"status\":\"BLOCKED\",\"reason\":\"$reason\",\"cache_status\":\"$cstat\",\"violations\":$vios}")
    status=2
  fi
  if [ $NO_CACHE -eq 0 ]; then
    cache_key="$tid:$tpath"
    cache_json="$(_dict_set "$cache_json" "$cache_key:reason" "\"$reason\"")"
    cache_json="$(_dict_set "$cache_json" "$cache_key:fhash" "\"$hfile\"")"
  fi
done

rm -rf "$tmpdir"

if [ "$FORMAT" = "json" ]; then
  printf '{"ready":%d,"blocked":%d,"tasks":[%s]}
' "$ready_count" "$blocked_count" "$(IFS=,; echo "${json_items[*]}")"
fi

if [ -n "$CACHE_FILE" ]; then
  echo "$cache_json" > "$CACHE_FILE"
fi

# append history (best effort)
python - <<'PY' "$HIST" "$ready_count" "$blocked_count" "$CACHE_FILE"
import sys, json, pathlib, datetime
hist, ready, blocked, cachef = sys.argv[1:5]
entry = {
    "ts": datetime.datetime.now().isoformat(timespec="seconds"),
    "kind": "queue",
    "exit": 0 if int(blocked)==0 else 2,
    "artifact": cachef or "",
    "ready": int(ready),
    "blocked": int(blocked)
}
p = pathlib.Path(hist)
data = json.loads(p.read_text()) if p.exists() and p.read_text().strip() else []
p.parent.mkdir(parents=True, exist_ok=True)
data.append(entry)
p.write_text(json.dumps(data, indent=2))
PY

exit $status
