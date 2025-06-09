package main

import (
	"flag"
	"log"

	"github.com/flohansen/documenter/internal/app"
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

	pool, err := pgxpool.New(ctx, flags.Database)
	if err != nil {
		log.Fatalf("could not create db pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	repo := repository.NewDocRepoPostgres(pool)

	svr := app.NewServer(repo)
	if err := svr.Run(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
