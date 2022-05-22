package adminapi

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/app"
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
func (a *AdminAPI) AddNewAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
