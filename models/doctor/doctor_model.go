package doctormodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	hospitalmodel "github.com/praveennagaraj97/online-consultation/models/hospital"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorQualificationEntity struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	DoctorId        primitive.ObjectID `json:"-" bson:"doctor_id"`
	Name            string             `json:"name" bson:"name"`
	InstituteName   string             `json:"institute_name" bson:"institute_name"`
	ProcurementYear primitive.DateTime `json:"procurement_year" bson:"procurement_year"`
}

type DoctorEntity struct {
	ID                primitive.ObjectID    `json:"id" bson:"_id"`
	Name              string                `json:"name" bson:"name"`
	Email             string                `json:"email" bson:"email"`
	Phone             *interfaces.PhoneType `json:"phone" bson:"phone"`
	ProfessionalTitle string                `json:"professional_title" bson:"professional_title"`
	Experience        uint8                 `json:"experience" bson:"experience"`
	ProfilePic        *interfaces.ImageType `json:"profile_pic" bson:"profile_pic"`
	RefreshToken      string                `json:"-" bson:"refresh_token"`

	// Populate fields
	ConsultationType string                        `json:"consultation_type,omitempty" bson:"consultation_type,omitempty"`
	Hospital         *hospitalmodel.HospitalEntity `json:"hospital" bson:"hospital,omitempty"`

	// reference fields
	TypeId     *primitive.ObjectID `json:"-" bson:"type"`
	HospitalId *primitive.ObjectID `json:"-" bson:"hospital"`
}
