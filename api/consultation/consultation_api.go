package consultationapi

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	consultationdto "github.com/praveennagaraj97/online-consultation/dto/consultation"
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

		var payload consultationdto.AddConsultationDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errors := payload.ValidateAddConsultationDTO(); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, errors.StatusCode, errors.Errors)
			return
		}

		file, err := ctx.FormFile("icon")
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		fileType := file.Header.Get("Content-Type")

		if !strings.Contains(fileType, "image") {
			api.SendErrorResponse(ctx, "Provided file is not acceptable", http.StatusUnprocessableEntity, nil)
			return
		}

		// Read the file buffer.
		multiPartFile, err := file.Open()
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, nil)
			return
		}
		buffer, err := io.ReadAll(multiPartFile)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, nil)
			return
		}

		defer multiPartFile.Close()

		_, err = a.appConf.AwsUtils.UploadAsset(bytes.NewBuffer(buffer), string(constants.ConsultationIcon), file.Filename, &fileType)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		// image := interfaces.ImageType{
		// 	OriginalImagePath: fmt.Sprintf("%s/%s",constants.ConsultationIcon,file.Filename),
		// }

		// fmt.Println(image)

	}
}

func (a *ConsultationAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *ConsultationAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
