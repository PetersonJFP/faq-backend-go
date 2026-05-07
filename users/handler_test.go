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

	// 1. Setup do Banco via Testcontainers
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	// 2. Inicializar o App (Injeção de Dependência)
	app := NewApp(dbConn)

	// 3. Preparar o payload JSON
	payload := RegisterRequest{
		Name:     "Novo Utilizador",
		Email:    "handler@test.com",
		Password: "senha_secreta",
	}
	body, _ := json.Marshal(payload)

	// 4. Criar Request e Recorder (Simulador de Response)
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// 5. Executar a chamada diretamente no Handler
	app.Register(rr, req)

	// 6. Asserções
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

	// 1. Setup: Registrar um utilizador diretamente no banco para poder logar
	password := "senha_secreta"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, err := app.Queries.CreateUser(ctx, db.CreateUserParams{
		Name:         "Utilizador Teste",
		Email:        "login@test.com",
		PasswordHash: string(hashedPassword),
	})
	assert.NoError(t, err)

	// 2. Preparar payload de Login
	loginPayload := LoginRequest{
		Email:    "login@test.com",
		Password: password,
	}
	body, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// 3. Executar Login
	app.Login(rr, req)

	// 4. Asserções
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar se recebemos o Token JWT e o nome do utilizador
	assert.NotEmpty(t, response["token"], "Deve retornar um token JWT")
	assert.Equal(t, "Utilizador Teste", response["name"])
}
