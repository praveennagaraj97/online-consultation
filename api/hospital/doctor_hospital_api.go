package hospitalapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
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
		var payload doctordto.AddDoctorHospitalDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		doc := doctormodel.DoctorHospitalEntity{
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

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*doctormodel.DoctorHospitalEntity]{
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

		ctx.JSON(http.StatusOK, serialize.DataResponse[*doctormodel.DoctorHospitalEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Hospital details retrieved successfully",
			},
		})

	}
}
