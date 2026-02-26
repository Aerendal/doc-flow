# Feedback Week — zasady i co zbieram

Cześć — prowadzę tydzień zbierania feedbacku (7 dni). Poniżej zasady i zakres.

Czas trwania: 7 dni (od <start> do <end>) — tu wklej dokładne daty.

Gdzie zgłaszać:
- Bugs → Issues (użyj formularza "Zgłoszenie błędu")
- Pytania i pomysły → Discussions (ten wątek i osobne wątki)
- Security → tylko wg [SECURITY.md](/SECURITY.md). NIE publikuj podatności publicznie.

In scope:
- Błędy, które mogę odtworzyć lokalnie lub na CI
- Propozycje dotyczące UX/CLI i integracji
- Problemy z dokumentacją związanymi z instalacją i użyciem

Out of scope:
- Propozycje dużych funkcjonalności wymagających projektu architektury (zamieszczaj, ale będą rozważone później)
- Prywatne dane i raporty bezpieczeństwa — użyj SECURITY.md

Zasady triage:
- Odpowiadam 2x dziennie przez 7 dni.
- Statusy: needs-repro, needs-info, accepted, wontfix, duplicate.
- Priorytety: P0 (crash/data loss/security), P1 (blokuje instalację/CI), P2 (pozostałe).

Co robię po tygodniu:
- Tworzę podsumowanie w docs/FEEDBACK_WEEK_YYYYMMDD.md z decyzjami i planem dalszych prac.

Dziękuję za pomoc — proszę o jasne kroki reprodukcji i logi tam, gdzie to możliwe.
