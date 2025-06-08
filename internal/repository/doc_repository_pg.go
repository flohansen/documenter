package repository

import (
	"context"

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
