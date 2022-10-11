package appointmentrepository

import (
	"context"
	"errors"
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

// For internal use only.
func (r *AppointmentRepository) FindById(id *primitive.ObjectID) (*appointmentmodel.AppointmentEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"_id": id})

	if res.Err() != nil {
		return nil, errors.New("couldn't find any appointment details")
	}

	var result appointmentmodel.AppointmentEntity

	res.Decode(&result)

	return &result, nil

}

func (r *AppointmentRepository) FindByIdAndUserId(id, userId *primitive.ObjectID) (*appointmentmodel.AppointmentEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"$and": bson.A{bson.M{"_id": id}, bson.M{"user_id": userId}}})

	if res.Err() != nil {
		return nil, errors.New("couldn't find any appointment details")
	}

	var result appointmentmodel.AppointmentEntity

	res.Decode(&result)

	return &result, nil

}

func (r *AppointmentRepository) UpdateById(userId, id *primitive.ObjectID, status appointmentmodel.AppointmentStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.UpdateOne(ctx, bson.M{"$and": bson.A{bson.M{"_id": id}, bson.M{"user_id": userId}}}, bson.M{"$set": bson.M{
		"status": status,
	}})

	return err
}

func (r *AppointmentRepository) CheckIfAppointmentExistsByDoctorId(id *primitive.ObjectID) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"doctor_id": id, "status": "upcoming"}

	return r.colln.CountDocuments(ctx, filter)

}
