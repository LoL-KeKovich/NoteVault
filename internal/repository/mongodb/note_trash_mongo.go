package mongodb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mc MongoClient) MoveNoteToTrash(id string) error {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("wrong id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$set", Value: bson.D{{Key: "is_deleted", Value: true}}}}

	_, err = mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return err
	}

	return nil
}

func (mc MongoClient) RestoreNoteFromTrash(id string) error {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("wrong id")
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$set", Value: bson.D{{Key: "is_deleted", Value: false}}}}

	_, err = mc.Client.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return err
	}

	return nil
}

func (mc MongoClient) GetTrashedNotes() ([]model.Note, error) {
	filter := bson.D{{Key: "is_deleted", Value: true}}

	cursor, err := mc.Client.Find(context.Background(), filter)
	if err != nil {
		return []model.Note{}, fmt.Errorf("error finding notes in trash")
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

	if len(notes) == 0 {
		return []model.Note{}, fmt.Errorf("empty slice")
	}

	return notes, nil
}
