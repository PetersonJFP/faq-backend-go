package users

// RegisterRequest representa o payload esperado para criar uma nova conta
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100" label:"Nome"`
	Email    string `json:"email" validate:"required,email" label:"E-mail"`
	Password string `json:"password" validate:"required,min=6,max=32" label:"Senha"`
}

// LoginRequest representa o payload esperado para autenticação
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" label:"E-mail"`
	Password string `json:"password" validate:"required" label:"Senha"`
}

// UserResponse define os dados públicos do utilizador que podem ser devolvidos à API
type UserResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
