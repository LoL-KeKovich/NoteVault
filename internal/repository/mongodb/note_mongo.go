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

func (mc MongoClient) CreateNote(note model.Note) (string, error) {
	res, err := mc.Client.InsertOne(context.Background(), note)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (mc MongoClient) GetNoteByID(id string) (model.Note, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Note{}, fmt.Errorf("wrong id")
	}

	var note model.Note
	filter := bson.D{{Key: "_id", Value: docId}}

	err = mc.Client.FindOne(context.Background(), filter).Decode(&note)
	if err == mongo.ErrNoDocuments {
		return model.Note{}, fmt.Errorf("note not found")
	} else if err != nil {
		return model.Note{}, err
	}

	return note, nil
}

func (mc MongoClient) GetNotes() ([]model.Note, error) {
	filter := bson.D{}

	cursor, err := mc.Client.Find(context.Background(), filter)
	if err != nil {
		return []model.Note{}, nil
	}
	defer cursor.Close(context.Background())

	var notes []model.Note

	for cursor.Next(context.Background()) {
		var note model.Note

		err := cursor.Decode(&note)
		if err != nil {
			slog.Error("error decoding notes", slog.String("error", err.Error()))
			continue
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (mc MongoClient) GetNotesByNoteBookID(id string) ([]model.Note, error) {
	return []model.Note{}, nil
}

func (mc MongoClient) UpdateNote(id, name, text, color, media string, order int) (int, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("wrong id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}

	//Необходимый костыль
	setDoc := bson.D{}
	if name != "" {
		setDoc = append(setDoc, bson.E{Key: "name", Value: name})
	}
	if text != "" {
		setDoc = append(setDoc, bson.E{Key: "text", Value: text})
	}
	if color != "" {
		setDoc = append(setDoc, bson.E{Key: "color", Value: color})
	}
	if media != "" {
		setDoc = append(setDoc, bson.E{Key: "media", Value: media})
	}
	if order != 0 {
		setDoc = append(setDoc, bson.E{Key: "order", Value: order})
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

func (mc MongoClient) DeleteNote(id string) (int, error) {
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
