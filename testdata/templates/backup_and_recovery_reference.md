---
title: Backup and Recovery Reference
status: needs_content
---

# Backup and Recovery Reference

## Metadane
- Właściciel: [Ops/SRE/DBA]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Zwięzła karta referencyjna backup/restore dla usług krytycznych: RPO/RTO, harmonogramy, lokalizacje/retencja, szyfrowanie/klucze, kontakty on‑call oraz linki do runbooków i ostatnich testów. Ma umożliwiać szybkie decyzje i audyt.

## Zakres i granice
- Obejmuje: parametry i lokalizacje backupów, retencję, szyfrowanie/klucze (KMS/HSM), kontakty i eskalacje, odnośniki do runbooków/testów DR, wyjątki/regulacje.
- Poza zakresem: projekt/strategia backupu, szczegółowe procedury restore (są w Guide/Procedure), wybór technologii.

## Wejścia i wyjścia
- Wejścia: matryca RPO/RTO, harmonogramy backupów, inwentarz lokacji (region/on‑prem/offsite), polityka retencji, polityka kluczy, wyniki testów restore/DR, lista kontaktów on‑call.
- Wyjścia: karta referencyjna (PDF/Confluence/MD) z parametrami aktualnymi, linkami i dowodami testów; sygnał dla audytu i operacji.

## Powiązania (meta)
- Key Documents: backup_and_recovery_strategy/design/guide/procedure; drp/bcp; security_key_management.
- Key Document Structures: parametry, lokalizacje, retencja, klucze, kontakty, dowody.
- Document Dependencies: CMDB/asset, KMS/HSM, monitoring backupów, runbooki usług.

## Zależności dokumentu
Wymaga aktualnych RPO/RTO per usługa, harmonogramów/retencji z systemu backup, polityki kluczy (custody/escrow), listy kontaktów on‑call oraz ostatnich wyników testów restore/DR. Bez tych danych DoR jest otwarte.

## Powiązania sekcja↔sekcja
- RPO/RTO ↔ Harmonogramy/retencja ↔ Lokalizacje: muszą być spójne z klasą krytyczności.
- Szyfrowanie/klucze ↔ Procedury restore: dostęp do kluczy wymagany do odtworzenia.
- Kontakty/eskalacje ↔ Runbooki: potrzebne ścieżki eskalacji i właściciele.

## Fazy cyklu życia
- Planowanie: uzupełnienie parametrów, zgodność z BCP/DRP.
- Operacje: utrzymanie aktualności po zmianach usług, rotacje kluczy, zmiany retencji.
- Testy DR: aktualizacja po każdym teście/restore.
- Audyt/Compliance: potwierdzenie zgodności retencji/lokalizacji/kluczy.
- Decommission: archiwizacja kart i kluczy, potwierdzenie usunięcia danych.

## Struktura sekcji
1) Matryca RPO/RTO i właściciele usług  
2) Harmonogramy/typy backupów (full/inc/diff/snapshot/log) + retencja  
3) Lokalizacje danych (region/on‑prem/offsite/tape) i wymagania regulacyjne  
4) Szyfrowanie i klucze (KMS/HSM, rotacja, escrow, procedura odzyskania)  
5) Dowody testów restore/DR (data, zakres, wynik, follow‑up)  
6) Kontakty i eskalacje (on‑call, DBA/infra/security, komunikacja kryzysowa)  
7) Linki: runbooki, procedury, raporty audytu  
8) Wyjątki/waivery i daty przeglądów  

## Wymagane rozwinięcia
- Matryca RPO/RTO z klasyfikacją usług i właścicielami.
- Pełne harmonogramy/retencja z narzędzia backup (per dataset).
- Dane KMS/HSM: właściciel kluczy, rotacja, procedura break‑glass.

## Wymagane streszczenia
- Top‑k krytyczne usługi: RPO/RTO, lokalizacje, retencja, klucze.
- Wyjątki/regulacje specyficzne (np. lokalizacja danych, minimalna retencja).

## Guidance (skrót)
- Karta ma być krótka, aktualna, jednolita; pełne instrukcje są w Guide/Procedure.
- Utrzymuj spójność RPO/RTO z DRP i klasą krytyczności usług.
- Linkuj do najnowszych dowodów testów restore/DR; zaznacz datę i wynik.
- Dodawaj właścicieli i kontakty on‑call dla każdej usługi i klucza.
- Oznacz wyjątki/waivery z datą wygaśnięcia i planem domknięcia.

## Szybkie powiązania
- linkage_index.jsonl (backup/dr/reference)
- backup_and_recovery_strategy/design/guide/procedure, drp/bcp, kms/hsm policy

## Jak używać dokumentu
1. Wpisz aktualne RPO/RTO, harmonogramy, retencję i lokalizacje.
2. Dodaj dane kluczy (KMS/HSM), właścicieli i procedury break‑glass.
3. Podlinkuj runbooki i dowody testów restore/DR.
4. Zaktualizuj kontakty i eskalacje; domknij checklisty DoR/DoD.

## Checklisty Definition of Ready (DoR)
- [ ] RPO/RTO i krytyczność usług zebrane.
- [ ] Harmonogramy/retencja dostępne; lokalizacje znane.
- [ ] Informacje o kluczach (KMS/HSM) i właścicielach dostępne.
- [ ] Ostatnie testy restore/DR zarejestrowane lub N/A z planem.
- [ ] Struktura sekcji wypełniona/N/A.

## Checklisty Definition of Done (DoD)
- [ ] Wszystkie sekcje wypełnione lub N/A z uzasadnieniem.
- [ ] Linki do runbooków i dowodów testów działają.
- [ ] Kontakty/eskalacje aktualne; wyjątki opisane z datą.
- [ ] Spójność RPO/RTO ↔ retencja ↔ lokalizacje potwierdzona.
- [ ] Wersja/data/właściciel zaktualizowane; dokument wpisany w linkage_index.

## Definicje robocze
- RPO (Recovery Point Objective) — maks. akceptowalna utrata danych.
- RTO (Recovery Time Objective) — maks. czas przywrócenia usługi.
- Break‑glass — procedura awaryjnego dostępu do kluczy.

## Przykłady użycia
- Szybkie sprawdzenie, które regiony mają offsite i jaką retencję.
- Audyt: dowody testów restore/DR i właściciele kluczy.

## Ryzyka i ograniczenia
- Nieaktualne RPO/RTO lub harmonogramy → fałszywe poczucie bezpieczeństwa.
- Brak dostępu do kluczy przy restore → wydłużone RTO lub utrata danych.

## Decyzje i uzasadnienia
- [Decyzja] Priorytety offsite per usługa — wg krytyczności i ryzyka.
- [Decyzja] Retencja per dataset — uzasadnienie regulacyjne/biznesowe/kosztowe.

## Założenia
- System backup raportuje sukces/fail i ma API/raporty do pobrania.
- KMS/HSM dostępne i posiada procedurę odzyskania.

## Otwarte pytania
- Czy wszystkie wyjątki lokalizacji danych mają zgodę compliance?
- Jak często rotujemy klucze i testujemy restore z nowym kluczem?

## Powiązania z innymi dokumentami
- Backup & Recovery Guide/Procedure, DR Plan, BCP, Security Key Management.

## Powiązania z sekcjami innych dokumentów
- DRP → RPO/RTO; Key Management → Szyfrowanie/klucze.

## Słownik pojęć w dokumencie
- Snapshot, Incremental, Differential, Offsite, WORM — dodaj definicje.

## Wymagane odwołania do standardów
- Wymagania regulatorów dot. retencji/lokalizacji, polityka szyfrowania.

## Mapa relacji sekcja→sekcja
- RPO/RTO → Harmonogramy/retencja → Lokalizacje → Klucze → Testy → Kontakty.

## Mapa relacji dokument→dokument
- Backup Reference → Guide/Procedure → DRP/BCP → Audit/Compliance.

## Ścieżki informacji
- Strategia/Design → Parametry → Karta referencyjna → Runbooki/DR → Audyt.

## Weryfikacja spójności
- [ ] Parametry spójne z DRP/BCP i realnymi konfiguracjami.
- [ ] Linki/dowody działają; daty są aktualne.
- [ ] Wyjątki mają datę wygaśnięcia i właściciela.

## Lista kontrolna spójności relacji
- [ ] Każdy dataset/usługa ma RPO/RTO, retencję, lokalizację, klucze, kontakty.
- [ ] Każdy test restore/DR jest powiązany z usługą i ma wynik.
- [ ] Relacje cross‑doc opisane z uzasadnieniem.

## Artefakty powiązane
- Raporty z systemu backup, testy restore/DR, linki do runbooków, lista kontaktów.

## Ścieżka decyzji
- [Decyzja] → [Uzasadnienie] → [Konsekwencje] → [Właściciel] → [Data].

## Użytkownicy i interesariusze
- SRE/Ops/DBA, Security/Compliance, Audyt, Właściciele usług.

## Ścieżka akceptacji
- Ops/DBA → Security/Compliance → Service Owners → Owner sign‑off.

## Kryteria ukończenia
- [ ] Karta kompletna i zgodna z DRP/BCP; linki działają.
- [ ] Dowody testów i klucze opisane; kontakty aktualne.
- [ ] Dokument powiązany w linkage_index i checklistach.

## Metryki jakości
- Aktualność danych (dni od ostatniej aktualizacji), odsetek usług z kompletną kartą, sukces testów restore/DR, czas dostępu do kluczy, liczba wyjątków otwartych.
