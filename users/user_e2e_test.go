package users

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_E2E_Scenarios(t *testing.T) {
	ctx := context.Background()
	// Setup do banco real isolado via Testcontainers
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	// Configuração do Roteador Real (igual ao main.go)
	r := chi.NewRouter()
	userApp := NewApp(dbConn)
	userApp.RegisterRoutes(r)

	t.Run("Registro de Novo Usuário (Caminho Feliz)", func(t *testing.T) {
		payload := map[string]string{
			"name":     "João Silva",
			"email":    "joao@exemplo.com",
			"password": "senha_segura_123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, payload["email"], resp.Email)
		assert.NotZero(t, resp.ID)
		assert.NotContains(t, w.Body.String(), "password_hash", "A resposta nunca deve vazar o hash da senha")
	})

	t.Run("Erro ao registrar email já existente", func(t *testing.T) {
		// Tentando registrar o mesmo e-mail do teste anterior
		payload := map[string]string{
			"name":     "Outro João",
			"email":    "joao@exemplo.com",
			"password": "outra_senha_123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Este e-mail já está em uso")
	})

	t.Run("Erro de validação (Senha muito curta)", func(t *testing.T) {
		payload := map[string]string{
			"name":     "Curto",
			"email":    "curto@exemplo.com",
			"password": "123", // Menor que o min=6 definido na struct
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "no mínimo 6 caracteres")
	})

	t.Run("Login com sucesso", func(t *testing.T) {
		payload := map[string]string{
			"email":    "joao@exemplo.com",
			"password": "senha_segura_123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["token"], "O login deve retornar um token JWT válido")
		assert.Equal(t, "João Silva", resp["name"])
	})

	t.Run("Erro no login (Senha incorreta)", func(t *testing.T) {
		payload := map[string]string{
			"email":    "joao@exemplo.com",
			"password": "senha_errada",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Credenciais inválidas")
	})

	t.Run("Erro no login (Usuário inexistente)", func(t *testing.T) {
		payload := map[string]string{
			"email":    "fantasma@exemplo.com",
			"password": "qualquer_senha",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
