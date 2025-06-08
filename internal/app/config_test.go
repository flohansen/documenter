package app_test

import (
	"strings"
	"testing"

	"github.com/flohansen/documenter/internal/app"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig_UnmarshalYAML(t *testing.T) {
	t.Run("section type", func(t *testing.T) {
		t.Run("should unmarshal git section type", func(t *testing.T) {
			// assign
			b := []byte(strings.Join([]string{
				"docs:",
				"  sections:",
				"    - name: Test",
				"      type: git",
				"      url: https://some.url.com/repo",
			}, "\n"))

			// act
			var config app.Config
			err := yaml.Unmarshal(b, &config)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, app.Config{
				Docs: app.DocsConfig{
					Sections: []app.SectionConfig{
						{Name: "Test", Type: app.SectionTypeGit, URL: "https://some.url.com/repo"},
					},
				},
			}, config)
		})

		t.Run("should return error when section type is unknown", func(t *testing.T) {
			// assign
			b := []byte(strings.Join([]string{
				"docs:",
				"  sections:",
				"    - name: Test",
				"      type: anything",
				"      url: https://some.url.com/repo",
			}, "\n"))

			// act
			var config app.Config
			err := yaml.Unmarshal(b, &config)

			// assert
			assert.Error(t, err)
		})
	})

	t.Run("logging format", func(t *testing.T) {
		t.Run("should unmarshal json logging format", func(t *testing.T) {
			// assign
			b := []byte(strings.Join([]string{
				"logging:",
				"  format: json",
			}, "\n"))

			// act
			var config app.Config
			err := yaml.Unmarshal(b, &config)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, app.Config{
				Logging: app.LoggingConfig{
					Format: app.LoggingFormatJSON,
				},
			}, config)
		})

		t.Run("should unmarshal text logging format", func(t *testing.T) {
			// assign
			b := []byte(strings.Join([]string{
				"logging:",
				"  format: text",
			}, "\n"))

			// act
			var config app.Config
			err := yaml.Unmarshal(b, &config)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, app.Config{
				Logging: app.LoggingConfig{
					Format: app.LoggingFormatText,
				},
			}, config)
		})

		t.Run("should return error when section type is unknown", func(t *testing.T) {
			// assign
			b := []byte(strings.Join([]string{
				"logging:",
				"  format: anything wrong",
			}, "\n"))

			// act
			var config app.Config
			err := yaml.Unmarshal(b, &config)

			// assert
			assert.Error(t, err)
		})
	})
}
