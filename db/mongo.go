package db

import (
	"context"
	"log"
	"time"

	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeMongoDatabase(uri *string) *mongo.Client {
	mgoClientOption := options.Client().ApplyURI(*uri)

	client, err := mongo.NewClient(mgoClientOption)

	if err != nil {
		logger.ErrorLogFatal(err)
		return nil
	}

	// provide connection timout limit for connection to get established.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	logger.PrintLog("Connected to MongoDB 🗄")

	return client

}

func OpenCollection(mgoClient *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return mgoClient.Database(dbName).Collection(collectionName)
}
