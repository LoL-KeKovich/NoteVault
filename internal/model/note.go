package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string               `bson:"name,omitempty" json:"name,omitempty"`
	Text       string               `bson:"text,omitempty" json:"text,omitempty"`
	Color      string               `bson:"color,omitempty" json:"color,omitempty"`
	Order      int                  `bson:"order,omitempty" json:"order,omitempty"`
	IsDeleted  *bool                `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
	IsArchived *bool                `bson:"is_archived,omitempty" json:"is_archived,omitempty"`
	CreatedAt  time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time            `bson:"updated_at" json:"updated_at"`
	NoteBookID primitive.ObjectID   `bson:"notebook_id,omitempty" json:"notebook_id,omitempty"`
	Tags       []primitive.ObjectID `bson:"tags,omitempty" json:"tags,omitempty"`
}
