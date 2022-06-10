package consultationdto

import (
	"fmt"
	"net/http"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

type AddConsultationDTO struct {
	Title       string  `json:"title" form:"title"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price"`
	ActionName  string  `json:"action_name" form:"action_name"`
	Type        string  `json:"type" form:"type"`
	IconWidth   uint64  `json:"icon_width" form:"icon_width"`
	IconHeight  uint64  `json:"icon_height" form:"icon_height"`
}

func (a *AddConsultationDTO) ValidateAddConsultationDTO() *serialize.ErrorResponse {
	errors := map[string]string{}

	if a.Title == "" {
		errors["title"] = "Consultation title cannot be empty"
	}

	if a.Description == "" {
		errors["description"] = "Consultation description cannot be empty"
	}

	if a.Price <= 0 {
		errors["price"] = "Price shoule be greater than zero"
	}

	if a.Type == "" || a.Type != consultationmodel.Instant && a.Type != consultationmodel.Schedule {
		errors["type"] = fmt.Sprintf("Type should be either %s or %s", consultationmodel.Instant, consultationmodel.Schedule)
	}

	if a.IconWidth <= 0 {
		errors["icon_width"] = "Icon width should be greater than zero"
	}

	if a.IconHeight <= 0 {
		errors["icon_height"] = "Icon height should be greater than zero"
	}

	if a.ActionName == "" {
		errors["action_name"] = "Action button name cannot be empty"
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

type EditConsultationDTO struct {
	Title       string                `json:"title,omitempty" form:"title,omitempty" bson:"title,omitempty"`
	Description string                `json:"description,omitempty" form:"description,omitempty" bson:"description,omitempty"`
	Price       float64               `json:"price,omitempty" form:"price,omitempty" bson:"price,omitempty"`
	ActionName  string                `json:"action_name,omitempty" form:"action_name,omitempty" bson:"action_name,omitempty"`
	IconWidth   uint64                `json:"icon_width,omitempty" form:"icon_width,omitempty" bson:"icon_width,omitempty"`
	IconHeight  uint64                `json:"icon_height,omitempty" form:"icon_height,omitempty" bson:"icon_height,omitempty"`
	Icon        *interfaces.ImageType `json:"-" form:"-" bson:"icon,omitempty"`
}
