package appointmentapi

import (
	"fmt"
	"log"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *AppointmentAPI) sendEmailAndSMSForBooking(ch chan bool, userId *primitive.ObjectID, apptTime time.Time, timezone *time.Location) {
	// Get User Details
	userRes, err := a.userRepo.FindById(userId)
	if err != nil {
		log.Default().Println("User", err.Error())
		ch <- true
		return
	}

	message := fmt.Sprintf("Your appointment has been scheduled on %s",
		apptTime.In(timezone).Format("Monday, January 2 2006 at 15:04:05 MST -07:00"))

	// Send Confirmation Mail and SMS to user
	if err := a.conf.EmailClient.SendNoReplyMail([]string{userRes.Email}, message, "appointment", "appointment",
		mailer.GetScheduledAppointmentBookingTemplateData()); err != nil {
		log.Default().Println(err.Error())
	}

	if _, err := a.conf.AwsUtils.SendTextSMS(&interfaces.SMSType{
		Message: message,
		To:      &userRes.PhoneNumber,
	}); err != nil {
		log.Default().Println(err.Error())
	}

	ch <- true

}
