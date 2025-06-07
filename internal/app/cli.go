package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/flohansen/documenter/internal/scraper"
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

			switch section.Type {
			case SectionTypeGit:
				startScraper(ctx, scraper.NewGitScraper(
					section.URL,
					scraper.WithSSHKey(section.SSHKey),
				))
			}
		}()
	}

	wg.Wait()
	return nil
}

type Scraper interface {
	Scrape(ctx context.Context) ([]byte, error)
}

func startScraper(ctx context.Context, scraper Scraper) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			md, err := scraper.Scrape(ctx)
			if err != nil {
				log.Printf("error scraping: %s", err)
				continue
			}

			log.Printf("received %d bytes", len(md))
		}
	}
}

func (c *Cli) readConfig() (Config, error) {
	f, err := os.Open(c.config.ConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("json decode error: %w", err)
	}

	return config, nil
}
