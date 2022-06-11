package doctorapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
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

		var payload doctordto.AddNewAppointmentSlotSetDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		fmt.Println(doctorId)

	}
}
