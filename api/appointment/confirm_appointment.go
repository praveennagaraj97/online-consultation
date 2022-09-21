package appointmentapi

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	appointmentmodel "github.com/praveennagaraj97/online-consultation/models/appointment"
	appointmentslotmodel "github.com/praveennagaraj97/online-consultation/models/appointment_slot"
	razorpaypayment "github.com/praveennagaraj97/online-consultation/pkg/razorpay"
	"github.com/praveennagaraj97/online-consultation/pkg/scheduler"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Once the Payment is confirmed add to scheduled reminder list and mark slot as booked
func (a *AppointmentAPI) ConfirmScheduledAppointmentFromWebhook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body razorpaypayment.RazorPayWebHook
		ctx.ShouldBind(&body)

		defer ctx.Request.Body.Close()

		switch body.Payload.Payment.Entity.Status {
		//  Handle Captured Event
		case "captured":
			details, err := razorpaypayment.GetOrderDetails(body.Payload.Payment.Entity.OrderId)
			if err != nil {
				// Save to error logs
				log.Default().Println(err.Error())
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}
			// For Scheduled Payment Capture
			if body.Payload.Payment.Entity.Notes.PayingFor == constants.ScheduledAppointment {
				a.confirmScheduledAppointment(ctx, details)
				return
			}

		case "failed":
			ctx.JSON(http.StatusNoContent, nil)

		//  Handle authorized Event
		case "authorized":
			ctx.JSON(http.StatusNoContent, nil)

		default:
			ctx.JSON(http.StatusNoContent, nil)
		}

	}
}

func (a *AppointmentAPI) confirmScheduledAppointment(
	ctx *gin.Context,
	details *razorpaypayment.RazorPayOrderDetails) {

	timeZone := utils.GetTimeZone(ctx)

	timeLoc, err := time.LoadLocation(timeZone)
	if err != nil {
		timeLoc = time.Now().Local().Location()
	}

	var refId = details.Notes.RefID

	appointmentSlotId, err := primitive.ObjectIDFromHex(refId)
	if err != nil {
		api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
			"reason": err.Error(),
		})
		return
	}

	// Get AppointmentSlot Details
	apptRes, err := a.apptRepo.FindById(&appointmentSlotId)
	if err != nil {
		api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
			"reason": err.Error(),
		})
		return
	}

	// If payment status has changed by webhook or cancelled by browser close.
	if apptRes.Status != appointmentmodel.Pending && apptRes.Status != appointmentmodel.Cancelled {
		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Appointment slot has been booked successfully",
		})
		return
	}

	// Change Slot Booked reason
	if err := a.apptSlotRepo.UpdateSlotAvailability(apptRes.AppointmentSlot, false, appointmentslotmodel.Confirmed); err != nil {
		api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
			"reason": err.Error(),
		})
		return
	}

	// Change Appointment Response
	if err := a.apptRepo.UpdateById(apptRes.UserId, &apptRes.ID, appointmentmodel.Upcoming); err != nil {
		api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
			"reason": err.Error(),
		})
		return
	}

	// Get Slot Info
	appSlotRes, err := a.apptSlotRepo.FindById(apptRes.AppointmentSlot)
	if err != nil {
		api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
			"reason": err.Error(),
		})
		return
	}

	invokeTime := primitive.NewDateTimeFromTime(appSlotRes.Start.Time().Add(-time.Minute * 5))

	// Create schedule reminder Docs
	apptSheduleDoc := appointmentmodel.AppointmentScheduleTaskEntity{
		ID:            primitive.NewObjectID(),
		InvokeTime:    &invokeTime,
		Date:          appSlotRes.Date,
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		AppointmentId: &apptRes.ID,
	}

	if err := a.apptReminderRepo.Create(&apptSheduleDoc); err != nil {
		a.apptRepo.DeleteById(apptRes.UserId, &apptRes.ID)
		api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
			"reason": err.Error(),
		})
		return
	}

	// Schedule if appointment is in current date.
	if apptSheduleDoc.InvokeTime.Time().Format("2006-01-02") == time.Now().Format("2006-01-02") {
		if err := a.conf.Scheduler.NewSchedule(apptSheduleDoc.InvokeTime.Time(), scheduler.AppointmentReminderTask); err != nil {
			a.apptRepo.DeleteById(apptRes.UserId, &apptRes.ID)
			a.apptReminderRepo.DeleteById(&apptSheduleDoc.ID)
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}
	}

	var ch chan bool = make(chan bool, 1)
	go a.sendEmailAndSMSForBooking(ch, apptRes.UserId, appSlotRes.Start.Time(), timeLoc)

	ctx.JSON(http.StatusOK, serialize.Response{
		StatusCode: http.StatusOK,
		Message:    "Appointment slot has been booked successfully",
	})
}
