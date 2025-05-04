package mongo

import "go.mongodb.org/mongo-driver/mongo"

type MongoClient struct {
	Client mongo.Collection
}
