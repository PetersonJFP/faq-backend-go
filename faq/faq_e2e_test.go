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
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	// Setup do roteador real
	r := chi.NewRouter()
	faqApp := NewApp(dbConn)
	faqApp.RegisterRoutes(r, users.AuthMiddleware)

	t.Run("Deve criar FAQ com sucesso passando por todo o stack", func(t *testing.T) {
		token, _ := users.GenerateToken(1)

		payload := map[string]interface{}{
			"question":   "O que é um teste E2E?",
			"answer":     "É um teste que valida o fluxo completo.",
			"is_premium": false,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/faqs", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})
}
