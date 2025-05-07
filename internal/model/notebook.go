package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteBook struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	IsActive    bool               `bson:"is_active,omitempty" json:"is_active,omitempty"`
}
