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
	razorpaypayment "github.com/praveennagaraj97/online-consultation/pkg/razorpay"
	"github.com/praveennagaraj97/online-consultation/pkg/scheduler"
	appointmentrepository "github.com/praveennagaraj97/online-consultation/repository/appointment"
	appointmentslotsrepo "github.com/praveennagaraj97/online-consultation/repository/appointment_slots"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentAPI struct {
	conf             *app.ApplicationConfig
	apptSlotRepo     *appointmentslotsrepo.AppointmentSlotsRepository
	apptRepo         *appointmentrepository.AppointmentRepository
	cnsltRepo        *consultationrepository.ConsultationRepository
	rltvRepo         *userrepository.UserRelativesRepository
	userRepo         *userrepository.UserRepository
	apptReminderRepo *appointmentrepository.AppointmentScheduleReminderRepository
	scheduler        *scheduler.Scheduler
}

func (a *AppointmentAPI) Initialize(conf *app.ApplicationConfig,
	apptSlotRepo *appointmentslotsrepo.AppointmentSlotsRepository,
	apptRepo *appointmentrepository.AppointmentRepository,
	cnsltRepo *consultationrepository.ConsultationRepository,
	rltvRepo *userrepository.UserRelativesRepository,
	apptReminderRepo *appointmentrepository.AppointmentScheduleReminderRepository,
	userRepo *userrepository.UserRepository) {

	a.conf = conf
	a.apptSlotRepo = apptSlotRepo
	a.apptRepo = apptRepo
	a.cnsltRepo = cnsltRepo
	a.rltvRepo = rltvRepo
	a.apptReminderRepo = apptReminderRepo
	a.userRepo = userRepo

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
		apptSlotRes, err := a.apptSlotRepo.FindById(payload.AppointmentSlotId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		if !apptSlotRes.IsAvailable {
			api.SendErrorResponse(ctx, "We are sorry requested slot is not available", http.StatusNotFound, &map[string]string{
				"reason": apptSlotRes.Reason,
			})
			return
		}

		if *apptSlotRes.Start <= primitive.NewDateTimeFromTime(time.Now()) {
			api.SendErrorResponse(ctx, "You are trying to book a slot from past time", http.StatusUnprocessableEntity, nil)
			return
		}

		// Get Consultation Info
		consRes, err := a.cnsltRepo.FindByType(consultationmodel.Schedule)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusNotFound, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		doc := appointmentmodel.AppointmentEntity{
			ID:              primitive.NewObjectID(),
			DoctorId:        apptSlotRes.DoctorId,
			ConsultationId:  &consRes.ID,
			UserId:          userId,
			ConsultingFor:   payload.RelativeId,
			AppointmentSlot: payload.AppointmentSlotId,
			BookedDate:      primitive.NewDateTimeFromTime(time.Now()),
			Status:          appointmentmodel.Pending,
		}

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Get User Details
		userRes, err := a.userRepo.FindById(userId)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		if err := a.apptRepo.Create(&doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Block the slot
		// if err = a.apptSlotRepo.UpdateSlotAvailability(payload.AppointmentSlotId, false, "Blocked for booking"); err != nil {
		// 	a.apptRepo.DeleteById(userId, &doc.ID)
		// 	api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
		// 		"reason": err.Error(),
		// 	})
		// 	return
		// }

		// paymentDescription := fmt.Sprintf("Pay Rs.%v for your appointment booking", consRes.Price)
		// email := userRes.Email

		// Create Payment Channel
		paymentRes, err := razorpaypayment.CreateOrder(consRes.Price, "INR", doc.ID.Hex())

		if err != nil {
			// Release blocked slot
			a.apptSlotRepo.UpdateSlotAvailability(payload.AppointmentSlotId, true, "")

			// Delete appointment if payment creation fails
			a.apptRepo.DeleteById(userId, &doc.ID)
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"user":       userRes,
			"paymentRes": paymentRes,
		})

		// ctx.JSON(http.StatusCreated, serialize.DataResponse[*paymentmodel.StripePaymentModel]{
		// 	Data: paymentRes,
		// 	Response: serialize.Response{
		// 		StatusCode: http.StatusCreated,
		// 		Message:    "Appointment slot has been blocked, pay to confirm the slot",
		// 	},
		// })

	}
}

// Once the Payment is confirmed add to scheduled reminder list and mark slot as booked
func (a *AppointmentAPI) ConfirmScheduledAppointment() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// apptSheduleDoc := appointmentmodel.AppointmentScheduleTaskEntity{
		// 	ID:            primitive.NewObjectID(),
		// 	InvokeTime:    *apptSlotRes.Start,
		// 	CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		// 	AppointmentId: &doc.ID,
		// }

		// 		if err := a.apptReminderRepo.Create(&apptSheduleDoc); err != nil {
		// 	a.apptRepo.DeleteById(userId, &doc.ID)
		// 	api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
		// 		"reason": err.Error(),
		// 	})
		// 	return
		// }

		// // Schedule if appointment is in current date.
		// if apptSheduleDoc.InvokeTime.Time().Format("2006-01-02") == time.Now().Format("2006-01-02") {
		// 	if err := a.scheduler.NewSchedule(apptSheduleDoc.InvokeTime.Time(), scheduler.AppointmentReminderTask); err != nil {
		// 		a.apptRepo.DeleteById(userId, &doc.ID)
		// 		a.apptReminderRepo.DeleteById(&apptSheduleDoc.ID)
		// 		api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
		// 		return
		// 	}
		// }
	}
}
