package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ensureUniqueEmailIndex(ctx context.Context) {

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order
		Options: options.Index().SetUnique(true),
	}

	_, err := UserCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
}
