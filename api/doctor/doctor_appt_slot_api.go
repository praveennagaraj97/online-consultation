package doctorapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *DoctorAPI) AddNewSlotSet() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var doctorId primitive.ObjectID

		id := ctx.Param("doctor_id")
		if id != "" {
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			doctorId = id
		} else {
			id, err := api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
			doctorId = *id
		}

		var payload doctordto.AddNewAppointmentSlotSetDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		count, err := a.apptSlotSetRepo.ExistingSetsCount(&doctorId)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		if count >= 12 {
			api.SendErrorResponse(ctx, "Slots limit reached, Please delete other slot and try again.", http.StatusNotAcceptable, nil)
			return
		}

		if payload.Title == "" {
			payload.Title = fmt.Sprintf("Slot %v", count+1)
		}

		doc := &doctormodel.AppointmentSlotSetEntity{
			ID:        primitive.NewObjectID(),
			DoctorId:  doctorId,
			Title:     payload.Title,
			SlotTime:  payload.SlotTimes,
			IsDefault: count == 0 || payload.IsDefault,
			AddedOn:   primitive.NewDateTimeFromTime(time.Now()),
		}

		if err := a.apptSlotSetRepo.CreateOne(doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*doctormodel.AppointmentSlotSetEntity]{
			Data: doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "New appointment slot set added successfully",
			},
		})

	}
}

func (a *DoctorAPI) GetSlotSetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Slot Id
		id := ctx.Param("id")
		slotId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var doctorId primitive.ObjectID

		// Doctor Id
		doctorIdHex := ctx.Param("doctor_id")
		if doctorIdHex != "" {
			id, err := primitive.ObjectIDFromHex(doctorIdHex)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			doctorId = id
		} else {
			id, err := api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
			doctorId = *id
		}

		res, err := a.apptSlotSetRepo.FindById(&doctorId, &slotId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*doctormodel.AppointmentSlotSetEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Retrieved slot set details successfully",
			},
		})
	}
}

func (a *DoctorAPI) UpdateSlotSetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Slot Id
		id := ctx.Param("id")
		slotId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var doctorId primitive.ObjectID

		// Doctor Id
		doctorIdHex := ctx.Param("doctor_id")
		if doctorIdHex != "" {
			id, err := primitive.ObjectIDFromHex(doctorIdHex)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			doctorId = id
		} else {
			id, err := api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
			doctorId = *id
		}

		var payload doctordto.UpdateAppointmentSlotSetDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if err := a.apptSlotSetRepo.UpdateById(&doctorId, &slotId, &payload); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}

func (a *DoctorAPI) DeleteSlotById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Slot Id
		id := ctx.Param("id")
		slotId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var doctorId primitive.ObjectID

		// Doctor Id
		doctorIdHex := ctx.Param("doctor_id")
		if doctorIdHex != "" {
			id, err := primitive.ObjectIDFromHex(doctorIdHex)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			doctorId = id
		} else {
			id, err := api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
			doctorId = *id
		}

		if err := a.apptSlotSetRepo.DeleteById(&doctorId, &slotId); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}
