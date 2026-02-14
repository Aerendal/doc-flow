package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DocsRoot       string        `yaml:"docs_root"`
	IgnorePatterns []string      `yaml:"ignore_patterns"`
	OutputDir      string        `yaml:"output_dir"`
	Cache          CacheConfig   `yaml:"cache"`
	Log            LogConfig     `yaml:"log"`
	DB             DBConfig      `yaml:"db"`
	Dependency     DependencyCfg `yaml:"dependency"`
	SectionAliases map[string][]string `yaml:"section_aliases"`
}

type CacheConfig struct {
	Enabled           bool   `yaml:"enabled"`
	Dir               string `yaml:"dir"`
	MaxSizeMB         int    `yaml:"max_size_mb"`
	MaxAgeDays        int    `yaml:"max_age_days"`
	ChecksumAlgorithm string `yaml:"checksum_algorithm"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	File   string `yaml:"file"`
}

type DBConfig struct {
	Path string `yaml:"path"`
}

type DependencyCfg struct {
	// promote_context_for: lista doc_type, dla których context_sources traktujemy jako zależność twardą.
	PromoteContextFor []string `yaml:"promote_context_for"`
	// family_rules_path: opcjonalny plik z regułami zależności per rodzina (doc_type).
	FamilyRulesPath string `yaml:"family_rules_path"`
}

const DefaultConfigFile = "docflow.yaml"

func Default() *Config {
	return &Config{
		DocsRoot: ".",
		IgnorePatterns: []string{
			".git",
			".docflow",
			"node_modules",
			"vendor",
			"_build",
		},
		OutputDir: ".docflow/output",
		Cache: CacheConfig{
			Enabled:           true,
			Dir:               ".docflow/cache",
			MaxSizeMB:         100,
			MaxAgeDays:        30,
			ChecksumAlgorithm: "sha256",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
			File:   "",
		},
		DB: DBConfig{
			Path: ".docflow/docflow.db",
		},
		Dependency: DependencyCfg{
			PromoteContextFor: []string{},
		},
		SectionAliases: map[string][]string{},
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("nie można odczytać pliku konfiguracji %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("błąd parsowania %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("błąd walidacji konfiguracji: %w", err)
	}

	return cfg, nil
}

func FindConfigFile() string {
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		return DefaultConfigFile
	}

	home, err := os.UserHomeDir()
	if err == nil {
		global := filepath.Join(home, ".config", "docflow", DefaultConfigFile)
		if _, err := os.Stat(global); err == nil {
			return global
		}
	}

	return ""
}

func (c *Config) Validate() error {
	if c.DocsRoot == "" {
		return fmt.Errorf("docs_root nie może być pusty")
	}

	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[c.Log.Level] {
		return fmt.Errorf("nieprawidłowy poziom logowania: %s (dozwolone: debug, info, warn, error)", c.Log.Level)
	}

	validFormats := map[string]bool{
		"text": true, "json": true,
	}
	if !validFormats[c.Log.Format] {
		return fmt.Errorf("nieprawidłowy format logowania: %s (dozwolone: text, json)", c.Log.Format)
	}

	if c.Cache.MaxSizeMB <= 0 {
		return fmt.Errorf("cache.max_size_mb musi być > 0")
	}

	if c.Cache.MaxAgeDays <= 0 {
		return fmt.Errorf("cache.max_age_days musi być > 0")
	}

	return nil
}
