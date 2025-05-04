package service

import (
	"net/http"

	"github.com/LoL-KeKovich/NoteVault/internal/repository"
)

type NoteService struct {
	DBClient repository.NoteRepo
}

func (srv NoteService) CreateNote(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) GetNoteByID(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) GetNotes(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) UpdateNote(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) DeleteNote(w http.ResponseWriter, r *http.Request) {

}
