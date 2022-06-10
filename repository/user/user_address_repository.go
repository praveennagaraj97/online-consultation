package userrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDeliveryAddressRepository struct {
	colln *mongo.Collection
}

func (r *UserDeliveryAddressRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln,
		bson.D{{Key: "_id", Value: 1}, {Key: "user_id", Value: 1}}, "UserIDAndDocIDIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "user_id", Value: 1}}, "UserIDIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "user_id", Value: 1}, {Key: "is_default", Value: 1}}, "ParentIDandDefaultAddressIndex", false)

}

func (r *UserDeliveryAddressRepository) CreateOne(payload *userdto.AddOrEditDeliveryAddressDTO) (*usermodel.UserDeliveryAddressEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc := usermodel.UserDeliveryAddressEntity{
		ID:       primitive.NewObjectID(),
		Name:     payload.Name,
		Address:  payload.Address,
		State:    payload.State,
		Locality: payload.Locality,
		PinCode:  payload.PinCode,
		PhoneNumber: interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		},
		IsDefault: payload.IsDefault,
		UserId:    *payload.UserId,
	}

	if _, err := r.colln.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	return &doc, nil

}

func (r *UserDeliveryAddressRepository) FindAll(pgnOpt *api.PaginationOptions,
	sortOpts *map[string]int8,
	filterOpts *map[string]primitive.M,
	keySetSortby string,
	userId *primitive.ObjectID,
) ([]usermodel.UserDeliveryAddressEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opt := &options.FindOptions{}

	filters := map[string]bson.M{
		"user_id": {"$eq": userId},
	}

	if pgnOpt != nil {
		opt.Limit = options.Find().SetLimit(int64(pgnOpt.PerPage)).Limit
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	if sortOpts != nil {
		opt.Sort = options.Find().SetSort(sortOpts).Sort
	} else {
		opt.Sort = options.Find().SetSort(bson.M{"_id": -1}).Sort
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

	var results []usermodel.UserDeliveryAddressEntity
	// check for errors in the conversion
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	return results, nil
}

func (r *UserDeliveryAddressRepository) GetDocumentsCount(UserId *primitive.ObjectID, filterOpts *map[string]primitive.M) (int64, error) {

	filters := map[string]bson.M{
		"user_id": {"$eq": UserId},
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

func (r *UserDeliveryAddressRepository) FindById(userId, id *primitive.ObjectID) (*usermodel.UserDeliveryAddressEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"_id": id, "user_id": userId}}}

	cur := r.colln.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any matching result for given id")
	}

	var data usermodel.UserDeliveryAddressEntity

	if err := cur.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UserDeliveryAddressRepository) UpdateById(userId, id *primitive.ObjectID, payload *userdto.AddOrEditDeliveryAddressDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"_id": id, "user_id": userId}}}

	if _, err := r.colln.UpdateOne(ctx, filter, bson.M{"$set": payload}); err != nil {
		return err
	}

	return nil
}

func (r *UserDeliveryAddressRepository) DeleteById(userId, id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.M{"$and": bson.A{bson.M{"_id": id, "user_id": userId}}}

	if _, err := r.colln.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (r *UserDeliveryAddressRepository) UpdateDefaultStatus(userId, id *primitive.ObjectID, status bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if status {
		filter := bson.D{{Key: "$and", Value: bson.A{bson.M{"user_id": userId}, bson.M{"is_default": true}}}}

		if _, err := r.colln.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"is_default": false}}); err != nil {
			return err
		}
	}

	updateFilter := bson.D{{Key: "$and", Value: bson.A{bson.M{"user_id": userId}, bson.M{"_id": id}}}}

	if _, err := r.colln.UpdateOne(ctx, updateFilter, bson.M{"$set": bson.M{"is_default": status}}); err != nil {
		return err
	}

	return nil
}
