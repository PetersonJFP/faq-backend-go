package faq

import (
	"app/faq/db"
	"app/internal/web" // Importando o novo utilitário
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type App struct {
	Queries *db.Queries
}

func NewApp(conn *sql.DB) *App {
	return &App{Queries: db.New(conn)}
}

func (a *App) List(w http.ResponseWriter, r *http.Request) {
	faqs, err := a.Queries.ListFAQs(r.Context())
	if err != nil {
		web.Error(w, http.StatusInternalServerError, "Erro ao buscar FAQs")
		return
	}
	web.JSON(w, http.StatusOK, faqs)
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	var f db.CreateFAQParams
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		web.Error(w, http.StatusBadRequest, "Dados inválidos")
		return
	}

	faq, err := a.Queries.CreateFAQ(r.Context(), f)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, "Erro ao criar FAQ")
		return
	}

	web.JSON(w, http.StatusCreated, faq)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var f db.UpdateFAQParams
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		web.Error(w, http.StatusBadRequest, "Dados inválidos")
		return
	}
	f.ID = int32(id)

	faq, err := a.Queries.UpdateFAQ(r.Context(), f)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, "Erro ao atualizar FAQ")
		return
	}

	web.JSON(w, http.StatusOK, faq)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	if err := a.Queries.DeleteFAQ(r.Context(), int32(id)); err != nil {
		web.Error(w, http.StatusInternalServerError, "Erro ao deletar FAQ")
		return
	}

	web.NoContent(w)
}
