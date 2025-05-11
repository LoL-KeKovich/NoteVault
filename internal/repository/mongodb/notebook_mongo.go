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

func (mc MongoClient) CreateNoteBook(notebook model.NoteBook) (string, error) {
	res, err := mc.Client.InsertOne(context.Background(), notebook)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (mc MongoClient) GetNoteBookByID(id string) (model.NoteBook, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.NoteBook{}, fmt.Errorf("wrong id")
	}

	var noteBook model.NoteBook

	filter := bson.D{{Key: "_id", Value: docId}}

	err = mc.Client.FindOne(context.Background(), filter).Decode(&noteBook)
	if err == mongo.ErrNoDocuments {
		return model.NoteBook{}, fmt.Errorf("notebook not found")
	} else if err != nil {
		return model.NoteBook{}, err
	}

	return noteBook, nil
}

func (mc MongoClient) GetNoteBooks() ([]model.NoteBook, error) {
	filter := bson.D{}

	cursor, err := mc.Client.Find(context.Background(), filter)
	if err != nil {
		return []model.NoteBook{}, fmt.Errorf("error finding notebooks")
	}
	defer cursor.Close(context.Background())

	var noteBooks []model.NoteBook

	for cursor.Next(context.Background()) {
		var noteBook model.NoteBook

		err := cursor.Decode(&noteBook)
		if err != nil {
			slog.Error("error decoding notebooks", slog.String("error", err.Error()))
			continue
		}

		noteBooks = append(noteBooks, noteBook)
	}

	if len(noteBooks) == 0 {
		return []model.NoteBook{}, fmt.Errorf("empty slice")
	}

	return noteBooks, nil
}

func (mc MongoClient) UpdateNoteBook(id, name, description string, isActive *bool) (int, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("wrong id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}

	setDoc := bson.D{}
	if name != "" {
		setDoc = append(setDoc, bson.E{Key: "name", Value: name})
	}
	if description != "" {
		setDoc = append(setDoc, bson.E{Key: "description", Value: description})
	}
	if isActive != nil {
		setDoc = append(setDoc, bson.E{Key: "is_active", Value: *isActive})
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

func (mc MongoClient) DeleteNoteBook(id string) (int, error) {
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
