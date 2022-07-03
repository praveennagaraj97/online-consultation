package userapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
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

		timeZone := ctx.Request.Header.Get(constants.TimeZoneHeaderKey)

		if err := payload.ValidateRelativeDTO(timeZone); err != nil {
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

		if err := payload.CompareAndValidateRelativeDTOWithUserData(user); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		payload.UserId = userId

		if exists := a.relativeRepo.CheckIfRelativeExist(payload.Email,
			interfaces.PhoneType{Code: payload.PhoneCode, Number: payload.PhoneNumber}, payload.UserId); exists {
			api.SendErrorResponse(ctx, "Relative account with given credentials already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		doc := &usermodel.RelativeEntity{
			ID:    primitive.NewObjectID(),
			Name:  payload.Name,
			Email: payload.Email,
			Phone: interfaces.PhoneType{
				Code:   payload.PhoneCode,
				Number: payload.PhoneNumber,
			},
			DateOfBirth: *payload.DateOfBirth,
			Gender:      payload.Gender,
			Relation:    payload.Relation,
			UserId:      payload.UserId,
		}

		if err := a.relativeRepo.CreateOne(doc); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*usermodel.RelativeEntity]{
			Data: doc,
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
		pgOpts := api.ParsePaginationOptions(ctx, "user_relative_account")
		srtOpts := map[string]int8{"_id": -1}
		filterOpts := api.ParseFilterByOptions(ctx)
		keySetSortby := "$gt"

		res, err := a.relativeRepo.FindAll(pgOpts, &srtOpts, filterOpts, keySetSortby, userId)

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

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId, "user_relative_account")

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

func (a *UserAPI) GetRelativeProfileById() gin.HandlerFunc {
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

		res, err := a.relativeRepo.FindById(userId, &docId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*usermodel.RelativeEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Retrieved relative profile successfully",
			},
		})
	}
}

func (a *UserAPI) UpdateRelativeProfileById() gin.HandlerFunc {
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

		var payload userdto.AddOrEditRelativeDTO

		if err = ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if payload.PhoneCode != "" && payload.PhoneNumber != "" {
			payload.Phone = &interfaces.PhoneType{
				Code:   payload.PhoneCode,
				Number: payload.PhoneNumber,
			}
		}

		timeZone := ctx.Request.Header.Get(constants.TimeZoneHeaderKey)

		if payload.DOBRef != "" {
			timeLoc, err := time.LoadLocation(timeZone)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			t, err := time.ParseInLocation("2006-01-02", payload.DOBRef, timeLoc)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			} else {
				dateOfBirth := primitive.NewDateTimeFromTime(t.UTC())
				payload.DateOfBirth = &dateOfBirth
			}
		}

		if err := a.relativeRepo.UpdateByID(userId, &docId, &payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (a *UserAPI) DeleteRelativeProfileById() gin.HandlerFunc {
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

		if err := a.relativeRepo.DeleteByID(userId, &docId); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}
