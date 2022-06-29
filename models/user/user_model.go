package usermodel

import (
	"time"

	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserEntity struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id"`
	Name          string               `json:"name" bson:"name"`
	Email         string               `json:"email" bson:"email"`
	PhoneNumber   interfaces.PhoneType `json:"phone_number" bson:"phone_number"`
	DateOfBirth   primitive.DateTime   `json:"date_of_birth" bson:"date_of_birth"`
	Gender        string               `json:"gender" bson:"gender"`
	RefreshToken  string               `json:"-" bson:"refresh_token"`
	EmailVerified bool                 `json:"email_verified" bson:"email_verified"`
}

func (u *UserEntity) GetAccessAndRefreshToken(acessExpires bool) (string, string, int, error) {

	var access, refresh string
	var err error
	var accessTime int = constants.CookieRefreshExpiryTime

	if acessExpires {
		accessTime = constants.CookieAccessExpiryTime
		access, err = tokens.GenerateTokenWithExpiryTimeAndType(u.ID.Hex(),
			time.Now().Local().Add(time.Minute*constants.JWT_AccessTokenExpiry).Unix(), "access", "user")
	} else {
		access, err = tokens.GenerateNoExpiryTokenWithCustomType(u.ID.Hex(), "access", "user")

	}
	refresh, err = tokens.GenerateNoExpiryTokenWithCustomType(u.ID.Hex(), "refresh", "user")

	return access, refresh, accessTime, err
}
