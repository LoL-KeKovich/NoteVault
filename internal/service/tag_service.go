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

type TagService struct {
	DBClient         repository.TagRepo
	HelperNoteClient repository.NoteRepo
}

func (srv TagService) HandleCreateTag(w http.ResponseWriter, r *http.Request) {
	response := dto.TagResponse{}
	var tagReq dto.TagRequest

	err := json.NewDecoder(r.Body).Decode(&tagReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tag := model.Tag{
		Name:  tagReq.Name,
		Color: tagReq.Color,
	}

	res, err := srv.DBClient.CreateTag(tag)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error inserting tag in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Created tag", slog.String("_id", res))
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv TagService) HandleGetTagByID(w http.ResponseWriter, r *http.Request) {
	response := dto.TagResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tag, err := srv.DBClient.GetTagByID(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding tag in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Tag found")
	response.Data = tag
	json.NewEncoder(w).Encode(response)
}

func (srv TagService) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	response := dto.TagResponse{}

	tags, err := srv.DBClient.GetTags()
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error finding tags in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Tags found")
	response.Data = tags
	json.NewEncoder(w).Encode(response)
}

func (srv TagService) HandleUpdateTag(w http.ResponseWriter, r *http.Request) {
	response := dto.TagResponse{}
	var tagReq dto.TagRequest

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&tagReq)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Wrong request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.UpdateTag(id, tagReq.Name, tagReq.Color)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error updating tag in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Tag updated")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}

func (srv TagService) HandleDeleteTag(w http.ResponseWriter, r *http.Request) {
	response := dto.TagResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("Empty id field")
		response.Error = "Wrong id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tag, err := srv.DBClient.GetTagByID(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Tag not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = srv.HelperNoteClient.UnlinkNotesFromTag(tag.Name)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error unlinking notes from tag"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	res, err := srv.DBClient.DeleteTag(id)
	if err != nil {
		slog.Error(err.Error())
		response.Error = "Error deleting tag in db"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	slog.Info("Tag deleted")
	response.Data = res
	json.NewEncoder(w).Encode(response)
}
