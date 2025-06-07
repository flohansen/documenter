package scraper_test

import (
	"context"
	"testing"

	"github.com/flohansen/documenter/internal/scraper"
	"github.com/stretchr/testify/assert"
)

func TestGitScraper_Scrape(t *testing.T) {
	t.Run("should return README.md content", func(t *testing.T) {
		// assign
		scpr := scraper.NewGitScraper("https://github.com/flohansen/documenter")

		// act
		md, err := scpr.Scrape(context.Background())

		// assert
		assert.NoError(t, err)
		assert.Greater(t, len(md), 0)
	})
}
