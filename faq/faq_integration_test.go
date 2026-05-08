package faq

import (
	"app/faq/db"
	"app/users"
	"bytes" // Adicionado para transformar []byte em io.Reader
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFAQListIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	app := NewApp(dbConn)

	req := httptest.NewRequest(http.MethodGet, "/api/faqs", nil)
	rr := httptest.NewRecorder()

	app.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFAQCreateWithAuthIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	app := NewApp(dbConn)

	// Gerar token real via módulo de users
	token, err := users.GenerateToken(1)
	assert.NoError(t, err)

	payload := db.CreateFAQParams{
		Question:  "Como funciona o TDD?",
		Answer:    "Primeiro o teste falha, depois ele passa.",
		IsPremium: false,
	}
	jsonBody, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/faqs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Simular o contexto que o Middleware colocaria
	userCtx := context.WithValue(req.Context(), users.UserIDKey, int32(1))
	app.Create(rr, req.WithContext(userCtx))

	assert.Equal(t, http.StatusCreated, rr.Code)
}
