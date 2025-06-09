package renderer

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownRenderer struct {
	renderer *html.Renderer
}

func NewMarkdownRenderer() *MarkdownRenderer {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}

	return &MarkdownRenderer{
		renderer: html.NewRenderer(opts),
	}
}

func (r *MarkdownRenderer) Render(md []byte, templateName string) ([]byte, error) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	return markdown.Render(doc, r.renderer), nil
}
