package mongodb

import (
	"context"
	"fmt"

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
