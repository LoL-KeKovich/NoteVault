package repository

import "github.com/LoL-KeKovich/NoteVault/internal/model"

type NoteRepo interface {
	CreateNote(model.Note) (string, error)
	GetNoteByID(string) (model.Note, error)
	GetNotes() ([]model.Note, error)
	GetNotesByNoteBookID(string) ([]model.Note, error)
	GetTrashedNotes() ([]model.Note, error)
	UpdateNote(string, string, string, string, string, int) (int, error)
	UpdateNoteNoteBook(string, string) (int, error)
	UnlinkNotesFromNoteBook(string) (int, error)
	UnlinkNotesFromTag(string) (int, error)
	AddTagToNote(string, string) (int, error)
	MoveNoteToTrash(string) error
	RestoreNoteFromTrash(string) error
	DeleteNote(string) (int, error)
}
