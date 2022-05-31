package specialitydto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/serialize"
)

type AddSpecialityDTO struct {
	Title           string `json:"title" form:"title"`
	Description     string `json:"description" form:"description"`
	ThumbnailWidth  uint64 `json:"thumbnail_width" form:"thumbnail_width"`
	ThumbnailHeight uint64 `json:"thumbnail_height" form:"thumbnail_height"`
}

func (s *AddSpecialityDTO) ValidateAddSpecialityDTO() *serialize.ErrorResponse {
	errrors := map[string]string{}

	if s.Title == "" {
		errrors["title"] = "Speciality title cannot be empty"
	}

	if s.ThumbnailHeight <= 170 {
		errrors["thumbnail_height"] = "thumbnail height should be atleast 170px"
	}

	if s.ThumbnailWidth <= 170 {
		errrors["thumbnail_width"] = "thumbnail width should be atleast 170px"
	}

	if len(errrors) > 0 {
		return &serialize.ErrorResponse{
			Errors: &errrors,
			Response: serialize.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Given data is invalid",
			},
		}
	}

	return nil

}
