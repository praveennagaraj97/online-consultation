package languageapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	languagedto "github.com/praveennagaraj97/online-consultation/dto/language"
	languagesmodel "github.com/praveennagaraj97/online-consultation/models/languages"
	languagerepo "github.com/praveennagaraj97/online-consultation/repository/language"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

type LanguageAPI struct {
	lngRepo *languagerepo.LanguageRepository
	appCong *app.ApplicationConfig
}

func (a *LanguageAPI) Initialize(lngrepo *languagerepo.LanguageRepository, conf *app.ApplicationConfig) {
	a.appCong = conf
	a.lngRepo = lngrepo

}

func (a *LanguageAPI) AddNewLanguage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload languagedto.AddLanguageDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if err := payload.Validate(); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		// Check If same language already exist
		if exist := a.lngRepo.ChechIfExistByName(payload.Name, payload.LocaleName); exist {
			api.SendErrorResponse(ctx, "Language with given name or locale name already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.lngRepo.CreateOne(&payload)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*languagesmodel.LanguageEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Language added successfully",
			},
		})

	}
}

func (a *LanguageAPI) GetAllLanguages() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *LanguageAPI) GetLanguageById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *LanguageAPI) UpdateLanguageById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *LanguageAPI) DeleteLanguageById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
