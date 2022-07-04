package app

import (
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/scheduler"

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
	EmailClient *mailer.Mailer
	AwsUtils    *awspkg.AWSConfiguration
	Scheduler   *scheduler.Scheduler
}
