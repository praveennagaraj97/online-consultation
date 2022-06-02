package doctorrepo

import "go.mongodb.org/mongo-driver/mongo"

type DoctorQualificationRepository struct {
	colln *mongo.Collection
}

func (r *DoctorQualificationRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}
