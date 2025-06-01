package repository

import "github.com/LoL-KeKovich/NoteVault/internal/model"

type TagRepo interface {
	CreateTag(model.Tag) (string, error)
	GetTagByID(string) (model.Tag, error)
	GetTagByName(string) (model.Tag, error)
	GetTags() ([]model.Tag, error)
	UpdateTag(string, string, string) (int, error)
	DeleteTag(string) (int, error)
}
