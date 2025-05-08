package main

import (
	"context"
	"log"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	mongoClient *mongo.Client
	UserCollection *mongo.Collection
	RevokedTokenCollection *mongo.Collection
)

func InitDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}
	dbName := os.Getenv("MONGO_DATABASE_NAME")
	if dbName == "" {
		log.Fatal("MONGO_DATABASE_NAME environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	mongoClient = client
	UserCollection = mongoClient.Database(dbName).Collection("users")
	RevokedTokenCollection = mongoClient.Database(dbName).Collection("revoked_tokens")

	ttlIndex := mongo.IndexModel{
		Keys:    bson.M{"expiresAt": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err = RevokedTokenCollection.Indexes().CreateOne(context.Background(), ttlIndex)
	if err != nil {
		log.Fatalf("Failed to create TTL index on revoked_tokens collection: %v", err)
	}

	log.Println("Successfully connected to MongoDB and database initialized.")
}