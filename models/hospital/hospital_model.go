package hospitalmodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HospitalEntity struct {
	ID       primitive.ObjectID                 `json:"id" bson:"_id"`
	Name     string                             `json:"name" bson:"name"`
	City     string                             `json:"city" bson:"city"`
	Country  string                             `json:"country" bson:"country"`
	Address  string                             `json:"address" bson:"address"`
	Location *interfaces.MongoPointLocationType `json:"location,omitempty" bson:"location,omitempty"`
}
