package hospitalapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	hospitaldto "github.com/praveennagaraj97/online-consultation/dto/hospital"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	hospitalmodel "github.com/praveennagaraj97/online-consultation/models/hospital"
	hospitalrepo "github.com/praveennagaraj97/online-consultation/repository/hospital"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HospitalAPI struct {
	hspRepo *hospitalrepo.HospitalRepository
	appConf *app.ApplicationConfig
}

func (a *HospitalAPI) Initialize(conf *app.ApplicationConfig, hspRepo *hospitalrepo.HospitalRepository) {
	a.appConf = conf
	a.hspRepo = hspRepo
}

func (a *HospitalAPI) AddNewHospital() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload hospitaldto.AddHospitalDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		doc := hospitalmodel.HospitalEntity{
			ID:      primitive.NewObjectID(),
			Name:    payload.Name,
			City:    payload.City,
			Country: payload.Country,
			Address: payload.Address,
			Location: &interfaces.MongoPointLocationType{
				Type:        "Point",
				Coordinates: []float64{payload.Longitude, payload.Latitude},
			},
		}

		if exist := a.hspRepo.CheckIfExistByName(doc.Name); exist {
			api.SendErrorResponse(ctx, "Hospital with given name already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		if err := a.hspRepo.CreateOne(&doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*hospitalmodel.HospitalEntity]{
			Data: &doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Hospital has been added successfully",
			},
		})

	}
}

func (a *HospitalAPI) GetHospitalById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.hspRepo.FindById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*hospitalmodel.HospitalEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Hospital details retrieved successfully",
			},
		})

	}
}

func (a *HospitalAPI) UpdateHospitalById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var payload hospitaldto.EditHospitalDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if payload.Name != "" {
			res, _ := a.hspRepo.FindById(&objectId)
			if payload.Name != res.Name {
				if exists := a.hspRepo.CheckIfExistByName(payload.Name); exists {
					api.SendErrorResponse(ctx, "Hosiptal with given name already exist", http.StatusUnprocessableEntity, nil)
					return
				}
			}
		}

		err = a.hspRepo.UpdateById(&objectId, &payload)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}

func (a *HospitalAPI) DeleteById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		err = a.hspRepo.DeleteById(&objectId)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}

func (a *HospitalAPI) GetAllHospitals() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pgOpts := api.ParsePaginationOptions(ctx, "hospitals")
		fltsOpts := api.ParseFilterByOptions(ctx)
		sortOpts := map[string]int8{"_id": -1}
		keySortBy := "$lt"

		res, err := a.hspRepo.FindAll(pgOpts, fltsOpts, &sortOpts, keySortBy)

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
			docCount, err = a.hspRepo.GetDocumentsCount(fltsOpts)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, docCount, lastResId, "hospitals")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]hospitalmodel.HospitalEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			PaginatedData: serialize.PaginatedData[[]hospitalmodel.HospitalEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of hospitals retrieved",
				},
			},
		})

	}
}
