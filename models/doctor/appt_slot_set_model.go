package doctormodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppointmentSlotSetEntity struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	DoctorId  primitive.ObjectID   `json:"doctor_id" bson:"doctor_id"`
	Title     string               `json:"title" bson:"title"`
	SlotTime  []primitive.DateTime `json:"slot_times" bson:"slot_times"`
	IsDefault bool                 `json:"is_default" bson:"is_default"`
	AddedOn   primitive.DateTime   `json:"added_on" bson:"added_on"`
}
