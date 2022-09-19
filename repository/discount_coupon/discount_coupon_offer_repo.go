package discountcouponrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	discountcoupondto "github.com/praveennagaraj97/online-consultation/dto/discount_coupon"
	discountcouponmodel "github.com/praveennagaraj97/online-consultation/models/discount_coupon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *DiscountCouponRespository) CreateOffer(payload *discountcouponmodel.DiscountCouponOfferEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.offerColln.InsertOne(ctx, payload)

	if mongo.IsDuplicateKeyError(err) {
		return errors.New("offer title should be unique")
	}

	return err
}

func (r *DiscountCouponRespository) UpdateOfferById(id *primitive.ObjectID, payload *discountcoupondto.UpdateOfferDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.offerColln.UpdateByID(ctx, id, bson.M{"$set": payload})

	return err

}

func (r *DiscountCouponRespository) FindAllOffers(
	pgOpts *api.PaginationOptions,
	keySortBy string,
	fltrOpts *map[string]primitive.M,
	sortOps *map[string]int8,
) ([]discountcouponmodel.DiscountCouponOfferEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opts := options.FindOptions{}

	opts.Limit = options.Find().SetLimit(int64(pgOpts.PerPage)).Limit

	if sortOps != nil {
		opts.Sort = options.Find().SetSort(sortOps).Sort
	}

	filters := make(map[string]primitive.M, 0)

	if fltrOpts != nil {
		for key, val := range *fltrOpts {
			filters[key] = val
		}
	}

	if pgOpts.PaginateId != nil {
		filters["_id"] = bson.M{keySortBy: pgOpts.PaginateId}
	} else if pgOpts != nil {
		opts.Skip = options.Find().SetSkip(int64(pgOpts.PerPage) * int64(pgOpts.PageNum-1)).Skip
	}

	cur, err := r.offerColln.Find(ctx, filters, &opts)
	if err != nil {
		return nil, err
	}

	var results []discountcouponmodel.DiscountCouponOfferEntity

	if err := cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *DiscountCouponRespository) CountOfferDocuments(filterOpts *map[string]primitive.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.offerColln.CountDocuments(ctx, filterOpts)

}

func (r *DiscountCouponRespository) FindOneOfferById(id *primitive.ObjectID) (*discountcouponmodel.DiscountCouponOfferEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.offerColln.FindOne(ctx, bson.M{"_id": id})

	if cur.Err() != nil {
		return nil, errors.New("couldn't find any matching offer with given ID")
	}

	var result discountcouponmodel.DiscountCouponOfferEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *DiscountCouponRespository) DeleteOfferById(id *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.offerColln.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
