package discountcouponmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type DiscountType string

const (
	Percentage DiscountType = "percentage"
	Flat       DiscountType = "flat"
)

// This entity represents model for creating coupons
type DiscountCouponOfferEntity struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id"`
	Title        string              `json:"title" bson:"title"`
	DiscountType DiscountType        `json:"discount_type" bson:"discount_type"`
	FlatRate     *uint64             `json:"flat_rate" bson:"falt_rate"`
	Percentage   *uint64             `json:"percentage" bson:"percentage"`
	OfferExpiry  *primitive.DateTime `json:"offer_expiry" bson:"offer_expiry"`
}

// If the coupon needs to tied to particular user.
// LinkedUserId *primitive.ObjectID `json:"-" bson:"linked_user_id"`
