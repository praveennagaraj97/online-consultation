package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	authvalidator "github.com/praveennagaraj97/online-consultation/pkg/validator/auth"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

// Register user.
func (a *UserAPI) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload userdto.RegisterDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, "Failed to parse payload", http.StatusUnprocessableEntity, nil)
		}
		defer ctx.Request.Body.Close()

		// Validate
		if err := authvalidator.ValidateRegisterDTO(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		// Check if OTP is verified
		objectId, phone, err := decodeVerificationID(payload.VerificationId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		otpRef, err := a.otpRepo.FindById(objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		_, refPhoneNumber, _ := decodeVerificationID(otpRef.VerificationID)

		if !(refPhoneNumber.Code == payload.PhoneCode && phone.Number == payload.PhoneNumber) {
			api.SendErrorResponse(ctx, "Provided verification ID is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		if !otpRef.Verified {
			api.SendErrorResponse(ctx, "Provided verification ID is not verified", http.StatusUnprocessableEntity, nil)
			return
		}

		// Store to database
		res, err := a.userRepo.CreateUser(&payload)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		access, refresh, err := res.GetAccessAndRefreshToken()

		a.userRepo.UpdateById(&res.ID, &userdto.UpdateUserDTO{
			RefreshToken: refresh,
		})

		// Set Access Token
		ctx.SetCookie(string(constants.AUTH_TOKEN),
			access,
			constants.CookieAccessExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		ctx.SetCookie(string(constants.REFRESH_TOKEN),
			access,
			constants.CookieAccessExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.AuthResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			DataResponse: serialize.DataResponse[*usermodel.UserEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusCreated,
					Message:    "Registered successfully",
				},
			},
		})

	}
}

// Send Login Email link to user.
func (a *UserAPI) SendLoginLink() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// Verify email after register
func (a *UserAPI) VerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
