package consultationapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
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
	return func(ctx *gin.Context) {

		file, err := ctx.FormFile("icon")
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		fmt.Println(file)

	}
}

func (a *ConsultationAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *ConsultationAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
