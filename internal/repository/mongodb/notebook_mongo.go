package mongodb

import (
	"github.com/LoL-KeKovich/NoteVault/internal/model"
)

func (mc MongoClient) CreateNoteBook(notebook model.NoteBook) (string, error) {
	return "", nil
}

func (mc MongoClient) GetNoteBookByID(id string) (model.NoteBook, error) {
	return model.NoteBook{}, nil
}

func (mc MongoClient) GetNoteBooks() ([]model.NoteBook, error) {
	return []model.NoteBook{}, nil
}

func (mc MongoClient) UpdateNoteBook(id, name, description string, isActive bool) (int, error) {
	return 0, nil
}

func (mc MongoClient) DeleteNoteBook(id string) (int, error) {
	return 0, nil
}
