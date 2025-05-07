package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type NoteRequest struct {
	Name       string             `json:"name,omitempty"`
	Text       string             `json:"text,omitempty"`
	Color      string             `json:"color,omitempty"`
	Media      string             `json:"media,omitempty"`
	Order      int                `json:"order,omitempty"`
	NoteBookID primitive.ObjectID `json:"note_book_id,omitempty"`
}

type NoteResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
