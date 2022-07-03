package userrepository

import (
	"context"
	"errors"
	"time"

	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	colln *mongo.Collection
}

// Method to initialize user repository
func (r *UserRepository) InitializeRepository(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{
		{Key: "phone_number.number", Value: 1},
		{Key: "phone_number.code", Value: 1}},
		"PhoneIndex", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "email", Value: 1}},
		"EmailIndex", true)
}

func (r *UserRepository) CreateUser(doc *usermodel.UserEntity) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)

	return err

}

func (r *UserRepository) UpdateById(id *primitive.ObjectID, payload *userdto.UpdateUserDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: payload}})
	return err

}

func (r *UserRepository) UpdateRefreshToken(id *primitive.ObjectID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: bson.M{"refresh_token": token}}})
	return err

}

func (r *UserRepository) CheckIfUserExistsWithEmailOrPhone(email, number, code string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$or": bson.A{
		bson.M{"email": email},
		bson.M{"$and": bson.A{bson.M{"phone_number.number": number, "phone_number.code": code}}},
	}}

	count, err := r.colln.CountDocuments(ctx, filter)

	if err != nil {
		return false
	}

	return count > 0

}

func (r *UserRepository) FindByPhoneNumber(number, code string) (*usermodel.UserEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"phone_number.number": number, "phone_number.code": code}}}

	doc := r.colln.FindOne(ctx, filter)

	if doc.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var payload usermodel.UserEntity

	if err := doc.Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil

}

func (r *UserRepository) FindByEmail(email string) (*usermodel.UserEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"email": email}

	doc := r.colln.FindOne(ctx, filter)

	if doc.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var payload usermodel.UserEntity

	if err := doc.Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil

}

func (r *UserRepository) FindById(id *primitive.ObjectID) (*usermodel.UserEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc := r.colln.FindOne(ctx, bson.M{"_id": id})

	if doc.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var payload usermodel.UserEntity

	if err := doc.Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil

}
