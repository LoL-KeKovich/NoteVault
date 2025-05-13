package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Media struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Data   string             `bson:"data,omitempty" json:"data,omitempty"`
	NoteID primitive.ObjectID `bson:"note_id,omitempty" json:"note_id,omitempty"`
}
