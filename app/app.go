package app

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationConfig struct {
	Port        string
	Environment string
	Domain      string
	DB          struct {
		MONGO_URI    string
		MONGO_DBNAME string
	}
	MongoClient *mongo.Client
}
