package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LoL-KeKovich/NoteVault/internal/dto"
	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"github.com/LoL-KeKovich/NoteVault/internal/repository"
	"github.com/go-chi/chi"
)

type NoteBookService struct {
	DBClient         repository.NoteBookRepo
	HelperNoteClient repository.NoteRepo
}

func (srv NoteBookService) HandleCreateNoteBook(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteBookResponse{}
	var noteBookReq dto.NoteBookRequest

	err := json.NewDecoder(r.Body).Decode(&noteBookReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	noteBook := model.NoteBook{
		Name:        noteBookReq.Name,
		Description: noteBookReq.Description,
		IsActive:    noteBookReq.IsActive,
	}

	res, err := srv.DBClient.CreateNoteBook(noteBook)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error inserting notebook in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Created notebook", slog.String("_id", res))
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv NoteBookService) HandleGetNoteBookByID(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteBookResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	noteBook, err := srv.DBClient.GetNoteBookByID(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding notebook in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notebook found")
	response.Data = noteBook
	json.NewEncoder(w).Encode(response)
}

func (srv NoteBookService) HandleGetNoteBooks(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteBookResponse{}

	noteBooks, err := srv.DBClient.GetNoteBooks()
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding notebooks in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notebooks found")
	response.Data = noteBooks
	json.NewEncoder(w).Encode(response)
}

func (srv NoteBookService) HandleUpdateNoteBook(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteBookResponse{}
	var noteBookReq dto.NoteBookRequest

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&noteBookReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.UpdateNoteBook(id, noteBookReq.Name, noteBookReq.Description, noteBookReq.IsActive)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error updating notebook in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notebook updated")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv NoteBookService) HandleDeleteNoteBook(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteBookResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err := srv.HelperNoteClient.UnlinkNotesFromNoteBook(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error unlinking notes from notebook"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.DeleteNoteBook(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error deleting notebook in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notebook deleted")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}
