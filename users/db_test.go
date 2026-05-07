package users

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// tLogger é um wrapper que redireciona logs do Docker para o sistema de testes do Go.
// Isto faz com que os logs só apareçam se o teste FALHAR.
type tLogger struct {
	t *testing.T
}

func (l tLogger) Printf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}

func SetupTestContainer(ctx context.Context, t *testing.T) (*sql.DB, func()) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		// Usamos o nosso logger customizado aqui
		testcontainers.WithLogger(tLogger{t: t}),
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

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	schemaPath := filepath.Join(dir, "schema.sql")

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("falha ao ler schema.sql em %s: %s", schemaPath, err)
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
