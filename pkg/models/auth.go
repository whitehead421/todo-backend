package models

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6,max=32,nefield=Name,nefield=Email"`
	Confirm  string `json:"confirm" validate:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Id        uint64 `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId uint64 `json:"user_id"`
}

type AuthorizeResponse struct {
	Message string `json:"message"`
	UserId  uint64 `json:"user_id"`
}
