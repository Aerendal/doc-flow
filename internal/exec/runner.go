package exec

import (
	"crypto/sha256"
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"sort"
	"strings"

	execache "docflow/internal/exec/cache"
	"docflow/internal/util"
	"docflow/pkg/index"
	"docflow/pkg/parser"
	"docflow/pkg/templatehint"
)

type RunnerOptions struct {
	Mode           string
	Root           string
	IgnorePatterns []string
	UseCache       bool
	CacheDir       string
	ChangedOnly    bool
	SinceRef       string
}

type RunnerStats struct {
	Files        int
	ChangedFiles int
	CacheHits    int
	CacheMisses  int
}

func CollectFacts(opts RunnerOptions) ([]execache.DocumentFacts, RunnerStats, error) {
	var stats RunnerStats

	absRoot, err := filepath.Abs(opts.Root)
	if err != nil {
		return nil, stats, err
	}

	files, err := util.WalkMarkdown(absRoot, opts.IgnorePatterns)
	if err != nil {
		return nil, stats, err
	}
	sort.Slice(files, func(i, j int) bool {
		return filepath.ToSlash(files[i].RelPath) < filepath.ToSlash(files[j].RelPath)
	})
	stats.Files = len(files)

	changedSet := map[string]bool{}
	if opts.SinceRef != "" {
		if set, err := changedMarkdownSince(absRoot, opts.SinceRef); err == nil {
			changedSet = set
		}
	}

	var (
		c      execache.Cache
		runID  string
		useDB  = opts.UseCache
		dbPath = filepath.Join(opts.CacheDir, "docflow_cache.sqlite")
	)
	if useDB {
		c, err = execache.OpenSQLite(dbPath)
		if err == nil {
			runID, _ = c.BeginRun(opts.Mode, absRoot)
		} else {
			useDB = false
		}
	}
	if c != nil {
		defer func() {
			_ = c.EndRun(runID, execache.RunStats{
				ChangedFiles: stats.ChangedFiles,
				CacheHits:    stats.CacheHits,
				CacheMisses:  stats.CacheMisses,
			})
			_ = c.Close()
		}()
	}

	out := make([]execache.DocumentFacts, 0, len(files))
	for _, f := range files {
		rel := filepath.ToSlash(f.RelPath)
		if changedSet[rel] {
			stats.ChangedFiles++
		}

		fp := execache.Fingerprint{
			Size:      f.Size,
			MTimeUnix: f.ModTime,
		}

		forceParse := opts.ChangedOnly && len(changedSet) > 0 && changedSet[rel]

		if useDB && !forceParse {
			if cached, hit, err := c.Get(rel, fp); err == nil && hit && cached != nil {
				cached.AbsPath = f.Path
				cached.Root = absRoot
				out = append(out, *cached)
				stats.CacheHits++
				continue
			}
		}

		facts := parseFileFacts(absRoot, f.Path, rel, f.Size, f.ModTime)
		out = append(out, facts)
		stats.CacheMisses++

		if useDB {
			putFP := execache.Fingerprint{
				Size:         facts.Size,
				MTimeUnix:    facts.MTimeUnix,
				MetaChecksum: facts.ChecksumMeta,
				BodyChecksum: facts.ChecksumBody,
			}
			_ = c.Put(rel, putFP, facts)
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, stats, nil
}

func BuildIndexFromFacts(root string, facts []execache.DocumentFacts) *index.DocumentIndex {
	records := make([]*index.DocumentRecord, 0, len(facts))
	for _, f := range facts {
		rec := &index.DocumentRecord{
			DocID:          f.DocID,
			Path:           index.NormalizePath(f.Path),
			Title:          f.Title,
			DocType:        f.DocType,
			Status:         f.Status,
			Version:        f.Version,
			Priority:       f.Priority,
			Owner:          f.Owner,
			Language:       f.Language,
			DependsOn:      append([]string(nil), f.DependsOn...),
			ContextSources: append([]string(nil), f.ContextSources...),
			Tags:           append([]string(nil), f.Tags...),
			TemplateSource: f.TemplateSource,
			HeadingCount:   f.HeadingCount,
			MaxDepth:       f.MaxDepth,
			Lines:          f.Lines,
			Size:           f.Size,
			ModTime:        "",
			Checksum:       f.ChecksumFull,
			MetaChecksum:   f.ChecksumMeta,
			BodyChecksum:   f.ChecksumBody,
			CodeBlocks:     f.CodeBlocks,
			Tables:         f.Tables,
		}
		if rec.DocID == "" {
			rec.DocID = index.NormalizeDocID(f.Path)
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
		records = append(records, rec)
	}
	return index.BuildIndexFromRecords(root, records)
}

func parseFileFacts(root, absPath, relPath string, size int64, mtime int64) execache.DocumentFacts {
	facts := execache.DocumentFacts{
		Path:      index.NormalizePath(relPath),
		AbsPath:   absPath,
		Root:      index.NormalizePath(root),
		Size:      size,
		MTimeUnix: mtime,
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		facts.ParseErrorCode = "READ_ERROR"
		facts.ParseErrorMessage = fmt.Sprintf("nie można odczytać pliku: %v", err)
		return facts
	}

	content := string(data)
	facts.Content = content
	facts.ChecksumFull = shortHash(content)

	fmResult, err := parser.ParseFrontmatterString(content, relPath)
	if err == nil && fmResult != nil {
		facts.HasFrontmatter = fmResult.HasFrontmatter
		facts.RawYAML = fmResult.RawYAML
		facts.Body = fmResult.Body
		facts.RawFields = fmResult.RawFields
		if fmResult.HasFrontmatter && fmResult.Frontmatter != nil {
			fm := fmResult.Frontmatter
			facts.DocID = fm.DocID
			facts.Title = fm.Title
			facts.DocType = fm.DocType
			facts.Status = fm.Status
			facts.Version = fm.Version
			facts.Priority = fm.Priority
			facts.Owner = fm.Owner
			facts.Language = fm.Language
			facts.DependsOn = append([]string(nil), fm.DependsOn...)
			facts.ContextSources = append([]string(nil), fm.ContextSources...)
			facts.Tags = append([]string(nil), fm.Tags...)
			facts.TemplateSource = fm.TemplateSource
			facts.ChecksumMeta = shortHash(fmResult.RawYAML)
			facts.ChecksumBody = shortHash(fmResult.Body)
		} else {
			if meta := parser.ExtractMetaBlock(content); meta != nil {
				facts.Title = meta.Title
				facts.Status = meta.Status
				facts.Owner = meta.Owner
				facts.Version = meta.Version
			}
			facts.ChecksumBody = shortHash(content)
		}
	} else {
		facts.ParseErrorCode = "INVALID_YAML"
		facts.ChecksumBody = shortHash(content)
		if perr, ok := err.(*parser.ParseError); ok {
			facts.ParseErrorLine = perr.Line
			facts.ParseErrorMessage = perr.Message
		} else if err != nil {
			facts.ParseErrorMessage = err.Error()
		}
	}

	if facts.DocID == "" {
		facts.DocID = index.NormalizeDocID(relPath)
	}
	if facts.DocType == "" {
		facts.DocType = "unknown"
	}
	if facts.Status == "" {
		facts.Status = "draft"
	}
	if facts.Language == "" {
		facts.Language = "pl"
	}

	ms := parser.ParseHeadingsString(content)
	facts.HeadingCount = len(ms.Headings)
	facts.MaxDepth = ms.MaxDepth()
	facts.Lines = ms.Lines
	hints := templatehint.Extract(content)
	facts.CodeBlocks = hints.CodeBlocks
	facts.Tables = hints.Tables
	return facts
}

func changedMarkdownSince(root, sinceRef string) (map[string]bool, error) {
	cmd := osexec.Command("git", "-C", root, "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	diffCmd := osexec.Command("git", "-C", root, "diff", "--name-only", sinceRef+"...HEAD")
	out, err := diffCmd.Output()
	if err != nil {
		return nil, err
	}

	changed := map[string]bool{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(line))
		if ext != ".md" && ext != ".markdown" {
			continue
		}
		changed[index.NormalizePath(line)] = true
	}
	return changed, nil
}

func shortHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:8])
}
