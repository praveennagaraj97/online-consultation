package uservalidator

import (
	"net/http"
	"strings"

	userdto "github.com/praveennagaraj97/online-consultation/dto"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

func ValidateRelativeDTO(payload *userdto.AddOrEditRelativeDTO) *serialize.ErrorResponse {

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

	if payload.Relation == "" {
		errors["relation"] = "Relation cannot be empty"
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

func CompareAndValidateRelativeDTOWithUserData(payload *userdto.AddOrEditRelativeDTO, user *usermodel.UserEntity) *serialize.ErrorResponse {

	errors := map[string]string{}

	if payload.Email == user.Email {
		errors["email"] = "Relative email cannot be same as your email"
	}

	if payload.PhoneCode == user.PhoneNumber.Code && payload.PhoneNumber == user.PhoneNumber.Number {
		errors["phone_number"] = "Relative phone number cannot be same as your number"
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
