#!/usr/bin/env bash
set -euo pipefail
NORMALIZE=0
FIX_FLOW=0
NO_OVERWRITE=0
MERGE_SECTIONS=0
EXPORT_PATCH=0
PATCH_PATH="/tmp/import_docx.patch"
STATUS_JSON=""
SECTION_DIFF=""

usage() {
  cat <<'USAGE'
import_docx.sh --in file.docx --out-dir DIR [--doc-type proposal] [--status draft] [--prefix PREFIX] [--normalize] [--fix-flow] [--no-overwrite] [--merge-sections] [--export-patch [PATH]] [--status-json PATH]

Opcje:
  --in         ścieżka do pliku DOCX (wymagane)
  --out-dir    katalog wyjściowy (wymagane)
  --doc-type   doc_type w frontmatter (domyślnie proposal)
  --status     status w frontmatter (domyślnie draft)
  --prefix     prefiks doc_id (np. proj1_); doc_id = <prefix><basename_without_ext>
  --normalize  po imporcie uruchom sections_normalize.sh --apply na pliku
  --fix-flow   po imporcie uruchom fix_flow.sh (TODO w _todo/)
  --no-overwrite  jeżeli plik istnieje, nie nadpisuj (exit 0, zapis statusu)
  --merge-sections jeżeli plik istnieje, dopisz brakujące sekcje bez usuwania istniejących
  --export-patch [PATH] wygeneruj patch (diff) pomiędzy stanem przed i po imporcie (domyślnie /tmp/import_docx.patch)
  --status-json PATH zapisz raport statusu w JSON
  --section-diff PATH zapis raportu sekcji (braki/dodane) do pliku MD/LOG

Wymagania: python3 (standard libs zipfile/xml.etree), brak pandoc. Media są pomijane (logowane ostrzeżenie w stderr).
USAGE
}

DOC_TYPE=proposal
STATUS=draft
PREFIX=""
IN=""
OUT_DIR=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --in) IN="$2"; shift 2;;
    --out-dir) OUT_DIR="$2"; shift 2;;
    --doc-type) DOC_TYPE="$2"; shift 2;;
    --status) STATUS="$2"; shift 2;;
    --prefix) PREFIX="$2"; shift 2;;
    --normalize) NORMALIZE=1; shift;;
    --fix-flow) FIX_FLOW=1; shift;;
    --no-overwrite) NO_OVERWRITE=1; shift;;
    --merge-sections) MERGE_SECTIONS=1; shift;;
    --export-patch)
      EXPORT_PATCH=1
      if [[ $# -ge 2 && "$2" != --* ]]; then
        PATCH_PATH="$2"; shift 2;
      else
        shift;
      fi;;
    --status-json) STATUS_JSON="$2"; shift 2;;
    --section-diff) SECTION_DIFF="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Nieznana opcja: $1" >&2; usage; exit 1;;
  esac
done

if [[ -z "$IN" || -z "$OUT_DIR" ]]; then
  usage; exit 1
fi

if [[ ! -f "$IN" ]]; then
  echo "Brak pliku: $IN" >&2; exit 1
fi

mkdir -p "$OUT_DIR"
if [[ -n "$STATUS_JSON" ]]; then
  mkdir -p "$(dirname "$STATUS_JSON")"
fi

BASEPATH=$(basename "$IN")
BASENAME=$(basename "$IN" .docx)
DOC_ID="${PREFIX}${BASENAME//[^A-Za-z0-9_-]/_}"
DOC_ID=$(echo "$DOC_ID" | tr '[:upper:]' '[:lower:]')
DOC_ID=${DOC_ID//__/_}
DOC_ID=${DOC_ID//-/_}
OUT_MD="$OUT_DIR/${DOC_ID}.md"

echo "[import_docx] source=$BASEPATH doc_id=$DOC_ID"

TMP_MD=$(mktemp)

python3 - "$IN" "$TMP_MD" "$DOC_ID" "$DOC_TYPE" "$STATUS" <<'PY'
import sys, zipfile, xml.etree.ElementTree as ET, textwrap, re, pathlib, yaml
src, out_md, doc_id, doc_type, status = sys.argv[1:6]
missing_media = []
try:
    with zipfile.ZipFile(src) as z:
        xml = z.read('word/document.xml')
        missing_media = [name for name in z.namelist() if name.startswith('word/media/')]
except Exception as e:
    sys.stderr.write(f"[import_docx] read error: {e}\n")
    sys.exit(1)

# load section_aliases from docflow.yaml (repo root)
repo_root = pathlib.Path(__file__).resolve().parents[1]
cfg_path = repo_root / "docflow.yaml"
rev = {}
if cfg_path.exists():
    cfg = yaml.safe_load(cfg_path.read_text(encoding="utf-8")) or {}
    for canon, vals in (cfg.get("section_aliases") or {}).items():
        rev[canon] = canon
        for v in vals:
            rev[v] = canon

ns = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
root = ET.fromstring(xml)
paras = []
headings = []
for p in root.findall('.//w:body/w:p', ns):
    texts = []
    level = None
    for child in p.iter():
        if child.tag == f"{{{ns['w']}}}pStyle":
            style = child.attrib.get(f"{{{ns['w']}}}val")
            if style and style.lower().startswith("heading"):
                try:
                    level = int(style[-1])
                except ValueError:
                    level = None
        if child.tag == f"{{{ns['w']}}}t":
            texts.append(child.text or '')
        elif child.tag == f"{{{ns['w']}}}tab":
            texts.append('\t')
        elif child.tag == f"{{{ns['w']}}}br":
            texts.append('\n')
    if texts:
        txt = ''.join(texts).strip()
        if level == 1:
            heading = txt
            canon = rev.get(heading, heading) if rev else heading
            paras.append(f"## {canon}")
            headings.append(canon)
        elif level == 2:
            heading = txt
            canon = rev.get(heading, heading) if rev else heading
            paras.append(f"### {canon}")
        else:
            paras.append(txt)
md_body = '\n\n'.join(paras)
front_matter = textwrap.dedent(f'''\
---
title: "{doc_id.replace('_',' ').title()}"
doc_id: {doc_id}
doc_type: {doc_type}
status: {status}
version: "v0.1.0"
owner: "doc_owner"
context_sources: ["docx_import"]
---

## Wstęp

''')
with open(out_md, 'w', encoding='utf-8') as f:
    f.write(front_matter)
    f.write(md_body)

if missing_media:
    sys.stderr.write(f"[import_docx] WARN: media not extracted ({len(missing_media)} files)\n")
sys.stderr.write(f"[import_docx] written {out_md}\n")
PY

STATUS="created"
ADDED_SECTIONS="[]"
MERGED_COUNT=0
SECTION_DIFF_LOG=""

# helper: extract heading names from file
extract_headings() {
  python3 - "$1" <<'PY'
import sys, re, json, pathlib
path = pathlib.Path(sys.argv[1])
text = path.read_text(encoding='utf-8') if path.exists() else ''
heads = []
for line in text.splitlines():
    m = re.match(r'^## +(.+)', line)
    if m:
        heads.append(m.group(1).strip())
print(json.dumps(heads))
PY
}

if [[ -f "$OUT_MD" ]]; then
  # conflict path
  if [[ "$NO_OVERWRITE" -eq 1 && "$MERGE_SECTIONS" -eq 0 ]]; then
    STATUS="skipped_existing"
    cp "$OUT_MD" "${OUT_MD}.orig" 2>/dev/null || true
  elif [[ "$MERGE_SECTIONS" -eq 1 ]]; then
    cp "$OUT_MD" "${OUT_MD}.orig" 2>/dev/null || true
    ADDED_JSON=$(python3 - "$OUT_MD" "$TMP_MD" <<'PY'
import sys, re, pathlib, json
orig_path = pathlib.Path(sys.argv[1])
new_path = pathlib.Path(sys.argv[2])
orig_txt = orig_path.read_text(encoding='utf-8')
new_txt = new_path.read_text(encoding='utf-8')

def split_sections(txt):
    lines = txt.splitlines()
    sections = {}
    current = None
    buf = []
    for line in lines:
        m = re.match(r'^## +(.+)', line)
        if m:
            if current is not None:
                sections[current] = '\n'.join(buf).strip('\n')
            current = m.group(1).strip()
            buf = [line]
        else:
            if current is None:
                current = '__preamble__'
                buf.append(line)
            else:
                buf.append(line)
    if current is not None:
        sections[current] = '\n'.join(buf).strip('\n')
    return sections

orig_sections = split_sections(orig_txt)
new_sections = split_sections(new_txt)

added = []
for name, body in new_sections.items():
    if name == '__preamble__':
        continue
    if any(name.lower() == k.lower() for k in orig_sections.keys() if k != '__preamble__'):
        continue
    added.append(name)
    orig_sections[name] = body

order = [k for k in orig_sections.keys() if k != '__preamble__']
out_lines = []
if '__preamble__' in orig_sections:
    out_lines.append(orig_sections['__preamble__'])
for name in order:
    out_lines.append(orig_sections[name])
out_txt = '\n\n'.join([p for p in out_lines if p.strip() != '']) + '\n'
orig_path.write_text(out_txt, encoding='utf-8')
print(json.dumps(added))
PY
)
    ADDED_SECTIONS="$ADDED_JSON"
    STATUS="merged"
    if [[ "$ADDED_SECTIONS" == "[]" ]]; then
      STATUS="unchanged"
    fi
    MERGED_COUNT=$(ADDED_SECTIONS="$ADDED_JSON" python3 - <<'PY'
import json, os
print(len(json.loads(os.environ.get('ADDED_SECTIONS','[]'))))
PY
)
  else
    echo "[import_docx] ERROR: plik już istnieje (${OUT_MD}). Użyj --no-overwrite lub --merge-sections" >&2
    rm -f "$TMP_MD"
    exit 1
  fi
else
  mv "$TMP_MD" "$OUT_MD"
fi

if [[ ! -f "$OUT_MD" ]]; then
  mv "$TMP_MD" "$OUT_MD"
fi

if [[ "$EXPORT_PATCH" -eq 1 ]]; then
  ORIG_TMP=$(mktemp)
  if [[ -f "${OUT_MD}.orig" ]]; then
    cp "${OUT_MD}.orig" "$ORIG_TMP"
  else
    cp /dev/null "$ORIG_TMP"
  fi
  diff -u "$ORIG_TMP" "$OUT_MD" > "$PATCH_PATH" || true
  rm -f "$ORIG_TMP" "${OUT_MD}.orig" 2>/dev/null || true
fi

if [[ -n "$STATUS_JSON" ]]; then
  MERGE_FLAG=$MERGE_SECTIONS NO_FLAG=$NO_OVERWRITE ADDED="$ADDED_SECTIONS" \
  python3 - <<PY
import json, pathlib, os
path = pathlib.Path("${STATUS_JSON}")
path.parent.mkdir(parents=True, exist_ok=True)
added = json.loads(os.environ.get("ADDED","[]"))
data = {
    "doc_id": "${DOC_ID}",
    "out_md": "${OUT_MD}",
    "status": "${STATUS}",
    "merge": bool(int(os.environ.get("MERGE_FLAG","0"))),
    "no_overwrite": bool(int(os.environ.get("NO_FLAG","0"))),
    "added_sections": added,
    "added_count": len(added),
    "unchanged": "${STATUS}" == "unchanged",
    "patch_path": "${PATCH_PATH}" if ${EXPORT_PATCH} else None
}
path.write_text(json.dumps(data, indent=2), encoding='utf-8')
print(f"[import_docx] status -> {path}")
PY
fi

# section diff report
if [[ -n "$SECTION_DIFF" ]]; then
  python3 - <<PY
import json, pathlib, re, sys
out_path = pathlib.Path("${SECTION_DIFF}")
out_path.parent.mkdir(parents=True, exist_ok=True)
status_path = pathlib.Path("${STATUS_JSON}") if "${STATUS_JSON}" else None
status = {}
if status_path and status_path.exists():
    try:
        status = json.loads(status_path.read_text())
    except Exception:
        status = {}

added = status.get("added_sections", []) if status else []
unchanged = status.get("status") == "unchanged" if status else False

lines = ["# Raport sekcji (import_docx)", ""]
lines.append(f"Plik: {status.get('out_md','')}" if status else "Plik: (brak danych)")
if unchanged:
    lines.append("Brak braków – dokument zgodny z kanonem (unchanged).")
elif added:
    lines.append("Dodane sekcje:")
    for s in added:
        lines.append(f"- {s}")
else:
    lines.append("Braki sekcji: brak danych lub nie wykryto dodanych sekcji (sprawdź status JSON).")

out_path.write_text("\\n".join(lines) + "\\n", encoding="utf-8")
print(f"[import_docx] section diff -> {out_path}")
PY
fi

echo "$OUT_MD"

if [[ "$NORMALIZE" -eq 1 ]]; then
  ./scripts/sections_normalize.sh --root "$(dirname "$OUT_MD")" --apply --log "${OUT_MD}.normalize.log"
fi

if [[ "$FIX_FLOW" -eq 1 ]]; then
  ./scripts/fix_flow.sh --root "$OUT_DIR" --out-dir "_todo" --log LOGS/fix_flow.md || true
fi

rm -f "$TMP_MD"
