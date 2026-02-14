---
title: Data Privacy Checklist
status: needs_content
---

# Data Privacy Checklist

## Metadane
- Właściciel: [Privacy/Legal/Security]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Lista kontrolna spełnienia wymogów prywatności dla produktu/systemu/feature – do szybkiej walidacji przed release lub audytem.

## Zakres i granice
- Obejmuje: dane i cele, podstawy prawne, notice/zgody, prawa osób (DSAR), rejestry (ROP/DPIA), transfery i umowy (SCC/BCR/DPA), retencję/usuwanie, bezpieczeństwo danych, logi/audyt, szkolenia, waivery i wyjątki.
- Poza zakresem: pełne raporty DPIA/PIA i polityka prywatności (oddzielne dokumenty).

## Struktura checklisty (przykładowe sekcje)
1) Dane i cele  
2) Podstawy prawne i notice/zgody  
3) Prawa osób (DSAR) i SLA  
4) Rejestry (ROP, DPIA)  
5) Transfery i umowy (SCC/BCR/DPA)  
6) Retencja i usuwanie/deidentyfikacja  
7) Bezpieczeństwo (IAM, szyfrowanie, DLP, logi, backup/DR)  
8) Logi/audyt i monitoring  
9) Szkolenia/komunikacja  
10) Waivery/wyjątki i kompensacje  

## Przykładowe punkty DoR/DoD (do oznaczenia ✔/N/A)
- [ ] Dane/kategorie i cele zidentyfikowane; podstawa prawna wskazana.  
- [ ] Notice/zgody wdrożone; UI/treści zatwierdzone.  
- [ ] DSAR: proces, narzędzia, SLA/logi działają.  
- [ ] ROP i DPIA aktualne; zmiany odnotowane.  
- [ ] Transfery/podmioty: SCC/BCR/DPA podpisane; lokalizacja danych znana.  
- [ ] Retencja/usuwanie/depers: polityka i implementacja, dowody.  
- [ ] Bezpieczeństwo: IAM, szyfrowanie, DLP, logi, backup/DR; kontrolki testowane.  
- [ ] Logi/audyt dostępne; monitoring KPI/KRI prywatności.  
- [ ] Szkolenia/awareness wykonane; rekordy szkoleń.  
- [ ] Waivery/wyjątki mają sunset i kompensacje.  

## Guidance (skrót)
- Używaj checklisty jako pre-release/audit gate; oznacz N/A z uzasadnieniem.  
- Każdy punkt powinien mieć dowód (link, ticket, dokument); jeśli brak – zapisz action item.  
- Aktualizuj checklistę po zmianach funkcji/danych/transferów i po incydentach.  

## Szybkie powiązania
- data_privacy_assessment, data_privacy_compliance_plan, data_privacy_compliance, records_of_processing, privacy_policy, data_retention_policy, access_control_policy, incident_response_runbook\n+\n+## Checklisty Definition of Ready (DoR)\n+- [ ] Zakres funkcji/systemu i kategorie danych znane; podstawa prawna i transfery zidentyfikowane.\n+- [ ] ROP/DPIA/umowy dostępne; ownerzy odpowiedzi na DSAR/notice/zgody wskazani.\n+\n+## Checklisty Definition of Done (DoD)\n+- [ ] Wszystkie punkty checklisty oznaczone ✔/N/A z dowodami; waivery z sunset/kompensacjami; metadane aktualne; dokument w linkage_index.\n+\n+## Artefakty powiązane\n+- Wypełniona checklista, dowody (linki), ROP/DPIA, umowy SCC/BCR/DPA, log DSAR/consent, raporty z testów bezpieczeństwa, waiver log.\n+\n+## Weryfikacja spójności\n+- [ ] Punkty oznaczone ✔ mają dowód; N/A mają uzasadnienie.\n+- [ ] Transfery/umowy/retencja spójne z ROP/DPIA; waivery mają sunset.\n+- [ ] DSAR/zgody, logi/audyt i bezpieczeństwo są pokryte i udokumentowane.\n*** End Patch to=functions.apply_patch  supervisory-tone##" style="background-color: #eee"/>"}쁘
