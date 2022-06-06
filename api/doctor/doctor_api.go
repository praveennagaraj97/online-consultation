package doctorapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	doctorrepo "github.com/praveennagaraj97/online-consultation/repository/doctor"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorAPI struct {
	authRepo *doctorrepo.DoctorAuthRepository
	appConf  *app.ApplicationConfig
}

func (a *DoctorAPI) Initialize(conf *app.ApplicationConfig, authRepo *doctorrepo.DoctorAuthRepository) {
	a.authRepo = authRepo
	a.appConf = conf

}

func (a *DoctorAPI) AddNewDoctor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload doctordto.AddNewDoctorDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		if exist := a.authRepo.CheckIfDoctorExistsByEmailOrPhone(payload.Email, interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}); exist {
			api.SendErrorResponse(ctx, "Doctor with given credentials already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.authRepo.CreateOne(&payload)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		if err := a.appConf.EmailClient.SendNoReplyMail([]string{res.Email},
			"Welcome to Online Consultation", "new-doctor", "welcome",
			mailer.GetNewDoctorAddedTemplateData(res.Name, res.ProfessionalTitle, env.GetEnvVariable("CLIENT_DOCTOR_LOGIN_LINK"))); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*doctormodel.DoctorEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "New doctor accounts has been added successfully",
			},
		})

	}
}

func (a *DoctorAPI) GetDoctorById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")

		objectId, _ := primitive.ObjectIDFromHex(id)

		res := a.authRepo.FindById(&objectId)

		ctx.JSON(200, map[string]interface{}{
			"result": res,
		})

	}
}
