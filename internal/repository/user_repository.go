package repository

import (
	"fmt"

	"github.com/epistax1s/photo-manul/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByChatID(chatID int64) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) FindByChatID(chatID int64) (*model.User, error) {
	var user model.User

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.UserChatIDColumn), chatID).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
