package users

import (
	"github.com/go-chi/chi/v5"
)

func (a *App) RegisterRoutes(r chi.Router) {
	r.Route("/api/users", func(router chi.Router) {
		router.Post("/register", a.Register)
		router.Post("/login", a.Login)
	})
}
