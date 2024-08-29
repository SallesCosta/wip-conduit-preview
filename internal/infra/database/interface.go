package database

import (
	userEntity "github.com/sallescosta/conduit-api/internal/entity/user"
)

type UserInterface interface {
	CreateUser(user *userEntity.User) error
	FindByEmail(email string) (*userEntity.User, error)
	FindById(id string) (*userEntity.User, error)
	GetAllUsers() ([]userEntity.User, error)
	UpdateUserDb(email, username, password, image, bio string) (*userEntity.User, error)
	GetProfileDb(userName string) (*ProfileWithId, error)
}
