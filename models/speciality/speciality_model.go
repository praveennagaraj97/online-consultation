package specialitymodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpecialityEntity struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id"`
	Title       string                `json:"title" bson:"title"`
	Description string                `json:"description" bson:"description"`
	Thumbnail   *interfaces.ImageType `json:"thumbnail" bson:"thumbnail"`
}
