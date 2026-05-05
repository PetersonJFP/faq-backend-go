package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"app/faq"
	"app/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Ficheiro .env não encontrado.")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro banco: %v\n", err)
	}
	defer dbConn.Close()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	faqApp := faq.NewApp(dbConn)
	userApp := users.NewApp(dbConn)

	userApp.RegisterRoutes(r)

	r.Route("/api/faqs", func(router chi.Router) {
		router.Get("/", faqApp.List)

		router.Group(func(protected chi.Router) {
			protected.Use(users.AuthMiddleware)
			protected.Post("/", faqApp.Create)
			protected.Put("/{id}", faqApp.Update)
			protected.Delete("/{id}", faqApp.Delete)
		})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Servidor rodando em http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
