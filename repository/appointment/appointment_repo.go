package appointmentrepository

import "go.mongodb.org/mongo-driver/mongo"

type AppointmentRepository struct {
	colln *mongo.Collection
}

func (r *AppointmentRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}
