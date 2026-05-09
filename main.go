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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Ficheiro .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	// 2. Banco de Dados
	dbConn := setupDatabase()
	defer dbConn.Close()

	// 3. Roteador e Middlewares Globais
	r := chi.NewRouter()
	setupMiddlewares(r)

	// 4. Inicialização dos Apps (Injeção de Dependência)
	faqApp := faq.NewApp(dbConn)
	userApp := users.NewApp(dbConn)

	// Rotas de Usuários
	userApp.RegisterRoutes(r)

	// Rotas de FAQ (Passando o Middleware de Usuários para proteção)
	faqApp.RegisterRoutes(r, users.AuthMiddleware)

	// 5. Partida
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Servidor Go rodando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupDatabase() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Não foi possível alcançar o banco de dados: %v\n", err)
	}

	return db
}

func setupMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)    // Loga cada requisição no terminal
	r.Use(middleware.Recoverer) // Evita que o servidor caia por pânicos no código

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
}
