package userdto

import (
	"net/http"
	"strings"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterDTO struct {
	Name           string             `json:"name" form:"name"`
	Email          string             `json:"email" form:"email"`
	PhoneCode      string             `json:"phone_code" form:"phone_code"`
	PhoneNumber    string             `json:"phone_number" form:"phone_number"`
	Gender         string             `json:"gender" form:"gender"`
	VerificationId string             `json:"verification_id" form:"verification_id"`
	DOB            primitive.DateTime `json:"-" form:"-"`
	DOBRef         string             `json:"date_of_birth" form:"date_of_birth"`
}

type VerifyCodeDTO struct {
	VerifyCode string `json:"verify_code" form:"verify_code"`
}

type UpdateUserDTO struct {
	Name          string              `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	Email         string              `json:"-" form:"-" bson:"email,omitempty"`
	PhoneCode     string              `json:"-" form:"-" bson:"phone_code,omitempty"`
	PhoneNumber   string              `json:"-" form:"-" bson:"phone_number,omitempty"`
	DOB           *primitive.DateTime `json:"-" form:"-" bson:"date_of_birth,omitempty"`
	DOBRef        string              `json:"date_of_birth" form:"date_of_birth" bson:"-"`
	Gender        string              `json:"gender,omitempty" form:"gender,omitempty" bson:"gender,omitempty"`
	EmailVerified bool                `json:"-" form:"-" bson:"email_verified,omitempty"`
}

type SignInWithPhoneDTO struct {
	VerificationId string `json:"verification_id" form:"verification_id"`
	PhoneCode      string `json:"phone_code" form:"phone_code"`
	PhoneNumber    string `json:"phone_number" form:"phone_number"`
}

type SignInWithEmailLinkDTO struct {
	Email      string `json:"email" form:"email"`
	RedirectTo string `json:"redirect_to" form:"redirect_to"`
}

type RequestEmailVerifyDTO struct {
	Email      string `json:"email" form:"email"`
	RedirectTo string `json:"redirect_to" form:"redirect_to"`
}

type AddOrEditRelativeDTO struct {
	Name        string                `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	Email       string                `json:"email,omitempty" form:"email,omitempty" bson:"email,omitempty"`
	PhoneCode   string                `json:"phone_code,omitempty" form:"phone_code,omitempty" bson:"phone_code,omitempty"`
	PhoneNumber string                `json:"phone_number,omitempty" form:"phone_number,omitempty" bson:"phone_number,omitempty"`
	DateOfBirth *primitive.DateTime   `json:"-" form:"-" bson:"date_of_birth,omitempty"`
	Gender      string                `json:"gender,omitempty" form:"gender,omitempty" bson:"gender,omitempty"`
	Relation    string                `json:"relation,omitempty" form:"relation,omitempty" bson:"relation,omitempty"`
	UserId      *primitive.ObjectID   `json:"-,omitempty" form:"-,omitempty"`
	Phone       *interfaces.PhoneType `json:"-" form:"-" bson:"phone,omitempty"`

	//
	DOBRef string `json:"date_of_birth,omitempty" form:"date_of_birth,omitempty" bson:"-"`
}

type AddOrEditDeliveryAddressDTO struct {
	Name        string                `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	Address     string                `json:"address,omitempty" form:"address,omitempty" bson:"address,omitempty"`
	State       string                `json:"state,omitempty" form:"state,omitempty" bson:"state,omitempty"`
	Locality    string                `json:"locality,omitempty" form:"locality,omitempty" bson:"locality,omitempty"`
	PinCode     string                `json:"pincode,omitempty" form:"pincode,omitempty" bson:"pincode,omitempty"`
	PhoneCode   string                `json:"phone_code,omitempty" form:"phone_code,omitempty"`
	PhoneNumber string                `json:"phone_number,omitempty" form:"phone_number,omitempty"`
	UserId      *primitive.ObjectID   `json:"-,omitempty" form:"-,omitempty" bson:"user_id,omitempty"`
	IsDefault   bool                  `bson:"is_default,omitempty"`
	Phone       *interfaces.PhoneType `form:"-" json:"-" bson:"phone,omitempty"`
}

func (payload *AddOrEditRelativeDTO) ValidateRelativeDTO(timeZone string) *serialize.ErrorResponse {

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

	if payload.DOBRef != "" {
		timeLoc, err := time.LoadLocation(timeZone)
		if err != nil {
			errors["time_zone_header"] = "Provided time zone is invalid"
		}
		t, err := time.ParseInLocation("2006-01-02", payload.DOBRef, timeLoc)
		if err != nil {
			errors["date_of_birth"] = err.Error()
		} else {
			dob := primitive.NewDateTimeFromTime(t.UTC())
			payload.DateOfBirth = &dob
		}
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

func (payload *AddOrEditRelativeDTO) CompareAndValidateRelativeDTOWithUserData(user *usermodel.UserEntity) *serialize.ErrorResponse {

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

func (payload *AddOrEditDeliveryAddressDTO) ValidateUserDeliveryAddressDTO() *serialize.ErrorResponse {

	var errors = map[string]string{}

	if payload.Name == "" {
		errors["name"] = "Name cannot be empty"
	}

	if payload.Address == "" {
		errors["address"] = "Address cannot be empty"
	}

	if payload.State == "" {
		errors["state"] = "state cannot be empty"
	}

	if payload.Locality == "" {
		errors["locality"] = "Locality / Town cannot be empty"
	}

	if payload.PinCode == "" {
		errors["pincode"] = "Pincode cannot be empty"
	}

	if payload.PhoneCode == "" {
		errors["phone_code"] = "Phone code cannot be empty"
	}

	if payload.PhoneNumber == "" {
		errors["phone_number"] = "Phone number cannot be empty"
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

func (payload *RegisterDTO) ValidateRegisterDTO(timeZone string) *serialize.ErrorResponse {

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

	if payload.DOBRef != "" {
		timeLoc, err := time.LoadLocation(timeZone)
		if err != nil {
			errors["time_zone_header"] = "Provided time zone is invalid"
		}
		t, err := time.ParseInLocation("2006-01-02", payload.DOBRef, timeLoc)
		if err != nil {
			errors["date_of_birth"] = err.Error()
		} else {
			payload.DOB = primitive.NewDateTimeFromTime(t.UTC())
		}
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

func (payload *SignInWithPhoneDTO) ValidateSignInWithPhoneDTO() *serialize.ErrorResponse {

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
