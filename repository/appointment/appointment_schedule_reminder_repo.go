package appointmentrepository

import (
	"context"
	"time"

	appointmentmodel "github.com/praveennagaraj97/online-consultation/models/appointment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentScheduleReminderRepository struct {
	colln *mongo.Collection
}

func (r *AppointmentScheduleReminderRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *AppointmentScheduleReminderRepository) Create(doc *appointmentmodel.AppointmentScheduleTaskEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)

	return err

}

func (r *AppointmentScheduleReminderRepository) DeleteById(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"_id": id})

	return err

}

// For Internal Use | Find AllReminders By Invoke Time
func (r *AppointmentScheduleReminderRepository) FindAllByInvokeTime(invoke_time *primitive.DateTime) ([]appointmentmodel.AppointmentScheduleTaskEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur, err := r.colln.Find(ctx, bson.M{"invoke_time": invoke_time})
	if err != nil {
		return nil, err
	}

	var results []appointmentmodel.AppointmentScheduleTaskEntity

	if err := cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// For Internal Use | Get All By date
func (r *AppointmentScheduleReminderRepository) FindAllByDate(date *primitive.DateTime) ([]appointmentmodel.AppointmentScheduleTaskEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur, err := r.colln.Find(ctx, bson.M{"date": date})
	if err != nil {
		return nil, err
	}

	var results []appointmentmodel.AppointmentScheduleTaskEntity

	if err := cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
