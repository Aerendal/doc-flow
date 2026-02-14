---
title: API Error Codes Reference
status: needs_content
---

# API Error Codes Reference

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zdefiniować i ujednolicić kody błędów API (REST/GraphQL/gRPC) wraz z formatami odpowiedzi i wskazówkami remediacji, aby poprawić DX i diagnostykę.

## Zakres i granice
- Obejmuje: katalog kodów (HTTP i domenowe), struktura odpowiedzi (trace id, timestamp, message, details, remediation), mapowanie do logów/alertów, lokalizacja/wersje językowe komunikatów, kompatybilność wsteczna, bezpieczeństwo (info disclosure), testy kontraktowe.
- Poza zakresem: pełna dokumentacja endpointów (osobne), kody UI (jeśli inne – linkować).

## Wejścia i wyjścia
- Wejścia: obecne kody z usług, standardy HTTP/REST/GraphQL/gRPC, wymagania DX, polityka lokalizacji, polityka logów/trace, SLO, wytyczne bezpieczeństwa.
- Wyjścia: referencja kodów i payloadów, guidance użycia, przykłady, mapa kod→warstwa, checklisty testów kontraktowych, plan rollout (zmiany wersjonowane).

## Powiązania (meta)
- Wymaga odniesienia do: Key Documents
- Wymaga odniesienia do: Key Document Structures
- Wymaga odniesienia do: Document Dependencies
- Wymaga odniesienia do: RACI i role
- Wymaga odniesienia do: Standardy i compliance

## Zależności dokumentu
Wskaż: katalog usług, log format/tracing, zasady lokalizacji, polityka informacji o błędach/security, kontrakty API; brak – odnotuj.

## Powiązania sekcja↔sekcja
Kody → payload → logi/alerty; bezpieczeństwo → maskowanie detali; kompatybilność → rollout.

## Fazy cyklu życia
Curacja kodów → Projekt payloadu → Review bezpieczeństwa/DX → Rollout → Utrzymanie/wersjonowanie.

## Struktura sekcji (szkielet)
- Założenia i cele (DX/diagnostyka/bezpieczeństwo).
- Format odpowiedzi (pola obowiązkowe/opcjonalne, trace id, localization).
- Katalog kodów (HTTP + domenowe) z przykładami i remediacją.
- Mapowanie kodów na logowanie/alerty/metryki.
- Kompatybilność i wersjonowanie (zmiany kodów/pól, deprecacje).
- Testy kontraktowe i walidacja (schema, golden tests, backward compat).
- Zasady bezpieczeństwa (info disclosure, PII w błędach, rate limit/429).
- Plan rollout i komunikacja.

## Wymagane rozwinięcia
- Katalog kodów → tabela per usługa.
- Testy → scenariusze kontraktowe i golden responses.

## Wymagane streszczenia
- Cheat-sheet top kodów i payload.

## Guidance
Cel: spójne i bezpieczne błędy. DoR: zebrane kody, policy security/DX, tracing/log format. DoD: format/payload/katalog/testy/rollout opisane; sekcje N/A uzasadnione; metadane aktualne.

## Szybkie powiązania
- Meta: Key Documents
- Meta: Key Document Structures
- Meta: Document Dependencies

## Jak używać dokumentu
- Ustal format; zmapuj kody; dodaj do kontraktów; przetestuj; rollout; utrzymuj wersje.

## Checklisty jakości (DoR/DoD skrót)
- DoR: [ ] Zebrane kody z usług; [ ] Polityka bezpieczeństwa/logów/tracingu; [ ] DX wymagania.
- DoD: [ ] Format/katalog/przykłady/testy gotowe; [ ] Rollout i wersjonowanie opisane; [ ] Sekcje N/A uzasadnione; metadane aktualne.

## Definicje robocze
- [Termin 1]
- [Termin 2]
- [Termin 3]

## Przykłady użycia
- [Przykład 1]
- [Przykład 2]

## Ryzyka i ograniczenia
- [Ryzyko 1]
- [Ryzyko 2]

## Decyzje i uzasadnienia
- [Decyzja 1]
- [Decyzja 2]

## Założenia
- [Założenie 1]
- [Założenie 2]

## Otwarte pytania
- [Pytanie 1]
- [Pytanie 2]

## Powiązania z innymi dokumentami
- [Dokument A] — [typ relacji] — [uzasadnienie]
- [Dokument B] — [typ relacji] — [uzasadnienie]

## Powiązania z sekcjami innych dokumentów
- [Dokument X → Sekcja Y] — [powód]
- [Dokument Z → Sekcja W] — [powód]

## Wymagane odwołania do standardów
- [Standard 1]
- [Standard 2]

## Mapa relacji sekcja→sekcja
- [Sekcja A] -> [Sekcja B] : [typ]
- [Sekcja C] -> [Sekcja D] : [typ]

## Mapa relacji dokument→dokument
- [Dokument A] -> [Dokument B] : [typ]
- [Dokument C] -> [Dokument D] : [typ]

## Ścieżki informacji
- [Wejście] → [Źródło] → [Rozwinięcie] → [Wyjście]
- [Wejście] → [Źródło] → [Streszczenie] → [Wyjście]

## Weryfikacja spójności
- [ ] Ścieżki informacji zamknięte
- [ ] Brak sprzecznych relacji
- [ ] Sekcje krytyczne mają źródła

## Lista kontrolna spójności relacji
- [ ] Relacje mają sekcje źródłowe
- [ ] Relacje nie są sprzeczne
- [ ] Cross-doc uzasadnione
- [ ] Rozwinięcia/streszczenia odnotowane

## Artefakty powiązane
- [Artefakt 1]
- [Artefakt 2]

## Ścieżka decyzji
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]

## Użytkownicy i interesariusze
- [Rola] — [potrzeby/odpowiedzialności]
- [Rola] — [potrzeby/odpowiedzialności]

## Ścieżka akceptacji
