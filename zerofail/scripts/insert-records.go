package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	mongoURI       = "mongodb://admin:admin@my-mongo-mongodb-headless.mongodb.svc.cluster.local:27017/?authSource=admin"
	dbName         = "appdb"
	collectionName = "records"
	orders         = 100_000
	uniquePairs    = 36
	batchSize      = 1_000
	workers        = 4
)

func createUniqueIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexModel1 := mongo.IndexModel{
		Keys:    map[string]any{"pairs.col1": 1},
		Options: options.Index().SetUnique(true).SetName("unique_col1"),
	}

	indexModel2 := mongo.IndexModel{
		Keys:    map[string]any{"pairs.col2": 1},
		Options: options.Index().SetUnique(true).SetName("unique_col2"),
	}

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{indexModel1, indexModel2})
	return err
}

func makeBatch(start, size int) []any {
	batch := make([]any, size)
	now := time.Now()
	counter := start * uniquePairs

	for i := 0; i < size; i++ {
		pairs := make([]bson.D, uniquePairs)
		for j := 0; j < uniquePairs; j++ {
			pairs[j] = bson.D{
				{Key: "col1", Value: fmt.Sprintf("col1_%02d", counter)},
				{Key: "col2", Value: fmt.Sprintf("col2_%02d", counter)},
			}
			counter++
		}

		batch[i] = map[string]any{
			"_id":       fmt.Sprintf("%d", start+i),
			"pairs":     pairs,
			"createdAt": now,
		}
	}

	return batch
}

func loadDataInParallel(ctx context.Context, collection *mongo.Collection) {
	insertOpts := options.InsertMany().SetOrdered(false)
	jobChan := make(chan []any, workers)

	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()
			for batch := range jobChan {
				_, err := collection.InsertMany(ctx, batch, insertOpts)
				if err != nil {
					log.Fatalf("Insert failed: %v", err)
				}
				// log.Println(batch...)
			}
		}()
	}

	start := time.Now()
	for i := 0; i < orders; i += batchSize {
		size := min(orders-i, batchSize)
		jobChan <- makeBatch(i, size)
	}

	close(jobChan)
	wg.Wait()
	fmt.Printf("All records inserted in %s\n", time.Since(start))
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(200).
		SetWriteConcern(&writeconcern.WriteConcern{})

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection(collectionName)
	// collection.Drop(ctx)

	fmt.Println("Creating unique indexes...")
	err = createUniqueIndexes(ctx, collection)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	fmt.Println("Starting bulk data load...")
	loadDataInParallel(ctx, collection)

	fmt.Println("Done.")
}
