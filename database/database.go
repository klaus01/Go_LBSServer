package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/klaus01/Go_LBSServer/config"
)

var db *mongo.Database

func Init() {
	c := config.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.GetString("mongo.uri")))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database(c.GetString("mongo.dbName"))
}

func Context() context.Context {
	return context.Background()
}

func GetDB() *mongo.Database {
	return db
}
