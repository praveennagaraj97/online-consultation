package doctordto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/serialize"
)

type AddDoctorHospitalDTO struct {
	Name      string  `json:"name" form:"name"`
	City      string  `json:"city" form:"city"`
	Country   string  `json:"country" form:"country"`
	Address   string  `json:"address" form:"address"`
	Latitude  float64 `json:"latitude" form:"latitude,omitempty"`
	Longitude float64 `json:"longitude" form:"longitude,omitempty"`
}

func (a *AddDoctorHospitalDTO) Validate() *serialize.ErrorResponse {
	errs := map[string]string{}

	if a.Name == "" {
		errs["name"] = "Hospital name cannot be empty"
	}

	if a.City == "" {
		errs["city"] = "City name cannot be empty"
	}

	if a.Country == "" {
		errs["country"] = "Country name cannot be empty"
	}

	if a.Address == "" {
		errs["address"] = "Hospital Address cannot be empty"
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