package languagerepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	languagedto "github.com/praveennagaraj97/online-consultation/dto/language"
	languagesmodel "github.com/praveennagaraj97/online-consultation/models/languages"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LanguageRepository struct {
	colln *mongo.Collection
}

func (r *LanguageRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "name", Value: 1}}, "Unique Name", true)
	utils.CreateIndex(colln, bson.D{{Key: "locale_name", Value: 1}}, "Unique locale name", true)
}

func (r *LanguageRepository) CreateOne(payload *languagedto.AddOrEditLanguageDTO) (*languagesmodel.LanguageEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc := languagesmodel.LanguageEntity{
		ID:         primitive.NewObjectID(),
		Name:       payload.Name,
		LocaleName: payload.LocaleName,
	}

	_, err := r.colln.InsertOne(ctx, doc)

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *LanguageRepository) CheckIfExistByName(name, localename string) bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.D{{Key: "$or", Value: bson.A{bson.M{"name": name}, bson.M{"locale_name": localename}}}}

	count, err := r.colln.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}

	return count > 0

}

func (r *LanguageRepository) Find(
	pgOpts *api.PaginationOptions,
	srtOpts *map[string]int8,
	fltrOpts *map[string]primitive.M,
	keySortBy string) ([]languagesmodel.LanguageEntity, error) {

	opt := options.FindOptions{}

	opt.Limit = options.Find().SetLimit(int64(pgOpts.PerPage)).Limit

	if srtOpts != nil {
		opt.Sort = options.Find().SetSort(srtOpts).Sort
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var filters map[string]primitive.M = make(map[string]primitive.M)
	if fltrOpts != nil {
		for key, value := range *fltrOpts {
			filters[key] = value
		}
	}

	if pgOpts.PaginateId != nil {
		filters["_id"] = bson.M{keySortBy: pgOpts.PaginateId}
	} else if pgOpts != nil {
		opt.Skip = options.Find().SetSkip(int64(pgOpts.PerPage) * int64(pgOpts.PageNum-1)).Skip
	}

	cur, err := r.colln.Find(ctx, filters, &opt)
	if err != nil {
		return nil, err
	}

	var result []languagesmodel.LanguageEntity

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	fmt.Println(result)

	return result, nil

}

func (r *LanguageRepository) GetDocumentsCount(filter *map[string]primitive.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.colln.CountDocuments(ctx, filter)
}

func (r *LanguageRepository) FindById(id *primitive.ObjectID) (*languagesmodel.LanguageEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"_id": id})

	if res.Err() != nil {
		return nil, errors.New("Couldn't find any language matching given Id")
	}

	var result languagesmodel.LanguageEntity

	if err := res.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *LanguageRepository) DeleteById(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"_id": id})

	return err

}
