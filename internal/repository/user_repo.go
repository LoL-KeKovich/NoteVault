package repository

import "github.com/LoL-KeKovich/NoteVault/internal/model"

type UserRepo interface {
	RegisterUser(model.User) (string, error)
	LoginUser(string) (model.User, error)
	GetProfile(string) (model.User, error)
}
