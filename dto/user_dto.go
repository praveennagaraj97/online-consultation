package userdto

type RegisterDTO struct {
	Name           string `json:"name" form:"name"`
	Email          string `json:"email" form:"email"`
	PhoneCode      string `json:"phone_code" form:"phone_code"`
	PhoneNumber    string `json:"phone_number" form:"phone_number"`
	DOB            string `json:"date_of_birth" form:"date_of_birth"`
	Gender         string `json:"gender" form:"gender"`
	VerificationId string `json:"verification_id" form:"verification_id"`
}

type VerifyCodeDTO struct {
	VerifyCode string `json:"verify_code" form:"verify_code"`
}

type UpdateUserDTO struct {
	Name          string `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	Email         string `json:"email,omitempty" form:"email,omitempty" bson:"email,omitempty"`
	PhoneCode     string `json:"phone_code,omitempty" form:"phone_code,omitempty" bson:"phone_code,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty" form:"phone_number,omitempty" bson:"phone_number,omitempty"`
	DOB           string `json:"date_of_birth,omitempty" form:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	Gender        string `json:"gender,omitempty" form:"gender,omitempty" bson:"gender,omitempty"`
	RefreshToken  string `json:"-" form:"-" bson:"refresh_token,omitempty"`
	EmailVerified bool   `json:"-" bson:"email_verified,omitempty"`
}

type SignInWithPhoneDTO struct {
	VerificationId string `json:"verify_code" form:"verify_code"`
	PhoneCode      string `json:"phone_code" form:"phone_code"`
	PhoneNumber    string `json:"phone_number" form:"phone_number"`
}

type SignInWithEmailLinkDTO struct {
	Email string `json:"email" form:"email"`
}

type RequestEmailVerifyDTO struct {
	Email      string `json:"email" form:"email"`
	RedirectTo string `json:"redirect_to" form:"redirect_to"`
}
