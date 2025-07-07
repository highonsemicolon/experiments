package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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
		Keys:    map[string]any{"col1": 1},
		Options: options.Index().SetUnique(true).SetName("unique_col1"),
	}

	indexModel2 := mongo.IndexModel{
		Keys:    map[string]any{"col2": 1},
		Options: options.Index().SetUnique(true).SetName("unique_col2"),
	}

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{indexModel1, indexModel2})
	return err
}

func makeBatch(start, size int) []any {
	batch := make([]any, 0, size)
	now := time.Now()
	counter := start * batchSize
	for i := range size {
		for j := range uniquePairs {
			document := map[string]any{
				"_id":       fmt.Sprintf("%d_%02d", start+i, j),
				"col1":      fmt.Sprintf("col1_%d", counter),
				"col2":      fmt.Sprintf("col2_%d", counter),
				"createdAt": now,
			}
			batch = append(batch, document)
			counter++
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

	fmt.Println("Starting bulk data load...")
	loadDataInParallel(ctx, collection)

	fmt.Println("Creating unique indexes...")
	err = createUniqueIndexes(ctx, collection)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	fmt.Println("Done.")
}
