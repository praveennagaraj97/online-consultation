package additionalapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	additionalmodel "github.com/praveennagaraj97/online-consultation/models/additional"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdditionalAPI struct {
	appConf *app.ApplicationConfig
}

func (a *AdditionalAPI) Initailize(appConf *app.ApplicationConfig) {
	a.appConf = appConf
}

// Check Auth Token Status without database comparison.
func (a *AdditionalAPI) CheckJWTStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		var token string

		// check if token exists in cookie
		cookie, err := c.Request.Cookie(string(constants.AUTH_TOKEN))
		if err != nil {
			// check in auth header as bearer
			authHeader := c.Request.Header.Get("Authorization")
			containsBearerToken := strings.HasPrefix(authHeader, "Bearer")
			if !containsBearerToken {
				api.SendErrorResponse(c, "Token is missing", http.StatusUnauthorized, nil)
				return
			} else {
				bearerToken := strings.Split(authHeader, "Bearer ")
				if len(bearerToken) == 2 {
					token = bearerToken[1]
				}
			}
		} else {
			token = cookie.Value
		}
		if token == "" {
			api.SendErrorResponse(c, "Token is missing", http.StatusUnprocessableEntity, nil)
			return
		}

		signedClaims, err := tokens.DecodeJSONWebToken(token)

		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		var expires *primitive.DateTime = nil

		if signedClaims.ExpiresAt != 0 {
			e := primitive.NewDateTimeFromTime(time.Unix(int64(signedClaims.ExpiresAt), 0))
			expires = &e
		}

		c.JSON(http.StatusOK, serialize.DataResponse[*additionalmodel.JWTStatus]{
			Data: &additionalmodel.JWTStatus{
				IsValid: true,
				Expires: expires,
			},
		})

	}
}
