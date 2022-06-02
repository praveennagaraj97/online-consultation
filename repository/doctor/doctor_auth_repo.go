package doctorrepo

import "go.mongodb.org/mongo-driver/mongo"

type DoctorAuthRepository struct {
	colln *mongo.Collection
}

func (r *DoctorAuthRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *DoctorAuthRepository) CreateOne() {

}
