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
		{Key: "parent_id", Value: 1},
		{Key: "phone.number", Value: 1},
		{Key: "phone.code", Value: 1}},
		"Phone", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "parent_id", Value: 1},
		{Key: "email", Value: 1}},
		"Email", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "parent_id", Value: 1}},
		"Parent ID", false)

	utils.CreateIndex(colln, bson.D{
		{Key: "parent_id", Value: 1}, {Key: "_id", Value: 1}},
		"ParentId and Doc Id", false)
}

func (r *UserRelativesRepository) CreateOne(payload *userdto.AddOrEditRelativeDTO) (*usermodel.RelativeEntity, error) {

	if exists := r.checkIfRelativeExist(payload.Email,
		interfaces.PhoneType{Code: payload.PhoneCode, Number: payload.PhoneNumber}, &payload.ParentId); exists {
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

	if _, err := r.colln.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil

}

func (r *UserRelativesRepository) checkIfRelativeExist(email string, phone interfaces.PhoneType, parentId *primitive.ObjectID) bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	emailExists := bson.D{{Key: "$and", Value: bson.A{bson.M{"parent_id": parentId}, bson.M{"email": email}}}}
	phoneExists := bson.D{{Key: "$and", Value: bson.A{bson.M{"parent_id": parentId},
		bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}}}}

	filter := bson.M{"$or": bson.A{emailExists, phoneExists}}

	fmt.Println(filter)

	count, err := r.colln.CountDocuments(ctx, filter)

	fmt.Println("Count", count)

	if err != nil {
		return false
	}

	return count > 0
}

func (r *UserRelativesRepository) FindAll(
	pgnOpt *api.PaginationOptions,
	sortOpts *map[string]int8,
	filterOpts *map[string]primitive.M,
	keySetSortby string, parentId *primitive.ObjectID) ([]usermodel.RelativeEntity, error) {

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

	if sortOpts != nil {
		opt.Sort = options.Find().SetSort(sortOpts).Sort
	} else {
		opt.Sort = options.Find().SetSort(bson.M{"created_at": -1}).Sort
	}

	if filterOpts != nil {
		for key, value := range *filterOpts {
			filters[key] = value
		}
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

	return results, nil

}

func (r *UserRelativesRepository) GetDocumentsCount(parentId *primitive.ObjectID, filterOpts *map[string]primitive.M) (int64, error) {

	filters := map[string]bson.M{
		"parent_id": {"$eq": parentId},
	}

	if filterOpts != nil {
		for key, value := range *filterOpts {
			filters[key] = value
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.colln.CountDocuments(ctx, filters)
}

func (r *UserRelativesRepository) FindById(parentId *primitive.ObjectID, id *primitive.ObjectID) (*usermodel.RelativeEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.D{{Key: "$and", Value: bson.A{bson.M{"parent_id": parentId}, bson.M{"_id": id}}}})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any matching result for given id")
	}

	var data usermodel.RelativeEntity

	if err := cur.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UserRelativesRepository) UpdateByID(parentId, id *primitive.ObjectID, payload *userdto.AddOrEditRelativeDTO) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"_id": id, "parent_id": parentId}}}

	if _, err := r.colln.UpdateOne(ctx, filter, bson.M{"$set": payload}); err != nil {
		return err
	}

	return nil
}

func (r *UserRelativesRepository) DeleteByID(parentId, id *primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"_id": id, "parent_id": parentId}}}

	if _, err := r.colln.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
