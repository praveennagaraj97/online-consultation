package adminvalidator

import (
	"net/http"
	"strings"

	admindto "github.com/praveennagaraj97/online-consultation/dto/admin"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

func ValidateNewAdminDTO(payload *admindto.AddNewAdminDTO) *serialize.ErrorResponse {
	errors := map[string]string{}

	if strings.Trim(payload.Name, "") == "" {
		errors["name"] = "Name cannot be empty"
	}

	if err := validator.ValidateEmail(payload.Email); err != nil {
		errors["email"] = "Provided email is not valid"
	}

	if payload.UserName == "" {
		errors["user_name"] = "Username cannot be empty"
	}

	if payload.Password == "" {
		errors["password"] = "Password cannot be empty"
	}

	if payload.Password != "" && len(payload.Password) < 6 {
		errors["password"] = "Password is too week, password should contain atleast 6 characters"
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

func ValidateAdminLoginDTO(payload *admindto.LoginDTO) *serialize.ErrorResponse {
	errors := map[string]string{}

	if payload.Email == "" && payload.UserName == "" {
		errors["email"] = "Email and username cannot be empty"
		errors["user_name"] = "Username and email cannot be empty"
	}

	if payload.Password == "" {
		errors["password"] = "Password cannnot be empty"
	}

	if len(errors) > 0 {
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

func ValidateUpdatePasswordDTO(payload *admindto.UpdatePasswordDTO) *serialize.ErrorResponse {

	errors := map[string]string{}

	if payload.NewPassword == "" {
		errors["new_password"] = "New Password cannot be empty"
	}

	if payload.ConfirmPassword == "" {
		errors["confirm_password"] = "Confirm Password cannot be empty"
	}

	if payload.Password == "" {
		errors["password"] = "Provide your existing password"
	}

	if len(errors) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errors,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	if payload.NewPassword != payload.ConfirmPassword {
		errors["confirm_password"] = "Password didn't match"
	}

	if payload.Password == payload.NewPassword {
		errors["new_password"] = "Current password and new password cannot be same"
	}

	if len(errors) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errors,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	if len(payload.NewPassword) < 6 {
		errors["new_password"] = "Password is too week, password should contain atleast 6 characters"
	}

	if len(errors) > 0 {
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
