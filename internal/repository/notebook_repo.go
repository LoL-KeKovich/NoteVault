package repository

import "github.com/LoL-KeKovich/NoteVault/internal/model"

type NoteBookRepo interface {
	CreateNoteBook(model.NoteBook) (string, error)
	GetNoteBookByID(string) (model.NoteBook, error)
	GetNoteBooks() ([]model.NoteBook, error)
	UpdateNoteBook(string, string, string, bool) (int, error)
	DeleteNoteBook(string) (int, error)
}
