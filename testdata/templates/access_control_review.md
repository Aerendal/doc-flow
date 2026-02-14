---
title: Access Control Review
status: needs_content
---

# Access Control Review

## Metadane
- Właściciel: [IAM/Security]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

> Powiązania: linkage_index.jsonl

## Cel dokumentu
Definiuje przeglądy kontroli dostępu (recertyfikacje) dla ról/uprawnień/SoD: zakres, częstotliwość, właścicieli, kryteria i dowody, aby utrzymać least privilege i zgodność (SOX/PCI/RODO).

## Zakres i granice
- Obejmuje: typy przeglądów (rola/user/access/SoD), częstotliwości, zakres systemów/danych, właścicieli i approverów, procedurę przeglądu, wyjątki/waivery z sunset, dowody i raporty, integracje z IAM/CMDB/ticketing.
- Poza zakresem: projekt macierzy AC (osobno) i operacyjne nadania (JML workflow).

## Wejścia i wyjścia
- Wejścia: macierze ról/uprawnień, SoD rules, listy użytkowników i dostępów, CMDB/asset, wymagania audyt/regulator, harmonogram JML/zmian, wcześniejsze odchylenia.
- Wyjścia: decyzje keep/remove/adjust, waivery z sunset/kompensacjami, raport przeglądu, action items (owner/ETA), metryki (completion, violations, findings), dowody dla audytu.

## Powiązania (meta)
- Key Documents: access_control_matrix_design, access_control_policy, access_control_improvement_plan, multi_factor_authentication_design, logging_and_audit_trail, security_controls_reference, risk_register.
- Dependencies: IdP/IAM, CMDB/asset, HR/JML, ticketing/workflow, SIEM/logi, SoD rules, audyt schedule.

## Zależności dokumentu
- Upstream: macierze ról, SoD rules, listy access, wymagania audyt/regulator, harmonogram JML/zmian.
- Downstream: zmiany access (remove/adjust), waivery, raporty audytowe, risk register aktualizacje, improvement backlog.
- Zewnętrzne: audytorzy/regulatorzy.

## Powiązania sekcja↔sekcja
- Zakres/frequency → procedura → decyzje → waivery/action items → raport/KPI.

## Fazy cyklu życia
- Planowanie (zakres, częstotliwość, ownerzy).
- Wykonanie przeglądu.
- Decyzje i wdrożenie zmian (remove/adjust/waiver).
- Raportowanie i retrospektywa.

## Struktura sekcji
1) Streszczenie (zakres, completion, violations, top findings)  
2) Zakres i częstotliwość przeglądów (systemy, role, SoD, user access)  
3) Procedura przeglądu (dane wejściowe, narzędzia, kroki, SLA)  
4) Rola/owner/approver, podział obowiązków (SoD)  
5) Decyzje i egzekucja (remove/adjust/waiver, ticketing, ETA)  
6) Waivery/wyjątki (powód, kompensacje, sunset, przegląd)  
7) Raportowanie i dowody (audyt, KPI: completion, violations, time-to-close)  
8) Ryzyka i zależności; decyzje (ADR) i otwarte pytania  

## Wymagane rozwinięcia
- Harmonogram przeglądów; wzór raportu; lista ownerów/approverów; procedura SoD; log waivers; checklisty recertyfikacji.
- Integracja z IAM/CMDB/ticketing; dowody (exporty, podpisy, logi).

## Wymagane streszczenia
- Executive summary: completion, violations, top findings, waivery, rekomendacje.
- One-pager: zakres, wyniki, waivery, terminy zamknięcia.

## Guidance (skrót)
- DoR: zakres/systemy/role/SoD zebrane; ownerzy/approverzy znani; narzędzia IAM/CMDB/ticketing dostępne; SLA i wymagania audytu/regulatora znane.
- DoD: przegląd wykonany; decyzje i zmiany w ticketingu; waivery z sunset; raport/KPI; metadane aktualne; dokument w linkage_index.
- Spójność: każde uprawnienie ma decyzję; waivery mają sunset/kompensacje; raport zawiera dowody i KPI.

## Szybkie powiązania
- access_control_matrix_design, access_control_improvement_plan, access_control_policy, multi_factor_authentication_design, logging_and_audit_trail, security_controls_reference, risk_register

## Checklisty Definition of Ready (DoR)
- [ ] Zakres/systemy/role/SoD i ownerzy/approverzy zebrani; narzędzia IAM/CMDB/ticketing dostępne; SLA/requirement audytu znane.

## Checklisty Definition of Done (DoD)
- [ ] Przegląd wykonany; decyzje wdrożone; waivery z sunset/kompensacjami; raport/KPI/dowody gotowe; dokument w linkage_index.

## Artefakty powiązane
- Harmonogram i raport przeglądów, exporty access/SoD, waiver log, ticket log, KPI dashboard, ADR log.

## Weryfikacja spójności
- [ ] Każdy dostęp/rola w scope ma decyzję (keep/remove/adjust/waiver) i dowód.
- [ ] Waivery mają sunset/kompensacje; KPI (completion, violations, time-to-close) obliczone.
- [ ] Raport zawiera dowody i jest zgodny z wymaganiami audytu/regulatora.
