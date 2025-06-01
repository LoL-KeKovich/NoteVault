package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mc MongoClient) CreateTag(tag model.Tag) (string, error) {
	res, err := mc.Client.InsertOne(context.Background(), tag)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (mc MongoClient) GetTagByID(id string) (model.Tag, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Tag{}, fmt.Errorf("wrong id")
	}

	var tag model.Tag

	filter := bson.D{{Key: "_id", Value: docId}}

	err = mc.Client.FindOne(context.Background(), filter).Decode(&tag)
	if err == mongo.ErrNoDocuments {
		return model.Tag{}, fmt.Errorf("tag not found")
	} else if err != nil {
		return model.Tag{}, err
	}

	return tag, nil
}

func (mc MongoClient) GetTagByName(tagName string) (model.Tag, error) {
	var tag model.Tag

	filter := bson.D{{Key: "name", Value: tagName}}

	err := mc.Client.FindOne(context.Background(), filter).Decode(&tag)
	if err == mongo.ErrNoDocuments {
		return model.Tag{}, fmt.Errorf("tag not found")
	} else if err != nil {
		return model.Tag{}, err
	}

	return tag, nil
}

func (mc MongoClient) GetTags() ([]model.Tag, error) {
	filter := bson.D{}

	cursor, err := mc.Client.Find(context.Background(), filter)
	if err != nil {
		return []model.Tag{}, fmt.Errorf("error finding tags")
	}
	defer cursor.Close(context.Background())

	var tags []model.Tag

	for cursor.Next(context.Background()) {
		var tag model.Tag

		err := cursor.Decode(&tag)
		if err != nil {
			slog.Error("error decoding notes", slog.String("error", err.Error()))
			continue
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (mc MongoClient) UpdateTag(id, name, color string) (int, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("wrong id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}

	setDoc := bson.D{}
	if name != "" {
		setDoc = append(setDoc, bson.E{Key: "name", Value: name})
	}
	if color != "" {
		setDoc = append(setDoc, bson.E{Key: "color", Value: color})
	}
	if len(setDoc) == 0 {
		return 0, nil
	}

	updateStmt := bson.D{{Key: "$set", Value: setDoc}}

	res, err := mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func (mc MongoClient) DeleteTag(id string) (int, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("wrong id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}

	res, err := mc.Client.DeleteOne(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return int(res.DeletedCount), nil
}
