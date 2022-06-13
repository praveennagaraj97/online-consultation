package doctorrepo

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *DoctorAppointmentSlotSetRepository) FindAll(doctorId *primitive.ObjectID,
	pgOpts *api.PaginationOptions,
	srtOpts *map[string]int8,
	keySortBy string) ([]doctormodel.AppointmentSlotSetEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opt := options.FindOptions{}

	opt.Limit = options.Find().SetLimit(int64(pgOpts.PerPage)).Limit

	if srtOpts != nil {
		opt.Sort = options.Find().SetSort(srtOpts).Sort
	}

	var filters map[string]primitive.M = make(map[string]primitive.M)
	filters["doctor_id"] = bson.M{"$eq": doctorId}

	if pgOpts.PaginateId != nil {
		filters["_id"] = bson.M{keySortBy: pgOpts.PaginateId}
	} else if pgOpts != nil {
		opt.Skip = options.Find().SetSkip(int64(pgOpts.PerPage) * int64(pgOpts.PageNum-1)).Skip
	}

	cur, err := r.colln.Find(ctx, filters, &opt)

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	var results []doctormodel.AppointmentSlotSetEntity

	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
