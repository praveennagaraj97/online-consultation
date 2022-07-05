package main

import (
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/db"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/pkg/scheduler"
	twiliopkg "github.com/praveennagaraj97/online-consultation/pkg/sms/twilio"
	stripepayment "github.com/praveennagaraj97/online-consultation/pkg/stripe"
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

	// Email Client
	app.EmailClient = &mailer.Mailer{}
	app.EmailClient.Initialize()

	// Initialize Database
	app.MongoClient = db.InitializeMongoDatabase(&app.DB.MONGO_URI)

	// AWS Utility Package
	app.AwsUtils = &awspkg.AWSConfiguration{}
	app.AwsUtils.Initialize()

	// Initialize Twilio SMS Package | Global instance is not available!
	twiliopkg.Initialize()

	// Stripe Package
	stripepayment.Initialize()

	// Scheduler Package
	app.Scheduler = &scheduler.Scheduler{}
	app.Scheduler.Initialize()
	defer app.Scheduler.Shutdown()

	// Start the server
	r := router.Router{}
	r.ListenAndServe(app)

}
