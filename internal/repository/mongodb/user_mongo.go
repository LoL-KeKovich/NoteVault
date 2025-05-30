package mongodb

import (
	"context"
	"fmt"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (mc MongoClient) GetProfile(id string) (model.User, error) {
	var user model.User

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, fmt.Errorf("wrong user id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}

	err = mc.Client.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
