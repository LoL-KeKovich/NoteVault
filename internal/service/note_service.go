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

type NoteService struct {
	DBClient repository.NoteRepo
}

func (srv NoteService) HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}
	var noteReq dto.NoteRequest

	err := json.NewDecoder(r.Body).Decode(&noteReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	note := model.Note{
		Name:  noteReq.Name,
		Text:  noteReq.Text,
		Color: noteReq.Color,
		Media: noteReq.Media,
		Order: noteReq.Order,
	}

	res, err := srv.DBClient.CreateNote(note)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error inserting note in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Created note", slog.String("_id", res))
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	note, err := srv.DBClient.GetNoteByID(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding note in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Note found")
	response.Data = note
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	notes, err := srv.DBClient.GetNotes()
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding notes in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notes found")
	response.Data = notes
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleGetNotesByNoteBookID(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) HandleUpdateNote(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}
	var noteReq dto.NoteRequest

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&noteReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.UpdateNote(id, noteReq.Name, noteReq.Text, noteReq.Color, noteReq.Media, noteReq.Order)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error updating note in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Note updated")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleDeleteNote(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.DeleteNote(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error deleting note in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Note deleted")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}
