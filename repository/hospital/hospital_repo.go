package hospitalrepo

import (
	"context"
	"errors"
	"time"

	hospitaldto "github.com/praveennagaraj97/online-consultation/dto/hospital"
	hospitalmodel "github.com/praveennagaraj97/online-consultation/models/hospital"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HospitalRepository struct {
	colln *mongo.Collection
}

func (r *HospitalRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "name", Value: 1}}, "Hospital name", true)
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
