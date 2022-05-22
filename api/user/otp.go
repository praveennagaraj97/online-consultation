package userapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	otpmodel "github.com/praveennagaraj97/online-consultation/models/otp"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Send Verification code to login via - SMS Gateway
func (a *UserAPI) SendVerificationCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload interfaces.PhoneType

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer ctx.Request.Body.Close()

		// Generate OTP
		verifyCode := utils.GenerateRandomCode(6)

		res, err := a.otpRepo.CreateOne(&payload, &verifyCode)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		// Send OTP

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*otpmodel.OneTimePasswordEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    fmt.Sprintf("%s is your verification code to register with Online Consultation", verifyCode),
			},
		})

	}
}

// Verify Code
func (a *UserAPI) VerifyCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		verification_query_str, exists := ctx.Params.Get("verification_id")
		if !exists {
			api.SendErrorResponse(ctx, "Verification ID is missing", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, _, err := decodeVerificationID(verification_query_str)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		var payload userdto.VerifyCodeDTO

		if err = ctx.ShouldBind(&payload); err != nil || payload.VerifyCode == "" {
			api.SendErrorResponse(ctx, "Verification code is missing", http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		// Check OTP with database
		res, err := a.otpRepo.FindById(objectId)

		if err != nil {
			api.SendErrorResponse(ctx, "Failed to get requested entity", http.StatusNotAcceptable, nil)
			return
		}

		if err = res.DecodeVerificationCode(payload.VerifyCode); err != nil {

			// Delete if expired
			if res.ExpiryTime <= primitive.NewDateTimeFromTime(time.Now()) {
				a.otpRepo.DeleteById(&res.ID)
				api.SendErrorResponse(ctx, "Verification code is expired", http.StatusNotAcceptable, nil)
				return
			}

			// Delete if attemps limit is reached.
			if (constants.VerifyCodeAcceptedAttempts+1)-res.Attempts == 0 {
				a.otpRepo.DeleteById(&res.ID)
				api.SendErrorResponse(ctx, "Too many attempts", http.StatusNotAcceptable, nil)
				return
			}

			a.otpRepo.UpdateAttemptsCount(objectId)
			ctx.JSON(http.StatusUnprocessableEntity, serialize.InvalidVerificationCodeErrorResponse{
				Response: serialize.Response{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Entered verification code is invalid",
				},
				RemainingAttempts: constants.VerifyCodeAcceptedAttempts - res.Attempts,
			})
			return
		}

		if res.Verified {
			api.SendErrorResponse(ctx, "Your code is already verified", http.StatusNotAcceptable, nil)
			return
		}

		// Delete if expired
		if res.ExpiryTime <= primitive.NewDateTimeFromTime(time.Now()) {
			a.otpRepo.DeleteById(&res.ID)
			api.SendErrorResponse(ctx, "Verification code is expired", http.StatusNotAcceptable, nil)
			return
		}

		if (constants.VerifyCodeAcceptedAttempts+1)-res.Attempts == 0 {
			api.SendErrorResponse(ctx, "Too many attempts", http.StatusNotAcceptable, nil)
			return
		}

		a.otpRepo.UpdateStatus(objectId)

		ctx.JSON(http.StatusOK, serialize.DataResponse[*otpmodel.OneTimePasswordEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Verified successfully",
			},
		})

	}
}
