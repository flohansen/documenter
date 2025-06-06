package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/flohansen/documenter/internal/docs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/yaml.v3"
)

type CliConfig struct {
	ConfigPath string
}

type Cli struct {
	config CliConfig
}

func NewCli(cfg CliConfig) *Cli {
	return &Cli{
		config: cfg,
	}
}

func (c *Cli) Run(ctx context.Context) error {
	config, err := c.readConfig()
	if err != nil {
		return fmt.Errorf("read config error: %w", err)
	}

	var wg sync.WaitGroup

	for _, section := range config.Docs.Sections {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					md, err := scrape(ctx, section)
					if err != nil {
						log.Printf("error scraping: %s", err)
					}

					log.Printf("Section '%s' (%s) has %d bytes", section.Name, section.URL, len(md))
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func scrape(ctx context.Context, config docs.SectionConfig) ([]byte, error) {
	cloneOptions, err := cloneOptionsForSection(config)
	if err != nil {
		return nil, fmt.Errorf("clone options setup error: %w", err)
	}

	repo, err := git.CloneContext(ctx, memory.NewStorage(), nil, &cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("clone error: %s", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("repository head error: %s", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("repository commit error: %s", err)
	}

	file, err := commit.File("README.md")
	if err != nil {
		return nil, fmt.Errorf("commit file error: %s", err)
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, fmt.Errorf("file reader error: %s", err)
	}

	b, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		return nil, fmt.Errorf("file reader error: %s", err)
	}

	return b, nil
}

func cloneOptionsForSection(section docs.SectionConfig) (git.CloneOptions, error) {
	cloneOptions := git.CloneOptions{
		URL:   section.URL,
		Depth: 1,
	}

	if len(section.SSHKey) > 0 {
		publicKeys, err := ssh.NewPublicKeysFromFile("git", section.SSHKey, "")
		if err != nil {
			return git.CloneOptions{}, fmt.Errorf("could not read ssh key: %w", err)
		}

		cloneOptions.Auth = publicKeys
	}

	return cloneOptions, nil
}

func (c *Cli) readConfig() (docs.Config, error) {
	f, err := os.Open(c.config.ConfigPath)
	if err != nil {
		return docs.Config{}, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	var config docs.Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return docs.Config{}, fmt.Errorf("json decode error: %w", err)
	}

	return config, nil
}
