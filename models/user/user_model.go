package usermodel

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserEntity struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id"`
	Name          string               `json:"name" bson:"name"`
	Email         string               `json:"email" bson:"email"`
	PhoneNumber   interfaces.PhoneType `json:"phone_number" bson:"phone_number"`
	DateOfBirth   string               `json:"date_of_birth" bson:"date_of_birth"`
	Gender        string               `json:"gender" bson:"gender"`
	RefreshToken  string               `json:"-" bson:"refresh_token"`
	EmailVerified bool                 `json:"email_verified" bson:"email_verified"`
}

func (u *UserEntity) GetAccessAndRefreshToken() (string, string, error) {
	access, err := tokens.GenerateTokenWithExpiryTimeAndType(u.ID.Hex(), 100, "access", "user")
	refresh, err := tokens.GenerateTokenWithExpiryTimeAndType(u.ID.Hex(), 100, "refresh", "user")

	return access, refresh, err

}
