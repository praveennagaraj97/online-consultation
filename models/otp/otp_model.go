package otpmodel

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type OneTimePasswordEntity struct {
	ID             primitive.ObjectID `json:"-" bson:"_id"`
	VerificationID string             `json:"verification_id" bson:"verification_id"`
	VerifyCode     string             `json:"-" bson:"verify_code"`
	ExpiryTime     primitive.DateTime `json:"-" bson:"expiry_time"`
	CreatedAt      primitive.DateTime `json:"created_at" bson:"created_at"`
	Attempts       uint8              `json:"-" bson:"attempts"`
	Verified       bool               `json:"-" bson:"verified"`
}

func (otp *OneTimePasswordEntity) Init(verifyCode *string, phoneNumber *interfaces.PhoneType) error {

	otp.ID = primitive.NewObjectID()

	encryptedCode, err := bcrypt.GenerateFromPassword([]byte(*verifyCode), 12)

	if err != nil {
		return err
	}

	otp.VerifyCode = string(encryptedCode)
	otp.VerificationID = base64.StdEncoding.
		EncodeToString([]byte(fmt.Sprintf("_id=%s&phone_code=%s&phone_number=%s",
			otp.ID.Hex(),
			phoneNumber.Code,
			phoneNumber.Number)))

	otp.Attempts = 0
	otp.ExpiryTime = primitive.NewDateTimeFromTime(time.Now().Add(time.Minute * 1))
	otp.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	return nil
}

func (otp *OneTimePasswordEntity) DecodeVerificationCode(enteredCode string) error {
	return bcrypt.CompareHashAndPassword([]byte(otp.VerifyCode), []byte(enteredCode))
}
