package dtos

import "github.com/mosmo1212312121/hexagonal_practice_go/internal/core/domain"

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}
}
