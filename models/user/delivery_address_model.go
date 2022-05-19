package usermodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDeliveryAddressEntity struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Address     string               `json:"address" bson:"address"`
	State       primitive.ObjectID   `json:"state" bson:"state"`
	Locality    string               `json:"locality" bson:"locality"`
	PinCode     string               `json:"pincode" bson:"pincode"`
	PhoneNumber interfaces.PhoneType `json:"phone" bson:"phone"`
	IsDefault   bool                 `json:"is_default" bson:"is_default"`
}
