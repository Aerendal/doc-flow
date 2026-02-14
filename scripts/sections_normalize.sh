#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
sections_normalize.sh --root DIR [--apply] [--doc-type-filter TYPE] [--log LOGFILE]

Normalizuje nazwy sekcji H2 w plikach .md do kanonu z docflow.yaml (section_aliases, YAML).

Opcje:
  --root DIR            katalog z .md (wymagane)
  --apply               wykonaj zmiany (domyślnie dry-run)
  --doc-type-filter T   przetwarzaj tylko pliki z doc_type==T (opcjonalne)
  --log LOGFILE         ścieżka do logu (opcjonalnie)

Zachowanie:
  - Zamienia aliasy na kanon (wg YAML section_aliases).
  - Nie zmienia kolejności sekcji.
  - Pomija pliki bez frontmattera/doc_type.

USAGE
}

ROOT=""
APPLY=0
DOC_FILTER=""
LOGFILE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --apply) APPLY=1; shift;;
    --doc-type-filter) DOC_FILTER="$2"; shift 2;;
    --log) LOGFILE="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
  esac
done

if [[ -z "$ROOT" ]]; then
  usage; exit 1
fi

python3 - <<'PY' "$ROOT" "$APPLY" "$DOC_FILTER" "$LOGFILE"
import sys, pathlib, re, yaml, io
root = pathlib.Path(sys.argv[1])
apply = sys.argv[2] == "1"
doc_filter = sys.argv[3]
logfile = sys.argv[4]

log_lines = []
def log(msg):
    print(msg)
    if logfile:
        log_lines.append(msg)

cfg_path = pathlib.Path.cwd() / "docflow.yaml"
if not cfg_path.exists():
    sys.stderr.write(f"missing docflow.yaml at {cfg_path}\n")
    sys.exit(1)
cfg = yaml.safe_load(cfg_path.read_text(encoding="utf-8")) or {}
aliases = cfg.get("section_aliases") or {}
rev = {}
for canon, vals in aliases.items():
    rev[canon] = canon
    for v in vals:
        rev[v] = canon

def split_frontmatter(text):
    if not text.startswith("---"):
        return None, text
    parts = text.split("---", 2)
    if len(parts) < 3:
        return None, text
    fm = parts[1]
    body = parts[2]
    return fm, body

def doc_type_from_fm(fm):
    for line in fm.splitlines():
        if line.startswith("doc_type:"):
            return line.split(":",1)[1].strip().strip('"')
    return None

changed = 0
processed = 0
files_changed = []
for md in sorted(root.rglob("*.md")):
    txt = md.read_text(encoding="utf-8")
    res = split_frontmatter(txt)
    if res is None:
        log(f"[SKIP] {md} (brak frontmattera)")
        continue
    fm, body = res
    doc_type = doc_type_from_fm(fm)
    if doc_filter and doc_type != doc_filter:
        continue
    if doc_type is None:
        log(f"[SKIP] {md} (brak doc_type)")
        continue
    processed += 1
    lines = body.splitlines()
    new_lines = []
    file_changed = False
    for ln in lines:
        m = re.match(r"^(#+)\\s+(.*)$", ln)
        if m and len(m.group(1)) == 2:  # H2
            name = m.group(2).strip()
            if name in rev and rev[name] != name:
                canon = rev[name]
                new_lines.append(f"## {canon}")
                file_changed = True
                log(f"[ALIAS] {md}: '{name}' -> '{canon}'")
            else:
                new_lines.append(ln)
        else:
            new_lines.append(ln)
    if file_changed:
        changed += 1
        files_changed.append(md)
        if apply:
            buf = io.StringIO()
            buf.write("---\n")
            buf.write(fm)
            buf.write("---")
            buf.write("\n".join(new_lines))
            md.write_text(buf.getvalue(), encoding="utf-8")

log(f"Processed={processed}, Changed={changed}, apply={apply}")
if logfile:
    pathlib.Path(logfile).write_text("\n".join(log_lines + [f"SUMMARY: processed={processed}, changed={changed}, apply={apply}"]), encoding="utf-8")
PY
