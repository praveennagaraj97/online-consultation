package doctordto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/pkg/validator"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddNewDoctorDTO struct {
	Name              string `json:"name" form:"name"`
	Email             string `json:"email" form:"email"`
	PhoneCode         string `json:"phone_code" form:"phone_code"`
	PhoneNumber       string `json:"phone_number" form:"phone_number"`
	ProfessionalTitle string `json:"professional_title" form:"professional_title"`
	Education         string `json:"education" form:"education"`
	Experience        uint8  `json:"experience" form:"experience"`
	ProfilePicWidth   uint64 `json:"profile_pic_width" form:"profile_pic_width"`
	ProfilePicHeight  uint64 `json:"profile_pic_height" form:"profile_pic_height"`

	// Input ref fields
	HospitalId         string   `json:"hospital_id" form:"hospital_id"`
	ConsultationTypeId string   `json:"consultation_type_id" form:"consultation_type_id"`
	SpecialityId       string   `json:"speciality_id" form:"speciality_id"`
	SpokenLanguagesIds []string `json:"spoken_language_id" form:"spoken_language_id"`

	Hospital         *primitive.ObjectID  `json:"-" bson:"-"`
	ConsultationType *primitive.ObjectID  `json:"-" form:"-"`
	Speciality       *primitive.ObjectID  `json:"-" form:"-"`
	SpokenLanguages  []primitive.ObjectID `json:"-" form:"-"`
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

	if a.ConsultationTypeId == "" {
		errs["consultation_type_id"] = "Consultation Id cannot be empty"
	}

	if a.HospitalId == "" {
		errs["hospital_id"] = "Hospital Id cannot be empty"
	}

	if a.ProfessionalTitle == "" {
		errs["professional_title"] = "Professional title cannot be empty"
	}

	if a.Education == "" {
		errs["education"] = "Education cannot be empty"
	}

	if len(a.SpokenLanguagesIds) == 0 {
		errs["spoken_language_id"] = "Spoken languages cannot be empty"
	}

	if a.ConsultationTypeId != "" {
		objectId, err := primitive.ObjectIDFromHex(a.ConsultationTypeId)
		if err != nil {
			errs["consultation_type_id"] = "Consultation Id should be valid primitive objectId"
		} else {
			a.ConsultationType = &objectId
		}
	}

	if a.ProfilePicWidth == 0 || a.ProfilePicHeight == 0 {
		a.ProfilePicWidth = 110
		a.ProfilePicHeight = 110
	}

	if a.HospitalId != "" {
		objectId, err := primitive.ObjectIDFromHex(a.HospitalId)
		if err != nil {
			errs["hospital_id"] = "Hospital Id should be valid primitive objectId"
		} else {
			a.Hospital = &objectId
		}
	}

	if a.SpecialityId != "" {
		objectId, err := primitive.ObjectIDFromHex(a.SpecialityId)
		if err != nil {
			errs["speciality_id"] = "Speciality Id should be valid primitive objectId"
		} else {
			a.Speciality = &objectId
		}
	}

	// Validate language ids
	for i := 0; i < len(a.SpokenLanguagesIds); i++ {
		objectId, err := primitive.ObjectIDFromHex(a.SpokenLanguagesIds[i])
		if err != nil {
			errs["spoken_language_id"] = "Spoken language Id should be valid primitive Object Id"
		}
		a.SpokenLanguages = append(a.SpokenLanguages, objectId)
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

type EditDoctorDTO struct {
	Name              string `json:"name,omitempty" form:"name,omitempty"  bson:"name,omitempty"`
	ProfessionalTitle string `json:"professional_title,omitempty" form:"professional_title,omitempty" bson:"professional_title,omitempty"`
	Education         string `json:"education,omitempty" form:"education,omitempty" bson:"education,omitempty"`
	Experience        uint8  `json:"experience,omitempty" form:"experience,omitempty" bson:"experience,omitempty"`
	ProfilePicWidth   uint64 `json:"profile_pic_width,omitempty" form:"profile_pic_width,omitempty" bson:"profile_pic_width,omitempty"`
	ProfilePicHeight  uint64 `json:"profile_pic_height,omitempty" form:"profile_pic_height,omitempty" bson:"profile_pic_height,omitempty"`

	// Input ref fields
	HospitalId         string   `json:"hospital_id,omitempty" form:"hospital_id,omitempty" bson:"-"`
	ConsultationTypeId string   `json:"consultation_type_id,omitempty" form:"consultation_type_id,omitempty" bson:"-"`
	SpecialityId       string   `json:"speciality_id,omitempty" form:"speciality_id,omitempty" bson:"-"`
	SpokenLanguagesIds []string `json:"spoken_language_id,omitempty" form:"spoken_language_id,omitempty" bson:"-"`

	Hospital         *primitive.ObjectID  `json:"-" form:"-" bson:"hospital_id,omitempty"`
	ConsultationType *primitive.ObjectID  `json:"-" form:"-" bson:"consultation_type_id,omitempty"`
	Speciality       *primitive.ObjectID  `json:"-" form:"-" bson:"speciality_id,omitempty"`
	SpokenLanguages  []primitive.ObjectID `json:"-" form:"-" bson:"languages_ids,omitempty"`
}

func (a *EditDoctorDTO) Validate() *serialize.ErrorResponse {
	errs := map[string]string{}

	if a.ConsultationTypeId != "" {
		objectId, err := primitive.ObjectIDFromHex(a.ConsultationTypeId)
		if err != nil {
			errs["consultation_type_id"] = "Consultation Id should be valid primitive objectId"
		} else {
			a.ConsultationType = &objectId
		}
	}

	if a.HospitalId != "" {
		objectId, err := primitive.ObjectIDFromHex(a.HospitalId)
		if err != nil {
			errs["hospital_id"] = "Hospital Id should be valid primitive objectId"
		} else {
			a.Hospital = &objectId
		}
	}

	if a.SpecialityId != "" {
		objectId, err := primitive.ObjectIDFromHex(a.SpecialityId)
		if err != nil {
			errs["speciality_id"] = "Speciality Id should be valid primitive objectId"
		} else {
			a.Speciality = &objectId
		}
	}

	// Validate language ids
	if len(a.SpokenLanguages) > 0 {
		for i := 0; i < len(a.SpokenLanguagesIds); i++ {
			objectId, err := primitive.ObjectIDFromHex(a.SpokenLanguagesIds[i])
			if err != nil {
				errs["spoken_language_id"] = "Spoken language Id should be valid primitive Object Id"
			}
			a.SpokenLanguages = append(a.SpokenLanguages, objectId)
		}
	} else {
		a.SpokenLanguages = nil
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
