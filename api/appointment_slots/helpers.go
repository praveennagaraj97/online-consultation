package appointmentslotsapi

import (
	"time"

	appointmentslotmodel "github.com/praveennagaraj97/online-consultation/models/appointment_slot"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateSlotDocuments(dates []primitive.DateTime,
	doctorId *primitive.ObjectID, startTimes []primitive.DateTime) []interface{} {

	var docs []interface{}

	for i := 0; i < len(dates); i++ {
		for j := 0; j < len(startTimes); j++ {

			// Set the start time by combining with date and slot time
			date := dates[i].Time()
			startTime := startTimes[j].Time()
			t, _ := time.Parse("2006-01-02", date.Format("2006-01-02"))
			t = t.Add(time.Hour*time.Duration(startTime.Hour()) +
				time.Minute*time.Duration(startTime.Minute()) +
				time.Second*time.Duration(startTime.Second()))

			st := primitive.NewDateTimeFromTime(t)

			docs = append(docs, appointmentslotmodel.AppointmentSlotEntity{
				ID:          primitive.NewObjectID(),
				DoctorId:    doctorId,
				Date:        &dates[i],
				Start:       &st,
				IsAvailable: true,
			})
		}
	}

	return docs

}
