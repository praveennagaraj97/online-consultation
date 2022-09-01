package consultationrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	consultationdto "github.com/praveennagaraj97/online-consultation/dto/consultation"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConsultationRepository struct {
	colln         *mongo.Collection
	imageBasePath string
}

func (r *ConsultationRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	r.imageBasePath = env.GetEnvVariable("S3_ACCESS_BASEURL")
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
	sortOpts *bson.D,
	filterOpts *map[string]primitive.M,
	keySetSortby string,
) ([]consultationmodel.ConsultationEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var pipeline mongo.Pipeline = make(mongo.Pipeline, 0)
	limitPipe := bson.D{{Key: "$limit", Value: pgnOpt.PerPage}}

	if len(*filterOpts) != 0 {
		matchPipe := bson.D{{Key: "$match", Value: *filterOpts}}
		pipeline = append(pipeline, matchPipe)
	}

	if sortOpts != nil {

		sortPipe := bson.D{{Key: "$sort", Value: sortOpts}}
		pipeline = append(pipeline, sortPipe)

	} else {
		sortPipe := bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}}
		pipeline = append(pipeline, sortPipe)
	}

	if pgnOpt.PaginateId != nil {
		// Key Set ID
		keySetPipe := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: bson.M{keySetSortby: pgnOpt.PaginateId}}}}}
		pipeline = append(pipeline, keySetPipe)
	} else {
		skipPipe := bson.D{{Key: "$skip", Value: (pgnOpt.PageNum - 1) * pgnOpt.PerPage}}
		pipeline = append(pipeline, skipPipe)
	}

	// Map Image fields
	addFieldsPipe := bson.D{{Key: "$addFields", Value: bson.M{"icon.image_src": bson.M{"$concat": bson.A{r.imageBasePath, "/", "$icon.original_image_path"}},
		"icon.blur_data_url": bson.M{"$concat": bson.A{r.imageBasePath, "/", "$icon.blur_image_path"}}},
	}}

	pipeline = append(pipeline, limitPipe, addFieldsPipe)
	cur, err := r.colln.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("Couldn't find any results")
	}

	var result []consultationmodel.ConsultationEntity

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

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

func (r *ConsultationRepository) FindByType(consType consultationmodel.ConsultationType) (*consultationmodel.ConsultationEntity, error) {
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

func (r *ConsultationRepository) FindById(id *primitive.ObjectID) (*consultationmodel.ConsultationEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"_id": id})

	if res.Err() != nil {
		return nil, errors.New("Couldn't find any results")
	}

	var result consultationmodel.ConsultationEntity

	if err := res.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ConsultationRepository) UpdateById(id *primitive.ObjectID, payload *consultationdto.EditConsultationDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": payload})
	return err
}

func (r *ConsultationRepository) DeleteById(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
