# Sprint Plan v0.2 (10 dni roboczych)

## Zadania (P0/P1)
1. Integracja realnych danych do recommend/template-sets/templates (usuń demo) — P0
2. Flagi CLI: `--strict-mode`, `--publish-strict`, `migrate-sections --dry-run` — P0
3. Usage store podpięte do generatora + konsumowane w recommend — P1
4. Raport brak przykładów w CLI — P1
5. Governance CLI raport + rozszerzone reguły — P1

## Kryteria sukcesu sprintu
- recommend/template-sets/templates działają na realnym indeksie (status/usage/quality z plików)
- generate zapisuje usage, recommend korzysta z usage store
- validate obsługuje flagi strict/publish, migrate-sections dostępne w CLI (dry-run)
- Raport brak przykładów wypisuje doc_id bez code/table
- Governance raport CLI pokazuje naruszenia i status exit=1, gdy błędy

## Rytm
- Daily sync: 15 min
- Review/demo: koniec sprintu (d) + log w LOGS/
- Retro: 30 min po demo

## Ryzyka sprintu
- Brak wypełnionego `template_source` i statusów w realnych dokumentach → trzeba uzupełnić lub dodać fallback
- Możliwy wzrost czasu skanu po dodaniu usage/governance raportów → monitorować
