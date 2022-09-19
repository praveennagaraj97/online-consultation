package discountcouponrepository

import (
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
