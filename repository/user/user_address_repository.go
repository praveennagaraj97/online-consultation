package userrepository

import (
	"context"
	"fmt"
	"time"

	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDeliveryAddressRepository struct {
	colln *mongo.Collection
}

func (r *UserDeliveryAddressRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln
}

func (r *UserDeliveryAddressRepository) CreateOne() (*usermodel.UserDeliveryAddressEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println(ctx)

	return nil, nil

}
