package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mc MongoClient) AddTagToNote(noteID, tagID string) (int, error) {
	docId, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return 0, fmt.Errorf("invalid note ID: %v", err)
	}

	objTagId, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		return 0, fmt.Errorf("invalid tag ID: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "tags", Value: objTagId}}}}

	res, err := mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func (mc MongoClient) RemoveTagFromNote(noteID, tagID string) (int, error) {
	docId, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return 0, fmt.Errorf("invalid note ID: %v", err)
	}

	objTagId, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		return 0, fmt.Errorf("invalid tag ID: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$pull", Value: bson.D{{Key: "tags", Value: objTagId}}}}

	res, err := mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func (mc MongoClient) GetNotesByTags(tagIDs []string) ([]model.Note, error) {
	var objTagIDs []primitive.ObjectID

	for _, id := range tagIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid tag id: %s", id)
		}
		objTagIDs = append(objTagIDs, objID)
	}

	filter := bson.D{{Key: "tags", Value: bson.D{{Key: "$all", Value: objTagIDs}}}}
	cursor, err := mc.Client.Find(context.Background(), filter)
	if err != nil {
		return []model.Note{}, fmt.Errorf("error finding notes by tags: %v", err)
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
