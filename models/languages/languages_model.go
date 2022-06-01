package languagesmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type LanguageEntity struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	LocaleName string             `json:"locale_name" bson:"locale_name"`
}
