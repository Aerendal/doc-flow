package index

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"docflow/internal/util"
	"docflow/pkg/parser"
	"docflow/pkg/templatehint"
)

type DocumentRecord struct {
	DocID          string   `json:"doc_id"`
	Path           string   `json:"path"`
	Title          string   `json:"title"`
	DocType        string   `json:"doc_type"`
	Status         string   `json:"status"`
	Version        string   `json:"version"`
	Priority       string   `json:"priority"`
	Owner          string   `json:"owner"`
	Language       string   `json:"language"`
	DependsOn      []string `json:"depends_on,omitempty"`
	ContextSources []string `json:"context_sources,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	TemplateSource string   `json:"template_source,omitempty"`
	HeadingCount   int      `json:"heading_count"`
	MaxDepth       int      `json:"max_depth"`
	Lines          int      `json:"lines"`
	Size           int64    `json:"size"`
	ModTime        string   `json:"mod_time"`
	Checksum       string   `json:"checksum,omitempty"`
	MetaChecksum   string   `json:"meta_checksum,omitempty"`
	BodyChecksum   string   `json:"body_checksum,omitempty"`
	CodeBlocks     int      `json:"code_blocks,omitempty"`
	Tables         int      `json:"tables,omitempty"`
	Cycle          []string `json:"cycle,omitempty"`
}

type DocumentIndex struct {
	Version   string            `json:"version"`
	CreatedAt string            `json:"created_at"`
	Count     int               `json:"count"`
	Checksum  string            `json:"checksum"`
	Root      string            `json:"root,omitempty"`
	Documents []*DocumentRecord `json:"documents"`

	mu     sync.RWMutex
	byID   map[string]*DocumentRecord
	byPath map[string]*DocumentRecord
}

type SaveOptions struct {
	NoTimestamps bool
}

func New() *DocumentIndex {
	return &DocumentIndex{
		Version:   "1.0.0",
		CreatedAt: time.Now().Format(time.RFC3339),
		byID:      make(map[string]*DocumentRecord),
		byPath:    make(map[string]*DocumentRecord),
	}
}

func (idx *DocumentIndex) Add(rec *DocumentRecord) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.Documents = append(idx.Documents, rec)
	if rec.DocID != "" {
		idx.byID[rec.DocID] = rec
	}
	idx.byPath[rec.Path] = rec
	idx.Count = len(idx.Documents)
}

func (idx *DocumentIndex) GetByID(docID string) *DocumentRecord {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.byID[docID]
}

func (idx *DocumentIndex) GetByPath(path string) *DocumentRecord {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.byPath[path]
}

func (idx *DocumentIndex) All() []*DocumentRecord {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	result := make([]*DocumentRecord, len(idx.Documents))
	copy(result, idx.Documents)
	return result
}

func (idx *DocumentIndex) computeChecksum() string {
	data, _ := json.Marshal(idx.Documents)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:8])
}

func (idx *DocumentIndex) SaveJSON(path string) error {
	return idx.SaveJSONWithOptions(path, SaveOptions{})
}

func (idx *DocumentIndex) SaveJSONWithOptions(path string, opts SaveOptions) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	sort.Slice(idx.Documents, func(i, j int) bool {
		return idx.Documents[i].Path < idx.Documents[j].Path
	})
	idx.Count = len(idx.Documents)
	idx.Checksum = idx.computeChecksum()

	createdAt := idx.CreatedAt
	if opts.NoTimestamps {
		createdAt = "0001-01-01T00:00:00Z"
	}

	out := struct {
		Version   string            `json:"version"`
		CreatedAt string            `json:"created_at"`
		Count     int               `json:"count"`
		Checksum  string            `json:"checksum"`
		Root      string            `json:"root,omitempty"`
		Documents []*DocumentRecord `json:"documents"`
	}{
		Version:   idx.Version,
		CreatedAt: createdAt,
		Count:     idx.Count,
		Checksum:  idx.Checksum,
		Root:      idx.Root,
		Documents: idx.Documents,
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("nie można utworzyć katalogu %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("błąd serializacji indeksu: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("błąd zapisu indeksu do %s: %w", path, err)
	}

	return nil
}

func LoadJSON(path string) (*DocumentIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nie można odczytać indeksu %s: %w", path, err)
	}

	idx := &DocumentIndex{
		byID:   make(map[string]*DocumentRecord),
		byPath: make(map[string]*DocumentRecord),
	}

	if err := json.Unmarshal(data, idx); err != nil {
		return nil, fmt.Errorf("błąd parsowania indeksu %s: %w", path, err)
	}

	for _, rec := range idx.Documents {
		if rec.DocID != "" {
			idx.byID[rec.DocID] = rec
		}
		idx.byPath[rec.Path] = rec
	}

	computed := idx.computeChecksum()
	if idx.Checksum != "" && computed != idx.Checksum {
		return nil, fmt.Errorf("indeks %s: checksum niezgodny (oczekiwano %s, obliczono %s) — plik może być uszkodzony", path, idx.Checksum, computed)
	}

	return idx, nil
}

var nonAlphanumRegex = regexp.MustCompile(`[^a-z0-9_-]+`)

func NormalizeDocID(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strings.ToLower(name)
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, name)
	name = nonAlphanumRegex.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	return name
}

func NormalizePath(path string) string {
	path = filepath.ToSlash(path)
	path = filepath.Clean(path)
	return path
}

func shortHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:8])
}

func BuildIndex(docsRoot string, ignorePatterns []string) (*DocumentIndex, error) {
	absRoot, err := filepath.Abs(docsRoot)
	if err != nil {
		return nil, fmt.Errorf("błąd skanowania: %w", err)
	}

	files, err := util.WalkMarkdown(absRoot, ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("błąd skanowania: %w", err)
	}

	idx := New()
	idx.Root = NormalizePath(absRoot)
	adj := map[string][]string{}

	for _, f := range files {
		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}

		content := string(data)
		rec := &DocumentRecord{
			Path:    NormalizePath(f.RelPath),
			Size:    f.Size,
			ModTime: time.Unix(f.ModTime, 0).Format(time.RFC3339),
		}

		fmResult, err := parser.ParseFrontmatterString(content, f.RelPath)
		if err == nil && fmResult.HasFrontmatter {
			fm := fmResult.Frontmatter
			rec.Title = fm.Title
			rec.DocType = fm.DocType
			rec.Status = fm.Status
			rec.Version = fm.Version
			rec.Priority = fm.Priority
			rec.Owner = fm.Owner
			rec.Language = fm.Language
			rec.DependsOn = fm.DependsOn
			rec.ContextSources = fm.ContextSources
			rec.Tags = fm.Tags
			rec.TemplateSource = fm.TemplateSource

			if fm.DocID != "" {
				rec.DocID = fm.DocID
			}
			rec.MetaChecksum = shortHash(fmResult.RawYAML)
			rec.BodyChecksum = shortHash(fmResult.Body)
		} else {
			meta := parser.ExtractMetaBlock(content)
			if meta != nil {
				rec.Title = meta.Title
				rec.Status = meta.Status
				rec.Owner = meta.Owner
				rec.Version = meta.Version
			}
			if rec.DocType == "" {
				rec.DocType = "unknown"
			}
			if rec.Status == "" {
				rec.Status = "draft"
			}
			if rec.Language == "" {
				rec.Language = "pl"
			}
			rec.BodyChecksum = shortHash(content)
		}

		if rec.DocID == "" {
			rec.DocID = NormalizeDocID(f.RelPath)
		}

		adj[rec.DocID] = append([]string{}, rec.DependsOn...)

		ms := parser.ParseHeadingsString(content)
		rec.HeadingCount = len(ms.Headings)
		rec.MaxDepth = ms.MaxDepth()
		rec.Lines = ms.Lines
		rec.Checksum = shortHash(content)
		hints := templatehint.Extract(content)
		rec.CodeBlocks = hints.CodeBlocks
		rec.Tables = hints.Tables

		idx.Add(rec)
	}

	markCycles(idx, adj)

	return idx, nil
}
