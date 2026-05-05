package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"faq-backend/faq" // Importa o nosso novo app

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Conexão com o Banco de Dados
	connStr := "user=root password=secret dbname=faq_db sslmode=disable host=localhost port=5432"
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

	// 3. Middlewares (CORS liberado para integração com o app mobile e web)
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
	faqApp.RegisterRoutes(r) // Registra as rotas em /api/faqs

	// 5. Inicia o Servidor
	fmt.Println("🚀 Servidor Go rodando na porta 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erro ao subir servidor: %v\n", err)
	}
}
