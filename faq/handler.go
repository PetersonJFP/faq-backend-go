package faq

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// App encapsula as dependências deste módulo (ex: acesso ao banco de dados).
type App struct {
	Queries *Queries // Gerado pelo SQLC
}

// NewApp é o construtor do nosso módulo FAQ.
func NewApp(db *sql.DB) *App {
	return &App{
		Queries: New(db),
	}
}

// List retorna todas as FAQs do banco.
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	faqs, err := a.Queries.ListFAQs(context.Background())
	if err != nil {
		http.Error(w, "Erro ao buscar as FAQs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faqs)
}

// Create recebe um JSON e cria uma nova FAQ.
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		IsPremium bool   `json:"is_premium"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	novoFaq, err := a.Queries.CreateFAQ(context.Background(), CreateFAQParams{
		Question:  body.Question,
		Answer:    body.Answer,
		IsPremium: body.IsPremium,
	})
	if err != nil {
		http.Error(w, "Erro ao salvar a FAQ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novoFaq)
}

// Update modifica uma FAQ existente usando o ID da URL.
func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	var body struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		IsPremium bool   `json:"is_premium"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	faqAtualizada, err := a.Queries.UpdateFAQ(context.Background(), UpdateFAQParams{
		ID:        int32(id),
		Question:  body.Question,
		Answer:    body.Answer,
		IsPremium: body.IsPremium,
	})
	if err != nil {
		http.Error(w, "Erro ao atualizar a FAQ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faqAtualizada)
}

// Delete remove uma FAQ do banco pelo ID da URL.
func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	err = a.Queries.DeleteFAQ(context.Background(), int32(id))
	if err != nil {
		http.Error(w, "Erro ao deletar a FAQ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
