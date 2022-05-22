package admindto

type AddNewAdminDTO struct {
	Name     string   `json:"name" form:"name"`
	Email    string   `json:"email" form:"email"`
	Password string   `json:"password" form:"password"`
	Role     []string `json:"role" form:"role"`
}
