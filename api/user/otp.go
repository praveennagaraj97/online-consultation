package userapi

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	otpmodel "github.com/praveennagaraj97/online-consultation/models/otp"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
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

		// Validate
		if errs := validator.ValidatePhoneNumber(payload); errs != nil {
			api.SendErrorResponse(ctx, "Given data is invalid", http.StatusUnprocessableEntity, errs)
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
		if _, err := a.appConf.AwsUtils.SendTextSMS(&interfaces.SMSType{
			Message: fmt.Sprintf("%s is your verification code for Online Consultation", verifyCode),
			To:      &payload,
		}); err != nil {
			log.Default().Println(err.Error())
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*otpmodel.OneTimePasswordEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "A text with verification code has been sent to your mobile number",
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

		objectId, _, err := utils.DecodeVerificationID(verification_query_str)

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

// Resend verification code once code is expired
func (a *UserAPI) ResendVerificationCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Verification ID
		vId, exists := ctx.Params.Get("id")
		if !exists {
			api.SendErrorResponse(ctx, "Verification Id is required", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, phone, err := utils.DecodeVerificationID(vId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
		}

		data, err := a.otpRepo.FindById(objectId)
		if err != nil {
			api.SendErrorResponse(ctx, "Couldn't find any matching reference with given verification ID", http.StatusNotFound, nil)
			return
		}

		if data.ExpiryTime > primitive.NewDateTimeFromTime(time.Now()) {
			api.SendErrorResponse(ctx, "Verification is not expired", http.StatusNotAcceptable, nil)
			return
		}

		// Delete the reference
		if err := a.otpRepo.DeleteById(objectId); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Generate OTP
		verifyCode := utils.GenerateRandomCode(6)

		res, err := a.otpRepo.CreateOne(phone, &verifyCode)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		// Send OTP
		if _, err := a.appConf.AwsUtils.SendTextSMS(&interfaces.SMSType{
			Message: fmt.Sprintf("%s is your verification code for Online Consultation", verifyCode),
			To:      phone,
		}); err != nil {
			log.Default().Println(err.Error())
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*otpmodel.OneTimePasswordEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "A text with verification code has been sent to your mobile number",
			},
		})

	}
}
