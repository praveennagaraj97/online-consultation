package additionalmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type JWTStatus struct {
	Expires *primitive.DateTime `json:"expires"`
	IsValid bool                `json:"is_valid"`
}
