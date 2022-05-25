package adminrepository

import (
	"context"
	"errors"
	"time"

	admindto "github.com/praveennagaraj97/online-consultation/dto/admin"
	adminmodel "github.com/praveennagaraj97/online-consultation/models/admin"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	colln *mongo.Collection
}

func (r *AdminRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{{Key: "email", Value: 1}}, "Email", true)
	utils.CreateIndex(colln, bson.D{{Key: "user_name", Value: 1}}, "User Name", true)

}

func (r *AdminRepository) CreateOne(payload *admindto.AddNewAdminDTO) (*adminmodel.AdminEntity, error) {

	if exists := r.checkIfUserExistsByEmailOrUserName(payload.Email, payload.UserName); exists {
		return nil, errors.New("User with given credentials already exist")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc := adminmodel.AdminEntity{
		ID:           primitive.NewObjectID(),
		Name:         payload.Name,
		Role:         payload.Role,
		Email:        payload.Email,
		UserName:     payload.UserName,
		Password:     payload.Password,
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		RefreshToken: "",
	}

	doc.EncodePassword()

	if _, err := r.colln.InsertOne(ctx, doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *AdminRepository) checkIfUserExistsByEmailOrUserName(email, userName string) bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.D{{Key: "$or", Value: bson.A{bson.M{"email": email}, bson.M{"user_name": userName}}}}

	count, err := r.colln.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *AdminRepository) FindByUserName(name string) (*adminmodel.AdminEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"user_name": name})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var result adminmodel.AdminEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *AdminRepository) FindByEmail(email string) (*adminmodel.AdminEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"email": email})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var result adminmodel.AdminEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *AdminRepository) UpdateById(id *primitive.ObjectID, payload *admindto.UpdateAdminDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": payload}); err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) UpdateRefreshToken(id *primitive.ObjectID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := r.colln.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: bson.M{"refresh_token": token}}}); err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) FindById(id *primitive.ObjectID) (*adminmodel.AdminEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"_id": id})

	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any user")
	}

	var result adminmodel.AdminEntity

	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r *AdminRepository) DeleteById(id *primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	return nil

}
