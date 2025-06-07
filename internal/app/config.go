package app

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Docs     DocsConfig     `yaml:"docs"`
	Scraping ScrapingConfig `yaml:"scraping"`
	Logging  LoggingConfig  `yaml:"logging"`
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

type ScrapingConfig struct {
	Interval time.Duration `yaml:"interval"`
}

type LoggingConfig struct {
	Format LoggingFormat `yaml:"format"`
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

type LoggingFormat int

const (
	LoggingFormatText LoggingFormat = iota
	LoggingFormatJSON
)

func (t *LoggingFormat) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "text":
		*t = LoggingFormatText
	case "json":
		*t = LoggingFormatJSON
	default:
		return fmt.Errorf("unknown logging format: %s", value.Value)
	}

	return nil
}
