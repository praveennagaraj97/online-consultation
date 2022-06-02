package languagerepo

import (
	"context"
	"time"

	languagedto "github.com/praveennagaraj97/online-consultation/dto/language"
	languagesmodel "github.com/praveennagaraj97/online-consultation/models/languages"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LanguageRepository struct {
	colln *mongo.Collection
}

func (r *LanguageRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "name", Value: 1}}, "Unique Name", true)
	utils.CreateIndex(colln, bson.D{{Key: "locale_name", Value: 1}}, "Unique locale name", true)
}

func (r *LanguageRepository) CreateOne(payload *languagedto.AddLanguageDTO) (*languagesmodel.LanguageEntity, error) {

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

func (r *LanguageRepository) ChechIfExistByName(name, localename string) bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.D{{Key: "$or", Value: bson.A{bson.M{"name": name}, bson.M{"locale_name": localename}}}}

	count, err := r.colln.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}

	return count > 0

}
