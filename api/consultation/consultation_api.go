package consultationapi

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/app"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
)

type ConsultationAPI struct {
	consultRepo *consultationrepository.ConsultationRepository
	appConf     *app.ApplicationConfig
}

func (a *ConsultationAPI) Initialize(appConf *app.ApplicationConfig, consultRepo *consultationrepository.ConsultationRepository) {
	a.appConf = appConf
	a.consultRepo = consultRepo
}

func (a *ConsultationAPI) AddNewConsultationType() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *ConsultationAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *ConsultationAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
