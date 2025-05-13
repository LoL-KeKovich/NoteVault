package mongodb

import (
	"context"

	"github.com/LoL-KeKovich/NoteVault/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mc MongoClient) CreateReminder(reminder model.Reminder) (string, error) {
	if reminder.IsActive == nil {
		isActive := true
		reminder.IsActive = &isActive
	}

	res, err := mc.Client.InsertOne(context.Background(), reminder)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
