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
		bson.D{{Key: "is_available", Value: 1}}, "AvailabilityIndex", false)

}

func (r *AppointmentSlotsRepository) CreateMany(docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := r.colln.InsertMany(ctx, docs)

	var e mongo.BulkWriteException

	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return errors.New("document with given slot already exist")
			}
		}
	}

	return err
}

func (r *AppointmentSlotsRepository) FindByIdAndDoctorID(docId, id *primitive.ObjectID) (*appointmentslotmodel.AppointmentSlotEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"$and": bson.A{bson.M{"doctor_id": docId}, bson.M{"_id": id}}})

	if cur.Err() != nil {
		return nil, errors.New("couldn't find appointment slot for given Id")
	}

	var result appointmentslotmodel.AppointmentSlotEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *AppointmentSlotsRepository) FindById(id *primitive.ObjectID) (*appointmentslotmodel.AppointmentSlotEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"_id": id})

	if cur.Err() != nil {
		return nil, errors.New("couldn't find appointment slot for given Id")
	}

	var result appointmentslotmodel.AppointmentSlotEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *AppointmentSlotsRepository) FindByDoctorIdAndDate(docId *primitive.ObjectID, date *primitive.DateTime) ([]appointmentslotmodel.AppointmentSlotEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur, err := r.colln.Find(ctx, bson.M{"$and": bson.A{bson.M{"doctor_id": docId}, bson.M{"date": date}}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var results []appointmentslotmodel.AppointmentSlotEntity

	if err := cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AppointmentSlotsRepository) UpdateSlotAvailability(id *primitive.ObjectID,
	isAvailable bool,
	reason appointmentslotmodel.BookingStatusMessages) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if !isAvailable && len(reason) == 0 {
		return errors.New("reason cannot be empty")
	}

	updateDoc := map[string]interface{}{}

	updateDoc["is_available"] = isAvailable
	updateDoc["reason_of_unavailablity"] = reason

	// Set release time if reason is payment processing
	if reason == appointmentslotmodel.PaymentProcessing {
		updateDoc["slot_release_at"] = time.Now().UTC().Add(time.Minute * 5)
	} else {
		updateDoc["slot_release_at"] = nil
	}

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": updateDoc})

	return err

}
