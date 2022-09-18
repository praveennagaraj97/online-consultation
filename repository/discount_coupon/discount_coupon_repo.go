package discountcouponrepository

import (
	"context"
	"errors"
	"time"

	discountcouponmodel "github.com/praveennagaraj97/online-consultation/models/discount_coupon"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscountCouponRespository struct {
	colln      *mongo.Collection
	offerColln *mongo.Collection
}

func (r *DiscountCouponRespository) Initialize(colln, offerColln *mongo.Collection) {
	r.colln = colln
	r.offerColln = offerColln

	// Unique index for offer
	utils.CreateIndex(offerColln, bson.D{{Key: "title", Value: 1}}, "Unique Offer Title", true)

}

func (r *DiscountCouponRespository) CreateOffer(payload *discountcouponmodel.DiscountCouponOfferEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.offerColln.InsertOne(ctx, payload)

	if mongo.IsDuplicateKeyError(err) {
		return errors.New("offer title should be unique")
	}

	return err
}
