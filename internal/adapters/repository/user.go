package repository

import (
	"github.com/google/uuid"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/domain"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserModel struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string    `gorm:"not null"`
	Email    string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uuid.UUID) (domain.User, error) {
	var userModel UserModel
	err := r.db.Where("id = ?", id).First(&userModel).Error
	if err != nil {
		return domain.User{}, err
	}
	var password *string
	if userModel.Password != "" {
		password = &userModel.Password
	}
	return domain.User{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: password,
	}, nil
}

func (r *userRepository) FindByEmail(email string) (domain.User, error) {
	var userModel UserModel
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		return domain.User{}, err
	}
	var password *string
	if userModel.Password != "" {
		password = &userModel.Password
	}
	return domain.User{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: password,
	}, nil
}

func (r *userRepository) Save(user domain.User) error {
	userModel := UserModel{
		Name:     user.Name,
		Email:    user.Email,
		Password: "",
	}
	return r.db.Create(&userModel).Error
}
