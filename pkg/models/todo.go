package models

type Status string

const (
	Pending    Status = "pending"
	InProgress Status = "in_progress"
	Completed  Status = "completed"
)

type TodoRequest struct {
	Description string `json:"description" example:"Buy milk" validate:"min=6"`
}

type TodoUpdateRequest struct {
	Description string `json:"description" example:"Buy milk" validate:"min=6"`
	Status      Status `json:"status" example:"pending" validate:"required"`
}

type TodoResponse struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
