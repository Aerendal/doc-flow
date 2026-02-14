---
title: Disaster Recovery Plan
status: needs_content
---

# Disaster Recovery Plan (DRP)

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Plan odtworzenia systemów/usług po katastrofie: RTO/RPO, role, procedury, testy.

## Struktura sekcji (szkielet)
1. Zakres i krytyczne systemy; RTO/RPO.
2. Scenariusze DR (DC loss, data corruption, ransomware, provider outage).
3. Strategie: backup/restore, standby/active-active, failover, dane (snapshots, replication).
4. Procedury DR: kroki per scenariusz, runbooki, kolejność usług.
5. Role i komunikacja (incident/DR lead, stakeholders, status page).
6. Testy/ćwiczenia i harmonogram, kryteria sukcesu.
7. Dokumentacja i audyt; post-mortem i poprawki.

## Szybkie powiązania
- Zależność: {'depends_on': 'Backup Verification Checklist', 'type': 'REFERENCES'}
- Zależność: {'depends_on': 'Backup and Recovery Strategy', 'type': 'REFERENCES'}
- Zależność: {'depends_on': 'Backup System Setup', 'type': 'REFERENCES'}

## Jak używać dokumentu
- Przeczytaj sekcje "Cel dokumentu" i "Zakres i granice" i upewnij się, że opisują Twój przypadek.
- Wypełniaj kolejne sekcje zgodnie z guidance i powiązaniami; korzystaj z kryteriów DoR/DoD w `reports/checklist_atomic.jsonl`.
- Aktualizuj statusy w checklistach (structure/clarity/links, DoR/DoD), gdy sekcje są gotowe lub oznaczone jako N/A.


## Checklisty jakości
- [ ] RTO/RPO dla krytycznych usług zdefiniowane.
- [ ] Procedury failover/restore per scenariusz opisane.
- [ ] Testy DR wykonywane wg harmonogramu.
- [ ] Komunikacja/role i statusy przygotowane.
