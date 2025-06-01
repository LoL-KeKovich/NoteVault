package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type NoteRequest struct {
	Name       string             `json:"name,omitempty"`
	Text       string             `json:"text,omitempty"`
	Color      string             `json:"color,omitempty"`
	Order      int                `json:"order,omitempty"`
	IsDeleted  *bool              `json:"is_deleted,omitempty"`
	IsArchived *bool              `json:"is_archived,omitempty"`
	NoteBookID primitive.ObjectID `json:"notebook_id,omitempty"`
	TagName    string             `json:"tag_name,omitempty"`
}

type NoteTagsRequest struct {
	TagNames []string `json:"tags,omitempty"`
}

type NoteResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
