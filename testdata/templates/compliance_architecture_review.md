---
title: Compliance Architecture Review
status: needs_content
---

# Compliance Architecture Review

## Metadane
- Właściciel: [CISO / Enterprise Architect / Compliance Lead]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Zapewnić, że architektura systemu spełnia wymagania regulacyjne (RODO/PII, PCI DSS, HIPAA/GxP, lokalne przepisy) oraz standardy bezpieczeństwa i audytu. Dokument identyfikuje luki, proponuje środki zaradcze i definiuje ścieżkę zgodności dla nowych i istniejących rozwiązań.

## Zakres i granice
- Obejmuje: architekturę aplikacji i danych, IAM/SSO/MFA, szyfrowanie w spoczynku i tranzycie, logging/audyt, segregację obowiązków, retencję/usuwanie danych, ciągłość działania, zarządzanie zmianą, ścieżki audytu.  
- Poza zakresem: szczegółowe procedury operacyjne SOC i IR (oddzielne runbooki), polityki HR, negocjacje umów z dostawcami.

## Wejścia i wyjścia
- Wejścia: katalog systemów i przepływów danych, klasyfikacja danych, mapa podmiotów przetwarzających, polityki bezpieczeństwa, wymagania branżowe (PCI/HIPAA/GxP), wyniki pen-testów, architektura referencyjna, matryca ról.  
- Wyjścia: raport luk zgodności, lista działań naprawczych z priorytetami, decyzje architektoniczne (ADR) dla kontroli, zaktualizowana mapa przepływów danych, checklisty DoR/DoD zgodności, linki do dowodów audytowych.

## Powiązania (meta)
- Key Documents: data_protection_compliance, security_controls_reference, logging_and_audit_trail, iam_strategy_document, retention_policy, business_continuity_plan.  
- Key Document Structures: przepływy danych, kontrola dostępu, szyfrowanie, monitoring/audyt, retencja, ciągłość, dostawcy.  
- Document Dependencies: SIEM/SOAR, secrets manager, DLP, backup/DR, CMDB, change management.  
- Standardy: ISO 27001/27701, SOC2, PCI DSS, HIPAA/GxP, lokalne akty prawne.

## Zależności dokumentu
Wymaga aktualnej mapy przepływów danych, klasyfikacji danych, inwentarza systemów i dostawców, matrycy ról/SoD, wyników testów bezpieczeństwa oraz polityk retencji. Brak któregokolwiek = blokery DoR.

## Powiązania sekcja↔sekcja
- Dane i klasyfikacja ↔ Szyfrowanie ↔ Retencja/usuwanie.  
- IAM/SoD ↔ Dostawcy ↔ Audyt/monitoring.  
- Ciągłość/DR ↔ Backup ↔ Recovery objectives (RPO/RTO).

## Fazy cyklu życia
- Scoping i identyfikacja regulacji.  
- Analiza architektury i przepływów danych.  
- Ocena kontroli i luk; plan działań.  
- Walidacja wdrożenia kontroli; dowody audytu.  
- Cykl przeglądów okresowych.

## Struktura sekcji
1) Kontekst systemu i regulacji  
2) Klasyfikacja danych i przepływy  
3) Kontrola dostępu (IAM/SoD/SSO/MFA)  
4) Szyfrowanie i ochrona danych  
5) Logging, audyt, monitorowanie  
6) Retencja i usuwanie danych  
7) Dostawcy i transfery transgraniczne  
8) Ciągłość działania i DR  
9) Luka → plan działania → dowody  
10) Ryzyka, decyzje, pytania

## Wymagane rozwinięcia
- Data flow diagram (DFD) z klasami danych.  
- Macierz SoD i dostępów uprzywilejowanych.  
- Plan retencji/usuwania i mechanizmy egzekucji.  
- Lista kontrolna szyfrowania (at-rest/in-transit, klucze/KMS).  
- Plan działań naprawczych z terminami i właścicielami.  
- Dowody audytowe: logi, konfiguracje, raporty testów.

## Wymagane streszczenia
- Executive summary: status zgodności, top 5 luk i ich ryzyko.  
- Snapshot dostawców z oceną ryzyka i miejscem przetwarzania danych.

## Guidance (skrót)
- Zacznij od przepływów danych i klasyfikacji; bez nich kontrola nie jest kompletna.  
- Preferuj SSO/MFA, least privilege, rotację kluczy, centralne KMS/HSM.  
- Logi: pełne, niezmienialne, z korelacją w SIEM; testuj alerty.  
- Retencja/usuwanie musi być egzekwowana automatycznie; dokumentuj wyjątki.  
- Każdą kontrolę wiąż z konkretną regulacją i dowodem audytowym.  
- Ustal cykl przeglądów (np. kwartalny) i własność luk.

## Szybkie powiązania
- linkage_index.jsonl (compliance/architecture/review)  
- security_controls_reference, data_protection_compliance, retention_policy, logging_and_audit_trail

## Jak używać dokumentu
1. Zbierz dane wejściowe (DFD, klasyfikacja, IAM, dostawcy).  
2. Oceń kontrole vs wymagania regulacyjne; spisz luki i ryzyko.  
3. Zdefiniuj działania naprawcze i właścicieli; wprowadź do backlogu.  
4. Waliduj wdrożenie i zarchiwizuj dowody audytowe; odhacz DoD.

## Checklisty Definition of Ready (DoR)
- [ ] Aktualne DFD i klasy danych.  
- [ ] Lista systemów, dostawców i lokalizacji danych.  
- [ ] Polityki: IAM, retencja, szyfrowanie, logging.  
- [ ] Wyniki ostatnich testów bezpieczeństwa dostępne.  
- [ ] Zidentyfikowane regulacje (RODO/PCI/HIPAA/GxP itp.).

## Checklisty Definition of Done (DoD)
- [ ] Luki zmapowane na regulacje; ryzyka ocenione.  
- [ ] Plan działań naprawczych z właścicielami i terminami.  
- [ ] Dowody wdrożenia kontroli zarchiwizowane.  
- [ ] Linkage_index/CMDB zaktualizowane; decyzje zapisane w ADR.  
- [ ] Harmonogram kolejnego przeglądu ustalony.

## Definicje robocze
- SoD (Segregation of Duties): rozdział uprawnień redukujący nadużycia.  
- DPIA: ocena skutków dla ochrony danych (RODO art. 35).  
- Evidence: materiał potwierdzający kontrolę (log, konfiguracja, ticket).

## Przykłady użycia
- Przegląd architektury e‑commerce pod kątem PCI DSS.  
- Walidacja systemu medycznego (PHI) wobec HIPAA/GxP.  
- Ocena SaaS z danymi UE i transferami poza EOG.

## Ryzyka i ograniczenia
- Niepełne DFD → ukryte przepływy danych.  
- Brak SoD → nadużycia i audytowe niezgodności.  
- Nieegzekwowana retencja → naruszenia RODO/PCI.  
- Dostawca bez SLA bezpieczeństwa → ryzyko transferu danych.

## Decyzje i uzasadnienia
- Przyjęte standardy i priorytety regulacyjne.  
- Model szyfrowania (KMS/HSM) i rotacja kluczy.  
- Zakres logowania i czas retencji logów.  
- Kryteria akceptacji ryzyka/wyjątków.

## Założenia
- Dane o przepływach i dostawcach są aktualne.  
- Dostępne są narzędzia SIEM/DLP/KMS i polityki bezpieczeństwa.  
- Zespoły produktowe dostarczą konfiguracje do audytu.

## Otwarte pytania
- Czy istnieją transfery transgraniczne wymagające SCC lub BCR?  
- Jakie minimalne wymagania SoD dla administratorów i developerów?  
- Jakie okresy retencji logów są wymagane przez regulatora/klientów?  
- Jak będzie weryfikowana skuteczność kontroli (testy okresowe)?
