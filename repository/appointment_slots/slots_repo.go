package appointmentslotsrepo

import (
	"context"
	"errors"
	"time"

	appointmentslotmodel "github.com/praveennagaraj97/online-consultation/models/appointment_slot"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentSlotsRepository struct {
	colln *mongo.Collection
}

func (r *AppointmentSlotsRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln,
		bson.D{{Key: "doctor_id", Value: 1}, {Key: "date", Value: 1}, {Key: "start", Value: 1}}, "UniqueSlotIndex", true)

	utils.CreateIndex(colln,
		bson.D{{Key: "doctor_id", Value: 1}, {Key: "date", Value: 1}}, "DoctorAndDateIndex", false)
}

func (r *AppointmentSlotsRepository) CreateMany(docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := r.colln.InsertMany(ctx, docs)

	var e mongo.BulkWriteException

	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return errors.New("Document with given slot already exist")
			}
		}
	}

	return err
}

func (r *AppointmentSlotsRepository) FindById(docId, id *primitive.ObjectID) (*appointmentslotmodel.AppointmentSlotEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"$and": bson.A{bson.M{"doctor_id": docId}, bson.M{"_id": id}}})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find appointment slot for given Id")
	}

	var result appointmentslotmodel.AppointmentSlotEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}
