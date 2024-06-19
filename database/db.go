package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	Port     = ":8080"
	MongoURL = "mongodb://localhost:27050"
	DB       *mongo.Client
	RmsDB    *mongo.Database

	UserCollection       *mongo.Collection
	RestaurantCollection *mongo.Collection
	DishCollection       *mongo.Collection
	OrderCollection      *mongo.Collection

	MongoCtx = context.Background()
)

func ConnectDatabase() {
	fmt.Println("Connecting to MongoDB...")
	clientOptions := options.Client().ApplyURI(MongoURL)
	client, err := mongo.Connect(MongoCtx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err := client.Ping(MongoCtx, nil); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	DB = client
	RmsDB = DB.Database("rmsDB")
	UserCollection = RmsDB.Collection("users")
	RestaurantCollection = RmsDB.Collection("restaurants")
	DishCollection = RmsDB.Collection("dishes")
	OrderCollection = RmsDB.Collection("orders")

	fmt.Println("Connected to MongoDB")

	// email unique in users
	ensureUniqueEmailIndex(MongoCtx)
}
