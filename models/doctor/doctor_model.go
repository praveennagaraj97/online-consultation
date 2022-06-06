package doctormodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	languagesmodel "github.com/praveennagaraj97/online-consultation/models/languages"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorHospital struct {
	ID       primitive.ObjectID       `json:"id" bson:"_id"`
	Name     string                   `json:"name" bson:"name"`
	City     string                   `json:"city" bson:"city"`
	Country  string                   `json:"country" bson:"country"`
	Location *interfaces.LocationType `json:"location,omitempty" bson:"location,omitempty"`
}

type DoctorQualificationEntity struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	DoctorId        primitive.ObjectID `json:"-" bson:"doctor_id"`
	Name            string             `json:"name" bson:"name"`
	InstituteName   string             `json:"institute_name" bson:"institute_name"`
	ProcurementYear primitive.DateTime `json:"procurement_year" bson:"procurement_year"`
}

type DoctorSpokenLanguagesEntity struct {
	ID       primitive.ObjectID             `json:"id" bson:"_id"`
	DoctorId primitive.ObjectID             `json:"doctor_id" bson:"doctor_id"`
	Language *languagesmodel.LanguageEntity `json:"language" bson:"language"`
}

type DoctorSpecialitiesEntity struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	DoctorId   primitive.ObjectID `json:"doctor_id" bson:"doctor_id"`
	Speciality primitive.ObjectID `json:"speciality_id" bson:"speciality_id"`
}

type DoctorEntity struct {
	ID                primitive.ObjectID    `json:"id" bson:"_id"`
	Name              string                `json:"name" bson:"name"`
	Email             string                `json:"email" bson:"email"`
	Phone             *interfaces.PhoneType `json:"phone" bson:"phone"`
	Type              *primitive.ObjectID   `json:"-" bson:"type"`
	ConsultationType  string                `json:"consultation_type,omitempty" bson:"consultation_type,omitempty"`
	ProfessionalTitle string                `json:"professional_title" bson:"professional_title"`
	Experience        uint8                 `json:"experience" bson:"experience"`
	ProfilePic        *interfaces.ImageType `json:"profile_pic" bson:"profile_pic"`
	Hospital          *DoctorHospital       `json:"hospital" bson:"hospital"`
	RefreshToken      string                `json:"-" bson:"refresh_token"`
}