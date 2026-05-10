package users

import (
	"app/internal/web"
	"app/users/db"
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Queries *db.Queries
}

func NewApp(conn *sql.DB) *App {
	return &App{Queries: db.New(conn)}
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := web.ReadJSON(r, &req); err != nil {
		web.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, "Erro ao processar segurança")
		return
	}

	user, err := a.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		web.Error(w, http.StatusConflict, "Este e-mail já está em uso")
		return
	}

	web.JSON(w, http.StatusCreated, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := web.ReadJSON(r, &req); err != nil {
		web.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		web.Error(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		web.Error(w, http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, "Erro ao gerar acesso")
		return
	}

	web.JSON(w, http.StatusOK, map[string]string{
		"token": token,
		"name":  user.Name,
	})
}
