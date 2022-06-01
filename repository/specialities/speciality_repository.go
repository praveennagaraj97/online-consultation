package specialityrepository

import (
	"context"
	"errors"
	"time"

	specialitydto "github.com/praveennagaraj97/online-consultation/dto/speciality"
	specialitymodel "github.com/praveennagaraj97/online-consultation/models/speciality"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpecialitysRepository struct {
	colln *mongo.Collection
}

func (r *SpecialitysRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

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
