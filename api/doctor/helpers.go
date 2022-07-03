package doctorapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/utils"
)

func (a *DoctorAPI) validateVerifyCode(ctx *gin.Context, verificationId string, phoneNumber interfaces.PhoneType) (int, error) {
	// Check if OTP is verified
	objectId, phone, err := utils.DecodeVerificationID(verificationId)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	otpRef, err := a.otpRepo.FindById(objectId)

	if err != nil {
		return http.StatusUnprocessableEntity, errors.New("couldn't find any reference to provided verification code")
	}

	_, refPhoneNumber, _ := utils.DecodeVerificationID(otpRef.VerificationID)

	if !(refPhoneNumber.Code == phoneNumber.Code && phone.Number == phoneNumber.Number) {
		return http.StatusUnprocessableEntity, errors.New("provided verification ID is invalid")
	}

	if !otpRef.Verified {
		return http.StatusUnprocessableEntity, errors.New("provided verification ID is not verified")
	}

	if err := a.otpRepo.DeleteById(&otpRef.ID); err != nil {
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	return 0, nil
}
