package faq

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes agora recebe o authMiddleware para proteger rotas específicas internamente
func (a *App) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/api/faqs", func(router chi.Router) {
		// Rotas Públicas
		router.Get("/", a.List)

		// Rotas Protegidas
		router.Group(func(protected chi.Router) {
			protected.Use(authMiddleware)
			protected.Post("/", a.Create)
			protected.Put("/{id}", a.Update)
			protected.Delete("/{id}", a.Delete)
		})
	})
}
