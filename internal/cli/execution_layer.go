package cli

import (
	"path/filepath"

	execore "docflow/internal/exec"
	execache "docflow/internal/exec/cache"
	"docflow/pkg/compliance"
	"docflow/pkg/config"
	"docflow/pkg/index"
)

type executionResult struct {
	Facts []execache.DocumentFacts
	Index *index.DocumentIndex
	Stats execore.RunnerStats
}

func effectiveCacheConfig(cfg *config.Config) (enabled bool, dir string) {
	enabled = cfg.Cache.Enabled && cacheEnabled && !noCache
	dir = cfg.Cache.Dir
	if dir == "" {
		dir = ".docflow/cache"
	}
	if cacheDir != "" {
		dir = cacheDir
	}
	return enabled, dir
}

func collectExecutionResult(mode string, cfg *config.Config) (*executionResult, error) {
	enabled, dir := effectiveCacheConfig(cfg)
	facts, stats, err := execore.CollectFacts(execore.RunnerOptions{
		Mode:           mode,
		Root:           cfg.DocsRoot,
		IgnorePatterns: cfg.IgnorePatterns,
		UseCache:       enabled,
		CacheDir:       dir,
		ChangedOnly:    changedOnly,
		SinceRef:       sinceRef,
	})
	if err != nil {
		return nil, err
	}
	idx := execore.BuildIndexFromFacts(cfg.DocsRoot, facts)
	return &executionResult{
		Facts: facts,
		Index: idx,
		Stats: stats,
	}, nil
}

func contentFactsByPath(facts []execache.DocumentFacts) map[string]compliance.ContentFact {
	out := make(map[string]compliance.ContentFact, len(facts))
	for _, f := range facts {
		if f.Path == "" {
			continue
		}
		out[filepath.ToSlash(f.Path)] = compliance.ContentFact{
			Content:   f.Content,
			ReadError: readErrorFromFacts(f),
		}
	}
	return out
}

func readErrorFromFacts(f execache.DocumentFacts) string {
	if f.ParseErrorCode == "READ_ERROR" {
		return f.ParseErrorMessage
	}
	return ""
}
