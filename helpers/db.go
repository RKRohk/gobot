package helpers

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var err error

//InitDb Initializes the database
func InitDb() {
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB\nExiting....")
	}
}

//SlapString is this
type SlapString struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Text string             `bson:"text,omitempty"`
}

//AddSlapToDB adds a slap string format to the mongodb database
func AddSlapToDB(slapString string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	slapCollection := client.Database("bot").Collection("slaps")
	var newslap = SlapString{Text: slapString}
	log.Println(newslap)
	res, err := slapCollection.InsertOne(ctx, &newslap)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(res)
	}
}

//GetSlapStrings gets a random slap string from document
func GetSlapStrings() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	slapCollection := client.Database("bot").Collection("slaps")

	slapArr, err := slapCollection.Aggregate(ctx, mongo.Pipeline{
		{primitive.E{Key: "$sample", Value: bson.D{primitive.E{Key: "size", Value: 1}}}},
	})

	if err != nil {
		log.Panic(err)
		return "", err
	}

	//Slap String from Database
	var result SlapString

	//Moving cursour
	slapArr.Next(ctx)
	error := slapArr.Decode(&result)
	log.Printf("Decode function has been run")

	if error != nil {
		log.Panic(error)
		return "", error
	}
	log.Println(result)

	return result.Text, nil

}
