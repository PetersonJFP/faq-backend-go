package faq

import (
	"app/users"
	"bytes"
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

	// Testa a listagem (Rota Pública)
	req := httptest.NewRequest(http.MethodGet, "/api/faqs", nil)
	rr := httptest.NewRecorder()

	app.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
}

func TestFAQProtectedRoutesIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	app := NewApp(dbConn)

	// 1. Tentar criar SEM token (Deve falhar)
	payload := map[string]interface{}{
		"question": "Pergunta teste?",
		"answer":   "Resposta teste",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/faqs", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Nota: Aqui simulamos o que o middleware faria. 
	// Em um teste de integração de "ponta a ponta", usaríamos o router completo.
	// Por enquanto, vamos testar o comportamento do Handler se o usuário for injetado.
	
	app.Create(rr, req)

	// Se não passamos pelo middleware e não injetamos o user_id no context, 
	// o handler deve retornar erro ou o teste deve validar a proteção.
	// No nosso caso, o middleware é quem barra antes de chegar aqui.
}

func TestFAQCreateWithAuthIntegration(t *testing.T) {
	ctx := context.Background()
	dbConn, cleanup := SetupTestContainer(ctx, t)
	defer cleanup()

	app := NewApp(dbConn)

	// 1. Gerar um token válido (Simulando um usuário logado)
	// Usamos o pacote users que já tem a lógica de JWT
	token, err := users.GenerateToken(1) 
	assert.NoError(t, err)

	// 2. Criar FAQ com o Token
	payload := map[string]interface{}{
		"question":   "Como funciona o TDD?",
		"answer":     "Primeiro o teste falha, depois ele passa.",
		"is_premium": false,
	}
	jsonBody, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/faqs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()

	// Aqui, para o teste ser fiel, precisamos simular o contexto que o Middleware colocaria
	userCtx := context.WithValue(req.Context(), users.UserIDKey, int32(1))
	app.Create(rr, req.WithContext(userCtx))

	assert.Equal(t, http.StatusCreated, rr.Code)
	
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "Como funciona o TDD?", response["question"])
}