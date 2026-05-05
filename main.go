package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"faq-backend/internal/db" // Importando o super pacote que o SQLC gerou para nós!

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq" // Importa o driver do Postgres silenciosamente
)

func main() {
	// 1. Criar a string de conexão (com os dados do nosso Docker)
	connStr := "user=root password=secret dbname=faq_db sslmode=disable host=localhost port=5432"

	// 2. Abrir a conexão com o banco
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}
	defer dbConn.Close()

	// 3. Entregar a conexão para o SQLC fazer a mágica dele
	queries := db.New(dbConn)

	// 4. Configurar o Roteador
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	// 5. Nossa rota agora bate no banco real!
	r.Get("/api/faqs", func(w http.ResponseWriter, r *http.Request) {
		// Chama a função ListFAQs que o SQLC gerou automaticamente
		faqs, err := queries.ListFAQs(context.Background())
		if err != nil {
			http.Error(w, "Erro ao buscar FAQs", http.StatusInternalServerError)
			return
		}

		// Se o banco retornar nulo (vazio), mandamos um array vazio [] para o React Native não quebrar
		if faqs == nil {
			faqs = []db.Faq{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(faqs)
	})

	r.Post("/api/faqs", func(w http.ResponseWriter, r *http.Request) {
		// 1. Criamos uma estrutura temporária para ler o JSON que vem na requisição
		var body struct {
			Question  string `json:"question"`
			Answer    string `json:"answer"`
			IsPremium bool   `json:"is_premium"`
		}

		// 2. Lemos o corpo (body) da requisição e convertemos para a nossa estrutura
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// 3. Passamos os dados para a função que o SQLC gerou para inserir no banco
		novoFaq, err := queries.CreateFAQ(context.Background(), db.CreateFAQParams{
			Question:  body.Question,
			Answer:    body.Answer,
			IsPremium: body.IsPremium,
		})

		if err != nil {
			http.Error(w, "Erro ao salvar a FAQ no banco de dados", http.StatusInternalServerError)
			return
		}

		// 4. Retornamos Status 201 (Created) e os dados da nova FAQ gerada (com o ID do banco)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(novoFaq)
	})

	r.Delete("/api/faqs/{id}", func(w http.ResponseWriter, r *http.Request) {
		// 1. Pegar o parâmetro {id} da URL
		idParam := chi.URLParam(r, "id")

		// 2. Converter de texto para número inteiro
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID inválido. Deve ser um número.", http.StatusBadRequest)
			return
		}

		// 3. Executar a deleção no banco (convertendo para int32 que o PostgreSQL espera)
		err = queries.DeleteFAQ(context.Background(), int32(id))
		if err != nil {
			http.Error(w, "Erro ao deletar a FAQ do banco de dados", http.StatusInternalServerError)
			return
		}

		// 4. Padrão REST: Retornar 204 (No Content) indicando sucesso sem corpo na resposta
		w.WriteHeader(http.StatusNoContent)
	})

	r.Put("/api/faqs/{id}", func(w http.ResponseWriter, r *http.Request) {
		// 1. Pegar o ID da URL
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID inválido.", http.StatusBadRequest)
			return
		}

		// 2. Ler os novos dados do JSON
		var body struct {
			Question  string `json:"question"`
			Answer    string `json:"answer"`
			IsPremium bool   `json:"is_premium"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// 3. Executar a atualização no banco usando os parâmetros gerados pelo SQLC
		faqAtualizada, err := queries.UpdateFAQ(context.Background(), db.UpdateFAQParams{
			ID:        int32(id),
			Question:  body.Question,
			Answer:    body.Answer,
			IsPremium: body.IsPremium,
		})

		// Se der erro (ex: tentou atualizar um ID que não existe)
		if err != nil {
			http.Error(w, "Erro ao atualizar a FAQ", http.StatusInternalServerError)
			return
		}

		// 4. Retornar a FAQ atualizada
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(faqAtualizada)
	})

	fmt.Println("🚀 Servidor Go com Postgres rodando na porta 8080...")
	http.ListenAndServe(":8080", r)
}
