package consultationrepository

import "go.mongodb.org/mongo-driver/mongo"

type ConsultationRepository struct {
	colln *mongo.Collection
}

func (r *ConsultationRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}
