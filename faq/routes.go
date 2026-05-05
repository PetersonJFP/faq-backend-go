package faq

import "github.com/go-chi/chi/v5"

// RegisterRoutes cria o agrupamento de rotas e liga aos handlers.
func (a *App) RegisterRoutes(r chi.Router) {
	// Cria um sub-roteador para agrupar tudo que começa com /api/faqs
	r.Route("/api/faqs", func(router chi.Router) {
		router.Get("/", a.List)
		router.Post("/", a.Create)
		router.Put("/{id}", a.Update)
		router.Delete("/{id}", a.Delete)
	})
}
