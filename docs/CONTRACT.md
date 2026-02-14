# CLI Contract

## Exit Codes

- `0`: success (including `--warn` mode)
- `1`: domain violation in `--strict` mode
- `2`: usage error (invalid/missing flags/args)
- `3`: runtime error (I/O, config load, baseline load, rules load)

## Machine Formats

### Validate

- `--format json`
- `--format sarif` (SARIF 2.1.0)

Common baseline flags:

- `--against <baseline.json>`
- `--fail-on any|new`
- `--show all|new|existing`

JSON report fields include:

- `schema_version`
- `identity_version` (`"2"` = comparator v2, niezależny od `message`)
- `issues[]` with `code`, `level`, `type`, `path`, `doc_id`, `line`, `message`
- `issues[].details` (stabilne pola maszynowe używane przez baseline identity v2)
- baseline metadata and counters:
  - `baseline`
  - `new_error_count`, `new_warn_count`
  - `existing_error_count`, `existing_warn_count`

### Compliance

- `--format json`
- `--format html`

Common baseline flags:

- `--against <baseline.json>`
- `--fail-on any|new`
- `--show all|new|existing`

JSON report fields include:

- `schema_version`
- `identity_version` (`"2"`)
- `rules_path`, `rules_checksum`
- aggregate counters and `docs[]`
- baseline metadata and counters:
  - `baseline`
  - `new_failed`, `existing_failed`
  - `new_violations_count`
- `duplicate_doc_ids`

### Health Bundle

- `docflow health --ci` jest kanonicznym one-command entrypointem dla CI.
- Bundle runu jest zapisywany w `.docflow/out/<run_id>/` i zawiera:
  - `validate.json`
  - `validate.sarif`
  - `compliance.json`
  - `summary.json`
  - `summary.md`
  - `meta.json`
- `summary.json` zawiera `overall_exit` wyliczony jako max priorytetu `3 > 2 > 1 > 0` z etapów validate/compliance.
- W trybie `--ci` i baseline dostępnym (`repo|artifact`) health używa gatingu `--fail-on new --show new`.

### Baseline Identity

- `identity_version=2` używa klucza strukturalnego (`code`, `path`, `doc_id`, `location`, `details`) i **nie zależy od `message`**.
- Baseline bez `identity_version` traktowany jest jako v1 (kompatybilność wstecz).
- Migracja v1 -> v2: `docflow baseline migrate --in <old.json> --out <new.json> --kind validate|compliance`.

## JSON Error Envelope

When `--format json` is used and command fails before producing a normal report, stdout contains:

```json
{
  "schema_version": "1.0",
  "error": {
    "kind": "usage|runtime",
    "code": "DOCFLOW....",
    "message": "...",
    "details": {}
  }
}
```

`stderr` still contains a short human-readable error line.
