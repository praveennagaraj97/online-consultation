package specialityrepository

import (
	"context"
	"time"

	specialitymodel "github.com/praveennagaraj97/online-consultation/models/speciality"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpecialitysRepository struct {
	colln *mongo.Collection
}

func (r *SpecialitysRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "title", Value: 1}}, "Unique title", true)
}

func (r *SpecialitysRepository) CreateOne(doc *specialitymodel.SpecialityEntity) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := r.colln.InsertOne(ctx, doc); err != nil {
		return err
	}

	return nil

}

// func (r *SpecialitysRepository) FindByTitle(title string) (*cons) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	cur:= r.colln.FindOne(ctx, bson.D{{Key: "$eq", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: title, Options: "i"}}}}}}})

// 	if cur.Err() != nil{
// 		return false
// 	}

// 	return true

// }
