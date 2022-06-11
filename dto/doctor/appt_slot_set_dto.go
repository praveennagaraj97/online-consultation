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

	errs := map[string]string{}

	if len(a.SlotTimeRef) == 0 {
		errs["slot_time"] = "Slot time should have atleast one time slot"
	} else {

		visited := make(map[string]bool, 0)

		for i := 0; i < len(a.SlotTimeRef); i++ {

			t, err := time.Parse("15:04:05", a.SlotTimeRef[i])

			if err != nil {
				errs["slot_time"] = "Invalid Slot time, " + err.Error()
			}

			// check for duplicates
			if visited[a.SlotTimeRef[i]] == true {
				errs[a.SlotTimeRef[i]] = "Time slot is found as duplicate"
			} else {
				visited[a.SlotTimeRef[i]] = true
			}

			a.SlotTimes = append(a.SlotTimes, primitive.NewDateTimeFromTime(t))
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
