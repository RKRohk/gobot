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

var (
	Client       *mongo.Client
	ctx          context.Context
	cancel       context.CancelFunc
	DATABASE_URI = os.Getenv("DATABASE")
)

func init() {
	log.Println("Connecting to database", DATABASE_URI)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var err error

	Client, err = mongo.NewClient(options.Client().ApplyURI(DATABASE_URI))

	for err != nil {
		log.Printf("Error initializing mongodb client %v", err)
		Client, err = mongo.NewClient(options.Client().ApplyURI(DATABASE_URI))
		fmt.Println("Trying in 5 seconds")
		time.Sleep(5 * time.Second)
	}

	err = Client.Connect(ctx)
	for err != nil {
		err = Client.Connect(ctx)
		fmt.Printf("Error connecting to mongodb %v", err)
		fmt.Println("Trying in 5 seconds")
		time.Sleep(5 * time.Second)
	}

	fmt.Println("Connected to mongodb")
}

func DisconnectDatabase() {
	log.Println("Disconnecting from database")
	if err := Client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}

	defer cancel()
}
