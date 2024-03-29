package appointmentslotmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookingStatusMessages string

const (
	PaymentProcessing BookingStatusMessages = "Blocked for payment processing"
	Confirmed         BookingStatusMessages = "Slot has been booked"
	Released          BookingStatusMessages = ""
)

type AppointmentSlotEntity struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id"`
	DoctorId    *primitive.ObjectID   `json:"-" bson:"doctor_id"`
	Date        *primitive.DateTime   `json:"date" bson:"date"`
	Start       *primitive.DateTime   `json:"start" bson:"start"`
	End         *primitive.DateTime   `json:"end" bson:"end"`
	IsAvailable bool                  `json:"is_available" bson:"is_available"`
	Reason      BookingStatusMessages `json:"reason_of_unavailablity,omitempty" bson:"reason_of_unavailablity"`

	// Blocked for booking
	SlotReleaseAt *primitive.DateTime `json:"slot_release_at,omitempty" bson:"slot_release_at"`
}
