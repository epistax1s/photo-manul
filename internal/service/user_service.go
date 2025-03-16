package service

import (
	"github.com/epistax1s/photo-manul/internal/model"
	"github.com/epistax1s/photo-manul/internal/repository"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUserByChatID(chatID int64) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) CreateUser(user *model.User) error {
	return service.repo.Create(user)
}

func (service *userService) GetUserByChatID(chatID int64) (*model.User, error) {
	return service.repo.FindByChatID(chatID)
}
