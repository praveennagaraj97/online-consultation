package doctorrepo

import (
	"context"
	"errors"
	"time"

	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorAppointmentSlotSetRepository struct {
	colln *mongo.Collection
}

func (r *DoctorAppointmentSlotSetRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "doctor_id", Value: 1}}, "DocterIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "_id", Value: 1}, {Key: "doctor_id", Value: 1}}, "DocIdAndDocterIndex", true)
}

func (r *DoctorAppointmentSlotSetRepository) CreateOne(doc *doctormodel.AppointmentSlotSetEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)

	return err
}

func (r *DoctorAppointmentSlotSetRepository) ExistingSetsCount(doctorId *primitive.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.colln.CountDocuments(ctx, bson.M{"doctor_id": doctorId})

}

func (r *DoctorAppointmentSlotSetRepository) FindById(doctorId, id *primitive.ObjectID) (*doctormodel.AppointmentSlotSetEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"$and": bson.A{bson.M{"_id": id}, bson.M{"doctor_id": doctorId}}})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any slot set.")
	}

	var result doctormodel.AppointmentSlotSetEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *DoctorAppointmentSlotSetRepository) UpdateById(doctorId, id *primitive.ObjectID, payload *doctordto.UpdateAppointmentSlotSetDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": payload})

	return err
}

func (r *DoctorAppointmentSlotSetRepository) DeleteById(doctorId, id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"$and": bson.A{bson.M{"_id": id}, bson.M{"doctor_id": doctorId}}})

	return err
}
