package repository

import "web3/models"

type DatabaseRepo interface {
	InsertPost(newPost models.Post) error
}
