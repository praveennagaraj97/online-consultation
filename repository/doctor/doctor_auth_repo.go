package doctorrepo

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorAuthRepository struct {
	colln         *mongo.Collection
	imageBasePath string
}

func (r *DoctorAuthRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	r.imageBasePath = env.GetEnvVariable("S3_ACCESS_BASEURL")

	utils.CreateIndex(colln, bson.D{
		{Key: "phone.number", Value: 1},
		{Key: "phone.code", Value: 1}},
		"Phone", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "email", Value: 1}},
		"Email", true)
}

func (r *DoctorAuthRepository) CreateOne(doc *doctormodel.DoctorEntity) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil

}

func (r *DoctorAuthRepository) CheckIfDoctorExistsByEmailOrPhone(email string, phone interfaces.PhoneType) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	emailFilter := bson.M{"email": email}
	phoneFilter := bson.M{"$and": bson.A{bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}}}

	filter := bson.M{"$or": bson.A{
		emailFilter,
		phoneFilter,
	}}

	count, err := r.colln.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *DoctorAuthRepository) FindById(id *primitive.ObjectID) (*doctormodel.DoctorEntity, error) {

	filterPipe := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}

	// Consultation ID Populate
	typeMatchPipe := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "consultation",
		"localField":   "type",
		"foreignField": "_id",
		"as":           "consultation_type",
	}}}
	unwindTypePipe := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$consultation_type",
		"preserveNullAndEmptyArrays": false,
	}}}
	setTypePipe := bson.D{{Key: "$set", Value: bson.M{"consultation_type": "$consultation_type.type"}}}

	pipeLine := mongo.Pipeline{
		filterPipe,
		typeMatchPipe,
		unwindTypePipe,
		setTypePipe,
	}

	// Add Prefix to image
	setImagePrefixPipe := bson.D{{Key: "$set",
		Value: bson.M{"profile_pic.image_src": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: bson.M{"$concat": bson.A{r.imageBasePath, "/", "$profile_pic.original_image_path"}}},
		}}}}}
	setBlurImagePrefixPipe := bson.D{{Key: "$set",
		Value: bson.M{"profile_pic.blur_data_url": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: bson.M{"$concat": bson.A{r.imageBasePath, "/", "$profile_pic.blur_image_path"}}},
		}}}}}
	resetNullImagePipe := bson.D{{Key: "$set", Value: bson.M{
		"profile_pic": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic.image_src", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: "$profile_pic"},
		}},
	}}}
	pipeLine = append(pipeLine, setImagePrefixPipe, setBlurImagePrefixPipe, resetNullImagePipe)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur, err := r.colln.Aggregate(ctx, pipeLine)
	if err != nil {
		return nil, errors.New("Couldn't find any doctor matching gived id")
	}

	var result []doctormodel.DoctorEntity

	defer cur.Close(context.TODO())

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 1 {
		return &result[0], nil
	}

	return nil, errors.New("Couldn't find any doctor matching gived id")

}
