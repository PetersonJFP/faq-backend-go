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
	// 0. Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Ficheiro .env não encontrado.")
	}

	// 1. Construção da string de conexão
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	// 2. Conexão com o banco de dados (Corrigido para sql.Open)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v\n", err)
	}
	defer dbConn.Close()

	// Verificar se a conexão está ativa
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Não foi possível alcançar o banco de dados: %v\n", err)
	}

	// 3. Configuração do Roteador Principal
	r := chi.NewRouter()

	// 4. Middlewares de Infraestrutura
	r.Use(middleware.Logger)    // Loga todas as requisições no terminal
	r.Use(middleware.Recoverer) // Evita que o servidor caia por causa de um Panic (ex: Nil pointer)

	// Middlewares de Segurança (CORS)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// 5. Inicialização dos Apps (Módulos)
	faqApp := faq.NewApp(dbConn)
	userApp := users.NewApp(dbConn)

	// 6. Registro de Rotas
	// Rotas de utilizadores (Públicas: Registro/Login)
	userApp.RegisterRoutes(r)

	// Rotas de FAQ (Passando o Middleware de Autenticação para proteção interna)
	faqApp.RegisterRoutes(r, users.AuthMiddleware)

	// 7. Inicia o Servidor
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Servidor Go rodando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
