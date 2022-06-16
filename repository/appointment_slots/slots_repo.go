package appointmentslotsrepo

import (
	"context"
	"time"

	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentSlotsRepository struct {
	colln *mongo.Collection
}

func (r *AppointmentSlotsRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln,
		bson.D{{Key: "doctor_id", Value: 1}, {Key: "date", Value: 1}, {Key: "start_time", Value: 1}}, "UniqueSlotIndex", true)

	utils.CreateIndex(colln,
		bson.D{{Key: "doctor_id", Value: 1}, {Key: "date", Value: 1}}, "DoctorAndDateIndex", false)
}

func (r *AppointmentSlotsRepository) CreateMany(docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := r.colln.InsertMany(ctx, docs)

	return err
}
