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

// Scraper defines the interface for documentation scrapers.
// Implementations should be able to scrape content from their respective sources
// and return the scraped data as bytes.
type Scraper interface {
	// Scrape extracts documentation content from the configured source.
	// It returns the scraped content as bytes or an error if scraping fails.
	Scrape(ctx context.Context) ([]byte, error)
}

//go:generate mockgen -destination=mocks/logger.go -package=mocks . Logger

// Logger defines the interface for application logging.
// It provides methods for different log levels with formatted output.
type Logger interface {
	// Warn logs a warning message with optional formatted arguments.
	Warn(format string, v ...any)
	// Info logs an informational message with optional formatted arguments.
	Info(format string, v ...any)
}

// Cli represents the main command-line interface application.
// It manages the configuration, scrapers, and logging for the documentation system.
type Cli struct {
	Config   Config    // Application configuration
	Scrapers []Scraper // List of active scrapers
	Logger   Logger    // Logger instance for application logging
}

// NewCli creates a new CLI instance with the provided configuration.
// It initializes scrapers based on the configuration sections and sets up
// the appropriate logger format. Scrapers are created for each configured
// documentation section, and unsupported section types are skipped.
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

// Run starts the CLI application and begins the scraping process.
// It launches each configured scraper in its own goroutine and waits for
// all scrapers to complete. The method blocks until the context is cancelled
// or all scrapers have finished execution.
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

// startScraper runs a single scraper in a continuous loop.
// It periodically executes the scraper based on the configured interval
// and handles scraping errors by logging warnings. The method respects
// context cancellation and will exit when the context is done.
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
