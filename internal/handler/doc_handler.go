package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/flohansen/documenter/internal/app"
	"github.com/flohansen/documenter/web"
)

type DocRepo interface {
	GetDocumentationNames(ctx context.Context) ([]string, error)
	GetDocumentationByName(ctx context.Context, name string) ([]byte, error)
}

type Renderer interface {
	Render(md []byte, templateName string) ([]byte, error)
}

type DocHandler struct {
	tmpl     *template.Template
	mux      *http.ServeMux
	repo     DocRepo
	renderer Renderer
	logger   app.Logger
}

func NewDocHandler(service DocRepo, renderer Renderer, logger app.Logger) *DocHandler {
	tmpl, err := template.ParseFS(web.Templates, "templates/index.gohtml")
	if err != nil {
		logger.Error("parsing templates failed", "error", err)
		panic(err)
	}

	h := &DocHandler{
		tmpl:     tmpl,
		mux:      http.NewServeMux(),
		repo:     service,
		renderer: renderer,
		logger:   logger,
	}

	h.mux.HandleFunc("GET /", h.GetRoot)
	h.mux.HandleFunc("GET /sections/{name}", h.GetSection)
	return h
}

type IndexModel struct {
	Navigation NavigationModel
	Content    template.HTML
}

type NavigationModel struct {
	Items []NavigationItemModel
}

type NavigationItemModel struct {
	Name string
}

func (h *DocHandler) GetRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	names, err := h.repo.GetDocumentationNames(ctx)
	if err != nil {
		h.logger.Error("retrieve documentations failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/sections/%s", names[0]))
	w.WriteHeader(http.StatusMovedPermanently)
}

func (h *DocHandler) GetSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PathValue("name")

	doc, err := h.repo.GetDocumentationByName(ctx, name)
	if err != nil {
		h.logger.Error("retrieve documentation failed", "error", err, "name", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	names, err := h.repo.GetDocumentationNames(ctx)
	if err != nil {
		h.logger.Error("retrieve documentations failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var model IndexModel
	for _, name := range names {
		model.Navigation.Items = append(model.Navigation.Items, NavigationItemModel{
			Name: name,
		})
	}

	content, err := h.renderer.Render(doc, "index")
	if err != nil {
		h.logger.Error("rendering failed", "error", err, "name", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	model.Content = template.HTML(content)

	w.Header().Add("Content-Type", "text/html")
	if err := h.tmpl.ExecuteTemplate(w, "index", model); err != nil {
		h.logger.Error("execute template failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DocHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
