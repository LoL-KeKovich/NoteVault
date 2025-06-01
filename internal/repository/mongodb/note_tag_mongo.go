package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mc MongoClient) AddTagToNote(noteID, tagName string) (int, error) {
	docId, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return 0, fmt.Errorf("invalid note ID: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "tags", Value: tagName}}}}

	res, err := mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func (mc MongoClient) RemoveTagFromNote(noteID, tagName string) (int, error) {
	docId, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return 0, fmt.Errorf("invalid note ID: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$pull", Value: bson.D{{Key: "tags", Value: tagName}}}}

	res, err := mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func (mc MongoClient) GetNotesByTags(tagNames []string) ([]model.Note, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "tags", Value: bson.D{{Key: "$all", Value: tagNames}}}},
			bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "is_deleted", Value: bson.D{{Key: "$exists", Value: false}}}},
					bson.D{{Key: "is_deleted", Value: false}},
				}},
			},
			bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "is_archived", Value: bson.D{{Key: "$exists", Value: false}}}},
					bson.D{{Key: "is_archived", Value: false}},
				}},
			},
		}},
	}
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
