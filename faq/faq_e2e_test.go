package faq

import (
	"app/users"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestFAQ_E2E_Flow(t *testing.T) {
	ctx := context.Background()
	
	// 1. Setup do Banco Real (via Testcontainers)
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	// 2. Setup do Roteador Real (Exatamente como no main.go)
	r := chi.NewRouter()
	faqApp := NewApp(dbConn)
	
	// Registramos as rotas injetando o middleware real de outro módulo
	faqApp.RegisterRoutes(r, users.AuthMiddleware)

	t.Run("Deve rejeitar criação sem token (Middleware Check)", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"question": "?", "answer": "!"})
		req := httptest.NewRequest(http.MethodPost, "/api/faqs", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		// Enviamos para o ROTEADOR (r), testando o middleware de proteção
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Deve criar FAQ com sucesso passando por todo o stack", func(t *testing.T) {
		// Geramos um token real usando a lógica do módulo de usuários
		token, _ := users.GenerateToken(1)

		payload := map[string]interface{}{
			"question":   "O que é um teste E2E?",
			"answer":     "É um teste que valida o fluxo completo, da rota ao banco.",
			"is_premium": false,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/faqs", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()

		// A mágica acontece aqui: ServeHTTP processa middlewares, roteamento e handler
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "O que é um teste E2E?", resp["question"])
	})
}