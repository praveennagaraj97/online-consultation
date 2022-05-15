package utils

import (
	"context"
	"time"

	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndex - creates an index for a collection
func CreateIndex(collection *mongo.Collection, keys bson.D, indexName string, unique bool) bool {

	var indexOptions *options.IndexOptions = &options.IndexOptions{}

	indexOptions.Unique = &unique
	indexOptions.Name = options.Index().SetName(indexName).Name

	// 1. Field key
	mod := mongo.IndexModel{
		Keys:    keys,
		Options: indexOptions,
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {

		logger.ErrorLogFatal(err.Error())
		return false
	}

	return true
}
