package doctorapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
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

		doc := doctormodel.DoctorEntity{
			ID:                primitive.NewObjectID(),
			Name:              payload.Name,
			Email:             payload.Email,
			Phone:             &interfaces.PhoneType{Code: payload.PhoneCode, Number: payload.PhoneNumber},
			TypeId:            payload.ConsultationType,
			ProfessionalTitle: payload.ProfessionalTitle,
			Experience:        payload.Experience,
			RefreshToken:      "",
		}

		multipartFile, _ := ctx.FormFile("profile_pic")
		if multipartFile != nil {
			var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
			defer close(ch)
			a.appConf.AwsUtils.UploadImageToS3(ctx, string(constants.DoctorProfilePic), doc.ID.Hex(), "profile_pic",
				payload.ProfilePicWidth, payload.ProfilePicHeight, ch)

			select {
			case value, ok := <-ch:
				if ok {
					if value.Err != nil {
						api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
						return
					} else {
						doc.ProfilePic = value.Result
					}
				}
			default:
			}

		}

		err := a.authRepo.CreateOne(&doc)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		if err := a.appConf.EmailClient.SendNoReplyMail([]string{doc.Email},
			"Welcome to Online Consultation", "new-doctor", "welcome",
			mailer.GetNewDoctorAddedTemplateData(doc.Name, doc.ProfessionalTitle, env.GetEnvVariable("CLIENT_DOCTOR_LOGIN_LINK"))); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		if multipartFile != nil {
			doc.ProfilePic.OriginalSrc = a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + doc.ProfilePic.OriginalImagePath
			doc.ProfilePic.BlurDataURL = a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + doc.ProfilePic.BlurImagePath
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*doctormodel.DoctorEntity]{
			Data: &doc,
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

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.authRepo.FindById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*doctormodel.DoctorEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Doctor details retrieved successfully",
			},
		})

	}
}

func (a *DoctorAPI) LinkHospitalToDoctor() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
