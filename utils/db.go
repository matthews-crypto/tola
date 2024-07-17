package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://itachi:Anoman2002@cluster0.aiemcbo.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	MongoClient = client
	log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoClient is not initialized")
	}
	return MongoClient.Database("tola").Collection(collectionName)
}
