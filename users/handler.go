package users

import (
	"app/users/db"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// App injeta as dependências do módulo de utilizadores
type App struct {
	Queries *db.Queries
}

// NewApp cria uma nova instância do módulo
func NewApp(conn *sql.DB) *App {
	return &App{
		Queries: db.New(conn),
	}
}

// UserResponse define quais campos do utilizador podem ser enviados ao frontend
type UserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao processar senha", http.StatusInternalServerError)
		return
	}

	user, err := a.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		http.Error(w, "Erro ao criar utilizador", http.StatusConflict)
		return
	}

	// Mapeia para a resposta segura (sem hash)
	response := UserResponse{
		ID:        user.ID, // Usando ID em maiúsculo
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// 1. Procurar utilizador
	user, err := a.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Utilizador ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 2. Validar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Utilizador ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 3. Gerar Token
	token, err := GenerateToken(user.ID) // Usando ID em maiúsculo
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"name":  user.Name,
	})
}
