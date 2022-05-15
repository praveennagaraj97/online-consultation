package onetimepasswordrepository

import (
	"context"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	otpmodel "github.com/praveennagaraj97/online-consultation/models/otp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OneTimePasswordRepository struct {
	col *mongo.Collection
}

func (r *OneTimePasswordRepository) InitializeRepository(colln *mongo.Collection) {
	if r.col == nil {
		r.col = colln
	}

}

func (r *OneTimePasswordRepository) CreateOne(phoneNumber *interfaces.PhoneType, verifyCode *string) (*otpmodel.OneTimePasswordEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc := otpmodel.OneTimePasswordEntity{}
	doc.Init(verifyCode, phoneNumber)

	if _, err := r.col.InsertOne(ctx, &doc); err != nil {
		return nil, err
	}

	return &doc, nil

}

func (r *OneTimePasswordRepository) DeleteById(id *primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := r.col.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

func (r *OneTimePasswordRepository) FindById(id *primitive.ObjectID) (*otpmodel.OneTimePasswordEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}})

	if result.Err() != nil {
		return nil, result.Err()
	}

	var payload otpmodel.OneTimePasswordEntity
	result.Decode(&payload)

	return &payload, nil
}

func (r *OneTimePasswordRepository) UpdateAttemptsCount(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := r.col.UpdateByID(ctx, id, bson.M{"$inc": bson.M{"attempts": 1}}); err != nil {
		return err
	}

	return nil
}

func (r *OneTimePasswordRepository) UpdateStatus(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := r.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"verified": true}}); err != nil {
		return err
	}

	return nil
}
