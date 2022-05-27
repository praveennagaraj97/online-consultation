package consultationapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	consultationdto "github.com/praveennagaraj97/online-consultation/dto/consultation"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		var payload consultationdto.AddConsultationDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errors := payload.ValidateAddConsultationDTO(); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, errors.StatusCode, errors.Errors)
			return
		}

		doc := &consultationmodel.ConsultationEntity{
			ID:          primitive.NewObjectID(),
			Title:       payload.Title,
			Description: payload.Description,
			Price:       payload.Price,
			ActionName:  payload.ActionName,
			Type:        payload.Type,
		}

		var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)

		a.appConf.AwsUtils.UploadImageToS3(ctx, string(constants.ConsultationIcon), doc.ID.Hex(), "icon", payload.IconWidth, payload.IconHeight, ch)

		select {
		case value, ok := <-ch:
			if ok {
				if value.Err != nil {
					api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
					return
				} else {
					doc.Icon = value.Result
				}
			}
		default:
		}

		if err := a.consultRepo.CreateOne(doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})

			a.appConf.AwsUtils.DeleteAsset(&doc.Icon.OriginalImagePath)
			a.appConf.AwsUtils.DeleteAsset(&doc.Icon.BlurImagePath)
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*consultationmodel.ConsultationEntity]{
			Data: doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Consultation type added successfully",
			},
		})

	}
}

func (a *ConsultationAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *ConsultationAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
