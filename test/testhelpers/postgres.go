package testhelpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testPgUser     = "test"
	testPgPassword = "test"
	testPgDatabase = "test"
)

type PostgresContainer struct {
	t         *testing.T
	container testcontainers.Container
}

func (c *PostgresContainer) Dsn() string {
	host, err := c.container.Host(context.Background())
	if err != nil {
		c.t.Fatal(err)
	}

	port, err := c.container.MappedPort(context.Background(), nat.Port("5432/tcp"))
	if err != nil {
		c.t.Fatal(err)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		testPgUser, testPgPassword, host, port.Int(), testPgDatabase)
}

func StartPostgresContainer(t *testing.T, opts ...PostgresContainerOption) *PostgresContainer {
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     testPgUser,
				"POSTGRES_PASSWORD": testPgPassword,
				"POSTGRES_DB":       testPgDatabase,
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
			},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort(nat.Port("5432/tcp")),
			),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		container.Terminate(context.Background())
	})

	c := &PostgresContainer{
		t:         t,
		container: container,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

type PostgresContainerOption func(*PostgresContainer)

func WithMigration(dir string) PostgresContainerOption {
	return func(c *PostgresContainer) {
		sourceURL := fmt.Sprintf("file://%s", dir)

		m, err := migrate.New(sourceURL, c.Dsn())
		if err != nil {
			c.t.Fatal(err)
		}

		if err := m.Up(); err != nil {
			c.t.Fatal(err)
		}
	}
}
