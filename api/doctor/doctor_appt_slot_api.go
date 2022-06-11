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

		} else {
			id, err := api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}

			doctorId = *id
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

		if payload.Title == "" {
			payload.Title = fmt.Sprintf("Slot %v", count+1)
		}

		doc := &doctormodel.AppointmentSlotSetEntity{
			ID:        primitive.NewObjectID(),
			DoctorId:  doctorId,
			Title:     payload.Title,
			SlotTime:  payload.SlotTimes,
			IsDefault: payload.IsDefault,
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
