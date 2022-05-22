package adminapi

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	adminrepository "github.com/praveennagaraj97/online-consultation/repository/admin"
)

type AdminAPI struct {
	appConf   *app.ApplicationConfig
	adminRepo *adminrepository.AdminRepository
}

func (a *AdminAPI) Initailize(conf *app.ApplicationConfig, adminRepo *adminrepository.AdminRepository) {
	a.appConf = conf
	a.adminRepo = adminRepo
}

// Create new user
func (a *AdminAPI) AddNewAdmin(role constants.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) ChangeRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *AdminAPI) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
