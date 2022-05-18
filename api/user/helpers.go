package userapi

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *UserAPI) validateVerifyCode(ctx *gin.Context, verificationId string, phoneNumber interfaces.PhoneType) (error, int) {
	// Check if OTP is verified
	objectId, phone, err := decodeVerificationID(verificationId)
	if err != nil {
		return err, http.StatusUnprocessableEntity
	}

	otpRef, err := a.otpRepo.FindById(objectId)

	if err != nil {
		return errors.New("Couldn't find any reference to provided verification code"), http.StatusUnprocessableEntity
	}

	_, refPhoneNumber, _ := decodeVerificationID(otpRef.VerificationID)

	if !(refPhoneNumber.Code == phoneNumber.Code && phone.Number == phoneNumber.Number) {
		return errors.New("Provided verification ID is invalid"), http.StatusUnprocessableEntity
	}

	if !otpRef.Verified {
		return errors.New("Provided verification ID is not verified"), http.StatusUnprocessableEntity
	}

	if err := a.otpRepo.DeleteById(&otpRef.ID); err != nil {
		return errors.New("Something went wrong"), http.StatusInternalServerError
	}

	return nil, 0
}

func decodeVerificationID(verification_query_str string) (*primitive.ObjectID, *interfaces.PhoneType, error) {

	decodedStr, err := base64.StdEncoding.DecodeString(verification_query_str)
	if err != nil {

		return nil, nil, err
	}

	parsedQuery, err := url.ParseQuery(string(decodedStr))
	if err != nil {
		return nil, nil, err
	}

	objectId, err := primitive.ObjectIDFromHex(parsedQuery.Get("_id"))
	if err != nil {
		return nil, nil, err
	}

	phone := interfaces.PhoneType{
		Code:   strings.Replace("+"+parsedQuery.Get("phone_code"), " ", "", 1),
		Number: parsedQuery.Get("phone_number"),
	}

	return &objectId, &phone, nil
}

// returns access, refresh and error
func (a *UserAPI) getAccessAndRefreshTokenFromRequest(c *gin.Context) (string, string, error) {
	var token string
	var refresh_token string

	// Get auth token either from cookie nor Header
	cookie, err := c.Request.Cookie(string(constants.AUTH_TOKEN))
	if err != nil {
		// check in auth header as bearer
		authHeader := c.Request.Header.Get("Authorization")
		containsBearerToken := strings.HasPrefix(authHeader, "Bearer")
		if !containsBearerToken {
			token = ""
		} else {
			token = strings.Split(authHeader, "Bearer ")[1]
		}
	} else {
		token = cookie.Value
	}

	// Get refresh token either from cookie nor params
	refreshCookie, err := c.Request.Cookie(string(constants.REFRESH_TOKEN))
	if err != nil {
		// check in auth header as bearer
		refresh_token = c.Request.URL.Query().Get("refresh_token")
		if refresh_token == "" {
			return "", "", errors.New("refresh token is missing")
		}

	} else {
		refresh_token = refreshCookie.Value
	}

	return token, refresh_token, nil
}
