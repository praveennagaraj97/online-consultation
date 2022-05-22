package userapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	authvalidator "github.com/praveennagaraj97/online-consultation/pkg/validator/auth"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register user.
func (a *UserAPI) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload userdto.RegisterDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
		}
		defer ctx.Request.Body.Close()

		// Validate
		if err := authvalidator.ValidateRegisterDTO(&payload); err != nil {
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
		res, err := a.userRepo.CreateUser(&payload)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, accessTime, err := res.GetAccessAndRefreshToken(!shouldExp)

		a.userRepo.UpdateById(&res.ID, &userdto.UpdateUserDTO{
			RefreshToken: refresh,
		})

		// Set Access Token
		ctx.SetCookie(string(constants.AUTH_TOKEN),
			access,
			accessTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		ctx.SetCookie(string(constants.REFRESH_TOKEN),
			access,
			constants.CookieRefreshExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		token, err := tokens.GenerateTokenWithExpiryTimeAndType(res.ID.Hex(),
			time.Now().Local().Add(time.Hour*48).Unix(),
			"verify_email", "user")

		emailLink := fmt.Sprintf("%s?verifyCode=%s",
			env.GetEnvVariable("CLIENT_VERIFY_EMAIL_LINK"), token)

		td := mailer.GetRegisterEmailTemplateData(res.Name, emailLink)
		if err = a.appConf.EmailClient.SendNoReplyMail([]string{res.Email},
			"Welcome to Online Consultation", "verify-email",
			"base", td); err != nil {
			log.Println("Register email failed to send")
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

// Login with Phone Number.
func (a *UserAPI) SignInWithPhoneNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload userdto.SignInWithPhoneDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if err := authvalidator.ValidateSignInWithPhoneDTO(&payload); err != nil {
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

		access, refresh, accessTime, err := res.GetAccessAndRefreshToken(!shouldExp)

		a.userRepo.UpdateById(&res.ID, &userdto.UpdateUserDTO{
			RefreshToken: refresh,
		})

		// Set Access Token
		ctx.SetCookie(string(constants.AUTH_TOKEN),
			access,
			accessTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		ctx.SetCookie(string(constants.REFRESH_TOKEN),
			access,
			constants.CookieRefreshExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

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

// Send Login Email link to user.
func (a *UserAPI) SignInWithEmailLink() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload userdto.SignInWithEmailLinkDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if payload.Email == "" {
			api.SendErrorResponse(ctx, "Email cannot be empty", http.StatusUnprocessableEntity, nil)
			return
		}

		if err := validator.ValidateEmail(payload.Email); err != nil {
			api.SendErrorResponse(ctx, "Provided email is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		// Check if user exists with email
		res, err := a.userRepo.FindByEmail(payload.Email)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// Check if email is verified.
		if !res.EmailVerified {
			api.SendErrorResponse(ctx, "Your email verification is still pending", http.StatusNotAcceptable, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		token, err := tokens.GenerateTokenWithExpiryTimeAndType(res.ID.Hex(),
			time.Now().Local().Add(time.Minute*5).Unix(),
			"sign-in-email", "user")

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		emailLink := fmt.Sprintf("%s?redirectTo=%s&remember_me=%v&verifyCode=%s",
			env.GetEnvVariable("CLIENT_EMAIL_SIGNIN_LINK"),
			payload.RedirectTo,
			shouldExp, token)

		td := mailer.GetSignWithEmailLinkTemplateData(res.Name, emailLink)
		if err = a.appConf.EmailClient.SendNoReplyMail([]string{res.Email}, "Sign in to your online consultation account", "verify-email", "base", td); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "An email with login link has been sent to your email address.",
		})

	}
}

// Login Credentials for email link
func (a UserAPI) SendLoginCredentialsForEmailLink() gin.HandlerFunc {
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

		if claimedInfo.Type != "sign-in-email" {
			api.SendErrorResponse(ctx, "Provided token is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(claimedInfo.ID)

		if err != nil {
			api.SendErrorResponse(ctx, "Token is malformed", http.StatusUnprocessableEntity, nil)
			return
		}

		// Find by phone number
		res, err := a.userRepo.FindById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, accessTime, err := res.GetAccessAndRefreshToken(!shouldExp)

		a.userRepo.UpdateById(&res.ID, &userdto.UpdateUserDTO{
			RefreshToken: refresh,
		})

		// Set Access Token
		ctx.SetCookie(string(constants.AUTH_TOKEN),
			access,
			accessTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		ctx.SetCookie(string(constants.REFRESH_TOKEN),
			access,
			constants.CookieRefreshExpiryTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

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

// Request email verify link
func (a *UserAPI) RequestEmailVerifyLink() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload userdto.RequestEmailVerifyDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if payload.Email == "" {
			api.SendErrorResponse(ctx, "Email cannot be empty", http.StatusUnprocessableEntity, nil)
			return
		}

		if err := validator.ValidateEmail(payload.Email); err != nil {
			api.SendErrorResponse(ctx, "Provided email is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		// Check if user exists with email
		res, err := a.userRepo.FindByEmail(payload.Email)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// Check if email is verified.
		if res.EmailVerified {
			api.SendErrorResponse(ctx, "Your email is already verified", http.StatusNotAcceptable, nil)
			return
		}

		token, err := tokens.GenerateTokenWithExpiryTimeAndType(res.ID.Hex(),
			time.Now().Local().Add(time.Hour*48).Unix(),
			"verify_email", "user")

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		emailLink := fmt.Sprintf("%s?redirectTo=%s&verifyCode=%s",
			env.GetEnvVariable("CLIENT_VERIFY_EMAIL_LINK"),
			payload.RedirectTo,
			token)

		td := mailer.GetVerifyEmailTemplateData(res.Name, emailLink)
		if err = a.appConf.EmailClient.SendNoReplyMail([]string{res.Email},
			"Verify email address", "verify-email", "base", td); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "An email with link to verify your email address has been sent to your email address.",
		})

	}
}

// Confirm Email address after register or manual verify request
func (a *UserAPI) ConfirmEmail() gin.HandlerFunc {
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

		if claimedInfo.Type != "verify_email" {
			api.SendErrorResponse(ctx, "Provided token is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(claimedInfo.ID)

		if err != nil {
			api.SendErrorResponse(ctx, "Token is malformed", http.StatusUnprocessableEntity, nil)
			return
		}

		// Get user By ID
		user, err := a.userRepo.FindById(&objectId)
		if err != nil {
			api.SendErrorResponse(ctx, "Provided token is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		if user.EmailVerified {
			api.SendErrorResponse(ctx, "Email is already verified", http.StatusBadRequest, nil)
			return
		}

		if err = a.userRepo.UpdateById(&objectId, &userdto.UpdateUserDTO{EmailVerified: true}); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Email verified successfully",
		})

	}
}

// Logout - Removes user token from Entity and disables token from used further.
func (a *UserAPI) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		err = a.userRepo.UpdateById(id, &userdto.UpdateUserDTO{
			RefreshToken: "",
		})

		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		c.SetCookie(string(constants.AUTH_TOKEN), "", 0, "/", a.appConf.Domain, false, true)
		c.SetCookie(string(constants.REFRESH_TOKEN), "", 0, "/", a.appConf.Domain, false, true)

		c.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Logged out successfully",
		})

	}
}

// Refresh access token using refresh token.
// Can be force refreshed by passing force param set to true.
func (a *UserAPI) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, refreshToken, err := a.getAccessAndRefreshTokenFromRequest(c)
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
		user, err := a.userRepo.FindById(&userId)
		if err != nil {
			api.SendErrorResponse(c, "Couldn't find any user for this refresh token", http.StatusNotFound, nil)
			return
		}

		if user.RefreshToken != refreshToken {
			api.SendErrorResponse(c, "Revalidate token Malformed", http.StatusUnauthorized, nil)
			return
		}

		access, refresh, accessTime, err := user.GetAccessAndRefreshToken(true)

		a.userRepo.UpdateById(&user.ID, &userdto.UpdateUserDTO{
			RefreshToken: refresh,
		})

		// Set Access Token
		c.SetCookie(string(constants.AUTH_TOKEN),
			access,
			accessTime, "/", a.appConf.Domain, a.appConf.Environment == "production", true)

		// Set Refresh Token
		c.SetCookie(string(constants.REFRESH_TOKEN),
			access,
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
			RefreshToken: refreshToken,
		})
	}
}
