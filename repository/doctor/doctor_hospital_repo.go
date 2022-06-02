package doctorrepo

import "go.mongodb.org/mongo-driver/mongo"

type DoctorHospitalRepository struct {
	colln *mongo.Collection
}

func (r *DoctorHospitalRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}
