package adminrepository

import (
	"context"
	"fmt"
	"time"

	adminmodel "github.com/praveennagaraj97/online-consultation/models/admin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	colln *mongo.Collection
}

func (r *AdminRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *AdminRepository) CreateOne() (*adminmodel.AdminEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println(ctx)

	return nil, nil
}
