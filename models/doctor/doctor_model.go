package doctormodel

import (
	"time"

	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	appointmentslotmodel "github.com/praveennagaraj97/online-consultation/models/appointment_slot"
	hospitalmodel "github.com/praveennagaraj97/online-consultation/models/hospital"
	languagesmodel "github.com/praveennagaraj97/online-consultation/models/languages"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorEntity struct {
	ID                primitive.ObjectID    `json:"id" bson:"_id"`
	Name              string                `json:"name" bson:"name"`
	Email             string                `json:"email" bson:"email"`
	Phone             *interfaces.PhoneType `json:"phone" bson:"phone"`
	ProfessionalTitle string                `json:"professional_title" bson:"professional_title"`
	Education         string                `json:"education" bson:"education"`
	Experience        uint8                 `json:"experience" bson:"experience"`
	ProfilePic        *interfaces.ImageType `json:"profile_pic" bson:"profile_pic"`
	RefreshToken      string                `json:"-" bson:"refresh_token"`
	IsActive          bool                  `json:"is_active" bson:"is_active"`

	// Populate fields
	Speciality        string                                      `json:"speciality,omitempty" bson:"speciality,omitempty"`
	ConsultationType  string                                      `json:"consultation_type,omitempty" bson:"consultation_type,omitempty"`
	Hospital          *hospitalmodel.HospitalEntity               `json:"hospital,omitempty" bson:"hospital,omitempty"`
	SpokenLanguages   []languagesmodel.LanguageEntity             `json:"spoken_languages,omitempty" bson:"spoken_languages,omitempty"`
	NextAvailableSlot *appointmentslotmodel.AppointmentSlotEntity `json:"next_available_slot" bson:"next_available_slot,omitempty"`

	// reference fields
	ConsultationTypeId *primitive.ObjectID  `json:"-" bson:"consultation_type_id"`
	HospitalId         *primitive.ObjectID  `json:"-" bson:"hospital_id"`
	SpecialityId       *primitive.ObjectID  `json:"-" bson:"speciality_id"`
	SpokenLanguagesIds []primitive.ObjectID `json:"-" bson:"languages_ids"`
}

func (a *DoctorEntity) GetAccessAndRefreshToken(acessExpires bool) (string, string, int, error) {

	var access, refresh string
	var err error
	var accessTime int = constants.CookieRefreshExpiryTime

	if acessExpires {
		accessTime = constants.CookieAccessExpiryTime
		access, err = tokens.GenerateTokenWithExpiryTimeAndType(a.ID.Hex(),
			time.Now().Add(time.Minute*constants.JWT_AccessTokenExpiry).Unix(), "access", "doctor")
	} else {
		access, err = tokens.GenerateNoExpiryTokenWithCustomType(a.ID.Hex(), "access", "doctor")

	}
	refresh, err = tokens.GenerateNoExpiryTokenWithCustomType(a.ID.Hex(), "refresh", "doctor")

	return access, refresh, accessTime, err
}
