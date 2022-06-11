package appointmentslotsrepo

import "go.mongodb.org/mongo-driver/mongo"

type AppointmentSlotsRepository struct {
	colln *mongo.Collection
}

func (r *AppointmentSlotsRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}
