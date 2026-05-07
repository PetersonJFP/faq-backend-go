package users

import (
	"app/internal/testutil"
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
)

func SetupTestContainer(ctx context.Context, t *testing.T) (*sql.DB, func()) {
	// Descobre o caminho do schema.sql deste módulo
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	schemaPath := filepath.Join(dir, "schema.sql")

	// Chama o utilitário global que já é silencioso
	return testutil.SetupPostgresContainer(ctx, t, schemaPath)
}
