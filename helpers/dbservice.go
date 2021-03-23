package helpers

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var ctx context.Context
var cancel context.CancelFunc

func init() {
	log.Println("Connecting to database")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	var err error
	if prod == "TRUE" {
		Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	} else {
		Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	}

	if err != nil {
		log.Panic("Unable to connect to database, exiting ....")
	}

	Client.Connect(ctx)
}

func DisconnectDatabase() {

	if err := Client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}

	defer cancel()
}
