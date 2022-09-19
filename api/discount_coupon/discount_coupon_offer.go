package discountcouponapi

import (
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

		var payload discountcoupondto.UpdateOfferDTO

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

		if err := a.repo.UpdateOfferById(&objectId, &payload); err != nil {
			api.SendErrorResponse(ctx, "Failed to update", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (a *DiscountCouponAPi) GetAllOffers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pgOpts := api.ParsePaginationOptions(ctx, "discount_offers")
		keySortBy := "$lt"
		fltrOpts := api.ParseFilterByOptions(ctx)
		sortOpts := map[string]int8{"_id": -1}

		res, err := a.repo.FindAllOffers(pgOpts, keySortBy, fltrOpts, &sortOpts)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}
		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.repo.CountOfferDocuments(fltrOpts)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, docCount, lastResId, "hospitals")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]discountcouponmodel.DiscountCouponOfferEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			PaginatedData: serialize.PaginatedData[[]discountcouponmodel.DiscountCouponOfferEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of discount offers plans retrieved",
				},
			},
		})

	}
}

func (a *DiscountCouponAPi) GetOfferById() gin.HandlerFunc {
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

		res, err := a.repo.FindOneOfferById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*discountcouponmodel.DiscountCouponOfferEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Discount offer details retrieved successfully",
			},
		})

	}
}

func (a *DiscountCouponAPi) DeleteOfferById() gin.HandlerFunc {
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

		if err := a.repo.DeleteOfferById(&objectId); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}
