package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/flohansen/documenter/internal/app"
	"github.com/flohansen/documenter/internal/handler"
	"github.com/flohansen/documenter/internal/renderer"
	"github.com/flohansen/documenter/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type flags struct {
	Database string
}

func main() {
	var flags flags
	flag.StringVar(&flags.Database, "database", "postgresql://localhost:5432/postgres", "The connection string used to connect to the PostgreSQL database")
	flag.Parse()

	ctx := app.SignalContext()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pool, err := pgxpool.New(ctx, flags.Database)
	if err != nil {
		logger.Error("could not create db pool", "error", err)
		os.Exit(1)
	}
	if err := pool.Ping(ctx); err != nil {
		logger.Error("could not ping database", "error", err)
		os.Exit(1)
	}

	docRepo := repository.NewDocRepoPostgres(pool)
	docRenderer := renderer.NewMarkdownRenderer()
	docHandler := handler.NewDocHandler(docRepo, docRenderer, logger)

	svr := app.NewServer(docHandler)
	if err := svr.Run(ctx); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
