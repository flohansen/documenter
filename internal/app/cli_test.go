package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/flohansen/documenter/internal/app"
	"github.com/flohansen/documenter/internal/app/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCli_NewCli(t *testing.T) {
	t.Run("should create git scraper using config", func(t *testing.T) {
		// assign
		config := app.Config{
			ScrapeInterval: 5 * time.Second,
			Docs: app.DocsConfig{
				Sections: []app.SectionConfig{
					{Name: "Section", Type: app.SectionTypeGit, URL: "https://doesnt-matter1.com"},
					{Name: "Section", Type: app.SectionTypeGit, URL: "https://doesnt-matter2.com"},
				},
			},
		}

		// act
		cli := app.NewCli(config)

		// assert
		assert.Equal(t, 5*time.Second, cli.Config.ScrapeInterval)
		assert.Len(t, cli.Scrapers, 2)
	})
}

func TestCli_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	scraperMock := mocks.NewMockScraper(ctrl)
	loggerMock := mocks.NewMockLogger(ctrl)

	t.Run("should periodically execute scaper", func(t *testing.T) {
		// assign
		ctx, cancel := context.WithCancel(context.Background())
		cli := app.Cli{
			Config: app.Config{
				ScrapeInterval: 10 * time.Millisecond,
			},
			Scrapers: []app.Scraper{scraperMock},
			Logger:   loggerMock,
		}

		scraperMock.EXPECT().
			Scrape(ctx).
			Return([]byte{}, nil).
			Times(1)
		scraperMock.EXPECT().
			Scrape(ctx).
			Do(func(_ context.Context) { cancel() }).
			Return([]byte{}, nil).
			Times(1)

		// act
		err := cli.Run(ctx)

		// assert
		assert.NoError(t, err)
	})

	t.Run("should log warning if scraping fails but continue", func(t *testing.T) {
		// assign
		ctx, cancel := context.WithCancel(context.Background())
		cli := app.Cli{
			Config: app.Config{
				ScrapeInterval: 10 * time.Millisecond,
			},
			Scrapers: []app.Scraper{scraperMock},
			Logger:   loggerMock,
		}

		loggerMock.EXPECT().
			Warn("scraping error: %s", errors.New("scrape error")).
			Times(1)

		scraperMock.EXPECT().
			Scrape(ctx).
			Return(nil, errors.New("scrape error")).
			Times(1)
		scraperMock.EXPECT().
			Scrape(ctx).
			Do(func(_ context.Context) { cancel() }).
			Return([]byte{}, nil).
			Times(1)

		// act
		err := cli.Run(ctx)

		// assert
		assert.NoError(t, err)
	})
}
