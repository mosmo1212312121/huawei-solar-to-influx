package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/domain"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/ports"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user domain.User) error {
	_, err := s.repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already exists")
	}
	return s.repo.Save(user)
}

func (s *userService) GetByID(id uuid.UUID) (domain.User, error) {
	return s.repo.FindByID(id)
}
