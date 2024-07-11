package entities

type TodoRequest struct {
	Description string `json:"description" example:"Buy milk" validate:"min=6"`
}

type TodoResponse struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
