package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

	if noteReq.IsArchived == nil {
		noteReq.IsArchived = new(bool)
	}

	*noteReq.IsArchived = false //При создании элемент не может попасть в архив

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		slog.Error("Failed to load Moscow location, using UTC", slog.String("error", err.Error()))
		location = time.UTC
	}

	now := time.Now().In(location)

	note := model.Note{
		Name:       noteReq.Name,
		Text:       noteReq.Text,
		Color:      noteReq.Color,
		Order:      noteReq.Order,
		IsDeleted:  noteReq.IsDeleted,
		IsArchived: noteReq.IsArchived,
		NoteBookID: noteReq.NoteBookID,
		CreatedAt:  now,
		UpdatedAt:  now,
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

func (srv NoteService) HandleGetArchivedNotes(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	notes, err := srv.DBClient.GetArchivedNotes()
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding archived notes in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Archived notes found")
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

func (srv NoteService) HandleGetNotesByTags(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}
	var noteReq dto.NoteTagsRequest

	err := json.NewDecoder(r.Body).Decode(&noteReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var tagNames []string

	tagNames = append(tagNames, noteReq.TagNames...)

	notes, err := srv.DBClient.GetNotesByTags(tagNames)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding notes by tags"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Notes by tags found")
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

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		slog.Error("Failed to load Moscow location, using UTC", slog.String("error", err.Error()))
		location = time.UTC
	}

	now := time.Now().In(location)

	res, err := srv.DBClient.UpdateNote(id, noteReq.Name, noteReq.Text, noteReq.Color, noteReq.Order, now)
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

	if noteReq.TagName == "" {
		slog.Error("Empty tag_name field")
		response.Error = "No tag name"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = srv.HelperTagClient.GetTagByName(noteReq.TagName)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "wrong tag name"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.AddTagToNote(id, noteReq.TagName)
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

func (srv NoteService) HandleMoveNoteToArchive(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := srv.DBClient.MoveNoteToArchive(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error moving note to archive"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Changed |is_archived| field to true")
	response.Data = "Successfully moved note to archive"
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

func (srv NoteService) HandleRestoreNoteFromArchive(w http.ResponseWriter, r *http.Request) {
	response := dto.NoteResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := srv.DBClient.RestoreNoteFromArchive(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error restoring note from archive"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Changed |is_archived| field to false")
	response.Data = "Successfully removed note from archive"
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

	if noteReq.TagName == "" {
		slog.Error("Empty tag_name field")
		response.Error = "No tag name"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = srv.HelperTagClient.GetTagByName(noteReq.TagName)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "wrong tag name"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.RemoveTagFromNote(id, noteReq.TagName)
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
