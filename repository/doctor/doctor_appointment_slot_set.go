package doctorrepo

import (
	"context"
	"time"

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
