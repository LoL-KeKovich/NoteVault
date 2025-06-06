package repository

import (
	"github.com/LoL-KeKovich/NoteVault/internal/model"
)

type NoteRepo interface {
	CreateNote(model.Note) (string, error)
	GetNoteByID(string) (model.Note, error)
	GetNotes() ([]model.Note, error)
	GetNotesByNoteBookID(string) ([]model.Note, error)
	GetTrashedNotes() ([]model.Note, error)
	GetArchivedNotes() ([]model.Note, error)
	GetNotesByTags([]string) ([]model.Note, error)
	UpdateNote(string, string, string, string, string, int) (int, error)
	UpdateNoteNoteBook(string, string) (int, error)
	RemoveNoteBookFromNote(string) (int, error)
	UnlinkNotesFromNoteBook(string) (int, error)
	UnlinkNotesFromTag(string) (int, error)
	AddTagToNote(string, string) (int, error)
	MoveNoteToTrash(string) error
	MoveNoteToArchive(string) error
	RestoreNoteFromTrash(string) error
	RestoreNoteFromArchive(string) error
	RemoveTagFromNote(string, string) (int, error)
	DeleteNote(string) (int, error)
}
