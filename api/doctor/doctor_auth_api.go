package doctorapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/serialize"
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
		if err, statusCode := a.validateVerifyCode(ctx, payload.VerificationId, interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}); err != nil {
			api.SendErrorResponse(ctx, err.Error(), statusCode, nil)
			return
		}

		// Find by phone number
		res, err := a.repo.FindByPhoneNumber(interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		})

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
