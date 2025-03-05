package repository

import (
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/database/connector"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/ports/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	conn *connector.Connector
}

// Функция для создания нового репозитория пользователей
func NewUserRepository(conn *connector.Connector) repository.UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) CreateUser(user *models.User) error {
	if err := r.conn.DB().Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.conn.DB().First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	if err := r.conn.DB().Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(id uint) error {
	if err := r.conn.DB().Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	if err := r.conn.DB().Save(user).Error; err != nil {
		return err
	}
	return nil
}
