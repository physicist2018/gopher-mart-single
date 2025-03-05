package repository

import (
	"github.com/physicist2018/gopher-mart-single/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}
