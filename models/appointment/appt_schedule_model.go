package appointmentmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentScheduleTaskEntity struct {
	ID            primitive.ObjectID  `bson:"_id"`
	InvokeTime    *primitive.DateTime `bson:"invoke_time"`
	Date          *primitive.DateTime `bson:"date"`
	CreatedAt     primitive.DateTime  `bson:"created_at"`
	AppointmentId *primitive.ObjectID `bson:""`
}
