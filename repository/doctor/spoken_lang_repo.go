package doctorrepo

import "go.mongodb.org/mongo-driver/mongo"

type DoctorSpokenLanguageRepository struct {
	colln *mongo.Collection
}

func (r *DoctorSpokenLanguageRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}
