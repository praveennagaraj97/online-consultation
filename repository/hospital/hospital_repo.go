package hospitalrepo

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	hospitaldto "github.com/praveennagaraj97/online-consultation/dto/hospital"
	hospitalmodel "github.com/praveennagaraj97/online-consultation/models/hospital"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HospitalRepository struct {
	colln *mongo.Collection
}

func (r *HospitalRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "name", Value: 1}}, "HospitalNameIndex", true)
}

func (r *HospitalRepository) CreateOne(doc *hospitalmodel.HospitalEntity) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)

	return err
}

func (r *HospitalRepository) CheckIfExistByName(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, err := r.colln.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false
	}

	return count > 0
}

func (r *HospitalRepository) FindById(id *primitive.ObjectID) (*hospitalmodel.HospitalEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"_id": id})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any matching hosiptal with given ID")
	}

	var result hospitalmodel.HospitalEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *HospitalRepository) UpdateById(id *primitive.ObjectID, payload *hospitaldto.EditHospitalDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": payload})

	return err

}

func (r *HospitalRepository) DeleteById(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *HospitalRepository) FindAll(pgOpts *api.PaginationOptions,
	fltrOpts *map[string]primitive.M,
	srtOpts *map[string]int8,
	keySoryBy string,
) ([]hospitalmodel.HospitalEntity, error) {

	opt := options.FindOptions{}

	opt.Limit = options.Find().SetLimit(int64(pgOpts.PerPage)).Limit

	if srtOpts != nil {
		opt.Sort = options.Find().SetSort(srtOpts).Sort
	}

	var filters map[string]primitive.M = make(map[string]primitive.M)

	if fltrOpts != nil {
		for key, val := range *fltrOpts {
			filters[key] = val
		}
	}

	if pgOpts.PaginateId != nil {
		filters["_id"] = bson.M{keySoryBy: pgOpts.PaginateId}
	} else if pgOpts != nil {
		opt.Skip = options.Find().SetSkip(int64(pgOpts.PerPage) * int64(pgOpts.PageNum-1)).Skip
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur, err := r.colln.Find(ctx, filters, &opt)

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	var results []hospitalmodel.HospitalEntity

	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *HospitalRepository) GetDocumentsCount(filter *map[string]primitive.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.colln.CountDocuments(ctx, filter)

}
