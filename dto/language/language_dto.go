package languagedto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/serialize"
)

type AddLanguageDTO struct {
	Name       string `json:"name" form:"name"`
	LocaleName string `json:"locale_name" form:"locale_name"`
}

func (a *AddLanguageDTO) Validate() *serialize.ErrorResponse {
	errors := map[string]string{}

	if a.Name == "" {
		errors["name"] = "Name of the language cannot be empty"
	}

	if a.LocaleName == "" {
		errors["locale_name"] = "A localized verison of the language name cannot be empty"
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

	return nil
}

type EditLanguageDTO struct {
	Name       string `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	LocaleName string `json:"locale_name,omitempty" form:"locale_name,omitempty" bson:"locale_name,omitempty"`
}
