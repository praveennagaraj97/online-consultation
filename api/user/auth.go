package userapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	authvalidator "github.com/praveennagaraj97/online-consultation/pkg/validator/auth"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

// Register user.
func (a *UserAPI) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload *userdto.RegisterDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
		}
		defer ctx.Request.Body.Close()

		// Validate
		if err := authvalidator.ValidateRegisterDTO(payload); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		if err, statusCode := a.validateVerifyCode(ctx, payload.VerificationId, interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}); err != nil {
			api.SendErrorResponse(ctx, err.Error(), statusCode, nil)
			return
		}

		// Store to database
		res, err := a.userRepo.CreateUser(payload)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, err := res.GetAccessAndRefreshToken(!shouldExp)

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

// Login with Phone Number.
func (a *UserAPI) SignInWithPhoneNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload *userdto.SignInWithPhoneDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if err := authvalidator.ValidateSignInWithPhoneDTO(payload); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		// Validate veriy code
		if err, statusCode := a.validateVerifyCode(ctx, payload.VerificationId, interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}); err != nil {
			api.SendErrorResponse(ctx, err.Error(), statusCode, nil)
			return
		}

		// Find by phone number
		res, err := a.userRepo.FindByPhoneNumber(payload.PhoneNumber, payload.PhoneCode)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, err := res.GetAccessAndRefreshToken(!shouldExp)

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

		ctx.JSON(http.StatusOK, serialize.AuthResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			DataResponse: serialize.DataResponse[*usermodel.UserEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "Logged in successfully",
				},
			},
		})

	}
}
