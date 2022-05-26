package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *UserAPI) AddNewAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		var payload userdto.AddOrEditDeliveryAddressDTO

		if err = ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errors := payload.ValidateUserDeliveryAddressDTO(); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, errors.StatusCode, errors.Errors)
			return
		}

		payload.UserId = *userId

		res, err := a.delvrAddrRepo.CreateOne(&payload)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*usermodel.UserDeliveryAddressEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Delivery address added successfully",
			},
		})
	}
}

func (a *UserAPI) GetAllAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		// get pagination/sort/filter options.
		pgOpts := api.ParsePaginationOptions(ctx, "user_delivery_address")
		srtOpts := api.ParseSortByOptions(ctx)
		filterOpts := api.ParseFilterByOptions(ctx)
		keySetSortby := "$gt"

		// Default options | sort by latest
		if len(*srtOpts) == 0 {
			srtOpts = &map[string]int8{"_id": -1}
		}

		if pgOpts.PaginateId != nil {
			for key, value := range *srtOpts {
				if value == -1 && key == "_id" {
					keySetSortby = "$lt"
				}
			}
		}

		res, err := a.delvrAddrRepo.FindAll(pgOpts, srtOpts, filterOpts, keySetSortby, userId)
		if err != nil {
			api.SendErrorResponse(ctx, "Couldn't find any relatives account assoiated with ypur account", http.StatusNotFound, nil)
			return
		}

		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.delvrAddrRepo.GetDocumentsCount(userId, filterOpts)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId, "user_delivery_address")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]usermodel.UserDeliveryAddressEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			DataResponse: serialize.DataResponse[[]usermodel.UserDeliveryAddressEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of delivery address retrieved successfully",
				},
			},
		})

	}
}

func (a *UserAPI) GetAddressById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		docParamId := ctx.Param("id")

		docId, err := primitive.ObjectIDFromHex(docParamId)
		if err != nil {
			api.SendErrorResponse(ctx, "Given ID is not valid", http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.delvrAddrRepo.FindById(userId, &docId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*usermodel.UserDeliveryAddressEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Delivery address retrieved successfully",
			},
		})

	}
}

func (a *UserAPI) UpdateAddressById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		docParamId := ctx.Param("id")

		docId, err := primitive.ObjectIDFromHex(docParamId)
		if err != nil {
			api.SendErrorResponse(ctx, "Given ID is not valid", http.StatusUnprocessableEntity, nil)
			return
		}

		var payload userdto.AddOrEditDeliveryAddressDTO

		if err = ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if err := a.delvrAddrRepo.UpdateById(userId, &docId, &payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (a *UserAPI) DeleteAddressById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		docParamId := ctx.Param("id")

		docId, err := primitive.ObjectIDFromHex(docParamId)
		if err != nil {
			api.SendErrorResponse(ctx, "Given ID is not valid", http.StatusUnprocessableEntity, nil)
			return
		}

		var payload userdto.AddOrEditDeliveryAddressDTO

		if err = ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if err := a.delvrAddrRepo.DeleteById(userId, &docId); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (a *UserAPI) ToggleDefaultAddress(status bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		docParamId := ctx.Param("id")

		docId, err := primitive.ObjectIDFromHex(docParamId)
		if err != nil {
			api.SendErrorResponse(ctx, "Given ID is not valid", http.StatusUnprocessableEntity, nil)
			return
		}

		if err = a.delvrAddrRepo.UpdateDefaultStatus(userId, &docId, status); err != nil {
			api.SendErrorResponse(ctx, "Failed to update status", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}
