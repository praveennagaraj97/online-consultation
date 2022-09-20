package appointmentmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentStatus string

const (
	Cancelled AppointmentStatus = "cancelled"
	Completed AppointmentStatus = "completed"
	Upcoming  AppointmentStatus = "upcoming"
	Pending   AppointmentStatus = "pending" // Payment pending.
	OnGoing   AppointmentStatus = "ongoing"
)

type AppointmentEntity struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id"`
	DoctorId        *primitive.ObjectID `json:"doctor_id" bson:"doctor_id"`
	ConsultationId  *primitive.ObjectID `json:"consultation_id" bson:"consultation_id"`
	UserId          *primitive.ObjectID `json:"user_id" bson:"user_id"`
	ConsultingFor   *primitive.ObjectID `json:"consulting_for" bson:"consulting_for"`
	AppointmentSlot *primitive.ObjectID `json:"appointment_slot" bson:"appointment_slot"`
	BookedDate      primitive.DateTime  `json:"booked_date" bson:"booked_date"`
	Status          AppointmentStatus   `json:"status" bson:"status"`
	RelativeId      *primitive.ObjectID `json:"relative_id" bson:"relative_id"`
}
