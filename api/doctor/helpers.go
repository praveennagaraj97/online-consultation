package doctorapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/utils"
)

func (a *DoctorAPI) validateVerifyCode(ctx *gin.Context, verificationId string, phoneNumber interfaces.PhoneType) (error, int) {
	// Check if OTP is verified
	objectId, phone, err := utils.DecodeVerificationID(verificationId)
	if err != nil {
		return err, http.StatusUnprocessableEntity
	}

	otpRef, err := a.otpRepo.FindById(objectId)

	if err != nil {
		return errors.New("Couldn't find any reference to provided verification code"), http.StatusUnprocessableEntity
	}

	_, refPhoneNumber, _ := utils.DecodeVerificationID(otpRef.VerificationID)

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
