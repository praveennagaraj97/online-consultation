package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

type EmailOrPhoneExists struct {
	IsAvailable bool `json:"is_available"`
}

func (a *UserAPI) CheckUserExistsByEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload struct {
			Email string `json:"email" form:"email"`
		}

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, "Given input is invalid", http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if payload.Email == "" {
			api.SendErrorResponse(ctx, "Email cannot be empty", http.StatusUnprocessableEntity, nil)
			return
		}

		if err := validator.ValidateEmail(payload.Email); err != nil {
			api.SendErrorResponse(ctx, "Provide email is not valid", http.StatusUnprocessableEntity, nil)
			return
		}

		if _, err := a.userRepo.FindByEmail(payload.Email); err != nil {
			ctx.JSON(http.StatusOK, serialize.DataResponse[EmailOrPhoneExists]{
				Data: EmailOrPhoneExists{
					IsAvailable: true,
				},
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "Email is available",
				},
			})
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[EmailOrPhoneExists]{
			Data: EmailOrPhoneExists{
				IsAvailable: false,
			},
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Email is already in use by another account",
			},
		})

	}
}

func (a *UserAPI) CheckUserExistsByPhoneNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload interfaces.PhoneType

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, "Given input is invalid", http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if payload.Code == "" || payload.Number == "" {
			api.SendErrorResponse(ctx, "Provide phone number is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		if _, err := a.userRepo.FindByPhoneNumber(payload.Number, payload.Code); err != nil {
			ctx.JSON(http.StatusOK, serialize.DataResponse[EmailOrPhoneExists]{
				Data: EmailOrPhoneExists{
					IsAvailable: true,
				},
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "Phone number is available",
				},
			})
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[EmailOrPhoneExists]{
			Data: EmailOrPhoneExists{
				IsAvailable: false,
			},
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Phone number is already in use by another account",
			},
		})

	}
}
