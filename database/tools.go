package database

import (
	"context"
	"errors"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	godotenv.Load("../.env")
	db_uri := os.Getenv("DB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db_uri))
	if err != nil {
		return new(mongo.Client), err
	}
	return client, nil
}

func DisconnectDB(client *mongo.Client) error {
	if err := client.Disconnect(context.Background()); err != nil {
		return errors.New("failed to close database connection")
	}
	return nil
}
