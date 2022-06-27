package appointmentrepository

import (
	"context"
	"time"

	appointmentmodel "github.com/praveennagaraj97/online-consultation/models/appointment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentRepository struct {
	colln *mongo.Collection
}

func (r *AppointmentRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *AppointmentRepository) Create(doc *appointmentmodel.AppointmentEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)

	return err

}

func (r *AppointmentRepository) DeleteById(userId, id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"$and": bson.A{bson.M{"user_id": userId}, bson.M{"_id": id}}})

	return err

}
