package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	userdto "github.com/praveennagaraj97/online-consultation/dto"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	uservalidator "github.com/praveennagaraj97/online-consultation/pkg/validator/user"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *UserAPI) AddRelative() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload userdto.AddOrEditRelativeDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		if err := uservalidator.ValidateRelativeDTO(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		// Get User Details.
		user, err := a.userRepo.FindById(userId)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		if err := uservalidator.CompareAndValidateRelativeDTOWithUserData(&payload, user); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		payload.ParentId = *userId

		res, err := a.relativeRepo.CreateOne(&payload)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*usermodel.RelativeEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Relative account added successfully",
			},
		})

	}
}

func (a *UserAPI) GetListOfRelatives() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		// get pagination/sort/filter options.
		pgOpts := api.ParsePaginationOptions(ctx)
		srtOpts := api.ParseSortByOptions(ctx)
		filterOpts := api.ParseFilterByOptions(ctx)
		keySetSortby := "$gt"

		// Default options | sort by latest
		if len(*srtOpts) == 0 {
			srtOpts = &map[string]int8{"_id": -1}
		}
		// Key Set fix for created_at desc
		if pgOpts.PaginateId != nil {
			for key, value := range *srtOpts {
				if value == -1 && key == "_id" {
					keySetSortby = "$lt"
				}
			}
		}

		res, err := a.relativeRepo.FindAll(pgOpts, srtOpts, filterOpts, keySetSortby, userId)

		if err != nil {
			api.SendErrorResponse(ctx, "Couldn't find any relatives account assoiated with ypur account", http.StatusNotFound, nil)
			return
		}

		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.relativeRepo.GetDocumentsCount(userId, filterOpts)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId)

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]usermodel.RelativeEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			DataResponse: serialize.DataResponse[[]usermodel.RelativeEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of relatives account retrieved successfully",
				},
			},
		})
	}
}
