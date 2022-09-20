package appointmentapi

import (
	"fmt"
	"log"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	twiliopkg "github.com/praveennagaraj97/online-consultation/pkg/sms/twilio"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *AppointmentAPI) sendEmailAndSMSForBooking(ch chan bool, userId *primitive.ObjectID, invokeTime time.Time, timezone *time.Location) {
	// Get User Details
	userRes, err := a.userRepo.FindById(userId)
	if err != nil {
		log.Default().Println("User", err.Error())
		ch <- true
		return
	}

	message := fmt.Sprintf("Your appointment has been scheduled for %s", invokeTime.In(timezone).Format(time.RFC1123))

	// Send Confirmation Mail and SMS to user
	if err := a.conf.EmailClient.SendNoReplyMail([]string{userRes.Email}, message, "appointment", "appointment",
		mailer.GetScheduledAppointmentBookingTemplateData()); err != nil {
		log.Default().Println(err.Error())
	}

	if err := twiliopkg.SendMessage(&interfaces.SMSType{
		Message: message,
		To:      fmt.Sprintf("%s%s", userRes.PhoneNumber.Code, userRes.PhoneNumber.Number),
	}); err != nil {
		log.Default().Println(err.Error())
	}

	ch <- true

}
