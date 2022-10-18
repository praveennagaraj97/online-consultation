package doctorapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *DoctorAPI) SignInWithPhoneNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload userdto.SignInWithPhoneDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		if err := payload.ValidateSignInWithPhoneDTO(); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		// Validate veriy code
		if statusCode, err := a.validateVerifyCode(ctx, payload.VerificationId, interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}); err != nil {
			api.SendErrorResponse(ctx, err.Error(), statusCode, nil)
			return
		}

		// Find by phone number
		res, err := a.repo.FindOne(nil, "", &interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}, false)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, accessTime, err := res.GetAccessAndRefreshToken(!shouldExp)

		if err := a.repo.UpdateRefreshToken(&res.ID, refresh); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}
		utils.SetAuthCookie(ctx, access, refresh, accessTime, a.appConf.Domain, a.appConf.Environment, constants.DOCTOR_AUTH_COOKIE)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.AuthResponse[*doctormodel.DoctorEntity]{
			AccessToken:  access,
			RefreshToken: refresh,
			DataResponse: serialize.DataResponse[*doctormodel.DoctorEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "Logged in successfully",
				},
			},
		})

	}
}

// Refresh access token using refresh token.
// Can be force refreshed by passing force param set to true.
func (a *DoctorAPI) RefreshToken() gin.HandlerFunc {
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
		user, err := a.repo.FindById(&userId)
		if err != nil {
			api.SendErrorResponse(c, "Couldn't find any user for this refresh token", http.StatusNotFound, nil)
			return
		}

		if user.RefreshToken != refreshToken {
			api.SendErrorResponse(c, "Revalidate token Malformed", http.StatusUnauthorized, nil)
			return
		}

		access, refresh, accessTime, err := user.GetAccessAndRefreshToken(true)

		a.repo.UpdateRefreshToken(&user.ID, refresh)

		utils.SetAuthCookie(c, access, refresh, accessTime, a.appConf.Domain, a.appConf.Environment, constants.DOCTOR_AUTH_COOKIE)

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

// Send Login Email link to user.
func (a *DoctorAPI) SignInWithEmailLink() gin.HandlerFunc {
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
		res, err := a.repo.FindByEmail(payload.Email)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// Check if email is verified.
		if !res.IsActive {
			api.SendErrorResponse(ctx, "Your account is not activated. Please check your email or contact admin.", http.StatusNotAcceptable, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		token, err := tokens.GenerateTokenWithExpiryTimeAndType(res.ID.Hex(),
			time.Now().Add(time.Minute*5).Unix(),
			"sign-in-email", "doctor")

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		emailLink := fmt.Sprintf("%s?redirectTo=%s&remember_me=%v&verifyCode=%s",
			env.GetEnvVariable("CLIENT_DOCTOR_SIGNIN_LINK"),
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
func (a *DoctorAPI) SendLoginCredentialsForEmailLink() gin.HandlerFunc {
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

		if claimedInfo.Type != "sign-in-email" || claimedInfo.Role != "doctor" {
			api.SendErrorResponse(ctx, "Provided token is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(claimedInfo.ID)

		if err != nil {
			api.SendErrorResponse(ctx, "Token is malformed", http.StatusUnprocessableEntity, nil)
			return
		}

		// Find by phone number
		res, err := a.repo.FindOne(&objectId, "", nil, false)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		shouldExp, err := strconv.ParseBool(ctx.Query("remember_me"))

		if err != nil {
			shouldExp = false
		}

		access, refresh, accessTime, err := res.GetAccessAndRefreshToken(!shouldExp)

		a.repo.UpdateRefreshToken(&res.ID, refresh)

		utils.SetAuthCookie(ctx, access, refresh, accessTime, a.appConf.Domain, a.appConf.Environment, constants.DOCTOR_AUTH_COOKIE)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.AuthResponse[*doctormodel.DoctorEntity]{
			AccessToken:  access,
			RefreshToken: refresh,
			DataResponse: serialize.DataResponse[*doctormodel.DoctorEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "Logged in successfully",
				},
			},
		})
	}
}

// Logout - Removes user token from Entity and disables token from used further.
func (a *DoctorAPI) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		err = a.repo.UpdateRefreshToken(id, "")

		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		utils.SetAuthCookie(c, "", "", 0, a.appConf.Domain, a.appConf.Environment, constants.DOCTOR_AUTH_COOKIE)

		c.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Logged out successfully",
		})

	}
}
