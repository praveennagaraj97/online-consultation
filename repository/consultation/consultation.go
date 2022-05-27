package consultationrepository

import (
	"context"
	"time"

	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConsultationRepository struct {
	colln *mongo.Collection
}

func (r *ConsultationRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *ConsultationRepository) CreateOne(payload *consultationmodel.ConsultationEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, payload)

	if err != nil {
		return err
	}

	return nil
}
