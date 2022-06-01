package specialitymodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpecialityEntity struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id"`
	Title       string                `json:"title" bson:"title"`
	Slug        string                `json:"slug" bson:"slug"`
	Description string                `json:"description,omitempty" bson:"description,omitempty"`
	Thumbnail   *interfaces.ImageType `json:"thumbnail" bson:"thumbnail"`
}
