package faq

import (
	"app/internal/testutil"
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
)

// SetupTestContainer inicializa o banco para o módulo FAQ usando o utilitário global
func SetupTestContainer(ctx context.Context, t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Descobre o caminho do schema.sql deste módulo
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	schemaPath := filepath.Join(dir, "schema.sql")

	return testutil.SetupPostgresContainer(ctx, t, schemaPath)
}
