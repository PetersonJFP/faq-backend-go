package users

import (
	"app/users/db"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandlerIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	app := NewApp(dbConn)

	payload := RegisterRequest{
		Name:     "Novo Utilizador",
		Email:    "handler@test.com",
		Password: "senha_secreta",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "O status deve ser 201 Created")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "handler@test.com", response["email"])
	assert.NotContains(t, rr.Body.String(), "password_hash", "O hash da senha nunca deve vazar no JSON")
}

func TestLoginHandlerIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	app := NewApp(dbConn)

	// 1. Criar utilizador diretamente no banco de teste
	password := "senha_secreta"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, err := app.Queries.CreateUser(ctx, db.CreateUserParams{
		Name:         "Utilizador Teste",
		Email:        "login@test.com",
		PasswordHash: string(hashedPassword),
	})
	assert.NoError(t, err)

	// 2. Tentar fazer login via Handler
	loginPayload := LoginRequest{
		Email:    "login@test.com",
		Password: password,
	}
	body, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.Login(rr, req)

	// CORREÇÃO: Voltando para 200 OK
	assert.Equal(t, http.StatusOK, rr.Code, "O login deve retornar 200 OK")

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"], "Deve retornar um token JWT")
}
