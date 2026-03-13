# Wsparcie

To repozytorium używa dwóch kanałów wsparcia:

1) **GitHub Discussions** (preferowane dla pytań, rozwiązywania problemów, pomysłów):
   https://github.com/Aerendal/doc-flow/discussions

2) **GitHub Issues** (tylko potwierdzone błędy i zadania z jasnym scope):
   https://github.com/Aerendal/doc-flow/issues

## Gdzie pisać?

**Discussions** — gdy:
- Masz pytanie o użycie, instalację, konfigurację
- Zastanawiasz się „jak to zrobić...?" i nie wiesz czy to błąd
- Masz pomysł na nową funkcję lub chcesz omówić design

**Issues** — gdy:
- Masz powtarzalny błąd (crash, niepoprawny wynik, zepsuty workflow)
- Zadanie ma jasne kryteria akceptacji

Zgłaszając błąd, dołącz:
- wersję / commit (`./build/docflow --version`)
- system operacyjny + architektura
- dokładne wywołanie(a) komendy
- oczekiwane vs. rzeczywiste zachowanie
- minimalne kroki do odtworzenia
- logi (bez sekretów — tokeny, klucze API, prywatne URL)

## Bezpieczeństwo

**Nie** zgłaszaj luk bezpieczeństwa przez Issues ani Discussions.
Użyj prywatnego kanału opisanego w SECURITY.md:
https://github.com/Aerendal/doc-flow/security/policy

## Gdy wklejasz logi

Usuń sekrety (tokeny, klucze API, prywatne URL, wrażliwe ścieżki) przed wklejeniem.
Jeśli przypadkowo wkleiłeś sekret — edytuj/usuń komentarz i natychmiast zrotuj sekret.

## Oczekiwania odpowiedzi

Projekt jest utrzymywany w miarę możliwości, bez gwarantowanych czasów odpowiedzi.
- Mogę poprosić o dodatkowe informacje lub minimalne odtworzenie.
- Zgłoszenia mogą być oznaczone jako `needs-info` / `needs-repro` / `duplicate` / `wontfix`.
- Issues i pull requesty mogą pozostawać bez odpowiedzi przez dłuższy czas.

