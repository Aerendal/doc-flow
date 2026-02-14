package governance

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type StatusRule struct {
	RequiredFields     []string `yaml:"required_fields"`
	AllowEmptySections bool     `yaml:"allow_empty_sections"`
	MinQuality         int      `yaml:"min_quality"`
	ApprovalsRequired  int      `yaml:"approvals_required"`
}

type FamilyRule struct {
	RequiredSections    []string `yaml:"required_sections"`
	AllowedStatus       []string `yaml:"allowed_status"`
	MinQualityPublished int      `yaml:"min_quality_published"`
}

type Rules struct {
	Statuses  map[string]StatusRule `yaml:"statuses"`
	Families  map[string]FamilyRule `yaml:"families"`
	Approvals map[string][]string   `yaml:"approvals"`
}

func Load(path string) (*Rules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read governance rules: %w", err)
	}
	var r Rules
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parse governance rules: %w", err)
	}
	return &r, nil
}
