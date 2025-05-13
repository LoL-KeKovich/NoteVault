package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reminder struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name,omitempty" json:"name,omitempty"`
	Message  string             `bson:"message,omitempty" json:"message,omitempty"`
	RemindAt time.Time          `bson:"remind_at,omitempty" json:"remind_at,omitempty"`
	IsActive *bool              `bson:"is_active,omitempty" json:"is_active,omitempty"`
	Repeat   string             `bson:"repeat,omitempty" json:"repeat,omitempty"`
	NoteID   primitive.ObjectID `bson:"note_id" json:"note_id"`
}
