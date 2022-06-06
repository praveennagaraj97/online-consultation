package doctordto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddNewDoctorDTO struct {
	Name              string              `json:"name" form:"name"`
	Email             string              `json:"email" form:"email"`
	PhoneCode         string              `json:"phone_code" form:"phone_code"`
	PhoneNumber       string              `json:"phone_number" form:"phone_number"`
	Type              string              `json:"type" form:"type"`
	ProfessionalTitle string              `json:"professional_title" form:"professional_title"`
	Experience        uint8               `json:"experience" form:"experience"`
	ConsultationType  *primitive.ObjectID `json:"-" form:"-"`
	ProfilePicWidth   uint64              `json:"profile_pic_width" form:"profile_pic_width"`
	ProfilePicHeight  uint64              `json:"profile_pic_height" form:"profile_pic_height"`
}

func (a *AddNewDoctorDTO) Validate() *serialize.ErrorResponse {
	errs := map[string]string{}

	if a.Name == "" {
		errs["name"] = "Doctor name cannot be empty"
	}

	if a.Email == "" {
		errs["email"] = "Email cannot be empty"
	}

	if a.Email != "" && validator.ValidateEmail(a.Email) != nil {
		errs["email"] = "Provided email is invalid"
	}

	if a.PhoneCode == "" {
		errs["phone_code"] = "Phone code cannot be empty"
	}

	if a.PhoneNumber == "" {
		errs["phone_number"] = "Phone number cannot be empty"
	}

	if a.Type == "" {
		errs["type"] = "Type cannot be empty"
	}

	if a.ProfessionalTitle == "" {
		errs["professional_title"] = "Professional title cannot be empty"
	}

	if a.Type != "" {
		objectId, err := primitive.ObjectIDFromHex(a.Type)
		if err != nil {
			errs["type"] = "Type should be valid consultation id"
		} else {
			a.ConsultationType = &objectId
		}
	}

	if a.ProfilePicWidth == 0 || a.ProfilePicHeight == 0 {
		a.ProfilePicWidth = 110
		a.ProfilePicHeight = 110
	}

	if len(errs) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errs,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			}}
	}

	return nil
}
