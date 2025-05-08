package service

import (
	"net/http"

	"github.com/LoL-KeKovich/NoteVault/internal/repository"
)

type NoteBookService struct {
	DBClient repository.NoteBookRepo
}

func (srv NoteBookService) HandleCreateNoteBook(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteBookService) HandleGetNoteBookByID(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteBookService) HandleGetNoteBooks(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteBookService) HandleUpdateNoteBook(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteBookService) HandleDeleteNoteBook(w http.ResponseWriter, r *http.Request) {

}
