package appointmentslotsapi

import (
	appointmentslotmodel "github.com/praveennagaraj97/online-consultation/models/appointment_slot"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateSlotDocuments(dates []primitive.DateTime,
	doctorId *primitive.ObjectID, startTimes []primitive.DateTime) []interface{} {

	var docs []interface{}

	for i := 0; i < len(dates); i++ {
		for j := 0; j < len(startTimes); j++ {
			docs = append(docs, appointmentslotmodel.AppointmentSlotEntity{
				ID:          primitive.NewObjectID(),
				DoctorId:    doctorId,
				Date:        &dates[i],
				StartTime:   &startTimes[j],
				IsAvailable: true,
			})
		}
	}

	return docs

}
