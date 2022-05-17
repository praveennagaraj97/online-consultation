package userrepository

import (
	"context"
	"errors"
	"time"

	userdto "github.com/praveennagaraj97/online-consultation/dto"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

// Method to initialize user repository
func (r *UserRepository) InitializeRepository(colln *mongo.Collection) {
	r.collection = colln

	utils.CreateIndex(colln, bson.D{
		{Key: "phone_number.number", Value: 1},
		{Key: "phone_number.code", Value: 1}},
		"Phone", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "email", Value: 1}},
		"Email", true)
}

func (r *UserRepository) CreateUser(payload *userdto.RegisterDTO) (*usermodel.UserEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Check if already exists
	exists := r.checkIfUserExistsWithEmailOrPhone(payload.Email, payload.PhoneNumber, payload.PhoneNumber)

	if exists {
		return nil, errors.New("User with given credentials already exist")
	}

	document := &usermodel.UserEntity{
		ID:    primitive.NewObjectID(),
		Name:  payload.Name,
		Email: payload.Email,
		PhoneNumber: interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		},
		DateOfBirth: payload.DOB,
		Gender:      payload.Gender,
	}

	_, err := r.collection.InsertOne(ctx, document)

	if err != nil {
		return nil, err
	}

	return document, nil

}

func (r *UserRepository) UpdateById(id *primitive.ObjectID, payload *userdto.UpdateUserDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: payload}})
	return err

}

func (r *UserRepository) checkIfUserExistsWithEmailOrPhone(email, number, code string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$or": bson.A{
		bson.M{"email": email},
		bson.M{"$and": bson.A{bson.M{"phone_number.number": number, "phone_number.code": code}}},
	}}

	count, err := r.collection.CountDocuments(ctx, filter)

	if err != nil {
		return false
	}

	return count > 0

}

func (r *UserRepository) FindByPhoneNumber(number, code string) (*usermodel.UserEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"phone_number.number": number, "phone_number.code": code}}}

	doc := r.collection.FindOne(ctx, filter)

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

	doc := r.collection.FindOne(ctx, filter)

	if doc.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var payload usermodel.UserEntity

	if err := doc.Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil

}
