package authvalidator

import (
	"net/http"
	"strings"

	userdto "github.com/praveennagaraj97/online-consultation/dto"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

func ValidateRegisterDTO(payload *userdto.RegisterDTO) *serialize.ErrorResponse {

	errors := map[string]string{}

	if strings.Trim(payload.Name, "") == "" {
		errors["name"] = "Name cannot be empty"
	}

	if err := validator.ValidateEmail(payload.Email); err != nil {
		errors["email"] = "Provided email is not valid"
	}

	if payload.PhoneCode == "" {
		errors["phone_code"] = "Phone code cannot be empty"
	}

	if payload.PhoneNumber == "" {
		errors["phone_number"] = "Phone number cannot be empty"
	}

	if payload.Gender == "" {
		errors["gender"] = "Gender cannot be empty"
	}

	if payload.VerificationId == "" {
		errors["verification_id"] = "Verification ID cannot be empty"
	}

	if len(errors) != 0 {
		return &serialize.ErrorResponse{
			Errors: &errors,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	return nil
}

func ValidateSignInWithPhoneDTO(payload *userdto.SignInWithPhoneDTO) *serialize.ErrorResponse {

	errors := map[string]string{}

	if payload.PhoneCode == "" {
		errors["phone_code"] = "Phone code cannot be empty"
	}

	if payload.PhoneNumber == "" {
		errors["phone_number"] = "Phone number cannot be empty"
	}

	if payload.VerificationId == "" {
		errors["verification_id"] = "Verification ID cannot be empty"
	}

	if len(errors) != 0 {
		return &serialize.ErrorResponse{
			Errors: &errors,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	return nil
}
