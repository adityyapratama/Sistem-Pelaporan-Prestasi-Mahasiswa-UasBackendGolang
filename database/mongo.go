package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func ConnectMongo() (*mongo.Client, *mongo.Database, error) {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	if uri == "" || dbName == "" {
		return nil, nil, fmt.Errorf("MONGO_URI atau MONGO_DB_NAME belum diset di .env")
	}

	// Setup opsi client
	clientOptions := options.Client().ApplyURI(uri)

	// Bikin context dengan timeout 10 detik (biar gak hang kalau error)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Test Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("✅ Berhasil terhubung ke MongoDB")
	return client, client.Database(dbName), nil
}