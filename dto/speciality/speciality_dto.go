package specialitydto

import (
	"net/http"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"github.com/praveennagaraj97/online-consultation/utils"
)

type AddSpecialityDTO struct {
	Title           string `json:"title" form:"title"`
	Description     string `json:"description" form:"description"`
	ThumbnailWidth  uint64 `json:"thumbnail_width" form:"thumbnail_width"`
	ThumbnailHeight uint64 `json:"thumbnail_height" form:"thumbnail_height"`
	Slug            string `json:"slug" form:"slug"`
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

func (s *AddSpecialityDTO) GenerateSlug() {
	if s.Slug == "" {
		s.Slug = utils.Slugify(s.Title)
	}
}

type EditSpecialityDTO struct {
	Title           string                `json:"title,omitempty" form:"title,omitempty" bson:"title,omitempty"`
	Description     string                `json:"description,omitempty" form:"description,omitempty" bson:"description,omitempty"`
	ThumbnailWidth  uint64                `json:"thumbnail_width,omitempty" form:"thumbnail_width,omitempty" bson:"thumbnail_width,omitempty"`
	ThumbnailHeight uint64                `json:"thumbnail_height,omitempty" form:"thumbnail_height,omitempty" bson:"thumbnail_height,omitempty"`
	Slug            string                `json:"slug,omitempty" form:"slug,omitempty" bson:"slug,omitempty"`
	Thumbnail       *interfaces.ImageType `json:"-" form:"-" bson:"thumbnail,omitempty"`
}
