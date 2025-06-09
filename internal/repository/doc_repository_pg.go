package repository

import (
	"context"
	"fmt"

	"github.com/flohansen/documenter/internal/database"
	"github.com/flohansen/documenter/internal/domain"
)

type DocRepoPostgres struct {
	q *database.Queries
}

func NewDocRepoPostgres(db database.DBTX) *DocRepoPostgres {
	return &DocRepoPostgres{
		q: database.New(db),
	}
}

func (r *DocRepoPostgres) UpsertDocumentation(ctx context.Context, doc domain.Documentation) error {
	return r.q.UpsertDocumentation(ctx, database.UpsertDocumentationParams{
		Name:    doc.Name,
		Content: doc.Content,
	})
}

func (r *DocRepoPostgres) GetDocumentationNames(ctx context.Context) ([]string, error) {
	names, err := r.q.GetDocumentations(ctx)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return names, nil
}

func (r *DocRepoPostgres) GetDocumentationByName(ctx context.Context, name string) ([]byte, error) {
	content, err := r.q.GetDocumentationByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return content, nil
}
