package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client       *mongo.Client
	ctx          context.Context
	cancel       context.CancelFunc
	DATABASE_URI = os.Getenv("DATABASE")
)

func init() {
	log.Println("Connecting to database",DATABASE_URI)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var err error

	Client, err = mongo.NewClient(options.Client().ApplyURI(DATABASE_URI))

	if err != nil {
		log.Panic("Unable to connect to database, exiting ....")
	}

	err = Client.Connect(ctx)
	if err != nil {
		log.Panic("Error connecting to DB")
	}
}

func DisconnectDatabase() {
	log.Println("Disconnecting from database")
	if err := Client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}

	defer cancel()
}
