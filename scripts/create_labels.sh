blank_issues_enabled: false
contact_links:
  - name: "Pytania i pomoc (Discussions)"
    url: "/discussions"
    about: "Używaj Discussions do pytań, pomocy i ogólnego wsparcia."
  - name: "Pomysły / Propozycje (Discussions)"
    url: "/discussions"
    about: "Używaj Discussions do pomysłów i dyskusji projektowych."
  - name: "Zgłoszenia bezpieczeństwa (nie publiczne)"
    url: "/security/policy"
    about: "Zgłaszaj podatności prywatnie zgodnie z polityką bezpieczeństwa. Nie ujawniaj publicznie."
  - name: "Dokumentacja"
    url: "/blob/main/README.md"
    about: "Tutaj: instalacja, użycie, przegląd projektu.","path":
".github/ISSUE_TEMPLATE/config.yml"},{"content":"name: "Zgłoszenie błędu / Bug report"
description: "Zgłoś powtarzalny problem (crash, nieprawidłowy output, błąd workflow). Jeśli to luka bezpieczeństwa — postępuj zgodnie z SECURITY.md (nie publikuj publicznie)."
title: "[BUG] ${{short_description}}"
labels:
  - "type/bug"
  - "status/needs-triage"
body:
  - type: markdown
    attributes:
      value: |
        Dziękuję za zgłoszenie. Proszę wypełnić poniższe pola jak najdokładniej.
        Jeśli podejrzewasz lukę bezpieczeństwa — NIE publikuj szczegółów publicznie; użyj /security/policy.
  - type: input
    id: short_description
    attributes:
      label: "Krótki opis (tytuł)"
      description: "Jasne, krótkie streszczenie problemu."
      placeholder: "np. 'błąd przy instalacji na macOS 13'"
      required: true
  - type: input
    id: version
    attributes:
      label: "Wersja"
      description: "Dokładna wersja lub commit/sha (np. wyjście `doc-flow --version`). Jeśli nieznana, wpisz 'unknown'."
      required: true
  - type: dropdown
    id: install_method
    attributes:
      label: "Sposób instalacji"
      options:
        - "Prebuilt binary (release)"
        - "Zbudowane ze źródeł (Go)"
        - "Docker / kontener"
        - "Package manager"
        - "Inne / nieznane"
      required: true
  - type: dropdown
    id: os
    attributes:
      label: "System operacyjny / środowisko"
      options:
        - "Linux"
        - "macOS"
        - "Windows"
        - "CI (GitHub Actions, inne)"
        - "Inne"
      required: true
  - type: input
    id: os_details
    attributes:
      label: "Szczegóły systemu"
      description: "Distro/wersja (Linux), wersja macOS, Windows build, architektura."
      placeholder: "np. Ubuntu 22.04 x86_64"
      required: true
  - type: input
    id: go_version
    attributes:
      label: "Wersja Go (jeśli budowane ze źródeł)"
      description: "Wyjście `go version` lub 'N/A' jeśli nie dotyczy."
      required: false
  - type: textarea
    id: command
    attributes:
      label: "Dokładne polecenie, które się nie powiodło"
      description: "Wklej dokładną komendę(y)."
      placeholder: "np. doc-flow <subcommand> --flag ..."
      required: true
  - type: textarea
    id: expected
    attributes:
      label: "Oczekiwane zachowanie"
      required: true
  - type: textarea
    id: actual
    attributes:
      label: "Rzeczywiste zachowanie / błąd"
      required: true
  - type: textarea
    id: repro
    attributes:
      label: "Kroki reprodukcji"
      description: "Krok po kroku, numerowane. Jeśli możesz, dołącz minimalne dane wejściowe."
      placeholder: |
        1) ...
        2) ...
        3) ...
      required: true
  - type: textarea
    id: logs
    attributes:
      label: "Logi / output (opcjonalne)"
      description: "Wklej istotne fragmenty logów. Usuń sekrety i prywatne dane."
      required: false
  - type: textarea
    id: additional
    attributes:
      label: "Dodatkowy kontekst (opcjonalne)"
      required: false
","path":".github/ISSUE_TEMPLATE/bug_report.yml"},{"content":"name: "Propozycja funkcji / Feature request"
description: "Zgłaszaj pomysły na nowe funkcje lub rozszerzenia. Dla szerszych dyskusji preferuję Discussions."
title: "[FEATURE] ${{short_summary}}"
labels:
  - "type/feature"
  - "status/needs-triage"
body:
  - type: markdown
    attributes:
      value: |
        Dziękuję za pomysł. Jeśli potrzebujesz szerszej dyskusji projektowej, utwórz Discussion i wklej link tutaj.
  - type: input
    id: short_summary
    attributes:
      label: "Krótki opis pomysłu"
      placeholder: "np. 'opcjonalne logowanie w formacie JSON'"
      required: true
  - type: textarea
    id: motivation
    attributes:
      label: "Motywacja / problem"
      description: "Dlaczego to jest potrzebne?"
      required: true
  - type: textarea
    id: proposal
    attributes:
      label: "Proponowane rozwiązanie"
      description: "Proszę opisać oczekiwane zachowanie, przykładowy UX/CLI/format wyjścia."
      required: true
  - type: textarea
    id: acceptance
    attributes:
      label: "Kryteria akceptacji"
      description: "Kiedy uznamy, że zadanie jest zamknięte? Podaj testowalne punkty."
      required: true
  - type: textarea
    id: alternatives
    attributes:
      label: "Alternatywy (opcjonalne)"
      description: "Co rozważałeś zamiast tego?"
      required: false
  - type: textarea
    id: scope
    attributes:
      label: "Zakres / wpływ (opcjonalne)"
      description: "Jakie obszary projektu to dotyczy (docs/cli/ci/install itd.)."
      required: false
","path":".github/ISSUE_TEMPLATE/feature_request.yml"},{"content":"## Opis zmian
- Co zmieniłem i dlaczego.
- Jeśli zamyka issue: #<numer>

## Rodzaj zmiany
- [ ] Poprawka błędu
- [ ] Nowa funkcja
- [ ] Dokumentacja
- [ ] Refactor / utrzymanie
- [ ] CI / release

## Jak przetestować
Podaj dokładne polecenia i oczekiwany rezultat.

## Kompatybilność
- [ ] Zmiany kompatybilne wstecz albo opisane jako breaking (RELEASE NOTES)

## Checklista
- [ ] Dodałem/zmodyfikowałem testy (jeśli dotyczy)
- [ ] Zaktualizowałem dokumentację (jeśli dotyczy)
- [ ] CI przechodzi lokalnie
","path":".github/PULL_REQUEST_TEMPLATE.md"},{"content":"# Wsparcie / Gdzie pytać

Używam dwóch głównych kanałów wsparcia:

1) GitHub Discussions — preferowane dla pytań, rozwiązywania problemów i pomysłów:
   https://github.com/Aerendal/doc-flow/discussions

2) GitHub Issues — wyłącznie dla reproducible bugów i konkretnych zadań:
   https://github.com/Aerendal/doc-flow/issues

## Co gdzie wstawiać?
Discussions:
- Pytania dotyczące użycia, instalacji, konfiguracji
- Brainstorming i projektowanie funkcji

Issues:
- Reproducible bugs (crashy, nieprawidłowy output, błędy instalacji)
- Konkretne zadania z kryteriami akceptacji

## Security
Nie zgłaszaj podatności przez Issues/Discussions publicznie.
Użyj: https://github.com/Aerendal/doc-flow/security/policy

## Logi
Usuń sekrety (tokeny, klucze). Jeśli przypadkowo opublikujesz sekret — zrotuj go natychmiast.

## Oczekiwania odpowiedzi
Podczas feedback week triage realizuję 2x dziennie. Mogę poprosić o dodatkowe informacje.
","path":"SUPPORT.md"},{"content":"# Zasady zachowania (Code of Conduct)

Chcę, żeby to repo było przyjazne i bezpieczne. Oto podstawowe zasady:

## Zasady
- Traktuję innych z szacunkiem.
- Unikam obraźliwego języka, zastraszania i dyskryminacji.
- Jestem konstruktywny i rzeczowy w feedbacku.

## Niedozwolone
- Nękanie, mowa nienawiści, personalne ataki.
- Trolling, celowe zakłócanie dyskusji.
- Publikowanie prywatnych danych (doxxing).
- Publiczne ujawnianie podatności.

## Zasięg
Dotyczy Issues, Pull Requests, Discussions i innych kanałów projektu.

## Zgłaszanie naruszeń
Proszę zgłaszać naruszenia prywatnie na e-mail: j.j.j.skoczylas@gmail.com
Możesz też użyć narzędzi GitHub do raportowania zachowań.
","path":"CODE_OF_CONDUCT.md"},{"content":"## Tydzień feedbacku (7 dni)

Przez 7 dni zbieram feedback w uporządkowany sposób:

- Błędy: Issues (użyj formularza "Zgłoszenie błędu")
- Pytania / Pomysły: Discussions
- Bezpieczeństwo: tylko wg https://github.com/Aerendal/doc-flow/security/policy (nie publikuj)

Przy zgłaszaniu błędu proszę podaj:
- wersję / commit
- system i architekturę
- dokładne polecenie
- expected vs actual
- kroki reprodukcji
- logi (zredagowane)
","path":"README_FEEDBACK_SNIPPET.md"},{"content":"# Szablony odpowiedzi (PL, 1. osoba)

## needs-repro
Dziękuję za zgłoszenie. Potrzebuję minimalnego przypadku reprodukcji (dokładne polecenia i dane wejściowe), wersji `doc-flow --version`, systemu operacyjnego i istotnych logów. Tymczasowo oznaczam jako `status/needs-repro`. Jeśli nie dostanę informacji w ciągu 14 dni, zamknę zgłoszenie.

## needs-info
Dziękuję — proszę o dodatkowe informacje:
- wersja programu / commit
- dokładne kroki (kopiuj/wklej polecenia)
- oczekiwane vs rzeczywiste zachowanie
- istotne logi
Oznaczam jako `status/needs-info` — odwołam po uzupełnieniu.

## use-discussions
Dziękuję — to wygląda na pytanie lub propozycję do dyskusji. Proszę kontynuuj w Discussions: https://github.com/Aerendal/doc-flow/discussions. Zamykam to Issue jako "converted to discussion".
","path":"REPLIES.md"},{"content":"# Feedback Week — zasady i co zbieram

Cześć — prowadzę tydzień zbierania feedbacku (7 dni). Poniżej zasady i zakres.

Gdzie zgłaszać:
- Bugs → Issues (użyj formularza "Zgłoszenie błędu")
- Pytania i pomysły → Discussions (użyj tego repozytorium Discussions)
- Security → tylko wg https://github.com/Aerendal/doc-flow/security/policy (nie publikuj publicznie)

In scope:
- Reproducible błędy, problemy instalacyjne, problemy z dokumentacją

Out of scope:
- Duże redesigny bez acceptance criteria

Zasady triage:
- Odpowiadam 2x dziennie
- Statusy: needs-repro, needs-info, accepted, wontfix, duplicate

Po tygodniu:
- Stworzę `docs/FEEDBACK_WEEK_YYYYMMDD.md` z krótkim podsumowaniem i decyzjami.
","path":"DISCUSSION_PINNED.md"},{"content":"#!/usr/bin/env bash
set -euo pipefail

repo="Aerendal/doc-flow"

ensure_label () {
  local name="$1" color="$2" desc="$3"
  if gh label list --repo "$repo" --search "$name" --json name --jq '.[].name' | grep -Fxq "$name"; then
    gh label edit "$name" --repo "$repo" --color "$color" --description "$desc" >/dev/null
  else
    gh label create "$name" --repo "$repo" --color "$color" --description "$desc" >/dev/null
  fi
  echo "OK: $name"
}

# type
ensure_label "type/bug"      "d73a4a" "Bug report (reproducible defect)"
ensure_label "type/feature"  "a2eeef" "Feature request / enhancement"
ensure_label "type/question" "d876e3" "Question / support request"
ensure_label "type/docs"     "0075ca" "Documentation issue/change"
ensure_label "type/security" "b60205" "Security-related (do not disclose)"

# status
ensure_label "status/needs-triage" "fbca04" "Not triaged yet"
ensure_label "status/needs-info"   "fef2c0" "Missing environment/logs/details"
ensure_label "status/needs-repro"  "fef2c0" "Need minimal reproduction steps/data"
ensure_label "status/accepted"     "0e8a16" "Accepted / confirmed"
ensure_label "status/duplicate"    "cfd3d7" "Duplicate of another issue"
ensure_label "status/wontfix"      "ffffff" "Not planned / out of scope"
ensure_label "status/blocked"      "e99695" "Blocked by dependency/decision"
ensure_label "status/in-progress"  "1d76db" "Work started"
ensure_label "status/review"       "1d76db" "PR/review needed"

# prio
ensure_label "prio/p0" "b60205" "Critical: security, data loss, crash, release blocker"
ensure_label "prio/p1" "d93f0b" "High: blocks common workflow / install / CI"
ensure_label "prio/p2" "fbca04" "Normal: important but not blocking"
ensure_label "prio/p3" "cfd3d7" "Low: nice-to-have"

# area
ensure_label "area/cli"      "5319e7" "CLI behavior, flags, UX"
ensure_label "area/install"  "0052cc" "Install / packaging / releases"
ensure_label "area/ci"       "0052cc" "CI workflows / pipelines"
ensure_label "area/indexing" "5319e7" "Indexing logic / ingestion"
ensure_label "area/sqlite"   "5319e7" "SQLite/FTS schema, queries, perf"
ensure_label "area/docs"     "0075ca" "Docs and examples"
ensure_label "area/security" "b60205" "Security policy/process"

# meta
ensure_label "meta/feedback-week"  "bfdadc" "Collected during feedback week"
ensure_label "meta/good-first-issue" "7057ff" "Good first issue"
ensure_label "meta/help-wanted"      "008672" "Help wanted / contributions welcome"


echo "Done creating/updating labels for $repo""