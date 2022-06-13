package doctordto

import (
	"net/http"
	"time"

	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddNewAppointmentSlotSetDTO struct {
	Title     string               `json:"title" form:"title"` // Optional
	IsDefault bool                 `json:"is_default" form:"is_default"`
	SlotTimes []primitive.DateTime `json:"-" form:"-"`

	// Input Refs
	SlotTimeRef []string `json:"slot_time" form:"slot_time"`
}

func (a *AddNewAppointmentSlotSetDTO) Validate() *serialize.ErrorResponse {

	errs, slots := validateSlots(a.SlotTimeRef)
	if errs != nil {
		return errs
	}

	a.SlotTimes = slots

	return nil
}

type UpdateAppointmentSlotSetDTO struct {
	Title     string               `json:"title,omitempty" form:"title,omitempty" bson:"title,omitempty"`
	SlotTimes []primitive.DateTime `json:"-" form:"-" bson:"slot_times,omitempty"`

	// Input Refs
	SlotTimeRef []string `json:"slot_time,omitempty" form:"slot_time,omitempty" bson:"-"`
}

func (a *UpdateAppointmentSlotSetDTO) Validate() *serialize.ErrorResponse {

	errs, slots := validateSlots(a.SlotTimeRef)
	if errs != nil {
		return errs
	}

	a.SlotTimes = slots

	return nil
}

func validateSlots(SlotTimeRef []string) (*serialize.ErrorResponse, []primitive.DateTime) {
	errs := map[string]string{}

	var SlotTimes []primitive.DateTime

	if len(SlotTimeRef) == 0 {
		errs["slot_time"] = "Slot time should have atleast one time slot"
	} else {

		visited := make(map[string]bool, 0)

		for i := 0; i < len(SlotTimeRef); i++ {

			t, err := time.Parse("15:04:05", SlotTimeRef[i])

			if err != nil {
				errs["slot_time"] = "Invalid Slot time, " + err.Error()
			}

			// check for duplicates
			if visited[SlotTimeRef[i]] == true {
				errs[SlotTimeRef[i]] = "Time slot is found as duplicate"
			} else {
				visited[SlotTimeRef[i]] = true
			}

			SlotTimes = append(SlotTimes, primitive.NewDateTimeFromTime(t))
		}

	}

	if len(errs) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errs,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}, nil
	}

	return nil, SlotTimes

}
