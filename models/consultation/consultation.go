package consultationmodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsultationType string

const (
	Instant  = "Instant"
	Schedule = "Schedule"
)

type ConsultationEntity struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id"`
	Title       string                `json:"title" bson:"title"`
	Icon        *interfaces.ImageType `json:"icon" bson:"icon"`
	Description string                `json:"description" bson:"description"`
	Price       float64               `json:"price" bson:"price"`
	ActionName  string                `json:"action_name" bson:"action_name"`
	Type        ConsultationType      `json:"type" bson:"type"`
}
