package consultationmodel

import (
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpecialityEntity struct {
	ID        primitive.ObjectID    `json:"_id" bson:"_id"`
	Name      string                `json:"name" bson:"name"`
	Order     int                   `json:"order" bson:"order"`
	Image     *interfaces.ImageType `json:"image" bson:"image"`
	CreatedAt primitive.DateTime    `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime    `json:"updated_at" bson:"updated_at"`
}

func (s *SpecialityEntity) Init() {
	s.ID = primitive.NewObjectID()
	s.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	s.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
}
