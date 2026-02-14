# PATHS — Struktura katalogów docflow

## Root projektu

Root projektu to katalog zawierający plik `docflow.yaml` lub `.docflow/`.
CLI szuka konfiguracji w bieżącym katalogu, następnie w `~/.config/docflow/`.

## Struktura katalogów

```
docflow/                        # Root projektu
├── cmd/docflow/main.go         # Punkt wejścia CLI
├── pkg/                        # Pakiety publiczne (API stabilne)
│   ├── config/config.go        # Loader konfiguracji docflow.yaml
│   └── parser/                 # Parsery dokumentów
│       ├── frontmatter.go      # Parser YAML frontmatter
│       └── markdown.go         # Parser nagłówków Markdown
├── internal/                   # Pakiety wewnętrzne
│   ├── cli/cli.go              # Definicja komend CLI (cobra)
│   ├── db/db.go                # Warstwa SQLite
│   ├── logger/logger.go        # Logger strukturalny
│   ├── model/model.go          # Modele danych
│   └── util/fileutil.go        # Narzędzia (walker plików)
├── docs/                       # Dokumentacja projektu
│   ├── _meta/                  # Kontrakty metadanych
│   │   ├── DOC_META_SCHEMA.md  # Schemat frontmatter
│   │   ├── DOC_DEPENDENCY_SPEC.md # Specyfikacja zależności
│   │   ├── DOC_TYPES.md        # Typy dokumentów
│   │   └── CACHE_SPEC.md       # Specyfikacja cache
│   ├── ADR/                    # Architecture Decision Records
│   ├── days/                   # Plan 90-dniowy (day_01..day_90)
│   └── PATHS.md                # Ten plik
├── LOGS/                       # Logi decyzji i badań
│   ├── DECISIONS.md            # Decyzje projektowe
│   └── RESEARCH.md             # Wyniki badań (day_02)
├── testdata/                   # Dane testowe
│   ├── templates/              # 829 szablonów Markdown
│   └── ground_truth/           # Przykłady referencyjne z pełnym frontmatter
├── tools/                      # Narzędzia pomocnicze
│   ├── bench/                  # Benchmarki
│   ├── mutations/              # Testy mutacyjne
│   └── scripts/                # Skrypty (analiza, CI)
├── docflow.yaml                # Konfiguracja projektu
├── go.mod                      # Zależności Go
└── .docflow/                   # Katalog roboczy (gitignored)
    ├── cache/                  # Cache checksums
    ├── output/                 # Wygenerowane dokumenty
    └── docflow.db              # Baza SQLite
```

## Konwencje ścieżek

- `pkg/` — pakiety z publicznym API, stabilne interfejsy
- `internal/` — pakiety wewnętrzne, nie importowane z zewnątrz
- `cmd/` — punkty wejścia binariów
- `testdata/` — dane testowe, nie kompilowane przez Go
- `.docflow/` — katalog roboczy, dodany do .gitignore
