package users

import (
	"app/users/db"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Queries *db.Queries
}

func NewApp(conn *sql.DB) *App {
	return &App{
		Queries: db.New(conn),
	}
}

// Structs para requests
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register cria um novo usuário com senha criptografada
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// 1. Gerar Hash da Senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao processar senha", http.StatusInternalServerError)
		return
	}

	// 2. Salvar no Banco
	user, err := a.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		http.Error(w, "Erro ao criar usuário (email já existe?)", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login verifica as credenciais
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// 1. Buscar usuário pelo email
	user, err := a.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 2. Comparar senha enviada com o hash do banco
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// Por enquanto, apenas confirmamos o sucesso.
	// O próximo passo será gerar um Token JWT aqui.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login realizado com sucesso",
		"user_id": string(rune(user.ID)),
		"name":    user.Name,
	})
}
