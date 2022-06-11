package doctorrepo

import (
	"context"
	"time"

	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorAppointmentSlotSetRepository struct {
	colln *mongo.Collection
}

func (r *DoctorAppointmentSlotSetRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *DoctorAppointmentSlotSetRepository) CreateOne(doc *doctormodel.AppointmentSlotSetEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)

	return err
}
