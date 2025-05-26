package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type NoteRequest struct {
	Name       string             `json:"name,omitempty"`
	Text       string             `json:"text,omitempty"`
	Color      string             `json:"color,omitempty"`
	Order      int                `json:"order,omitempty"`
	IsDeleted  *bool              `json:"is_deleted,omitempty"`
	NoteBookID primitive.ObjectID `json:"notebook_id,omitempty"`
	TagID      primitive.ObjectID `json:"tag_id,omitempty"`
}

type NoteTagsRequest struct {
	TagIDs []primitive.ObjectID `json:"tags,omitempty"`
}

type NoteResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
