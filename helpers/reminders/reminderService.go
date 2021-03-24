package reminders

import (
	"context"
	"log"
	"time"

	"github.com/rkrohk/gobot/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reminders Reminders

//InitializeReminders fetches the closest reminder from the database
func InitializeReminders() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	cursor, err := database.Client.Database("bot").Collection("reminders").Aggregate(ctx, mongo.Pipeline{
		{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "date", Value: 1}}}},
		{primitive.E{Key: "$limit", Value: 1}}})

	if err != nil {
		log.Println("Error fetching reminder from database", err)
		return
	}

	err = cursor.All(ctx, &reminders)
	if err != nil {
		log.Println("Error", err)
	}
	log.Println(reminders)

}
