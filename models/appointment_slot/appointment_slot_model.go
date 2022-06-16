package appointmentslotmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppointmentSlotEntity struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	DoctorId    *primitive.ObjectID `json:"-" bson:"doctor_id"`
	Date        *primitive.DateTime `json:"date" bson:"date"`
	StartTime   *primitive.DateTime `json:"start_time" bson:"start_time"`
	EndTime     *primitive.DateTime `json:"end_time" bson:"end_time"`
	IsAvailable bool                `json:"is_available" bson:"is_available"`
}
