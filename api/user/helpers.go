package userapi

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
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
