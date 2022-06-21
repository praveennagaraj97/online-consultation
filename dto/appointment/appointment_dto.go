package appointmentdto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookScheduledAppointmentDTO struct {
	DoctorId          *primitive.ObjectID `json:"-" bson:"-"`
	AppointmentSlotId *primitive.ObjectID `json:"-" form:"-"`
	RelativeId        *primitive.ObjectID `json:"-,omitempty" form:"-,omitempty"`

	DoctorIdRef          string `json:"doctor_id" form:"doctor_id"`
	AppointmentSlotIdRef string `json:"appointment_slot_id" form:"appointment_slot_id"`
	RelativeIdRef        string `json:"relative_id" form:"relative_id"`
}

func (a *BookScheduledAppointmentDTO) Validate() *serialize.ErrorResponse {
	errs := map[string]string{}

	if a.DoctorIdRef == "" {
		errs["doctor_id"] = "Doctor Id cannot be empty"
	} else {
		objectId, err := primitive.ObjectIDFromHex(a.DoctorIdRef)
		if err != nil {
			errs["doctor_id"] = "Doctor Id should be valid primitive ObjectId"
		} else {
			a.DoctorId = &objectId
		}
	}

	if a.AppointmentSlotIdRef == "" {
		errs["appointment_slot_id"] = "Appointment slot Id cannot be empty"
	} else {
		objectId, err := primitive.ObjectIDFromHex(a.AppointmentSlotIdRef)
		if err != nil {
			errs["appointment_slot_id"] = "Appointment slot Id should be valid primitive ObjectId"
		} else {
			a.AppointmentSlotId = &objectId
		}
	}

	// Validate if appointment is for relative
	if a.RelativeIdRef != "" {
		objectId, err := primitive.ObjectIDFromHex(a.RelativeIdRef)
		if err != nil {
			errs["relative_id"] = "Relative Id should be valid primitive ObjectId"
		} else {
			a.RelativeId = &objectId
		}
	}

	if len(errs) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errs,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	return nil
}
