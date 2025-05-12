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
	DBClient             repository.NoteRepo
	HelperNoteBookClient repository.NoteBookRepo
	HelperTagClient      repository.TagRepo
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

	if noteReq.IsDeleted == nil {
		noteReq.IsDeleted = new(bool)
	}

	*noteReq.IsDeleted = false //При создании элемент не может попасть в корзину

	note := model.Note{
		Name:       noteReq.Name,
		Text:       noteReq.Text,
		Color:      noteReq.Color,
		Media:      noteReq.Media,
		Order:      noteReq.Order,
		IsDeleted:  noteReq.IsDeleted,
		NoteBookID: noteReq.NoteBookID,
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

func (srv NoteService) HandleGetTrashedNotes(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	notes, err := srv.DBClient.GetTrashedNotes()
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding trashed notes in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Trashed notes found")
	response.Data = notes
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleGetNotesByNoteBookID(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err := srv.HelperNoteBookClient.GetNoteBookByID(id)
	if err != nil {
		slog.Error("No such notebook")
		response.Error = "Wrong notebook id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	notes, err := srv.DBClient.GetNotesByNoteBookID(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding notes from notebook"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notes from notebook found")
	response.Data = notes
	json.NewEncoder(w).Encode(response)
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

func (srv NoteService) HandleUpdateNoteNoteBook(w http.ResponseWriter, r *http.Request) {
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

	_, err = srv.HelperNoteBookClient.GetNoteBookByID(noteReq.NoteBookID.Hex())
	if err != nil {
		slog.Error(err.Error())
		response.Error = "error: wrong group id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.UpdateNoteNoteBook(id, noteReq.NoteBookID.Hex())
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error changing notebook for note"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notebook for note changed")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleAddTagToNote(w http.ResponseWriter, r *http.Request) {
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

	if noteReq.TagID.IsZero() {
		slog.Error("Empty tag_id field")
		response.Error = "No tag id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = srv.HelperTagClient.GetTagByID(noteReq.TagID.Hex())
	if err != nil {
		slog.Error(err.Error())
		response.Error = "wrong tag id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.AddTagToNote(id, noteReq.TagID.Hex())
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error adding tag to note"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Added tag to notebook")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleMoveNoteToTrash(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := srv.DBClient.MoveNoteToTrash(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error moving note to trash"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Changed |is_deleted| field to true")
	response.Data = "Successfully moved note to trash"
	json.NewEncoder(w).Encode(response)
}

func (srv NoteService) HandleRestoreNoteFromTrash(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := srv.DBClient.RestoreNoteFromTrash(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error restoring note from trash"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Changed |is_deleted| field to false")
	response.Data = "Successfully removed note from trash"
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

func (srv NoteService) HandleRemoveTagFromNote(w http.ResponseWriter, r *http.Request) {
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

	if noteReq.TagID.IsZero() {
		slog.Error("Empty tag_id field")
		response.Error = "No tag id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = srv.HelperTagClient.GetTagByID(noteReq.TagID.Hex())
	if err != nil {
		slog.Error(err.Error())
		response.Error = "wrong tag id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.RemoveTagFromNote(id, noteReq.TagID.Hex())
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error removing tag from note"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Removed tag from notebook")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}
