package discountcouponapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	discountcoupondto "github.com/praveennagaraj97/online-consultation/dto/discount_coupon"
	discountcouponmodel "github.com/praveennagaraj97/online-consultation/models/discount_coupon"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create An Offer
func (a *DiscountCouponAPi) CreateNewOffer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload discountcoupondto.NewOfferDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		timeZone := ctx.Request.Header.Get(constants.TimeZoneHeaderKey)

		if errs := payload.Validate(timeZone); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		doc := discountcouponmodel.DiscountCouponOfferEntity{
			ID:           primitive.NewObjectID(),
			Title:        payload.Title,
			DiscountType: payload.DiscountType,
			FlatRate:     payload.FlatRate,
			Percentage:   payload.Percentage,
			OfferExpiry:  payload.OfferExpiry,
		}

		if err := a.repo.CreateOffer(&doc); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[discountcouponmodel.DiscountCouponOfferEntity]{
			Data: doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Offer plan created successfully",
			},
		})

	}
}

func (a *DiscountCouponAPi) UpdateOfferById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, exits := ctx.Params.Get("id")

		if !exits {
			api.SendErrorResponse(ctx, "Offer ID is missing", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		fmt.Println(objectId)

	}
}
