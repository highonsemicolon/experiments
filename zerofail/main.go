package main

import (
	server "github.com/highonsemicolon/experiments/zerofail/server"
)

const (
	mongoURI       = "mongodb://admin:admin@my-mongo-mongodb-headless.mongodb.svc.cluster.local:27017/?authSource=admin"
	dbName         = "appdb"
	collectionName = "records"
	serverPort     = "8080"
)

func main() {
	server.StartServer(serverPort, mongoURI, dbName)
}
