package testutil

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type silentLogger struct{}

func (l silentLogger) Printf(format string, v ...interface{}) {}

// SetupPostgresContainer inicia um container único para o teste do módulo e aplica o schema
func SetupPostgresContainer(ctx context.Context, t *testing.T, schemaPath string) (*sql.DB, func()) {
	t.Helper()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		// Aplicando o silenciador globalmente aqui
		testcontainers.WithLogger(silentLogger{}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("falha ao iniciar container: %s", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("falha ao obter connection string: %s", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("falha ao conectar no banco: %s", err)
	}

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("falha ao ler schema em %s: %s", schemaPath, err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("falha ao aplicar schema: %s", err)
	}

	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}
