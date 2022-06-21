package appointmentapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	appointmentdto "github.com/praveennagaraj97/online-consultation/dto/appointment"
	appointmentmodel "github.com/praveennagaraj97/online-consultation/models/appointment"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	"github.com/praveennagaraj97/online-consultation/pkg/scheduler"
	appointmentrepository "github.com/praveennagaraj97/online-consultation/repository/appointment"
	appointmentslotsrepo "github.com/praveennagaraj97/online-consultation/repository/appointment_slots"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentAPI struct {
	conf         *app.ApplicationConfig
	apptSlotRepo *appointmentslotsrepo.AppointmentSlotsRepository
	apptRepo     *appointmentrepository.AppointmentRepository
	cnsltRepo    *consultationrepository.ConsultationRepository
	rltvRepo     *userrepository.UserRelativesRepository
	scheduler    *scheduler.Scheduler
}

func (a *AppointmentAPI) Initialize(conf *app.ApplicationConfig,
	apptSlotRepo *appointmentslotsrepo.AppointmentSlotsRepository,
	apptRepo *appointmentrepository.AppointmentRepository,
	cnsltRepo *consultationrepository.ConsultationRepository,
	rltvRepo *userrepository.UserRelativesRepository) {

	a.conf = conf
	a.apptSlotRepo = apptSlotRepo
	a.apptRepo = apptRepo
	a.cnsltRepo = cnsltRepo
	a.rltvRepo = rltvRepo

	// Task Scheduler
	a.scheduler = &scheduler.Scheduler{}
	a.scheduler.Initialize(a.conf)

}

// Takes input and create payment intent
func (a *AppointmentAPI) BookAnScheduledAppointment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		var payload appointmentdto.BookScheduledAppointmentDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		// Get Slot Details
		apptSlotRes, err := a.apptSlotRepo.FindById(payload.DoctorId, payload.AppointmentSlotId)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		if !apptSlotRes.IsAvailable {
			api.SendErrorResponse(ctx, "We are sorry this slot has been booked", http.StatusConflict, nil)
			return
		}

		// if *apptSlotRes.Start <= primitive.NewDateTimeFromTime(time.Now()) {
		// 	api.SendErrorResponse(ctx, "You are trying to book a slot from past time", http.StatusUnprocessableEntity, nil)
		// 	return
		// }

		// Get Consultation Info
		consRes, err := a.cnsltRepo.FindByType(consultationmodel.Schedule)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Create Payment Info
		// ** Pending **
		//

		doc := appointmentmodel.AppointmentEntity{
			ID:              primitive.NewObjectID(),
			DoctorId:        payload.DoctorId,
			ConsultationId:  &consRes.ID,
			UserId:          userId,
			ConsultingFor:   payload.RelativeId,
			AppointmentSlot: payload.AppointmentSlotId,
			BookedDate:      primitive.NewDateTimeFromTime(time.Now()),
			Status:          appointmentmodel.Pending,
		}

		apptSheduleDoc := appointmentmodel.AppointmentScheduleTaskEntity{
			ID:            primitive.NewObjectID(),
			InvokeTime:    *apptSlotRes.Start,
			Type:          scheduler.ReminderTask,
			CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
			AppointmentId: &doc.ID,
		}

		if err := a.scheduler.NewSchedule(apptSheduleDoc.InvokeTime.Time(), scheduler.ReminderTask); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusCreated, map[string]interface{}{
			"res":             doc,
			"apptScheduleDoc": apptSheduleDoc,
		})
	}
}
