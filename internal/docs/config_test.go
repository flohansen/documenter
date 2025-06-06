package docs_test

import (
	"strings"
	"testing"

	"github.com/flohansen/documenter/internal/docs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig_UnmarshalYAML(t *testing.T) {
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
		var config docs.Config
		err := yaml.Unmarshal(b, &config)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, docs.Config{
			Docs: docs.DocsConfig{
				Sections: []docs.SectionConfig{
					{Name: "Test", Type: docs.SectionTypeGit, URL: "https://some.url.com/repo"},
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
		var config docs.Config
		err := yaml.Unmarshal(b, &config)

		// assert
		assert.Error(t, err)
	})
}
