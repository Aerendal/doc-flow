#!/usr/bin/env bash
set -euo pipefail

ROOT="tmp/localfold_import"
POLICY="docs/_meta/IMPORT_POLICY.yaml"
DOCFLOW="docflow.yaml"
LOG="LOGS/import_policy_lint.md"
SUGGEST=0

usage() {
  echo "usage: $0 [--root dir] [--policy file] [--docflow file] [--log file]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --policy) POLICY="$2"; shift 2;;
    --docflow) DOCFLOW="$2"; shift 2;;
    --log) LOG="$2"; shift 2;;
    --suggest) SUGGEST=1; shift;;
    -h|--help) usage;;
    *) usage;;
  esac
done

[ -f "$POLICY" ] || { echo "missing policy $POLICY" >&2; exit 2; }
[ -f "$DOCFLOW" ] || { echo "missing docflow $DOCFLOW" >&2; exit 2; }
[ -d "$ROOT" ] || { echo "missing root dir $ROOT" >&2; exit 2; }

mkdir -p "$(dirname "$LOG")"

fail_count=$(python - "$POLICY" "$DOCFLOW" "$ROOT" "$LOG" "$SUGGEST" <<'PY'
import sys, pathlib, yaml, re, json, textwrap
policy_path, docflow_path, root, log_path, suggest_flag = sys.argv[1:6]
suggest_flag = int(suggest_flag)
policy = yaml.safe_load(pathlib.Path(policy_path).read_text()) or {}
docflow = yaml.safe_load(pathlib.Path(docflow_path).read_text()) or {}
aliases = docflow.get("section_aliases", {}) or {}

required = policy.get("required_sections", []) or []
forbidden = set(s.lower() for s in (policy.get("forbidden_sections") or []))

norm = {}
for canon in required + list(aliases.keys()):
    norm[canon.lower()] = canon
    for a in aliases.get(canon, []) or []:
        norm[a.lower()] = canon

def normalize(name: str):
    key = name.strip().strip("#").strip()
    low = key.lower()
    return norm.get(low, key)

def headings(md_path: pathlib.Path):
    hs=[]
    for line in md_path.read_text().splitlines():
        if line.startswith("#"):
            title=line.lstrip("#").strip()
            if title:
                hs.append(title)
    return hs

def doc_type(md_path: pathlib.Path):
    txt = md_path.read_text()
    if not txt.startswith("---"):
        return None
    parts = txt.split("---",2)
    if len(parts)<3:
        return None
    fm = parts[1]
    for line in fm.splitlines():
        if line.startswith("doc_type:"):
            return line.split(":",1)[1].strip()
    return None

violations=[]
files = sorted(pathlib.Path(root).glob("*.md"))
for md in files:
    dt = doc_type(md)
    if dt == "tasklist":
        violations.append({
            "file": md.as_posix(),
            "status": "PASS",
            "missing": [],
            "forbidden": [],
            "headings": headings(md),
        })
        continue
    hs = headings(md)
    normalized = [normalize(h) for h in hs]
    found = set(normalized)
    missing = [r for r in required if r not in found]
    forbid_hits = [h for h in hs if h.lower() in forbidden]
    status = "PASS" if not missing and not forbid_hits else "FAIL"
    violations.append({
        "file": md.as_posix(),
        "status": status,
        "missing": missing,
        "forbidden": forbid_hits,
        "headings": hs,
    })

fail = sum(1 for v in violations if v["status"] != "PASS")

with open(log_path, "w", encoding="utf-8") as f:
    f.write("# Import policy lint\n")
    f.write(f"Root: {root}\n")
    f.write(f"Policy: {policy_path}\n\n")
    for v in violations:
        f.write(f"## {v['file']}\n")
        f.write(f"- Status: {v['status']}\n")
        f.write(f"- Missing: {', '.join(v['missing']) if v['missing'] else 'none'}\n")
        f.write(f"- Forbidden: {', '.join(v['forbidden']) if v['forbidden'] else 'none'}\n")
        f.write(f"- Headings: {', '.join(v['headings']) if v['headings'] else 'none'}\n\n")
        if suggest_flag and v["missing"]:
            f.write("### Hints\n```markdown\n")
            limit = min(len(v["missing"]), 10)
            for i, m in enumerate(v["missing"][:limit], start=1):
                f.write(f"{i}. {m}\n- TODO: uzupełnij sekcję {m} (1-2 zdania, kluczowe decyzje/ryzyka)\n\n")
            if len(v["missing"]) > limit:
                f.write(f"+{len(v['missing'])-limit} więcej...\n")
            f.write("```\n\n")
    f.write(f"Summary: PASS={len(violations)-fail}, FAIL={fail}\n")

print(fail)
PY
)

if [ "$fail_count" -gt 0 ]; then
  exit 2
fi
