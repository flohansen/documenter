package app

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main application configuration structure.
// It contains settings for documentation scraping, scraping intervals, and logging.
type Config struct {
	Docs     DocsConfig     `yaml:"docs"`     // Documentation configuration
	Scraping ScrapingConfig `yaml:"scraping"` // Scraping behavior configuration
	Logging  LoggingConfig  `yaml:"logging"`  // Logging output configuration
}

// DocsConfig contains configuration for documentation sections to be processed.
type DocsConfig struct {
	Sections []SectionConfig `yaml:"sections"` // List of documentation sections
}

// SectionConfig defines configuration for a single documentation section.
// Each section represents a source of documentation with its type, location, and access credentials.
type SectionConfig struct {
	Name   string      `yaml:"name"`   // Human-readable name for the section
	Type   SectionType `yaml:"type"`   // Type of the documentation source (e.g., git)
	URL    string      `yaml:"url"`    // URL or path to the documentation source
	SSHKey string      `yaml:"sshKey"` // SSH key for authentication (if required)
}

// ScrapingConfig defines how frequently the application should scrape documentation sources.
type ScrapingConfig struct {
	Interval time.Duration `yaml:"interval"` // Time interval between scraping operations
}

// LoggingConfig specifies the format for application log output.
type LoggingConfig struct {
	Format LoggingFormat `yaml:"format"` // Format for log messages
}

// SectionType represents the different types of documentation sources supported.
type SectionType int

const (
	// SectionTypeGit represents a Git repository as a documentation source
	SectionTypeGit SectionType = iota
)

// UnmarshalYAML implements yaml.Unmarshaler to parse SectionType from YAML.
// It converts string values from YAML into the appropriate SectionType constant.
func (t *SectionType) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "git":
		*t = SectionTypeGit
	default:
		return fmt.Errorf("unknown section type: %s", value.Value)
	}

	return nil
}

// LoggingFormat represents the available formats for log output.
type LoggingFormat int

const (
	// LoggingFormatText represents plain text log format
	LoggingFormatText LoggingFormat = iota
	// LoggingFormatJSON represents structured JSON log format
	LoggingFormatJSON
)

// UnmarshalYAML implements yaml.Unmarshaler to parse LoggingFormat from YAML.
// It converts string values from YAML into the appropriate LoggingFormat constant.
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
