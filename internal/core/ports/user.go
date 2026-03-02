package ports

import (
	"github.com/google/uuid"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/domain"
)

type UserRepository interface {
	Save(user domain.User) error
	FindByEmail(email string) (domain.User, error)
	FindByID(id uuid.UUID) (domain.User, error)
}

type UserService interface {
	Register(user domain.User) error
	GetByID(id uuid.UUID) (domain.User, error)
}
