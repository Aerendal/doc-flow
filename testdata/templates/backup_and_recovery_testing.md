---
title: Backup and Recovery Testing
status: needs_content
---

# Backup and Recovery Testing

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


> Powiązania: linkage_index.jsonl

## Cel dokumentu
Zaplanować i wykonywać testy backupu/odtwarzania, aby potwierdzić możliwość przywrócenia danych w założonych RPO/RTO, wykrywać luki i zapewnić dowody zgodności (audyt/regulator).

## Zakres i granice
- Obejmuje: scenariusze testowe (plik/DB/konfiguracja, partial/full, ransomware), kryteria wejścia/wyjścia, role, metryki i raportowanie, action items po testach.
- Nie obejmuje: pełnych procedur DR (są w DRP) ani runbooków specyficznych dla aplikacji.

## Wejścia i wyjścia
- Wejścia: strategia backupu i replikacji, matryca RPO/RTO, inwentarz danych/backupów, polityka kluczy (KMS/HSM), okna serwisowe, lista właścicieli systemów, wymagania regulatorów/klientów.
- Wyjścia: plan testu (scenariusz, kroki, środowisko), log przebiegu, wyniki (RPO/RTO zmierzone), błędy i incydenty, action items z właścicielami/terminami, raport audytowy.

## Powiązania (meta)
- Backup and Recovery Strategy (cele, częstotliwości, technologia)
- DRP/BCP (spójność z RPO/RTO, priorytety usług)
- Polityka bezpieczeństwa danych (szyfrowanie, klucze, lokalizacja)
- Standardy: ISO 27001 A.12/A.17, ISO 22301, SOC2 CC7, wytyczne branżowe

## Zależności dokumentu
- Dostępne kopie/replicy + metadane retencji
- Polityka kluczy KMS/HSM i procedury odzyskiwania kluczy
- Runbooki usług (do walidacji po restore)
- Okna serwisowe i zgody właścicieli systemów

## Powiązania sekcja↔sekcja
- Scenariusze testowe ↔ RPO/RTO ↔ metryki sukcesu
- Plan komunikacji ↔ IR/BCP (kto informuje o teście i wynikach)
- Action items ↔ rejestr zmian/ryzyk (śledzenie wdrożenia)

## Fazy cyklu życia
- Faza 1: Koncepcja i Wizja: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 2: Analiza Wymagań: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 3: Projekt / Design: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 4: Planowanie: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 5: Implementacja: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 6: Testowanie / QA: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 7: Bezpieczeństwo / Compliance: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 8: Wdrożenie / Deployment: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 9: Operacje / Maintenance: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 10: Incident Management: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 11: Monitoring / Observability: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 12: Dokumentacja referencyjna: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 13: Szkolenie / Onboarding: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 14: Komunikacja stakeholders: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 15: Knowledge Management: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 16: Postmortem / Retrospektywa: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 17: Budżetowanie / Cost Management: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 18: Vendor Management: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 19: Governance / Compliance: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 20: Decommission / Sunset: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 21: DR / BCP: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 22: Change Management: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.
- Faza 23: Capacity Planning: Określ czy w tej fazie dokument powstaje, jest aktualizowany, przeglądany lub archiwizowany; podaj uzasadnienie i odpowiedzialnych.

## Struktura sekcji (szkielet)
1. Zakres testu i cele (dane/usługi, RPO/RTO do weryfikacji).
2. Scenariusze i środowiska (partial/full, ransomware, utrata DC; środowisko testowe/produkcyjne z oknem serwisowym).
3. Zespół i role (koordynator testu, właściciele systemów, DBA/infra, bezpieczeństwo, obserwator audytu).
4. Plan testu (kroki backup/restore, weryfikacja, pomiar czasu, kontrola integralności, dekrypt/klucze).
5. Kryteria wejścia/wyjścia i akceptacji (gotowość danych/środowisk, progi RPO/RTO, bezpieczeństwo danych).
6. Walidacja po odtworzeniu (kontrole techniczne, smoke testy aplikacyjne, akceptacja biznesowa).
7. Raportowanie i action items (wyniki, odchylenia, właściciele, terminy, priorytety).
8. Harmonogram testów (częstotliwość per system/klasa danych, wymagania regulacyjne).
9. Lessons learned i aktualizacja strategii/DRP/runbooków.

## Wymagane rozwinięcia
- Macierz RPO/RTO i sposób pomiaru w testach.
- Lista kluczy i procedur KMS/HSM użytych przy restore.
- Plan komunikacji testowej (odbiorcy, kanały, szablony).
- Rejestr action items z priorytetem, właścicielem, terminem.

## Wymagane streszczenia
- Streszczenie wyników testów (RPO/RTO osiągnięte, główne incydenty).
- Streszczenie lessons learned i zmian do wdrożenia (dla zarządu/audytu).

## Guidance
Cel: skrócone wskazówki do wypełniania szablonów dokumentów (core/satellite).

- Cel dokumentu: 2–3 zdania o decyzjach, ryzykach i wartości dokumentu.
- Zakres i granice: co obejmuje (systemy/procesy/zespoły) i czego nie obejmuje; zaznacz granice odpowiedzialności.
- Wejścia: dane, wymagania, standardy, zależności potrzebne przed startem.
- Wyjścia: artefakty/rezultaty, kto je konsumuje, format (link/plik).
- Zależności dokumentu: wymagane dokumenty lub decyzje; właściciel; wpływ na kolejność prac.
- Powiązania sekcja↔sekcja: które sekcje się rozwijają/streszczają; podaj uzasadnienie.
- Struktura sekcji: utrzymuj układ logiczny; sekcje bez treści oznacz jako N/A z krótkim uzasadnieniem.
- Fazy cyklu życia: zaznacz, w których fazach dokument powstaje/aktualizuje się/archiwizuje; kto odpowiada.
- DoR (Definition of Ready): zakres, wejścia, role, zależności, kryteria akceptacji gotowe.
- DoD (Definition of Done): sekcje uzupełnione lub N/A, powiązania wpisane, checklisty jakości sprawdzone, wersja/data/właściciel, linki/artefakty działają.
- Język: polski; nazwy własne pozostają bez zmian; liczby w nazwach plików usunięte już w szablonach.
- Filozofia: optymalizuj przez rozwój, nie ucinanie — dodawaj, nie kasuj; elementy „satelitarne” zostają.


## Szybkie powiązania
- Zależność: {'depends_on': 'Database Schema Design', 'type': 'REFERENCES'}
- Zależność: {'depends_on': 'Schema Implementation', 'type': 'REFERENCES'}
- Zależność: {'depends_on': 'Database Schema Reference', 'type': 'REFERENCES'}

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] Czy cel dokumentu jest jednoznaczny?
- [ ] Czy zakres i granice są jasno określone?
- [ ] Czy wszystkie zależności są opisane?
- [ ] Czy wskazano wymagane rozwinięcia i streszczenia?
- [ ] Czy powiązania sekcja↔sekcja są spójne?

## Definicje robocze
- placeholder-term-1 — [definicja robocza]
- placeholder-term-2 — [definicja robocza]
- placeholder-term-3 — [definicja robocza]

## Przykłady użycia
- [Przykład 1 — krótki opis sytuacji i zastosowania]
- [Przykład 2 — krótki opis sytuacji i zastosowania]

## Ryzyka i ograniczenia
- [Ryzyko 1 — wpływ i sposób ograniczenia]
- [Ryzyko 2 — wpływ i sposób ograniczenia]

## Decyzje i uzasadnienia
- [Decyzja 1 — uzasadnienie]
- [Decyzja 2 — uzasadnienie]

## Założenia
- [Założenie 1]
- [Założenie 2]

## Otwarte pytania
- [Pytanie 1]
- [Pytanie 2]

## Powiązania z innymi dokumentami
- [odniesienie do dokumentu] — [typ relacji] — [uzasadnienie]
- [odniesienie do dokumentu] — [typ relacji] — [uzasadnienie]

## Powiązania z sekcjami innych dokumentów
- [odniesienie do dokumentu] — [powód powiązania]
- [odniesienie do dokumentu] — [powód powiązania]

## Słownik pojęć w dokumencie
- [Pojęcie 1] — [definicja i źródło]
- [Pojęcie 2] — [definicja i źródło]
- [Pojęcie 3] — [definicja i źródło]

## Wymagane odwołania do standardów
- [Standard 1] — - [miejsce na wpisanie konkretnej relacji/elementu]
- [Standard 2] — - [miejsce na wpisanie konkretnej relacji/elementu]

## Mapa relacji sekcja→sekcja
- - [miejsce na wpisanie konkretnej relacji/elementu] -> - [miejsce na wpisanie konkretnej relacji/elementu] : [typ relacji]
- - [miejsce na wpisanie konkretnej relacji/elementu] -> - [miejsce na wpisanie konkretnej relacji/elementu] : [typ relacji]

## Mapa relacji dokument→dokument
- [odniesienie do dokumentu] -> [odniesienie do dokumentu] : [typ relacji]
- [odniesienie do dokumentu] -> [odniesienie do dokumentu] : [typ relacji]

## Ścieżki informacji
- - Ścieżka informacji: zdefiniuj wejście → sekcja źródłowa → sekcja docelowa → wyjście (opisz konkretnie).
- - Ścieżka informacji: zdefiniuj wejście → sekcja źródłowa → sekcja docelowa → wyjście (opisz konkretnie).

## Weryfikacja spójności
- [ ] Czy wszystkie ścieżki informacji są zamknięte?
- [ ] Czy istnieją pętle lub sprzeczne relacje?
- [ ] Czy sekcje krytyczne mają wskazane źródła i rozwinięcia?

## Lista kontrolna spójności relacji
- [ ] Czy każda sekcja z relacją ma wskazaną sekcję źródłową?
- [ ] Czy relacje nie tworzą sprzecznych wymagań (np. wzajemne wykluczanie)?
- [ ] Czy relacje cross‑doc mają uzasadnienie i są zgodne z fazą?
- [ ] Czy relacje wymagają rozwinięć lub streszczeń są odnotowane?

## Artefakty powiązane
- [Artefakt 1] — [opis i relacja do dokumentu]
- [Artefakt 2] — [opis i relacja do dokumentu]

## Ścieżka decyzji
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]
- [Decyzja] → [Uzasadnienie] → [Konsekwencje]

## Użytkownicy i interesariusze
- [Rola / interesariusz] — [potrzeby i odpowiedzialności]
- [Rola / interesariusz] — [potrzeby i odpowiedzialności]

## Ścieżka akceptacji
- [Kto zatwierdza] → [kryteria akceptacji] → [status]
- [Kto zatwierdza] → [kryteria akceptacji] → [status]

## Kryteria ukończenia
- [ ] Kryterium 1 — [opis]
- [ ] Kryterium 2 — [opis]

## Metryki jakości
- [Metryka 1] — [cel / próg]
- [Metryka 2] — [cel / próg]

## Monitoring i utrzymanie
- [Co monitorujemy] — [narzędzie / częstotliwość]
- [Kto utrzymuje] — [rola]

## Kontrola zmian
- [Zmiana] — [powód] — [data] — [akceptacja]

## Wymogi prawne i regulacyjne
- [Wymóg 1] — [źródło / akt prawny / standard]
- [Wymóg 2] — [źródło / akt prawny / standard]

## Zasady bezpieczeństwa informacji
- [Zasada 1] — [opis i wpływ na dokument]
- [Zasada 2] — [opis i wpływ na dokument]

## Ochrona danych i prywatność
- [Wymaganie 1] — [opis i sekcja docelowa]
- [Wymaganie 2] — [opis i sekcja docelowa]

## Wersjonowanie treści
- [Wersja] — [zmiana] — [autor] — [data]
- [Wersja] — [zmiana] — [autor] — [data]

## Historia zmian sekcji
- - [miejsce na wpisanie konkretnej relacji/elementu] — [zmiana] — [data]
- - [miejsce na wpisanie konkretnej relacji/elementu] — [zmiana] — [data]

## Wymagane aktualizacje
- - [miejsce na wpisanie konkretnej relacji/elementu] — [powód aktualizacji] — - [miejsce na wpisanie konkretnej relacji/elementu]
- - [miejsce na wpisanie konkretnej relacji/elementu] — [powód aktualizacji] — - [miejsce na wpisanie konkretnej relacji/elementu]

## Integracje i interfejsy
- [System / API] — [zakres integracji] — [wymagania]
- [System / API] — [zakres integracji] — [wymagania]

## Wymagania danych
- [Dane wejściowe] — [format] — [walidacja]
- [Dane wyjściowe] — [format] — [walidacja]

## Logowanie i audyt
- [Zdarzenie] — [poziom] — [retencja]
- [Zdarzenie] — [poziom] — [retencja]

## Utrzymanie i operacje
- [Procedura] — [cel] — [częstotliwość]
- [Procedura] — [cel] — [częstotliwość]

## KPI i SLA
- [KPI] — [cel] — [pomiar]
- [SLA] — [cel] — [pomiar]

## Scenariusze awaryjne
- [Scenariusz] — [objawy] — [reakcja]
- [Scenariusz] — [objawy] — [reakcja]

## Wpływ na inne systemy
- [System] — [rodzaj wpływu] — [ryzyko]
- [System] — [rodzaj wpływu] — [ryzyko]

## Zależności danych między systemami
- [Źródło danych] → [Odbiorca] — [opis]
- [Źródło danych] → [Odbiorca] — [opis]

## Harmonogram przeglądów
- [Obszar] — [częstotliwość] — [właściciel]
- [Obszar] — [częstotliwość] — [właściciel]

## Wymagania wydajnościowe
- [Wymaganie] — [metryka] — [próg]
- [Wymaganie] — [metryka] — [próg]

## Wymagania dostępnościowe
- [Wymaganie] — [SLA] — [metoda pomiaru]
- [Wymaganie] — [SLA] — [metoda pomiaru]

## Wymagania skalowalności
- [Wymaganie] — [cel] — [warunki]
- [Wymaganie] — [cel] — [warunki]

## Wymagania dostępności danych
- [Dane] — [częstotliwość dostępu] — [SLA]
- [Dane] — [częstotliwość dostępu] — [SLA]

## Retencja i archiwizacja
- [Dane] — [retencja] — [archiwizacja]
- [Dane] — [retencja] — [archiwizacja]

## Dostępność w sytuacjach awaryjnych
- [Scenariusz] — [zachowanie] — [priorytet]
- [Scenariusz] — [zachowanie] — [priorytet]

## Testy i weryfikacja
- [Test] — [cel] — [wynik oczekiwany]
- [Test] — [cel] — [wynik oczekiwany]

## Walidacja zgodności
- [Wymóg] — [metoda weryfikacji]
- [Wymóg] — [metoda weryfikacji]

## Audyty i przeglądy
- [Audyty] — [częstotliwość] — [odpowiedzialny]
- [Audyty] — [częstotliwość] — [odpowiedzialny]
