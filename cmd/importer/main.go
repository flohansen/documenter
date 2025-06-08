package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/flohansen/documenter/internal/app"
	"github.com/flohansen/documenter/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

type flags struct {
	ConfigPath string
	Database   string
}

func main() {
	var flags flags
	flag.StringVar(&flags.ConfigPath, "config", "documenter.config.yaml", "The path to the configuration file")
	flag.StringVar(&flags.Database, "database", "postgresql://localhost:5432/postgres", "The connection string used to connect to the PostgreSQL database")
	flag.Parse()

	config, err := readConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalf("could not read config: %v", err)
	}

	ctx := app.SignalContext()
	pool, err := pgxpool.New(ctx, flags.Database)
	if err != nil {
		log.Fatalf("could not create db pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	repo := repository.NewDocRepoPostgres(pool)

	cli := app.NewImporter(repo, config)
	if err := cli.Run(ctx); err != nil {
		log.Fatalf("cli error: %v", err)
	}
}

func readConfig(name string) (app.Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return app.Config{}, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	var config app.Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return app.Config{}, fmt.Errorf("json decode error: %w", err)
	}

	return config, nil
}
