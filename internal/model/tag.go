package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string             `bson:"name,omitempty" json:"name,omitempty"`
	Color string             `bson:"color,omitempty" json:"color,omitempty"`
}
