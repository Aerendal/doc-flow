---
title: Test Data Preparation
status: needs_content
---

# Test Data Preparation

## Metadane
- Właściciel: [QA/Data Engineering]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Opisuje zasady przygotowania danych testowych: pozyskanie, maskowanie/syntetyka, wersjonowanie, odświeżanie, zgodność z privacy oraz dostępność dla zespołów. Ma zapewnić powtarzalne, legalne i wysokiej jakości dane do testów.

## Zakres i granice
- Obejmuje: źródła danych (prod/anonymized/synthetic), maskowanie/pseudonimizację, generatory danych syntetycznych, zestawy scenariuszowe, dane do testów negatywnych i edge, dane do wydajności/load, polityki privacy/RODO/PII/PCI/HIPAA, kontrolę dostępu, wersjonowanie/odświeżanie, zarządzanie konfliktami i data drift w testach, dokumentację i artefakty.
- Poza zakresem: pełny data governance (odniesienie), pipeline ETL produkcyjny (poza reuse).

## Wejścia i wyjścia
- Wejścia: wymagania testów (funkcjonalne, perf, bezpieczeństwo), klasy danych/PII, polityki privacy/compliance, źródła (prod snapshots, data warehouse, syntetyczne generatory), schematy, kryteria pokrycia, narzędzia maskowania/syntetyki, harmonogram testów.
- Wyjścia: zestawy danych testowych (maskowane/syntetyczne), instrukcje pozyskania i odświeżenia, katalog zestawów i wersji, kontrola dostępu, logi zgodności, checklisty jakości, linki do repo/artefaktów.

## Powiązania (meta)
- Key Documents: test_data_management_policy, privacy_policy, data_masking_standards, synthetic_data_guidelines, performance_testing_plan, security_testing_plan.
- Key Document Structures: źródła, maskowanie/syntetyka, wersje, dostęp, dokumentacja.
- Document Dependencies: DB/storage, masking tools, synthetic generators, IAM, audit/logging, CI/CD.

## Zależności dokumentu
Wymaga polityk privacy/PII, narzędzi maskowania/syntetyki, schematów danych, listy testów i kryteriów pokrycia, repo/storage na dane, kontroli dostępu i logów. Brak = DoR otwarte.

## Powiązania sekcja↔sekcja
- Wymagania testów → Zestawy danych → Maskowanie/syntetyka → Wersje/odświeżanie.
- Privacy/PII → Maskowanie → Dostęp/audyt.
- Performance/load → Oddzielne zestawy i hygiene (cleanup).

## Fazy cyklu życia
- Planowanie: wymagania, źródła, privacy, narzędzia, kryteria pokrycia.
- Przygotowanie: maskowanie/syntetyka, walidacja jakości, wersjonowanie, dostęp.
- Utrzymanie: odświeżanie, drift check, cleanup, audyt, aktualizacja katalogu.

## Struktura sekcji
1) Wymagania i zakres (typy testów, PII, compliance)  
2) Źródła danych i wybór podejścia (maskowane vs. syntetyczne vs. mieszane)  
3) Maskowanie/pseudonimizacja (metody, narzędzia, kontrola jakości, logi)  
4) Dane syntetyczne (generatory, parametry, seed, walidacja pokrycia)  
5) Zestawy scenariuszowe (happy/edge/negative, dane do regresji)  
6) Dane performance/load (skala, profile, cleanup)  
7) Wersjonowanie i odświeżanie (kalendarz, diff, checksumy)  
8) Dostęp i bezpieczeństwo (IAM, least privilege, audyt)  
9) Dokumentacja i katalog (metadata, lineage, linki)  
10) Ryzyka, decyzje, open issues\n\n## Wymagane rozwinięcia\n- Polityki maskowania (kolumny PII/PCI/PHI, techniki, nieodwracalność) i walidacje.\n- Spec danych syntetycznych (parametry/seed) i kryteria pokrycia scenariuszy.\n- Kalendarz odświeżania, wersjonowanie i logi audytu/dostępu.\n\n## Wymagane streszczenia\n- Podsumowanie źródeł i podejścia (maskowane/syntetyczne), kluczowe ryzyka privacy.\n- Najważniejsze zestawy i ich wersje, dostęp/owner, harmonogram odświeżeń.\n\n## Guidance (skrót)\n- Preferuj dane syntetyczne, gdy to możliwe; gdy korzystasz z prod snapshotów — obowiązkowe maskowanie/PII removal.\n- Ustal pokrycie scenariuszy; dodaj negative/edge cases i dane do perf.\n- Utrzymuj wersje i checksumy; automatyzuj odświeżanie w CI/CD.\n- Loguj dostęp i operacje; zapewnij audyt zgodności.\n- Czyść środowiska po testach (szczególnie load/perf) i minimalizuj retencję.\n\n## Szybkie powiązania\n- linkage_index.jsonl (qa/test_data)\n- test_data_management_policy, privacy_policy, data_masking_standards, synthetic_data_guidelines, performance_testing_plan, security_testing_plan\n\n## Jak używać dokumentu\n1. Zbierz wymagania testów i klasy danych; wybierz podejście (maskowane/syntetyczne).\n2. Przygotuj zestawy, zamaskuj/ wygeneruj dane, waliduj i wersjonuj.\n3. Ustaw dostępy, logi i kalendarz odświeżeń; zarejestruj w katalogu.\n4. Aktualizuj po zmianach schematów/testów; zamknij DoR/DoD.\n\n## Checklisty Definition of Ready (DoR)\n- [ ] Polityki privacy/PII i maskowania dostępne; klasy danych zidentyfikowane.\n- [ ] Schematy i wymagania testów zebrane; narzędzia maskowania/syntetyki gotowe.\n- [ ] Storage/repo i IAM ustawione; struktura sekcji wypełniona/N/A.\n\n## Checklisty Definition of Done (DoD)\n- [ ] Zestawy danych utworzone, zamaskowane/wygen., zweryfikowane; pokrycie scenariuszy spełnione.\n- [ ] Wersje/checksumy, kalendarz odświeżeń i logi audytu zapisane.\n- [ ] Dostępy/IAM i privacy zgodne; dokument w linkage_index.\n- [ ] Wersja/data/właściciel zaktualizowane.\n\n## Definicje robocze\n- Maskowanie, Pseudonimizacja, Synthetic data, Seed, Checksum, Data drift (test data).\n\n## Przykłady użycia\n- Funkcjonalne: maskowany snapshot + syntetyczne edge case’y; refresh co sprint.\n- Performance: syntetyczne dane skalowane z profilem ruchu; cleanup po testach.\n\n## Ryzyka i ograniczenia\n- Wycieki PII przy słabym maskowaniu; nieodpowiednie pokrycie scenariuszy; zbyt stare dane → błędne wyniki.\n\n## Decyzje i uzasadnienia\n- [Decyzja] Podejście maskowane vs. syntetyczne vs. mieszane — uzasadnienie ryzyk/zasobów.\n- [Decyzja] Kalendarz odświeżeń i retencja — uzasadnienie privacy/kosztów.\n\n## Założenia\n- Narzędzia maskowania/syntetyki dostępne; polityki privacy/PII obowiązują.\n\n## Otwarte pytania\n- Czy wymagane są dane produkcyjne w perf/load? Jeśli tak, jakie gwarancje maskowania?\n- Jakie SLA odświeżania dla krytycznych zestawów?\n\n## Powiązania z innymi dokumentami\n- Test Data Mgmt Policy, Privacy Policy, Data Masking Standards, Synthetic Data Guidelines, Performance/Security Testing Plans.\n\n## Powiązania z sekcjami innych dokumentów\n- Privacy → maskowanie/retencja; Perf → dane load; Security → dane do testów bezpieczeństwa.\n\n## Słownik pojęć w dokumencie\n- Maskowanie, Pseudonimizacja, Synthetic data, Seed, Checksum, Data drift.\n\n## Wymagane odwołania do standardów\n- RODO/PII/PHI, PCI/HIPAA jeśli dotyczy, polityki privacy i data retention.\n\n## Mapa relacji sekcja→sekcja\n- Wymagania → Dane → Maskowanie/Syntetyka → Wersje → Dostęp → Odświeżanie.\n\n## Mapa relacji dokument→dokument\n- Test Data Preparation → Privacy/Masking/Synthetic → Perf/Security Testing → CI/CD.\n\n## Ścieżki informacji\n- Wymagania → Źródła → Transformacje → Walidacja → Wersjonowanie → Użycie → Odświeżenie.\n\n## Weryfikacja spójności\n- [ ] Dane zgodne z privacy; maskowanie/syntetyka udokumentowane i przetestowane.\n- [ ] Pokrycie scenariuszy kompletne; wersje/odświeżenia kontrolowane.\n- [ ] Relacje cross‑doc opisane; dokument w linkage_index.\n\n## Lista kontrolna spójności relacji\n- [ ] Każdy zestaw ma źródło, maskowanie/syntetykę, wersję, właściciela, harmonogram.\n- [ ] Każdy dostęp/log zgodny z IAM/privacy; każda zmiana ma wersję i checksum.\n- [ ] Relacje cross‑doc opisane z uzasadnieniem.\n\n## Artefakty powiązane\n- Masking config, synthetic generators, checklisty walidacji, katalog zestawów, logi audytu, checksumy/wersje, harmonogramy refresh.\n\n## Ścieżka decyzji\n- [Decyzja] → [Uzasadnienie] → [Konsekwencje] → [Właściciel] → [Data].\n\n## Użytkownicy i interesariusze\n- QA, Data Engineering, Security/Privacy, Product/Dev, Compliance.\n\n## Ścieżka akceptacji\n- QA/Data → Security/Privacy → Product/Compliance → Owner sign‑off.\n\n## Kryteria ukończenia\n- [ ] Zestawy danych gotowe i zgodne; wersje/odświeżenia/retencja opisane.\n- [ ] Dokument w linkage_index/checklistach; wersja/data/właściciel aktualne.\n\n## Metryki jakości\n- % pokrycia scenariuszy, defekt rate związany z danymi, czas przygotowania danych, liczba incydentów privacy, świeżość danych, sukces refresh w CI/CD.\n*** End Patch
