package users

import (
	"app/users/db"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserEmailUniquenessIntegration(t *testing.T) {
	ctx := context.Background()

	// 1. Setup do Banco via Testcontainers
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	queries := db.New(dbConn)
	email := "unique@test.com"

	// 2. Primeira Inserção: Deve funcionar
	params := db.CreateUserParams{
		Name:         "Primeiro User",
		Email:        email,
		PasswordHash: "hash_test",
	}

	user, err := queries.CreateUser(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NotZero(t, user.ID)

	// 3. Segunda Inserção com o mesmo Email: Deve falhar
	params.Name = "Segundo User"
	_, err2 := queries.CreateUser(ctx, params)

	assert.Error(t, err2, "O banco deveria retornar um erro de violação de UNIQUE")
	assert.Contains(t, err2.Error(), "unique constraint", "A mensagem de erro deve mencionar a restrição UNIQUE")
}
