package repository

import "web3/models"

type DatabaseRepo interface {
	InsertPost(newPost models.Post) error

	GetUserByID(id int) (models.User, error)

	UpdateUser(user models.User) error

	AuthenticateUser(email, testPassword string) (int, string, error)
}
