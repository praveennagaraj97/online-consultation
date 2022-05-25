package admindto

import "github.com/praveennagaraj97/online-consultation/constants"

type AddNewAdminDTO struct {
	Name     string             `json:"name" form:"name"`
	UserName string             `json:"user_name" form:"user_name"`
	Email    string             `json:"email" form:"email"`
	Password string             `json:"password" form:"password"`
	Role     constants.UserType `json:"-" form:"-"`
}

type LoginDTO struct {
	UserName string `json:"user_name" form:"user_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UpdateAdminDTO struct {
	Name     string `json:"name,omitempty" form:"name,omitempty" bson:"name,omitempty"`
	Password string `bson:"password,omitempty"`
}

type UpdatePasswordDTO struct {
	Password        string `json:"password" form:"password"`
	NewPassword     string `json:"new_password" form:"new_password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

type ForgorPasswordDTO struct {
	Email string `json:"email" form:"email"`
}

type AdminIdDTO struct {
	ID string `json:"id" form:"id"`
}

type ResetPasswordDTO struct {
	NewPassword     string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}
