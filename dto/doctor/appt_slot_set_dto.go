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

func (a *AddNewAppointmentSlotSetDTO) Validate(timeZone string) *serialize.ErrorResponse {

	errors := map[string]string{}
	if timeZone == "" {
		errors["time_zone"] = "Time zone cannot be empty"
	}
	if len(a.SlotTimeRef) == 0 {
		errors["slot_time"] = "Slot time should have atleast one time slot"
	}

	timeLoc, err := time.LoadLocation(timeZone)
	if err != nil {
		errors["time_zone"] = "Invlaid time zone " + err.Error()
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

	errs, slots := validateSlots(a.SlotTimeRef, timeLoc)
	if errs != nil {
		return errs
	}

	a.SlotTimes = slots

	return nil
}

type UpdateAppointmentSlotSetDTO struct {
	Title     string               `json:"title,omitempty" form:"title,omitempty" bson:"title,omitempty"`
	IsDefault *bool                `json:"is_default,omitempty" form:"is_default,omitempty" bson:"is_default,omitempty"`
	SlotTimes []primitive.DateTime `json:"-" form:"-" bson:"slot_times,omitempty"`

	// Input Refs
	SlotTimeRef []string `json:"slot_time,omitempty" form:"slot_time,omitempty" bson:"-"`
}

func (a *UpdateAppointmentSlotSetDTO) Validate(timeZone string) *serialize.ErrorResponse {
	errors := map[string]string{}
	if timeZone == "" {
		errors["time_zone"] = "Time zone cannot be empty"
	}
	if len(a.SlotTimeRef) == 0 {
		errors["slot_time"] = "Slot time should have atleast one time slot"
	}

	timeLoc, err := time.LoadLocation(timeZone)
	if err != nil {
		errors["time_zone"] = "Invlaid time zone " + err.Error()
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

	errs, slots := validateSlots(a.SlotTimeRef, timeLoc)
	if errs != nil {
		return errs
	}

	a.SlotTimes = slots

	return nil
}

func validateSlots(SlotTimeRef []string, timeZone *time.Location) (*serialize.ErrorResponse, []primitive.DateTime) {
	errs := map[string]string{}

	var SlotTimes []primitive.DateTime

	visited := make(map[string]bool, 0)

	for i := 0; i < len(SlotTimeRef); i++ {

		// Set Random Date to align time zones
		t, err := time.ParseInLocation("2006-01-02 15:04:05", "1970-01-01 "+SlotTimeRef[i], timeZone)

		if err != nil {
			errs["slot_time"] = "Invalid Slot time, " + err.Error()
		}

		// check for duplicates
		if visited[SlotTimeRef[i]] {
			errs[SlotTimeRef[i]] = "Time slot is found as duplicate"
		} else {
			visited[SlotTimeRef[i]] = true
		}
		// Store time in UTC
		SlotTimes = append(SlotTimes, primitive.NewDateTimeFromTime(t.UTC()))
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
