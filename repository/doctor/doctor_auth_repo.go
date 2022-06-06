package doctorrepo

import (
	"context"
	"fmt"
	"time"

	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorAuthRepository struct {
	colln *mongo.Collection
}

func (r *DoctorAuthRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	utils.CreateIndex(colln, bson.D{
		{Key: "phone.number", Value: 1},
		{Key: "phone.code", Value: 1}},
		"Phone", true)

	utils.CreateIndex(colln, bson.D{
		{Key: "email", Value: 1}},
		"Email", true)
}

func (r *DoctorAuthRepository) CreateOne(dto *doctordto.AddNewDoctorDTO) (*doctormodel.DoctorEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc := doctormodel.DoctorEntity{
		ID:                primitive.NewObjectID(),
		Name:              dto.Name,
		Email:             dto.Email,
		Phone:             &interfaces.PhoneType{Code: dto.PhoneCode, Number: dto.PhoneNumber},
		Type:              dto.ConsultationType,
		ProfessionalTitle: dto.ProfessionalTitle,
		Experience:        dto.Experience,
		RefreshToken:      "",
	}

	_, err := r.colln.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil

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

func (r *DoctorAuthRepository) FindById(id *primitive.ObjectID) interface{} {

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, _ := r.colln.Aggregate(ctx, pipeLine)

	var result []doctormodel.DoctorEntity

	if err := res.All(ctx, &result); err != nil {
		fmt.Println(err)
		return err.Error()
	}

	if len(result) == 1 {
		return result[0]
	}

	return nil

}
