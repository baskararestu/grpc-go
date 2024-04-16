package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB            *mongo.Collection
	DBCollections ColsType
	MongoClient   *mongo.Client
	err           error
)

func ConnectMongo() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DATABASE_NAME")
	productCols := os.Getenv("COLLECTION_PRODUCTS")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error while connecting to db: %v", err)
	}

	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DBCollections.Products = MongoClient.Database(dbName).Collection(productCols)
}

func GetMongoDB() *ColsType {
	return &DBCollections
}
