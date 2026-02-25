#!/usr/bin/env bash
# scripts/sync-labels.sh
# Idempotentny skrypt synchronizujący labele GitHuba z .github/labels.yml
# Wymaga: gh CLI (zalogowany), python3 lub yq
#
# Użycie:
#   ./scripts/sync-labels.sh                        # używa bieżącego repo
#   REPO=owner/repo ./scripts/sync-labels.sh        # wskaż inne repo
#   DRY_RUN=1 ./scripts/sync-labels.sh              # podgląd bez zmian

set -euo pipefail

REPO="${REPO:-$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null)}"
LABELS_FILE="$(git rev-parse --show-toplevel)/.github/labels.yml"
DRY_RUN="${DRY_RUN:-0}"

if [[ -z "$REPO" ]]; then
  echo "ERROR: nie można ustalić repo. Ustaw REPO=owner/repo lub uruchom z katalogu z git remote." >&2
  exit 1
fi

if [[ ! -f "$LABELS_FILE" ]]; then
  echo "ERROR: nie znaleziono $LABELS_FILE" >&2
  exit 1
fi

echo "Repo: $REPO"
echo "Plik: $LABELS_FILE"
[[ "$DRY_RUN" == "1" ]] && echo "Tryb DRY-RUN (bez zmian)"

# Parsuj YAML przez python3 (bez dodatkowych zależności)
python3 - "$LABELS_FILE" <<'PYEOF'
import sys, json, re

def parse_labels_yml(path):
    labels = []
    current = {}
    with open(path) as f:
        for line in f:
            line = line.rstrip('\n')
            m = re.match(r'^  - name:\s+"?([^"]+)"?', line)
            if m:
                if current:
                    labels.append(current)
                current = {"name": m.group(1)}
            m = re.match(r'^    color:\s+"?([^"]+)"?', line)
            if m and current:
                current["color"] = m.group(1).lstrip('#')
            m = re.match(r'^    description:\s+"?([^"]+)"?', line)
            if m and current:
                current["description"] = m.group(1).rstrip('"')
    if current:
        labels.append(current)
    print(json.dumps(labels))

parse_labels_yml(sys.argv[1])
PYEOF
) | python3 -c "
import sys, json, subprocess, os

labels = json.load(sys.stdin)
repo = os.environ['REPO']
dry = os.environ.get('DRY_RUN', '0') == '1'

# Pobierz istniejące labele
existing_raw = subprocess.check_output(['gh', 'label', 'list', '--repo', repo, '--json', 'name,color,description', '--limit', '200'])
existing = {l['name']: l for l in json.loads(existing_raw)}

created = updated = skipped = 0
for label in labels:
    name = label['name']
    color = label.get('color', 'ededed')
    desc = label.get('description', '')
    if name in existing:
        e = existing[name]
        if e['color'].lower() == color.lower() and e.get('description','') == desc:
            print(f'  skip   {name}')
            skipped += 1
        else:
            print(f'  update {name}')
            if not dry:
                subprocess.run(['gh', 'label', 'edit', name, '--repo', repo,
                                '--color', color, '--description', desc], check=True)
            updated += 1
    else:
        print(f'  create {name}')
        if not dry:
            subprocess.run(['gh', 'label', 'create', name, '--repo', repo,
                            '--color', color, '--description', desc], check=True)
        created += 1

print(f'\nGotowe: {created} utworzono, {updated} zaktualizowano, {skipped} bez zmian.')
"
