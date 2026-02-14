---
title: Backup Verification Checklist
status: needs_content
---

# Backup Verification Checklist

## Metadane
- Właściciel: [osoba/rola]
- Wersja: v0.1
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved


## Cel dokumentu
Operacyjna checklista, aby potwierdzić, że test odtwarzania backupów spełnia wymagania RPO/RTO, integralności i zgodności (audyt/regulator).

## Checklist
- [ ] Zakres backupów (systemy/dane) zdefiniowany, retencja i lokalizacje znane.
- [ ] RPO/RTO powiązane z testem; cele potwierdzone przez właścicieli.
- [ ] Scenariusze full/partial/PITR (w tym ransomware/tampering detection) wykonane zgodnie z planem.
- [ ] Środowisko testowe izolowane; dostępy/klucze (KMS/HSM) zweryfikowane.
- [ ] Dane po przywróceniu: integralność/spójność potwierdzona (checksums, DB consistency).
- [ ] Czas odtwarzania zmierzony i ≤ RTO; RPO osiągnięte.
- [ ] Dowody z testu (logi, checksumy, zrzuty, ticket) zapisane i dostępne do audytu/regulatora.
- [ ] Action items/remediacja zidentyfikowane, przypisani ownerzy i terminy; follow-up zaplanowany.

## Szybkie powiązania
- Zależność: {'depends_on': 'Database Schema Design', 'type': 'REFERENCES'}
- Zależność: {'depends_on': 'Schema Implementation', 'type': 'REFERENCES'}
- Zależność: {'depends_on': 'Database Schema Reference', 'type': 'REFERENCES'}
