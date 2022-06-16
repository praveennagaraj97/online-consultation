package appointmentslotsdto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddAppointmentSlotDTO struct {
	Dates       []primitive.DateTime `json:"-" form:"-"`
	SlotSetId   *primitive.ObjectID  `json:"" form:"-"`
	Days        uint16               `json:"days" form:"days"`
	ExcludeDays []string             `json:"exclude_days" form:"exclude_days"`

	SlotSetIdRef string   `json:"slot_set_id" form:"slot_set_id"`
	DatesRef     []string `json:"date" form:"date"`
}

func (a *AddAppointmentSlotDTO) Validate() *serialize.ErrorResponse {
	errs := map[string]string{}

	if a.SlotSetIdRef != "" {
		objectId, err := primitive.ObjectIDFromHex(a.SlotSetIdRef)
		if err != nil {
			errs["slot_set_id"] = "Slot Id should be valid primitive ObjectId"
		}
		a.SlotSetId = &objectId
	}

	if len(a.DatesRef) == 0 && a.Days == 0 {
		errs["date"] = "Date and days cannot be empty"
		errs["days"] = "Days and date Cannot be empty"
	}

	if a.Days > 365 {
		errs["days"] = "Days should be less than 365"
	}

	if len(errs) != 0 {
		return &serialize.ErrorResponse{
			Errors: &errs,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	if a.Days != 0 {
		dates, ers := validateExcludeDaysAndGetDates(a.ExcludeDays, a.Days)
		if ers != nil {
			errs = *ers
		}
		a.Dates = dates
	} else {
		dates, ers := validateAndGetDates(a.DatesRef)
		errs = *ers
		a.Dates = dates
	}

	if len(errs) != 0 {
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
