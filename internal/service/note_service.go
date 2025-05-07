package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LoL-KeKovich/NoteVault/internal/dto"
	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"github.com/LoL-KeKovich/NoteVault/internal/repository"
)

type NoteService struct {
	DBClient repository.NoteRepo
}

func (srv NoteService) CreateNote(w http.ResponseWriter, r *http.Request) {
	resp := dto.NoteResponse{}
	var noteReq dto.NoteRequest

	err := json.NewDecoder(r.Body).Decode(&noteReq)
	if err != nil {
		slog.Error(err.Error())
		resp.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
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
		resp.Error = "Error inserting note in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	slog.Info("Created note", slog.String("_id", res))
	resp.Data = res
	json.NewEncoder(w).Encode(resp)
}

func (srv NoteService) GetNoteByID(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) GetNotes(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) UpdateNote(w http.ResponseWriter, r *http.Request) {

}

func (srv NoteService) DeleteNote(w http.ResponseWriter, r *http.Request) {

}
