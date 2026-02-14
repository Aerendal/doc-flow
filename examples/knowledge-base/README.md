# Example: Knowledge Base

Use-case: mieszane treści (how-to, FAQ) z aliasami sekcji.

Steps:
1. `./docflow scan`
2. `./docflow validate --status-aware --governance docs/_meta/GOVERNANCE_RULES.yaml`
3. `./docflow migrate-sections` (zobacz aliasy Overview→Przegląd)

Files:
- `getting_started.md` — guide
- `faq.md` — FAQ (uses alias Overview)
