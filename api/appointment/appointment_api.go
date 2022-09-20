package appointmentapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	appointmentdto "github.com/praveennagaraj97/online-consultation/dto/appointment"
	appointmentmodel "github.com/praveennagaraj97/online-consultation/models/appointment"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	razorpaypayment "github.com/praveennagaraj97/online-consultation/pkg/razorpay"
	appointmentrepository "github.com/praveennagaraj97/online-consultation/repository/appointment"
	appointmentslotsrepo "github.com/praveennagaraj97/online-consultation/repository/appointment_slots"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
	"github.com/praveennagaraj97/online-consultation/serialize"
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
	a.conf.Scheduler.InitializeAppointmentRemainderPersistRepo(apptReminderRepo)

}

// Takes input and create payment intent
func (a *AppointmentAPI) BookScheduledAppointment() gin.HandlerFunc {
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

		// Check if relative belongs to user booking.
		if payload.RelativeId != nil {
			if _, err := a.rltvRepo.FindById(userId, payload.RelativeId); err != nil {
				api.SendErrorResponse(ctx, "Provided relative ID is not valid", http.StatusUnprocessableEntity, nil)
				return
			}
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
			RelativeId:      payload.RelativeId,
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

		// Create appointment booking history for user.
		if err := a.apptRepo.Create(&doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Block the slot
		if err = a.apptSlotRepo.UpdateSlotAvailability(payload.AppointmentSlotId, false, "Blocked for booking"); err != nil {
			a.apptRepo.DeleteById(userId, &doc.ID)
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Discount Check

		orderData := razorpaypayment.CreateRazorPayOrder{
			Amount:         uint64(consRes.Price - consRes.Discount),
			Currency:       "INR",
			Receipt:        doc.ID.Hex(),
			PartialPayment: false,
			Notes: struct {
				PayingFor   constants.PaymentFor
				ReferenceId string
			}{
				PayingFor:   constants.ScheduledAppointment,
				ReferenceId: doc.ID.Hex(),
			},
		}

		// Create Payment Channel
		orderId, err := razorpaypayment.CreateOrder(orderData)
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

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*razorpaypayment.RazorPayPaymentOutput]{
			Data: &razorpaypayment.RazorPayPaymentOutput{
				AppointmentId: doc.ID.Hex(),
				Amount:        orderData.Amount,
				Currency:      orderData.Currency,
				OrderId:       orderId,
				Prefill: razorpaypayment.PrefillData{
					Name:    userRes.Name,
					Email:   userRes.Email,
					Contact: fmt.Sprintf("%s %s", userRes.PhoneNumber.Code, userRes.PhoneNumber.Number),
				},
				Name:        "Online Consultation | Schedule Booking",
				Description: fmt.Sprintf("Pay Rs. %.2f for your appointment booking", float64(orderData.Amount)),
			},
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Appointment slot has been blocked, pay to confirm the slot",
			},
		})

	}
}

// Called from frontend - by default invoked by webhook. - Just a fallback api if webhook fails or delays.
func (a *AppointmentAPI) ConfirmScheduledAppointmentPaymentIntent() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// Release blocked slot if user cancels payment or leaves the browser (Recieved as beacon)
func (a *AppointmentAPI) CancelApponintmentBooking() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		id := ctx.Param("id")

		apptId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		// get Appointment details.
		apptRes, err := a.apptRepo.FindByIdAndUserId(&apptId, userId)
		if err != nil {
			api.SendErrorResponse(ctx, "Couldn't find any appointment matching", http.StatusNotFound, nil)
			return
		}

		if apptRes.Status == appointmentmodel.Cancelled {
			api.SendErrorResponse(ctx, "Appointment has been already cancelled", http.StatusNotAcceptable, nil)
			return
		}

		if apptRes.Status != appointmentmodel.Pending {
			api.SendErrorResponse(ctx, "This appointment cannot be cancelled", http.StatusNotAcceptable, nil)
			return
		}

		// Mark appointent entity as cancelled
		if err := a.apptRepo.UpdateById(userId, &apptId, appointmentmodel.Cancelled); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Release slot for booking
		if err := a.apptSlotRepo.UpdateSlotAvailability(apptRes.AppointmentSlot, true, ""); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}
