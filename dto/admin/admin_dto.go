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
