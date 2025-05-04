package repository

import "github.com/LoL-KeKovich/NoteVault/internal/model"

type NoteRepo interface {
	CreateNote(model.Note) (string, error)
	GetNoteByID(string) (model.Note, error)
	GetNotes() ([]model.Note, error)
	UpdateNote(string, string, string, string, string, int) (int, error)
	DeleteNote(string) (int, error)
}
