package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	userdto "github.com/praveennagaraj97/online-consultation/dto"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	uservalidator "github.com/praveennagaraj97/online-consultation/pkg/validator/user"
	"github.com/praveennagaraj97/online-consultation/serialize"
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
		keySetSortby := "$gt"

		res, _ := a.relativeRepo.FindAll(pgOpts, userId, keySetSortby)

		ctx.JSON(200, map[string]interface{}{
			"res": res,
		})

	}
}
