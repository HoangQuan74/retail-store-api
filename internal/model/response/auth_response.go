package response

import "time"

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int64     `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthAPIResponse struct {
	Status  int          `json:"status" example:"200"`
	Message string       `json:"message" example:"Success"`
	Data    AuthResponse `json:"data"`
}
