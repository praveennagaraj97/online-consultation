package usermodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RelativeEntity struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Email       string               `json:"email" bson:"email"`
	Phone       interfaces.PhoneType `json:"phone" bson:"phone"`
	DateOfBirth string               `json:"date_of_birth" bson:"date_of_birth"`
	Gender      string               `json:"gender" bson:"gender"`
	Relation    string               `json:"relation" bson:"relation"`
	UserId      primitive.ObjectID   `json:"user_id" bson:"user_id"`
}
