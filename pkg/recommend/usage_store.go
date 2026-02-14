package recommend

import (
	"encoding/json"
	"fmt"
	"os"
)

type UsageStore struct {
	Path  string
	Usage map[string]int `json:"usage"`
}

func LoadUsage(path string) (*UsageStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &UsageStore{Path: path, Usage: map[string]int{}}, nil
		}
		return nil, fmt.Errorf("cannot read usage %s: %w", path, err)
	}
	var store UsageStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("cannot parse usage %s: %w", path, err)
	}
	if store.Usage == nil {
		store.Usage = map[string]int{}
	}
	store.Path = path
	return &store, nil
}

func (s *UsageStore) Inc(path string) {
	if s.Usage == nil {
		s.Usage = map[string]int{}
	}
	s.Usage[path]++
}

func (s *UsageStore) Save() error {
	if s.Path == "" {
		return fmt.Errorf("usage store path is empty")
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal usage: %w", err)
	}
	if err := os.WriteFile(s.Path, data, 0o644); err != nil {
		return fmt.Errorf("write usage: %w", err)
	}
	return nil
}
