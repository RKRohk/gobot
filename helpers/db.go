package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var prod = os.Getenv("PROD")

//InitDb Initializes the database
func newClient() (*mongo.Client, error) {
	if prod == "TRUE" {
		return mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	}
	return mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
}

//SlapString is this
type SlapString struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Text string             `bson:"text,omitempty"`
}

//SlapSticker denotes a sticker in the database
type SlapSticker struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	FileID string             `bson:"fileid,omitempty"`
}

//AddSlapToDB adds a slap string format to the mongodb database
func AddSlapToDB(slapString string) {
	client, err := newClient()
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
		log.Panic(err)
	} else {
		log.Println(res)
	}
}

//GetAllSlapStrings returns a list of all the slap documents
func GetAllSlapStrings() ([]SlapString, error) {
	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//Reference to the collection in the database
	slapCollection := client.Database("bot").Collection("slaps")

	documentCursor, err := slapCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Panic(err)
		return make([]SlapString, 0), err
	}

	//Array of documents returned from db
	var documentArr []SlapString

	//Decoding all documents to struct
	if err := documentCursor.All(ctx, &documentArr); err != nil {
		log.Panic(err)
		return make([]SlapString, 0), nil
	}

	//Printing for my reference
	log.Println(documentArr)

	defer documentCursor.Close(ctx)

	return documentArr, nil
}

//GetSlapStrings gets a random slap string from document
func GetSlapStrings() (string, error) {
	client, err := newClient()
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

	defer slapArr.Close(ctx)
	if err != nil {
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

//DeleteSlapFromDb removes a slap string with a given ID from the database
func DeleteSlapFromDb(documentID string) error {
	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	slapCollection := client.Database("bot").Collection("slaps")
	fmt.Println(documentID)
	objectID, _ := primitive.ObjectIDFromHex(documentID)
	fmt.Println(objectID)
	res, err := slapCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: objectID}})

	if err == nil {
		log.Println(res)
	}
	return err
}

//AddSlapStickerToDb Adds a slap string to database
func AddSlapStickerToDb(fileID string) error {
	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	slapCollection := client.Database("bot").Collection("slapsstickers")
	var newslap = SlapSticker{FileID: fileID}
	res, err := slapCollection.InsertOne(ctx, &newslap)
	if err != nil {
		log.Panic(err)
		return err
	}
	log.Println(res)

	return nil
}

//GetSlapStickers gets a random slap sticker from document
func GetSlapStickers() (string, error) {
	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	slapCollection := client.Database("bot").Collection("slapsstickers")

	slapArr, err := slapCollection.Aggregate(ctx, mongo.Pipeline{
		{primitive.E{Key: "$sample", Value: bson.D{primitive.E{Key: "size", Value: 1}}}},
	})

	defer slapArr.Close(ctx)
	if err != nil {
		log.Panic(err)
		return "", err
	}

	//Slap String from Database
	var result []SlapSticker

	//Moving cursour
	error := slapArr.All(ctx, &result)

	if error != nil {
		log.Panic(error)
		return "", error
	}
	log.Println(result)

	return result[0].FileID, nil

}
