package doctorrepo

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *DoctorRepository) FindById(id *primitive.ObjectID) (*doctormodel.DoctorEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"_id": id})
	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any doctor with given id")
	}

	var result doctormodel.DoctorEntity
	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *DoctorRepository) FindByPhoneNumber(phone interfaces.PhoneType) (*doctormodel.DoctorEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	phoneFilter := bson.M{"$and": bson.A{bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}}}

	cur := r.colln.FindOne(ctx, phoneFilter)
	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any doctor with given number")
	}

	var result doctormodel.DoctorEntity
	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *DoctorRepository) FindByEmail(email string) (*doctormodel.DoctorEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur := r.colln.FindOne(ctx, bson.M{"email": email})
	if cur.Err() != nil {
		return nil, errors.New("Couldn't find any doctor with given email")
	}

	var result doctormodel.DoctorEntity
	if err := cur.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *DoctorRepository) CheckIfDoctorExistsByEmail(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	emailFilter := bson.M{"email": email}

	count, err := r.colln.CountDocuments(ctx, emailFilter)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *DoctorRepository) CheckIfDoctorExistsByPhone(phone interfaces.PhoneType) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	phoneFilter := bson.M{"$and": bson.A{bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}}}

	count, err := r.colln.CountDocuments(ctx, phoneFilter)
	if err != nil {
		return false
	}

	return count > 0
}
