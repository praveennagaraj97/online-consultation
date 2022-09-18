package discountcoupondto

import (
	"net/http"
	"time"

	discountcouponmodel "github.com/praveennagaraj97/online-consultation/models/discount_coupon"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewOfferDTO struct {
	Title          string                           `json:"title" form:"title"`
	DiscountType   discountcouponmodel.DiscountType `json:"discount_type" form:"discount_type"`
	FlatRate       *uint64                          `json:"flat_rate,omitempty" form:"flat_rate,omitempty"`
	Percentage     *uint64                          `json:"percentage,omitempty" form:"percentage,omitempty"`
	OfferExpiryRef string                           `json:"offer_expiry,omitempty" form:"offer_expiry,omitempty"`

	OfferExpiry *primitive.DateTime `json:"-" form:"-"`
}

func (d *NewOfferDTO) Validate(timeZone string) *serialize.ErrorResponse {
	errs := map[string]string{}

	if d.Title == "" {
		errs["title"] = "Offer title cannot be empty"
	}

	if d.DiscountType == "" {
		errs["discount_type"] = "Discount type cannot be empty"
	} else if d.DiscountType != discountcouponmodel.Flat && d.DiscountType != discountcouponmodel.Percentage {
		errs["discount_type"] = "Discount type should be either 'flat or percentage'"
	}

	if d.DiscountType == discountcouponmodel.Flat && d.FlatRate == nil {
		errs["flat_rate"] = "Flat rate cannot be empty"
	}

	if d.DiscountType == discountcouponmodel.Percentage && d.Percentage == nil {
		errs["percentage"] = "Percentage cannot be empty"
	}

	if d.OfferExpiryRef != "" {
		timeLoc, err := time.LoadLocation(timeZone)
		if err != nil {
			errs["time_zone_header"] = "Provided time zone is invalid"
		}
		t, err := time.ParseInLocation("2006-01-02", d.OfferExpiryRef, timeLoc)
		if err != nil {
			errs["offer_expiry"] = err.Error()
		} else {
			offerExpiry := primitive.NewDateTimeFromTime(t.UTC())
			if offerExpiry <= primitive.NewDateTimeFromTime(time.Now()) {
				errs["offer_expiry"] = "Offer expiry cannot be set in past date"
			}
			d.OfferExpiry = &offerExpiry
		}
	}

	if len(errs) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errs,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Provided data is not valid",
			},
		}
	}

	return nil

}
