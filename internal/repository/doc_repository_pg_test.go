package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/flohansen/documenter/internal/database"
	"github.com/flohansen/documenter/internal/domain"
	"github.com/flohansen/documenter/internal/repository"
	"github.com/flohansen/documenter/test/testhelpers"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestDocRepoPostgres_Integration(t *testing.T) {
	container := testhelpers.StartPostgresContainer(t, testhelpers.WithMigration("../../sql/migrations"))

	db, err := pgx.Connect(context.Background(), container.Dsn())
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	getDoc := getDoc(t, db)
	insertDoc := insertDoc(t, db)
	beforeEach := beforeEach(t, db)

	repo := repository.NewDocRepoPostgres(db)

	t.Run("UpsertDocumentation", func(t *testing.T) {
		t.Run("should insert new documentation", func(t *testing.T) {
			beforeEach()

			// assign
			// act
			err := repo.UpsertDocumentation(context.Background(), domain.Documentation{
				Name:    "name",
				Content: []byte("content"),
			})

			// assert
			assert.NoError(t, err)
			assert.Equal(t, database.Documentation{
				Name:    "name",
				Content: []byte("content"),
			}, getDoc("name"))
		})

		t.Run("should update existing documentation", func(t *testing.T) {
			beforeEach()

			// assign
			insertDoc(database.Documentation{
				ID:      0,
				Name:    "name",
				Content: []byte("change me"),
			})

			// act
			err := repo.UpsertDocumentation(context.Background(), domain.Documentation{
				Name:    "name",
				Content: []byte("content"),
			})

			// assert
			assert.NoError(t, err)
			assert.Equal(t, database.Documentation{
				ID:      0,
				Name:    "name",
				Content: []byte("content"),
			}, getDoc("name"))
		})
	})
}

func beforeEach(t *testing.T, db *pgx.Conn) func() {
	return func() {
		if _, err := db.Exec(context.Background(), "DELETE FROM documentations"); err != nil {
			t.Fatal(err)
		}
	}
}

func insertDoc(t *testing.T, db *pgx.Conn) func(doc database.Documentation) {
	return func(doc database.Documentation) {
		if _, err := db.Exec(context.Background(),
			"INSERT INTO documentations (id, name, content) VALUES ($1, $2, $3)", doc.ID, doc.Name, doc.Content,
		); err != nil {
			t.Fatal(err)
		}
	}
}

func getDoc(t *testing.T, db *pgx.Conn) func(name string) database.Documentation {
	return func(name string) database.Documentation {
		row := db.QueryRow(context.Background(),
			"SELECT name, content FROM documentations WHERE name = $1 LIMIT 1", name)

		var doc database.Documentation
		if err := row.Scan(
			&doc.Name,
			&doc.Content,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return database.Documentation{}
			}

			t.Fatal(err)
		}

		return doc
	}
}
