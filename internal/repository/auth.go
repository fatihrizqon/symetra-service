package repository

import (
	"errors"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	Login(username string) (entity.User, error)
}

type AuthRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) IAuthRepository {
	return &AuthRepository{Db: Db}
}

// Login implements IAuthRepository.
func (e *AuthRepository) Login(email string) (entity.User, error) {
	var entity entity.User
	if err := e.Db.Where("email = ?", email).First(&entity).Error; err != nil {
		return entity, errors.New("credentials does not matches our record")
	}
	return entity, nil
}
