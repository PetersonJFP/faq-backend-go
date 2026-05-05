package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"faq-backend/faq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 0. Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de sistema.")
	}

	// 1. Construção da Conexão com o Banco de Dados via Env
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v\n", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Não foi possível alcançar o banco de dados: %v\n", err)
	}

	// 2. Configuração do Roteador Principal
	r := chi.NewRouter()

	// 3. Middlewares (CORS) - Mantido exatamente como você enviou
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 4. Inicialização dos Módulos (Apps)
	faqApp := faq.NewApp(dbConn)
	faqApp.RegisterRoutes(r)

	// 5. Inicia o Servidor usando a porta do .env
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Valor padrão caso não esteja no .env
	}

	fmt.Printf("🚀 Servidor Go rodando na porta %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Erro ao subir servidor: %v\n", err)
	}
}
