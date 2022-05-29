package consultationrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConsultationRepository struct {
	colln *mongo.Collection
}

func (r *ConsultationRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *ConsultationRepository) CreateOne(payload *consultationmodel.ConsultationEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, payload)

	if err != nil {
		return err
	}

	return nil
}

func (r *ConsultationRepository) CheckIfConsultationTypeExists(consultationType string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count, err := r.colln.CountDocuments(ctx, bson.M{"type": consultationType})
	if err != nil {
		return false
	}

	return count > 0
}

func (r *ConsultationRepository) FindAll(
	pgnOpt *api.PaginationOptions,
	sortOpts *map[string]int8,
	filterOpts *map[string]primitive.M,
	keySetSortby string,
) ([]consultationmodel.ConsultationEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opt := &options.FindOptions{}

	filters := map[string]bson.M{}

	if pgnOpt != nil {
		opt.Limit = options.Find().SetLimit(int64(pgnOpt.PerPage)).Limit
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	if sortOpts != nil {
		opt.Sort = options.Find().SetSort(sortOpts).Sort
	} else {
		opt.Sort = options.Find().SetSort(bson.M{"created_at": -1}).Sort
	}

	if filterOpts != nil {
		for key, value := range *filterOpts {
			filters[key] = value
		}
	}

	if pgnOpt.PaginateId != nil {
		filters["_id"] = bson.M{keySetSortby: pgnOpt.PaginateId}
	} else {
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	cur, err := r.colln.Find(ctx, filters, opt)
	if err != nil {
		return nil, errors.New("Couldn't find any results")
	}

	var result []consultationmodel.ConsultationEntity

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ConsultationRepository) GetDocumentsCount(filterOpts *map[string]primitive.M) (int64, error) {

	filters := map[string]bson.M{}

	if filterOpts != nil {
		for key, value := range *filterOpts {
			filters[key] = value
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.colln.CountDocuments(ctx, filters)

}

func (r *ConsultationRepository) FindByType(consType string) (*consultationmodel.ConsultationEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"type": consType})

	if res.Err() != nil {
		return nil, errors.New("Couldn't find any results")
	}

	var result consultationmodel.ConsultationEntity

	if err := res.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
