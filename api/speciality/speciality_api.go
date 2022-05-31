package specialityapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	specialitydto "github.com/praveennagaraj97/online-consultation/dto/speciality"
	specialitymodel "github.com/praveennagaraj97/online-consultation/models/speciality"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	specialityrepository "github.com/praveennagaraj97/online-consultation/repository/specialities"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpecialityAPI struct {
	conf    *app.ApplicationConfig
	splRepo *specialityrepository.SpecialitysRepository
}

func (a *SpecialityAPI) Initialize(conf *app.ApplicationConfig, splRepo *specialityrepository.SpecialitysRepository) {
	a.conf = conf
	a.splRepo = splRepo
}

func (a *SpecialityAPI) AddNewSpeciality() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload specialitydto.AddSpecialityDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if err := payload.ValidateAddSpecialityDTO(); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		doc := &specialitymodel.SpecialityEntity{
			ID:          primitive.NewObjectID(),
			Title:       payload.Title,
			Description: payload.Description,
		}

		var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
		defer close(ch)

		a.conf.AwsUtils.UploadImageToS3(ctx, string(constants.SpecialityThumbnail), doc.ID.Hex(),
			"thumbnail", payload.ThumbnailWidth, payload.ThumbnailHeight, ch)

		select {
		case value, ok := <-ch:
			if ok {
				if value.Err != nil {
					api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
					return
				} else {
					doc.Thumbnail = value.Result
				}
			}
		default:
		}

		if err := a.splRepo.CreateOne(doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})

			a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.OriginalImagePath)
			a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.BlurImagePath)
			return
		}

		doc.Thumbnail.OriginalSrc = fmt.Sprintf("%s/%s", a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL, doc.Thumbnail.OriginalImagePath)
		doc.Thumbnail.BlurDataURL = fmt.Sprintf("%s/%s", a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL, doc.Thumbnail.BlurImagePath)

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*specialitymodel.SpecialityEntity]{
			Data: doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Speciality added successfully",
			},
		})

	}
}

func (a *SpecialityAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
