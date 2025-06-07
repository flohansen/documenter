package app

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Docs DocsConfig `yaml:"docs"`
}

type DocsConfig struct {
	Sections []SectionConfig `yaml:"sections"`
}

type SectionConfig struct {
	Name   string      `yaml:"name"`
	Type   SectionType `yaml:"type"`
	URL    string      `yaml:"url"`
	SSHKey string      `yaml:"sshKey"`
}

type SectionType int

const (
	SectionTypeGit SectionType = iota
)

func (t *SectionType) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "git":
		*t = SectionTypeGit
	default:
		return fmt.Errorf("unknown section type: %s", value.Value)
	}

	return nil
}
