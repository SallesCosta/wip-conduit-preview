package entity

import (
	"github.com/sallescosta/conduit-api/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        entity.ID `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Bio       string    `json:"bio"`
	Image     string    `json:"image"`
	Following []entity.ID
}

func DoHash(password string) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func NewUser(username, email, password string) (*User, error) {

	hash, err := DoHash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        entity.NewID(),
		UserName:  username,
		Email:     email,
		Password:  string(hash),
		Bio:       "",
		Image:     "",
		Following: []entity.ID{},
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
