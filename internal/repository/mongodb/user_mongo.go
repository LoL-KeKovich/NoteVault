package mongodb

import (
	"context"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (mc MongoClient) LoginUser(email string) (model.User, error) {
	var user model.User

	filter := bson.D{{Key: "email", Value: email}}

	err := mc.Client.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
