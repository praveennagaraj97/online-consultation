package userrepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	userdto "github.com/praveennagaraj97/online-consultation/dto"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRelativesRepository struct {
	colln *mongo.Collection
}

// Method to initialize user repository
func (r *UserRelativesRepository) InitializeRepository(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{
		{Key: "phone.number", Value: 1},
		{Key: "phone.code", Value: 1}},
		"Phone", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "email", Value: 1}},
		"Email", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "parent_id", Value: 1}},
		"Parent ID", false)
}

func (r *UserRelativesRepository) CreateOne(payload *userdto.AddOrEditRelativeDTO) (*usermodel.RelativeEntity, error) {

	if exists := r.checkIfRelativeExist(payload.Email,
		interfaces.PhoneType{Code: payload.PhoneCode, Number: payload.PhoneNumber}); exists {
		return nil, errors.New("Relative account with given credentials already exist")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	doc := &usermodel.RelativeEntity{
		ID:    primitive.NewObjectID(),
		Name:  payload.Name,
		Email: payload.Email,
		Phone: interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		},
		DateOfBirth: payload.DateOfBirth,
		Gender:      payload.Gender,
		Relation:    payload.Relation,
		ParentId:    payload.ParentId,
	}

	fmt.Println(doc)

	if _, err := r.colln.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil

}

func (r *UserRelativesRepository) checkIfRelativeExist(email string, phone interfaces.PhoneType) bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$or": bson.A{
		bson.M{"email": email},
		bson.M{"$and": bson.A{bson.M{"phone_number.number": phone.Number, "phone_number.code": phone.Code}}},
	}}

	count, err := r.colln.CountDocuments(ctx, filter)

	if err != nil {
		return false
	}

	return count > 0
}

func (r *UserRelativesRepository) FindAll(
	pgnOpt *api.PaginationOptions,
	parentId *primitive.ObjectID,
	keySetSortby string) (*[]usermodel.RelativeEntity, error) {

	opt := &options.FindOptions{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filters := map[string]bson.M{
		"parent_id": {"$eq": parentId},
	}

	if pgnOpt != nil {
		opt.Limit = options.Find().SetLimit(int64(pgnOpt.PerPage)).Limit
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	if pgnOpt.PaginateId != nil {
		filters["_id"] = bson.M{keySetSortby: pgnOpt.PaginateId}
	} else {
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	cur, err := r.colln.Find(ctx, filters, opt)
	if err != nil {
		return nil, err
	}

	var results []usermodel.RelativeEntity
	// check for errors in the conversion
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	return &results, nil

}
