package services

import "github.com/croisade/chimichanga/models"

type UserService interface {
	CreateUser(*models.User) (*models.User, error)
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}
