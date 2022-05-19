package repository

import (
	"github.com/praveennagaraj97/online-consultation/db"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
	onetimepasswordrepository "github.com/praveennagaraj97/online-consultation/repository/onetimepassword"
	specialityrepo "github.com/praveennagaraj97/online-consultation/repository/speciality"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	userRepo                *userrepository.UserRepository
	userRelativeRepo        *userrepository.UserRelativesRepository
	userDeliveryAddressRepo *userrepository.UserDeliveryAddressRepository
	otpRepo                 *onetimepasswordrepository.OneTimePasswordRepository
	specialityRepo          *specialityrepo.SpecialityRepository
}

func (r *Repository) Initialize(mgoClient *mongo.Client, dbName string) {

	// One Time Password Repo
	r.otpRepo = &onetimepasswordrepository.OneTimePasswordRepository{}
	r.otpRepo.InitializeRepository(db.OpenCollection(mgoClient, dbName, "otp"))

	// User Repo
	r.userRepo = &userrepository.UserRepository{}
	r.userRepo.InitializeRepository(db.OpenCollection(mgoClient, dbName, "user"))

	// User relatives repo
	r.userRelativeRepo = &userrepository.UserRelativesRepository{}
	r.userRelativeRepo.InitializeRepository(db.OpenCollection(mgoClient, dbName, "user_relative"))

	// UserDelivery Address Repo
	r.userDeliveryAddressRepo = &userrepository.UserDeliveryAddressRepository{}
	r.userDeliveryAddressRepo.Initialize(db.OpenCollection(mgoClient, dbName, "user_delivery_address"))

	// Specialities Repo
	r.specialityRepo = &specialityrepo.SpecialityRepository{}
	r.specialityRepo.Init(db.OpenCollection(mgoClient, dbName, "speciality"))

	logger.PrintLog("Repositories initialized 📜")
}

func (r *Repository) GetUserRepository() *userrepository.UserRepository {
	return r.userRepo
}

func (r *Repository) GetUserRelativeRepository() *userrepository.UserRelativesRepository {
	return r.userRelativeRepo
}

func (r *Repository) GetUserDeliveryAddressRepository() *userrepository.UserDeliveryAddressRepository {
	return r.userDeliveryAddressRepo
}

func (r *Repository) GetOneTimePasswordRepository() *onetimepasswordrepository.OneTimePasswordRepository {
	return r.otpRepo
}

func (r *Repository) GetSpecialityRepository() *specialityrepo.SpecialityRepository {
	return r.specialityRepo
}
