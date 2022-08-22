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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LanguageAPI struct {
	lngRepo *languagerepo.LanguageRepository
	appConf *app.ApplicationConfig
}

func (a *LanguageAPI) Initialize(lngrepo *languagerepo.LanguageRepository, conf *app.ApplicationConfig) {
	a.appConf = conf
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
		if exist := a.lngRepo.CheckIfExistByName(payload.Name, payload.LocaleName); exist {
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
	return func(ctx *gin.Context) {

		pgOpts := api.ParsePaginationOptions(ctx, "languages_list")
		sortOpts := map[string]int8{"_id": -1}
		filterOptions := api.ParseFilterByOptions(ctx)
		keySortById := "$lt"

		res, err := a.lngRepo.Find(pgOpts, &sortOpts, filterOptions, keySortById)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.lngRepo.GetDocumentsCount(filterOptions)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId, "languages_list")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]languagesmodel.LanguageEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			PaginatedData: serialize.PaginatedData[[]languagesmodel.LanguageEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of languages retrieved successfully",
				},
			},
		})

	}
}

func (a *LanguageAPI) GetLanguageById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.lngRepo.FindById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*languagesmodel.LanguageEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Language details retrieved successfully",
			},
		})

	}
}

func (a *LanguageAPI) DeleteLanguageById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		err = a.lngRepo.DeleteById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (a *LanguageAPI) UpdateLanguageById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var payload languagedto.EditLanguageDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if payload.Name == "" && payload.LocaleName == "" {
			api.SendErrorResponse(ctx, "Formdata is empty", http.StatusUnprocessableEntity, nil)
			return
		}

		err = a.lngRepo.UpdateById(&objectId, &payload)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}
