package global

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds database connection
var DB mongo.Database

func connectToDB() {

	ctx, cancel := NewDBContext(10 * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		fmt.Println("Error connecting to db", err.Error())
	}
	DB = *client.Database(dbname)
}

// NewDBContext return a mew Context according to app performance
func NewDBContext(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t*performance/100)
}

// ConnectToTestDB connects to a test Database
func ConnectToTestDB() {
	ctx, cancel := NewDBContext(10 * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		fmt.Println("Error connecting to db", err.Error())
	}
	DB = *client.Database(dbname + "_test")
	ctx, cancel = NewDBContext(30 * time.Second)
	defer cancel()
	collections, _ := DB.ListCollectionNames(ctx, bson.M{})
	for _, collection := range collections {
		ctx, cancel = NewDBContext(30 * time.Second)
		defer cancel()
		DB.Collection(collection).Drop(ctx)
	}
}
