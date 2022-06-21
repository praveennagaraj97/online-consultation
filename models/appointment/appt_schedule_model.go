package appointmentmodel

import (
	"github.com/praveennagaraj97/online-consultation/pkg/scheduler"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentScheduleTaskEntity struct {
	ID            primitive.ObjectID   `bson:"_id"`
	InvokeTime    primitive.DateTime   `bson:"invoke_time"`
	Type          scheduler.TasksTypes `bson:"type"`
	CreatedAt     primitive.DateTime   `bson:"created_at"`
	AppointmentId *primitive.ObjectID  `bson:""`
}
