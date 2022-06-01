package specialityrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	specialitydto "github.com/praveennagaraj97/online-consultation/dto/speciality"
	specialitymodel "github.com/praveennagaraj97/online-consultation/models/speciality"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpecialitysRepository struct {
	colln         *mongo.Collection
	imageBasePath string
}

func (r *SpecialitysRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
	r.imageBasePath = env.GetEnvVariable("S3_ACCESS_BASEURL")

	utils.CreateIndex(colln, bson.D{{Key: "title", Value: 1}}, "Unique title", true)
	utils.CreateIndex(colln, bson.D{{Key: "slug", Value: 1}}, "Unique slug", true)
}

func (r *SpecialitysRepository) CreateOne(doc *specialitymodel.SpecialityEntity) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := r.colln.InsertOne(ctx, doc); err != nil {
		return err
	}

	return nil

}

func (r *SpecialitysRepository) CheckIfExists(title, slug string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	regex := primitive.Regex{Pattern: title, Options: "i"}

	count, _ := r.colln.CountDocuments(ctx, bson.M{"$or": bson.A{bson.M{"slug": slug}, bson.M{"title": bson.M{"$regex": regex}}}})

	return count > 0
}

func (r *SpecialitysRepository) CheckIfExistsByTitle(title string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, _ := r.colln.CountDocuments(ctx, bson.M{"title": title})

	return count > 0
}

func (r *SpecialitysRepository) CheckIfExistsBySlug(slug string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, _ := r.colln.CountDocuments(ctx, bson.M{"slug": slug})

	return count > 0
}

func (r *SpecialitysRepository) FindById(id *primitive.ObjectID) (*specialitymodel.SpecialityEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"_id": id})

	if res.Err() != nil {
		return nil, errors.New("Couldn't find any speciality with given id")
	}

	var result specialitymodel.SpecialityEntity

	if err := res.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *SpecialitysRepository) DeleteById(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := r.colln.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

func (r *SpecialitysRepository) UpdateById(id *primitive.ObjectID, payload *specialitydto.EditSpecialityDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": payload}); err != nil {
		return err
	}

	return nil
}

func (r *SpecialitysRepository) FindBySlug(slug string) (*specialitymodel.SpecialityEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res := r.colln.FindOne(ctx, bson.M{"slug": slug})

	if res.Err() != nil {
		return nil, errors.New("Couldn't find any speciality with given slug")
	}

	var result specialitymodel.SpecialityEntity

	if err := res.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SpecialitysRepository) FindAll(
	pgnOpt *api.PaginationOptions,
	sortOpts *map[string]int8,
	filterOpts *map[string]primitive.M,
	keySetSortby string) ([]specialitymodel.SpecialityEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var pipeline mongo.Pipeline = make(mongo.Pipeline, 0)
	limitPipe := bson.D{{Key: "$limit", Value: pgnOpt.PerPage}}

	if filterOpts != nil {
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
	addFieldsPipe := bson.D{{Key: "$addFields", Value: bson.M{"thumbnail.image_src": bson.M{"$concat": bson.A{r.imageBasePath, "/", "$thumbnail.original_image_path"}},
		"thumbnail.blur_data_url": bson.M{"$concat": bson.A{r.imageBasePath, "/", "$thumbnail.blur_image_path"}}},
	}}

	pipeline = append(pipeline, limitPipe, addFieldsPipe)
	cur, err := r.colln.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("Couldn't find any results")
	}

	var result []specialitymodel.SpecialityEntity

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil

}

func (r *SpecialitysRepository) GetDocumentsCount(filterOpts *map[string]primitive.M) (int64, error) {

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
