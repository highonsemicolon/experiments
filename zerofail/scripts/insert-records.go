package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb://admin:admin@my-mongo-mongodb-headless.mongodb.svc.cluster.local:27017/?authSource=admin"
	dbName         = "appdb"
	collectionName = "records"
	totalRecords   = 5_000_000
	batchSize      = 1_000
)

func createUniqueIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexModel1 := mongo.IndexModel{
		Keys:    map[string]interface{}{"col1": 1},
		Options: options.Index().SetUnique(true).SetName("unique_col1"),
	}

	indexModel2 := mongo.IndexModel{
		Keys:    map[string]interface{}{"col2": 1},
		Options: options.Index().SetUnique(true).SetName("unique_col2"),
	}

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{indexModel1, indexModel2})
	return err
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection(collectionName)

	err = createUniqueIndexes(ctx, collection)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	fmt.Println("Unique indexes created (if not already present).")

	fmt.Printf("Inserting %d records in batches of %d...\n", totalRecords, batchSize)

	start := time.Now()

	for i := 0; i < totalRecords; i += batchSize {
		var batch []any

		for j := 0; j < batchSize && (i+j) < totalRecords; j++ {
			id := i + j
			batch = append(batch, map[string]any{
				"col1":      fmt.Sprintf("col1_value_%d", id),
				"col2":      fmt.Sprintf("col2_value_%d", id),
				"createdAt": time.Now(),
			})
		}

		_, err := collection.InsertMany(ctx, batch)
		if err != nil {
			log.Fatalf("Insert failed at record %d: %v", i, err)
		}

		if (i+batchSize)%100000 == 0 {
			fmt.Printf("Inserted %d records...\n", i+batchSize)
		}
	}

	fmt.Printf("Successfully inserted %d records in %s\n", totalRecords, time.Since(start))
}
