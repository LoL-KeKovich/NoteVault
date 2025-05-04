package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Note struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Text  string             `bson:"text,omitempty"`
	Color string             `bson:"color,omitempty"`
	Media []byte             `bson:"media,omitempty"` //???
	Order int                `bson:"order,omitempty"`
}
