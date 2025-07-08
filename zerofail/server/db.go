package server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var RecordCollection *mongo.Collection
var DeletedCollection *mongo.Collection
var Client *mongo.Client

func InitMongoDB(uri string, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMinPoolSize(2).SetMaxPoolSize(10))
	if err != nil {
		log.Fatal(err)
	}

	db := Client.Database(dbName)
	RecordCollection = db.Collection("records")
	DeletedCollection = db.Collection("deleted_records")

	_, _ = RecordCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "col1", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "col2", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

}
