#!/usr/bin/env bash
# scripts/sync-labels.sh
# Idempotentny skrypt synchronizujacy labele GitHuba z .github/labels.yml
#
# Uzycie:
#   ./scripts/sync-labels.sh                   # uzywa biezacego repo
#   REPO=owner/repo ./scripts/sync-labels.sh   # wskaż inne repo
#   DRY_RUN=1 ./scripts/sync-labels.sh         # podglad bez zmian

set -euo pipefail

REPO="${REPO:-$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null)}"
DRY_RUN="${DRY_RUN:-0}"

if [[ -z "$REPO" ]]; then
  echo "ERROR: nie mozna ustalic repo. Ustaw REPO=owner/repo lub uruchom z katalogu z git remote." >&2
  exit 1
fi

echo "Repo: $REPO"
[[ "$DRY_RUN" == "1" ]] && echo "Tryb DRY-RUN (bez zmian)"

ensure_label() {
  local name="$1" color="$2" desc="$3"
  if gh label list --repo "$REPO" --search "$name" --json name --jq '.[].name' | grep -Fxq "$name"; then
    if [[ "$DRY_RUN" != "1" ]]; then
      gh label edit "$name" --repo "$REPO" --color "$color" --description "$desc" >/dev/null
    fi
    echo "  update: $name"
  else
    if [[ "$DRY_RUN" != "1" ]]; then
      gh label create "$name" --repo "$REPO" --color "$color" --description "$desc" >/dev/null
    fi
    echo "  create: $name"
  fi
}

# type/*
ensure_label "type/bug"      "d73a4a" "Potwierdzony blad lub niepoprawne zachowanie"
ensure_label "type/feature"  "a2eeef" "Prosba o nowa funkcjonalnosc lub ulepszenie"
ensure_label "type/question" "d876e3" "Pytanie - preferuj GitHub Discussions"
ensure_label "type/docs"     "0075ca" "Zmiana lub blad w dokumentacji"
ensure_label "type/security" "b60205" "Bezpieczenstwo - nie ujawniaj publicznie"

# status/*
ensure_label "status/needs-triage" "fbca04" "Nowe zgloszenie, wymaga oceny"
ensure_label "status/needs-info"   "fef2c0" "Brakuje srodowiska, logow lub szczegolów"
ensure_label "status/needs-repro"  "fef2c0" "Potrzebne minimalne kroki reprodukcji"
ensure_label "status/accepted"     "0e8a16" "Zaakceptowane do implementacji"
ensure_label "status/duplicate"    "cfd3d7" "Duplikat innego zgloszenia"
ensure_label "status/wontfix"      "ffffff" "Nie bedzie naprawiane / poza zakresem"
ensure_label "status/blocked"      "e99695" "Zablokowane przez zaleznosc lub decyzje"
ensure_label "status/in-progress"  "1d76db" "Praca w toku"
ensure_label "status/review"       "1d76db" "Oczekuje na review PR"

# prio/*
ensure_label "prio/p0" "b60205" "Krytyczny: bezpieczenstwo, utrata danych, crash"
ensure_label "prio/p1" "d93f0b" "Wysoki: blokuje instalacje lub CI"
ensure_label "prio/p2" "fbca04" "Normalny: wazny, ale jest obejscie"
ensure_label "prio/p3" "cfd3d7" "Niski: nice-to-have"

# area/*
ensure_label "area/cli"      "5319e7" "Komendy CLI, flagi, UX terminala"
ensure_label "area/install"  "0052cc" "Instalacja, pakowanie, release"
ensure_label "area/ci"       "0052cc" "GitHub Actions, workflows, CI/CD"
ensure_label "area/indexing" "5319e7" "Logika indeksowania i ingestion"
ensure_label "area/sqlite"   "5319e7" "SQLite, FTS, schematy, zapytania"
ensure_label "area/docs"     "0075ca" "Dokumentacja i przyklady"
ensure_label "area/security" "b60205" "Polityka i procesy bezpieczenstwa"

# meta/*
ensure_label "meta/feedback-week"    "bfdadc" "Zgloszenie z tygodnia feedbacku"
ensure_label "meta/good-first-issue" "7057ff" "Dobre pierwsze zadanie dla nowego kontrybutora"
ensure_label "meta/help-wanted"      "008672" "Poszukujemy pomocy / mile widziane PR"

echo ""
echo "Gotowe."

