package main

import (
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/db"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	twiliopkg "github.com/praveennagaraj97/online-consultation/pkg/sms/twilio"
	"github.com/praveennagaraj97/online-consultation/router"
)

func main() {

	// Initialize application with configs
	app := &app.ApplicationConfig{
		Port:        env.GetEnvVariable("PORT"),
		Environment: env.GetEnvVariable("ENVIRONMENT"),
		Domain:      env.GetEnvVariable("DOMAIN"),
		DB: struct {
			MONGO_URI    string
			MONGO_DBNAME string
		}{
			MONGO_URI:    env.GetEnvVariable("MONGO_URI"),
			MONGO_DBNAME: "online-consultation",
		},
	}

	// Initialize Database
	app.MongoClient = db.InitializeMongoDatabase(&app.DB.MONGO_URI)

	// Initialize Twilio Package
	twiliopkg.Initialize()

	// Start the server
	r := router.Router{}
	r.ListenAndServe(app)
}
