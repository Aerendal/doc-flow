# CACHE_SPEC — Specyfikacja systemu cache

## Wersja

`v1.0.0`

## Cel

Cache przechowuje wyniki parsowania i analizy dokumentów, eliminując konieczność ponownego przetwarzania niezmienionych plików.

## Lokalizacja

```
.docflow/
  cache/
    meta/           - Cache sparsowanych metadanych (frontmatter)
    graph/          - Cache grafu zależności
    index/          - Cache indeksu dokumentów
    checksums.json  - Sumy kontrolne plików źródłowych
```

Domyślna lokalizacja: `.docflow/cache/` w katalogu roboczym.
Konfigurowalna przez `docflow.yaml` → `cache_dir`.

## Format danych

### checksums.json

Mapuje ścieżkę pliku na jego hash (SHA-256), rozmiar i czas modyfikacji:

```json
{
  "version": "1.0.0",
  "entries": {
    "docs/api_documentation.md": {
      "sha256": "abc123...",
      "size": 4096,
      "mod_time": "2026-02-09T10:00:00Z"
    }
  }
}
```

### Cache metadanych (meta/)

Jeden plik JSON per dokument, nazwany `{doc_id}.json`:

```json
{
  "doc_id": "api_documentation",
  "title": "API Documentation",
  "doc_type": "specification",
  "status": "draft",
  "depends_on": ["system_architecture"],
  "context_sources": ["product_vision_statement"],
  "source_path": "docs/api_documentation.md",
  "source_hash": "abc123...",
  "parsed_at": "2026-02-09T10:00:00Z"
}
```

### Cache grafu (graph/)

Plik `dependency_graph.json`:

```json
{
  "version": "1.0.0",
  "generated_at": "2026-02-09T10:00:00Z",
  "nodes": [...],
  "edges": [...],
  "stats": {
    "total_nodes": 100,
    "total_edges": 250,
    "components": 5,
    "max_depth": 8
  }
}
```

## Strategia invalidacji

### Reguły

1. **Zmiana pliku źródłowego** — jeśli SHA-256 pliku różni się od cache → invalidacja cache tego dokumentu.
2. **Zmiana zależności** — jeśli dokument z `depends_on` się zmienił → invalidacja cache dokumentu zależnego.
3. **Zmiana schematu** — jeśli wersja schematu metadanych się zmieniła → invalidacja całego cache.
4. **Ręczna invalidacja** — komenda `docflow cache clear`.

### Szybka detekcja zmian

Kolejność sprawdzeń (od najtańszej):
1. `mod_time` — jeśli czas modyfikacji się nie zmienił → skip.
2. `size` — jeśli rozmiar się zmienił → invalidacja.
3. `sha256` — oblicz hash i porównaj → definitywna decyzja.

## Komendy CLI

| Komenda | Opis |
|---------|------|
| `docflow cache status` | Pokaż statystyki cache (rozmiar, liczba wpisów, wiek) |
| `docflow cache clear` | Wyczyść cały cache |
| `docflow cache clear --doc <doc_id>` | Wyczyść cache konkretnego dokumentu |
| `docflow cache rebuild` | Przebuduj cache od zera |

## Limity

- Maksymalny rozmiar cache: konfigurowalne, domyślnie 100 MB.
- Maksymalny wiek wpisu: konfigurowalne, domyślnie 30 dni.
- Cache nie jest współdzielony między maszynami (lokalny).

## Konfiguracja w docflow.yaml

```yaml
cache:
  enabled: true
  dir: ".docflow/cache"
  max_size_mb: 100
  max_age_days: 30
  checksum_algorithm: "sha256"
```
