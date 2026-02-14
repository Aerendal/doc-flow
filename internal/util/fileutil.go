package util

import (
	"os"
	"path/filepath"
	"strings"
)

type WalkResult struct {
	Path    string
	RelPath string
	Size    int64
	ModTime int64
}

func WalkMarkdown(root string, ignorePatterns []string) ([]WalkResult, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var results []WalkResult

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)

		if info.IsDir() {
			if shouldIgnore(relPath, info.Name(), ignorePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIgnore(relPath, info.Name(), ignorePatterns) {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".md" && ext != ".markdown" {
			return nil
		}

		results = append(results, WalkResult{
			Path:    path,
			RelPath: relPath,
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})

		return nil
	})

	return results, err
}

func shouldIgnore(relPath, name string, patterns []string) bool {
	if name == "." {
		return false
	}

	for _, pattern := range patterns {
		if name == pattern {
			return true
		}

		if strings.HasPrefix(name, ".") && pattern == ".*" {
			return true
		}

		matched, err := filepath.Match(pattern, name)
		if err == nil && matched {
			return true
		}

		parts := strings.Split(filepath.ToSlash(relPath), "/")
		for _, part := range parts {
			if part == pattern {
				return true
			}
			matched, err := filepath.Match(pattern, part)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}
