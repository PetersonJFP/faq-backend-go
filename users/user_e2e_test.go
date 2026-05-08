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
)

func TestUser_E2E_AuthFlow(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	r := chi.NewRouter()
	userApp := NewApp(dbConn)
	userApp.RegisterRoutes(r)

	t.Run("Fluxo: Registrar -> Login -> Obter Token", func(t *testing.T) {
		// 1. Registro
		regPayload := map[string]string{
			"name":     "User E2E",
			"email":    "e2e@test.com",
			"password": "password123",
		}
		regBody, _ := json.Marshal(regPayload)
		
		reqReg := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(regBody))
		rrReg := httptest.NewRecorder()
		r.ServeHTTP(rrReg, reqReg)

		assert.Equal(t, http.StatusCreated, rrReg.Code)

		// 2. Login
		loginPayload := map[string]string{
			"email":    "e2e@test.com",
			"password": "password123",
		}
		loginBody, _ := json.Marshal(loginPayload)
		
		reqLogin := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(loginBody))
		rrLogin := httptest.NewRecorder()
		r.ServeHTTP(rrLogin, reqLogin)

		assert.Equal(t, http.StatusOK, rrLogin.Code)

		var resp map[string]string
		json.Unmarshal(rrLogin.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["token"], "O login E2E deve retornar um token válido")
	})
}