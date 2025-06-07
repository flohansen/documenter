package app

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/flohansen/documenter/internal/scraper"
)

//go:generate mockgen -destination=mocks/scraper.go -package=mocks . Scraper

type Scraper interface {
	Scrape(ctx context.Context) ([]byte, error)
}

//go:generate mockgen -destination=mocks/logger.go -package=mocks . Logger

type Logger interface {
	Warn(format string, v ...any)
	Info(format string, v ...any)
}

type Cli struct {
	Config   Config
	Scrapers []Scraper
	Logger   Logger
}

func NewCli(cfg Config) *Cli {
	var scrapers []Scraper
	for _, section := range cfg.Docs.Sections {
		var s Scraper
		switch section.Type {
		case SectionTypeGit:
			s = scraper.NewGitScraper(section.URL, scraper.WithSSHKey(section.SSHKey))
		default:
			continue
		}

		scrapers = append(scrapers, s)
	}

	var loggerHandler slog.Handler
	switch cfg.Logging.Format {
	case LoggingFormatJSON:
		loggerHandler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		loggerHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(loggerHandler)

	return &Cli{
		Config:   cfg,
		Scrapers: scrapers,
		Logger:   logger,
	}
}

func (c *Cli) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, s := range c.Scrapers {
		wg.Add(1)

		go func() {
			defer wg.Done()
			c.startScraper(ctx, s)
		}()
	}

	wg.Wait()
	return nil
}

func (c *Cli) startScraper(ctx context.Context, scraper Scraper) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.Config.Scraping.Interval):
			md, err := scraper.Scrape(ctx)
			if err != nil {
				c.Logger.Warn("scraping error", "error", err)
				continue
			}

			c.Logger.Info("scrape successful", "received bytes", len(md))
		}
	}
}
