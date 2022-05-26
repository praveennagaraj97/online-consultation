package admindto

import (
	"errors"
	"net/http"
	"strings"

	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

type AddNewAdminDTO struct {
	Name     string             `json:"name" form:"name"`
	UserName string             `json:"user_name" form:"user_name"`
	Email    string             `json:"email" form:"email"`
	Password string             `json:"password" form:"password"`
	Role     constants.UserType `json:"-" form:"-"`
}

type LoginDTO struct {
	UserName string `json:"user_name" form:"user_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UpdateAdminDTO struct {
	Name     string `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	Password string `bson:"password,omitempty"`
}

type UpdatePasswordDTO struct {
	Password        string `json:"password" form:"password"`
	NewPassword     string `json:"new_password" form:"new_password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

type ForgorPasswordDTO struct {
	Email string `json:"email" form:"email"`
}

type AdminIdDTO struct {
	ID string `json:"id" form:"id"`
}

type ResetPasswordDTO struct {
	NewPassword     string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

func (payload *AddNewAdminDTO) ValidateNewAdminDTO() *serialize.ErrorResponse {
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

func (payload *LoginDTO) ValidateAdminLoginDTO() *serialize.ErrorResponse {
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

func (payload *UpdatePasswordDTO) ValidateUpdatePasswordDTO() *serialize.ErrorResponse {

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

func (payload *ResetPasswordDTO) ValidateResetPasswordDTO() error {

	if payload.NewPassword == "" {
		return errors.New("Password cannot be empty")

	}

	if payload.ConfirmPassword == "" {
		return errors.New("Confirm Password cannot be empty")

	}

	if payload.ConfirmPassword != payload.NewPassword {
		return errors.New("Password didn't match")

	}

	if len(payload.NewPassword) < 6 {
		return errors.New("Password is too week, password should contain atleast 6 characters")
	}

	return nil
}
