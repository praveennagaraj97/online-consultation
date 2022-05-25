package adminapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	admindto "github.com/praveennagaraj97/online-consultation/dto/admin"
	adminmodel "github.com/praveennagaraj97/online-consultation/models/admin"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	adminvalidator "github.com/praveennagaraj97/online-consultation/pkg/validator/admin"
	adminrepository "github.com/praveennagaraj97/online-consultation/repository/admin"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminAPI struct {
	appConf   *app.ApplicationConfig
	adminRepo *adminrepository.AdminRepository
}

func (a *AdminAPI) Initailize(conf *app.ApplicationConfig, adminRepo *adminrepository.AdminRepository) {
	a.appConf = conf
	a.adminRepo = adminRepo
}

// Create new user
func (a *AdminAPI) AddNewAdmin(role constants.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload admindto.AddNewAdminDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		if errors := adminvalidator.ValidateNewAdminDTO(&payload); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, http.StatusUnprocessableEntity, errors.Errors)
			return
		}

		payload.Role = role

		user, err := a.adminRepo.CreateOne(&payload)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*adminmodel.AdminEntity]{
			Data: user,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    fmt.Sprintf("New user with %s role has been added successfully", string(role)),
			},
		})
	}
}

func (a *AdminAPI) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload admindto.LoginDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errors := adminvalidator.ValidateAdminLoginDTO(&payload); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, errors.StatusCode, errors.Errors)
			return
		}

		var user *adminmodel.AdminEntity
		var err error
		if payload.Email != "" {
			user, err = a.adminRepo.FindByEmail(payload.Email)
		} else {
			user, err = a.adminRepo.FindByUserName(payload.UserName)

		}
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err = user.DecodePassword(payload.Password); err != nil {
			api.SendErrorResponse(ctx, "Entered password is not valid", http.StatusUnauthorized, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, accessTime, err := user.GetAccessAndRefreshToken(!shouldExp, string(user.Role))

		a.adminRepo.UpdateRefreshToken(&user.ID, refresh)

		// Set Access Token
		ctx.SetCookie(string(constants.AUTH_TOKEN),
			access,
			accessTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		ctx.SetCookie(string(constants.REFRESH_TOKEN),
			refresh,
			constants.CookieRefreshExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.AdminAuthResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			DataResponse: serialize.DataResponse[*adminmodel.AdminEntity]{
				Data: user,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "Logged in successfully",
				},
			},
		})

	}
}

func (a *AdminAPI) UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload admindto.UpdatePasswordDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		if errors := adminvalidator.ValidateUpdatePasswordDTO(&payload); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, errors.StatusCode, errors.Errors)
			return
		}

		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		user, err := a.adminRepo.FindById(userId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if err := user.DecodePassword(payload.Password); err != nil {
			api.SendErrorResponse(ctx, "Current password is not valid", http.StatusUnauthorized, nil)
			return
		}

		user.Password = payload.NewPassword
		user.EncodePassword()

		if err = a.adminRepo.UpdateById(&user.ID, &admindto.UpdateAdminDTO{
			Password: user.Password,
		}); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (a *AdminAPI) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload admindto.ForgorPasswordDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		if err := validator.ValidateEmail(payload.Email); err != nil {
			api.SendErrorResponse(ctx, "Entered email is not valid", http.StatusUnprocessableEntity, nil)
			return
		}

		user, err := a.adminRepo.FindByEmail(payload.Email)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		token, err := tokens.GenerateTokenWithExpiryTimeAndType(user.ID.Hex(),
			time.Now().Local().Add(time.Hour*48).Unix(),
			"reset-email", "user")

		emailLink := fmt.Sprintf("%s?verifyCode=%s",
			env.GetEnvVariable("CLIENT_VERIFY_FORGOT_PASSWORD_LINK"), token)

		td := mailer.GetForgotEmailLinkTemplateData(user.Name, emailLink)
		if err = a.appConf.EmailClient.SendNoReplyMail([]string{user.Email},
			"Reset your password", "verify-email",
			"base", td); err != nil {
			log.Println("Reset email failed to send")
		}

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "An email with reset link has been sent to your email.",
		})

	}
}

func (a *AdminAPI) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, exists := ctx.Params.Get("token")
		if !exists {
			api.SendErrorResponse(ctx, "Couldn't find any token", http.StatusUnprocessableEntity, nil)
			return
		}

		claimedInfo, err := tokens.DecodeJSONWebToken(token)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if claimedInfo.Type != "reset-email" {
			api.SendErrorResponse(ctx, "Provided token is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(claimedInfo.ID)

		if err != nil {
			api.SendErrorResponse(ctx, "Token is malformed", http.StatusUnprocessableEntity, nil)
			return
		}

		var payload admindto.ResetPasswordDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if err := adminvalidator.ValidateResetPasswordDTO(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		user, err := a.adminRepo.FindById(&objectId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		user.Password = payload.NewPassword
		user.EncodePassword()

		if err = a.adminRepo.UpdateById(&user.ID, &admindto.UpdateAdminDTO{
			Password: user.Password,
		}); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Password changed successfully",
		})

	}
}

func (a *AdminAPI) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, refreshToken, err := tokens.GetAccessAndRefreshTokenFromRequest(c)
		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// check if force refresh is requested
		isForce, _ := strconv.ParseBool(c.Request.URL.Query().Get("force"))

		// parse auth Token
		_, err = tokens.DecodeJSONWebToken(token)
		if err == nil && !isForce {
			api.SendErrorResponse(c, "Token is not expired", http.StatusNotAcceptable, nil)
			return
		}

		// parse refresh token
		claimedRefreshToken, err := tokens.DecodeJSONWebToken(refreshToken)
		if err != nil {
			api.SendErrorResponse(c, "Revalidate token malformed", http.StatusNotAcceptable, nil)
			return
		}

		userId, err := primitive.ObjectIDFromHex(claimedRefreshToken.ID)
		if err != nil {
			api.SendErrorResponse(c, "Something went wrong", http.StatusNotAcceptable, nil)
			return
		}

		// cross check refresh token with db.
		user, err := a.adminRepo.FindById(&userId)
		if err != nil {
			api.SendErrorResponse(c, "Couldn't find any user for this refresh token", http.StatusNotFound, nil)
			return
		}

		if user.RefreshToken != refreshToken {
			api.SendErrorResponse(c, "Revalidate token Malformed", http.StatusUnauthorized, nil)
			return
		}

		access, refresh, accessTime, err := user.GetAccessAndRefreshToken(true, string(user.Role))

		if err = a.adminRepo.UpdateRefreshToken(&user.ID, refresh); err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusBadGateway, nil)
			return
		}

		// Set Access Token
		c.SetCookie(string(constants.AUTH_TOKEN),
			access,
			accessTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		c.SetCookie(string(constants.REFRESH_TOKEN),
			refresh,
			constants.CookieRefreshExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		if err != nil {
			api.SendErrorResponse(c, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, &serialize.RefreshResponse{
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Token refreshed successfully",
			},
			Token:        token,
			RefreshToken: refresh,
		})
	}
}

func (a *AdminAPI) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload admindto.AdminIdDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		adminId, err := primitive.ObjectIDFromHex(payload.ID)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if err := a.adminRepo.DeleteById(&adminId); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(204, nil)

	}
}

func (a *AdminAPI) ChangeRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"info": "Method not implemented",
		})
	}
}

func (a *AdminAPI) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		err = a.adminRepo.UpdateRefreshToken(id, "")

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		ctx.SetCookie(string(constants.AUTH_TOKEN), "", 0, "/", a.appConf.Domain, false, true)
		ctx.SetCookie(string(constants.REFRESH_TOKEN), "", 0, "/", a.appConf.Domain, false, true)

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Logged out successfully",
		})
	}
}
