package users

import (
	"app/internal/web"
	"app/users/db"
	"database/sql"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Queries *db.Queries
}

func NewApp(conn *sql.DB) *App {
	return &App{Queries: db.New(conn)}
}

// UserResponse define o que enviamos para o frontend (Segurança)
type UserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest define os dados necessários para cadastro
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest define as credenciais de acesso (Necessário para os testes)
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := web.ReadJSON(r, &req); err != nil {
		web.Error(w, http.StatusBadRequest, "Dados inválidos")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user, err := a.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		web.Error(w, http.StatusConflict, "Este email já está em uso")
		return
	}

	web.JSON(w, http.StatusCreated, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	})
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := web.ReadJSON(r, &req); err != nil {
		web.Error(w, http.StatusBadRequest, "Dados inválidos")
		return
	}

	user, err := a.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		web.Error(w, http.StatusUnauthorized, "Email ou senha incorretos")
		return
	}

	token, _ := GenerateToken(user.ID)

	web.JSON(w, http.StatusOK, map[string]string{
		"token": token,
		"name":  user.Name,
	})
}
