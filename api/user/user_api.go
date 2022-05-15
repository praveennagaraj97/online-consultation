package userapi

import (
	"github.com/praveennagaraj97/online-consultation/app"

	onetimepasswordrepository "github.com/praveennagaraj97/online-consultation/repository/onetimepassword"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
)

type UserAPI struct {
	appConf  *app.ApplicationConfig
	userRepo *userrepository.UserRepository
	otpRepo  *onetimepasswordrepository.OneTimePasswordRepository
}

func (a *UserAPI) Initialize(appConf *app.ApplicationConfig, repo *userrepository.UserRepository, otpRepo *onetimepasswordrepository.OneTimePasswordRepository) {
	a.appConf = appConf
	a.userRepo = repo
	a.otpRepo = otpRepo

}
