package adminmodel

import (
	"time"

	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AdminEntity struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	UserName     string             `json:"user_name" bson:"user_name"`
	Role         constants.UserType `json:"role" bson:"role"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"-" bson:"password"`
	CreatedAt    primitive.DateTime `json:"joined_on" bson:"created_at"`
	RefreshToken string             `json:"-" bson:"refresh_token"`
}

func (a *AdminEntity) EncodePassword() error {
	passcode, err := bcrypt.GenerateFromPassword([]byte(a.Password), 12)
	if err != nil {
		return err
	}
	a.Password = string(passcode)

	return nil
}

func (a *AdminEntity) DecodePassword(passcode string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(passcode))
}

func (a *AdminEntity) GetAccessAndRefreshToken(acessExpires bool, userType string) (string, string, int, error) {

	var access, refresh string
	var err error
	var accessTime int = constants.CookieRefreshExpiryTime

	if acessExpires {
		accessTime = constants.CookieAccessExpiryTime
		access, err = tokens.GenerateTokenWithExpiryTimeAndType(a.ID.Hex(),
			time.Now().Add(time.Minute*constants.JWT_AccessTokenExpiry).Unix(), "access", userType)
	} else {
		access, err = tokens.GenerateNoExpiryTokenWithCustomType(a.ID.Hex(), "access", userType)

	}
	refresh, err = tokens.GenerateNoExpiryTokenWithCustomType(a.ID.Hex(), "refresh", userType)

	return access, refresh, accessTime, err
}
